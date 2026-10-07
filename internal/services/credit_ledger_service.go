package services

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditLedgerService interface {
	GetLedgers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetLedgerByID(ctx context.Context, id uuid.UUID) (*dto.CreditLedgerDTO, error)
	GetLedgerByUserID(ctx context.Context, userID uuid.UUID) (*dto.CreditLedgerDTO, error)
	CreateLedger(ctx context.Context, req dto.CreateCreditLedgerDTO) (*dto.CreditLedgerDTO, error)
	AdjustCreditLimit(ctx context.Context, id uuid.UUID, req dto.AdjustCreditLimitDTO) (*dto.CreditLedgerDTO, error)
	GetPaymentHistory(ctx context.Context, ledgerID uuid.UUID, req dto.PaginationRequest) (dto.PaginatedResponse, error)
}

type creditLedgerService struct {
	db            *gorm.DB
	ledgerRepo    repositories.CreditLedgerRepository
	creditPayRepo repositories.CreditPaymentRepository
}

func NewCreditLedgerService(
	db *gorm.DB,
	ledgerRepo repositories.CreditLedgerRepository,
	creditPayRepo repositories.CreditPaymentRepository,
) CreditLedgerService {
	return &creditLedgerService{
		db:            db,
		ledgerRepo:    ledgerRepo,
		creditPayRepo: creditPayRepo,
	}
}

func (s *creditLedgerService) GetLedgers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	ledgers, total, err := s.ledgerRepo.FindAll(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.CreditLedgerDTO, len(ledgers))
	for i, l := range ledgers {
		dtos[i] = mapLedgerToDTO(&l)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *creditLedgerService) GetLedgerByID(ctx context.Context, id uuid.UUID) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByID(ctx, id)
	if err != nil || l == nil {
		return nil, errors.New("credit ledger account not found")
	}
	res := mapLedgerToDTO(l)
	return &res, nil
}

func (s *creditLedgerService) GetLedgerByUserID(ctx context.Context, userID uuid.UUID) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByUserID(ctx, userID)
	if err != nil || l == nil {
		return nil, errors.New("B2B credit ledger not found for this customer")
	}
	res := mapLedgerToDTO(l)
	return &res, nil
}

func (s *creditLedgerService) CreateLedger(ctx context.Context, req dto.CreateCreditLedgerDTO) (*dto.CreditLedgerDTO, error) {
	existing, err := s.ledgerRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("credit ledger already exists for this user")
	}

	ledger := &models.CreditLedger{
		UserID:             req.UserID,
		CreditLimit:        req.CreditLimit,
		UsedCredit:         0,
		RemainingCredit:    req.CreditLimit,
		OutstandingBalance: 0,
	}

	if err := s.ledgerRepo.Create(ctx, nil, ledger); err != nil {
		return nil, err
	}

	return s.GetLedgerByID(ctx, ledger.ID)
}

func (s *creditLedgerService) AdjustCreditLimit(ctx context.Context, id uuid.UUID, req dto.AdjustCreditLimitDTO) (*dto.CreditLedgerDTO, error) {
	l, err := s.ledgerRepo.FindByID(ctx, id)
	if err != nil || l == nil {
		return nil, errors.New("credit ledger not found")
	}

	l.CreditLimit = req.NewCreditLimit
	l.RemainingCredit = l.CreditLimit - l.UsedCredit
	if l.RemainingCredit < 0 {
		l.RemainingCredit = 0
	}

	if err := s.ledgerRepo.Update(ctx, nil, l); err != nil {
		return nil, err
	}

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
			Amount:           p.Amount,
			PaymentDate:      p.PaymentDate,
			RemainingBalance: p.RemainingBalance,
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func mapLedgerToDTO(l *models.CreditLedger) dto.CreditLedgerDTO {
	userName := "B2B Customer"
	if l.User != nil {
		userName = l.User.FullName
	}
	return dto.CreditLedgerDTO{
		ID:                 l.ID,
		UserID:             l.UserID,
		UserName:           userName,
		CreditLimit:        l.CreditLimit,
		UsedCredit:         l.UsedCredit,
		RemainingCredit:    l.RemainingCredit,
		OutstandingBalance: l.OutstandingBalance,
		UpdatedAt:          l.UpdatedAt,
	}
}
