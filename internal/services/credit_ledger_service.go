package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrLedgerNotFound = errors.New("akun kasbon tidak ditemukan")

type CreditLedgerService interface {
	GetLedgers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetLedgerByID(ctx context.Context, id uuid.UUID) (*dto.CreditLedgerDTO, error)
	GetLedgerByUserID(ctx context.Context, userID uuid.UUID) (*dto.CreditLedgerDTO, error)
	CreateLedger(ctx context.Context, req dto.CreateCreditLedgerDTO) (*dto.CreditLedgerDTO, error)
	AdjustCreditLimit(ctx context.Context, id uuid.UUID, req dto.AdjustCreditLimitDTO) (*dto.CreditLedgerDTO, error)
	SettleCredit(ctx context.Context, id uuid.UUID, req dto.SettleCreditDTO) (*dto.CreditLedgerDTO, error)
	GetPaymentHistory(ctx context.Context, ledgerID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error)
}

type creditLedgerService struct {
	db            *gorm.DB
	ledgerRepo    repositories.CreditLedgerRepository
	creditPayRepo repositories.CreditPaymentRepository
	userRepo      repositories.UserRepository
	bus           eventbus.EventBus
}

func NewCreditLedgerService(
	db *gorm.DB,
	ledgerRepo repositories.CreditLedgerRepository,
	creditPayRepo repositories.CreditPaymentRepository,
	userRepo repositories.UserRepository,
	bus eventbus.EventBus,
) CreditLedgerService {
	return &creditLedgerService{
		db:            db,
		ledgerRepo:    ledgerRepo,
		creditPayRepo: creditPayRepo,
		userRepo:      userRepo,
		bus:           bus,
	}
}

func (s *creditLedgerService) GetLedgers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	ledgers, total, err := s.ledgerRepo.FindAll(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.CreditLedgerDTO, len(ledgers))
	for i := range ledgers {
		dtos[i] = mapLedgerToDTO(&ledgers[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *creditLedgerService) GetLedgerByID(ctx context.Context, id uuid.UUID) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrLedgerNotFound
	}
	res := mapLedgerToDTO(l)
	return &res, nil
}

func (s *creditLedgerService) GetLedgerByUserID(ctx context.Context, userID uuid.UUID) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrLedgerNotFound
	}
	res := mapLedgerToDTO(l)
	return &res, nil
}

func (s *creditLedgerService) CreateLedger(ctx context.Context, req dto.CreateCreditLedgerDTO) (*dto.CreditLedgerDTO, error) {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Role != constants.RoleB2BCustomer {
		return nil, errors.New("akun kasbon hanya bisa dibuat untuk pelanggan mitra toko (B2B)")
	}

	existing, err := s.ledgerRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("pelanggan ini sudah punya akun kasbon")
	}

	ledger := &models.CreditLedger{
		UserID:          req.UserID,
		CreditLimit:     req.CreditLimit,
		RemainingCredit: req.CreditLimit,
	}
	if err := s.ledgerRepo.Create(ctx, nil, ledger); err != nil {
		return nil, err
	}

	return s.GetLedgerByID(ctx, ledger.ID)
}

func (s *creditLedgerService) AdjustCreditLimit(ctx context.Context, id uuid.UUID, req dto.AdjustCreditLimitDTO) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrLedgerNotFound
	}

	l.CreditLimit = req.NewCreditLimit
	l.RemainingCredit = l.CreditLimit - l.UsedCredit
	if l.RemainingCredit < 0 {
		l.RemainingCredit = 0
	}
	if err := s.ledgerRepo.Update(ctx, nil, l); err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&models.B2BProfile{}).Where("user_id = ?", l.UserID).
		Update("approved_credit_limit", l.CreditLimit).Error; err != nil {
		return nil, err
	}

	return s.GetLedgerByID(ctx, id)
}

// SettleCredit records a kasbon repayment received by the store.
func (s *creditLedgerService) SettleCredit(ctx context.Context, id uuid.UUID, req dto.SettleCreditDTO) (*dto.CreditLedgerDTO, error) {
	var userID uuid.UUID
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ledger models.CreditLedger
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&ledger).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrLedgerNotFound
			}
			return err
		}
		if req.Amount > ledger.OutstandingBalance+0.005 {
			return fmt.Errorf("jumlah pelunasan melebihi sisa hutang (%s)", formatRupiah(ledger.OutstandingBalance))
		}
		userID = ledger.UserID

		ledger.OutstandingBalance -= req.Amount
		ledger.UsedCredit -= req.Amount
		if ledger.OutstandingBalance < 0 {
			ledger.OutstandingBalance = 0
		}
		if ledger.UsedCredit < 0 {
			ledger.UsedCredit = 0
		}
		ledger.RemainingCredit = ledger.CreditLimit - ledger.UsedCredit
		if ledger.RemainingCredit < 0 {
			ledger.RemainingCredit = 0
		}
		if err := tx.Model(&models.CreditLedger{}).Where("id = ?", ledger.ID).Updates(map[string]interface{}{
			"used_credit":         ledger.UsedCredit,
			"outstanding_balance": ledger.OutstandingBalance,
			"remaining_credit":    ledger.RemainingCredit,
		}).Error; err != nil {
			return err
		}

		notes := strings.TrimSpace(req.Notes)
		if notes == "" {
			notes = "Pelunasan kasbon diterima toko"
		}
		return tx.Create(&models.CreditPayment{
			LedgerID:         ledger.ID,
			Type:             string(constants.CreditTxSettlement),
			Amount:           req.Amount,
			PaymentDate:      time.Now(),
			RemainingBalance: ledger.OutstandingBalance,
			Notes:            notes,
		}).Error
	})
	if err != nil {
		return nil, err
	}

	publishEvent(s.bus, "CREDIT_SETTLED", "CREDIT_LEDGER", id, map[string]interface{}{
		"user_id": userID.String(),
		"title":   "Pembayaran kasbon diterima",
		"message": "Pembayaran kasbon sebesar " + formatRupiah(req.Amount) + " sudah dicatat. Terima kasih!",
	})

	return s.GetLedgerByID(ctx, id)
}

func (s *creditLedgerService) GetPaymentHistory(ctx context.Context, ledgerID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	payments, total, err := s.creditPayRepo.FindByLedgerID(ctx, ledgerID, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.CreditPaymentHistoryDTO, len(payments))
	for i, p := range payments {
		dtos[i] = dto.CreditPaymentHistoryDTO{
			ID:               p.ID,
			LedgerID:         p.LedgerID,
			PaymentID:        p.PaymentID,
			OrderID:          p.OrderID,
			Type:             p.Type,
			Amount:           p.Amount,
			PaymentDate:      p.PaymentDate,
			RemainingBalance: p.RemainingBalance,
			Notes:            p.Notes,
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func mapLedgerToDTO(l *models.CreditLedger) dto.CreditLedgerDTO {
	res := dto.CreditLedgerDTO{
		ID:                 l.ID,
		UserID:             l.UserID,
		UserName:           "Pelanggan",
		CreditLimit:        l.CreditLimit,
		UsedCredit:         l.UsedCredit,
		RemainingCredit:    l.RemainingCredit,
		OutstandingBalance: l.OutstandingBalance,
		UpdatedAt:          l.UpdatedAt,
	}
	if l.User != nil {
		res.UserName = l.User.FullName
		res.UserPhone = l.User.PhoneNumber
		if l.User.B2BProfile != nil {
			res.StoreName = l.User.B2BProfile.StoreName
		}
	}
	return res
}
