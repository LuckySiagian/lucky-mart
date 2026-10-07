package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound     = errors.New("payment record not found")
	ErrInvoiceAlreadyPaid  = errors.New("invoice is already paid")
	ErrPaymentExpired      = errors.New("payment has expired")
	ErrOnlyPaidCanRefund   = errors.New("refund can only be performed on paid payments")
	ErrExceededCreditLimit = errors.New("requested transaction exceeds your remaining B2B credit limit")
	ErrUnauthorizedPayment = errors.New("you do not have permission to view or manage this payment")
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreatePaymentDTO) (*dto.PaymentResponseDTO, error)
	GetPaymentByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.PaymentResponseDTO, error)
	GetPayments(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	UploadProof(ctx context.Context, userID uuid.UUID, req dto.UploadProofDTO) (*dto.PaymentResponseDTO, error)
	GetWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	ApprovePayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error)
	RejectPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error)
	CancelPayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, paymentID uuid.UUID) (*dto.PaymentResponseDTO, error)
	RefundPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, req dto.RefundDTO) (*dto.PaymentResponseDTO, error)
}

type paymentService struct {
	db             *gorm.DB
	paymentRepo    repositories.PaymentRepository
	methodRepo     repositories.PaymentMethodRepository
	invoiceRepo    repositories.InvoiceRepository
	orderRepo      repositories.OrderRepository
	txRepo         repositories.PaymentTransactionRepository
	logRepo        repositories.PaymentLogRepository
	ledgerRepo     repositories.CreditLedgerRepository
	creditPayRepo  repositories.CreditPaymentRepository
}

func NewPaymentService(
	db *gorm.DB,
	paymentRepo repositories.PaymentRepository,
	methodRepo repositories.PaymentMethodRepository,
	invoiceRepo repositories.InvoiceRepository,
	orderRepo repositories.OrderRepository,
	txRepo repositories.PaymentTransactionRepository,
	logRepo repositories.PaymentLogRepository,
	ledgerRepo repositories.CreditLedgerRepository,
	creditPayRepo repositories.CreditPaymentRepository,
) PaymentService {
	return &paymentService{
		db:            db,
		paymentRepo:   paymentRepo,
		methodRepo:   methodRepo,
		invoiceRepo:  invoiceRepo,
		orderRepo:    orderRepo,
		txRepo:       txRepo,
		logRepo:      logRepo,
		ledgerRepo:   ledgerRepo,
		creditPayRepo: creditPayRepo,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreatePaymentDTO) (*dto.PaymentResponseDTO, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, req.InvoiceID)
	if err != nil || invoice == nil {
		return nil, errors.New("invoice not found")
	}

	if invoice.Status == "PAID" {
		return nil, ErrInvoiceAlreadyPaid
	}

	method, err := s.methodRepo.FindByID(ctx, req.PaymentMethodID)
	if err != nil || method == nil || !method.IsActive {
		return nil, errors.New("invalid or inactive payment method")
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	payNum := fmt.Sprintf("PAY-%s-%04d", now.Format("20060102"), now.UnixNano()%10000)

	adminFee := 0.00
	grandTotal := invoice.TotalAmount + adminFee
	initialStatus := "PENDING"

	var createdPayment *models.Payment

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// B2B CREDIT PAYMENT METHOD HANDLING
		if method.Type == "B2B_CREDIT" {
			ledger, err := s.ledgerRepo.FindByUserID(ctx, userID)
			if err != nil || ledger == nil {
				return errors.New("B2B credit ledger account not found")
			}

			if grandTotal > ledger.RemainingCredit {
				return fmt.Errorf("%w: required %.2f, remaining credit limit %.2f", ErrExceededCreditLimit, grandTotal, ledger.RemainingCredit)
			}

			// Update B2B Credit Ledger
			ledger.UsedCredit += grandTotal
			ledger.RemainingCredit = ledger.CreditLimit - ledger.UsedCredit
			ledger.OutstandingBalance += grandTotal

			if err := s.ledgerRepo.Update(ctx, tx, ledger); err != nil {
				return err
			}

			initialStatus = "PAID"
			paidTime := now
			createdPayment = &models.Payment{
				PaymentNumber:   payNum,
				OrderID:         req.OrderID,
				InvoiceID:       req.InvoiceID,
				UserID:          userID,
				PaymentMethodID: req.PaymentMethodID,
				Amount:          invoice.TotalAmount,
				AdminFee:        adminFee,
				GrandTotal:      grandTotal,
				Status:          initialStatus,
				PaidAt:          &paidTime,
				ExpiredAt:       expiredAt,
				Notes:           req.Notes,
			}

			if err := s.paymentRepo.Create(ctx, tx, createdPayment); err != nil {
				return err
			}

			// Mark Invoice & Order as PAID
			if err := s.invoiceRepo.UpdateStatus(ctx, tx, invoice.ID, "PAID"); err != nil {
				return err
			}
			order, _ := s.orderRepo.FindByID(ctx, req.OrderID)
			if order != nil {
				order.Status = constants.OrderStatusConfirmed
				_ = s.orderRepo.Update(ctx, tx, order)
				_ = s.orderRepo.AddHistory(ctx, tx, &models.OrderHistory{
					OrderID: order.ID,
					Status:  constants.OrderStatusConfirmed,
					Notes:   "Pembayaran via B2B Credit Tempo disetujui. Pesanan dikonfirmasi.",
				})
			}

			// Log Credit Payment & Payment Log
			_ = s.creditPayRepo.Create(ctx, tx, &models.CreditPayment{
				LedgerID:         ledger.ID,
				PaymentID:        createdPayment.ID,
				Amount:           grandTotal,
				PaymentDate:      now,
				RemainingBalance: ledger.OutstandingBalance,
			})

			_ = s.logRepo.Create(ctx, tx, &models.PaymentLog{
				PaymentID: createdPayment.ID,
				Action:    "B2B_CREDIT_AUTO_PAY",
				OldStatus: "PENDING",
				NewStatus: "PAID",
				Notes:     "Pembayaran berhasil dipotong dari limit kredit B2B",
			})

			return nil
		}

		// REGULAR PAYMENT METHODS (CASH, COD, MANUAL TRANSFER)
		createdPayment = &models.Payment{
			PaymentNumber:   payNum,
			OrderID:         req.OrderID,
			InvoiceID:       req.InvoiceID,
			UserID:          userID,
			PaymentMethodID: req.PaymentMethodID,
			Amount:          invoice.TotalAmount,
			AdminFee:        adminFee,
			GrandTotal:      grandTotal,
			Status:          initialStatus,
			ExpiredAt:       expiredAt,
			Notes:           req.Notes,
		}

		if err := s.paymentRepo.Create(ctx, tx, createdPayment); err != nil {
			return err
		}

		_ = s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: createdPayment.ID,
			Action:    "CREATE_PAYMENT",
			OldStatus: "NONE",
			NewStatus: initialStatus,
			Notes:     "Tagihan pembayaran berhasil dibuat",
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, userID, userRole, createdPayment.ID)
}

func (s *paymentService) GetPaymentByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	if userRole != constants.RoleSuperAdmin && userRole != constants.RoleStoreManager {
		if payment.UserID != userID {
			return nil, ErrUnauthorizedPayment
		}
	}

	res := mapPaymentToDTO(payment)
	return &res, nil
}

func (s *paymentService) GetPayments(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	var payments []models.Payment
	var total int64
	var err error

	if userRole == constants.RoleSuperAdmin || userRole == constants.RoleStoreManager {
		payments, total, err = s.paymentRepo.FindAll(ctx, req)
	} else {
		payments, total, err = s.paymentRepo.FindByUserID(ctx, userID, req)
	}

	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.PaymentResponseDTO, len(payments))
	for i, p := range payments {
		dtos[i] = mapPaymentToDTO(&p)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *paymentService) UploadProof(ctx context.Context, userID uuid.UUID, req dto.UploadProofDTO) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, req.PaymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	if payment.UserID != userID {
		return nil, ErrUnauthorizedPayment
	}

	if payment.Status == "PAID" {
		return nil, ErrInvoiceAlreadyPaid
	}

	oldStatus := payment.Status
	payment.ProofImageURL = &req.ProofImageURL
	payment.Status = "WAITING_CONFIRMATION"
	payment.Notes = req.Notes

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}

		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID,
			Action:    "UPLOAD_TRANSFER_PROOF",
			OldStatus: oldStatus,
			NewStatus: "WAITING_CONFIRMATION",
			Notes:     "Pelanggan mengunggah bukti transfer manual",
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, userID, constants.RoleB2CCustomer, payment.ID)
}

func (s *paymentService) GetWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	payments, total, err := s.paymentRepo.FindWaitingConfirmation(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.PaymentResponseDTO, len(payments))
	for i, p := range payments {
		dtos[i] = mapPaymentToDTO(&p)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *paymentService) ApprovePayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	now := time.Now()
	oldStatus := payment.Status

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = "PAID"
		payment.PaidAt = &now
		payment.Notes = notes

		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}

		if err := s.invoiceRepo.UpdateStatus(ctx, tx, payment.InvoiceID, "PAID"); err != nil {
			return err
		}

		order, _ := s.orderRepo.FindByID(ctx, payment.OrderID)
		if order != nil {
			order.Status = constants.OrderStatusConfirmed
			_ = s.orderRepo.Update(ctx, tx, order)
			_ = s.orderRepo.AddHistory(ctx, tx, &models.OrderHistory{
				OrderID: order.ID,
				Status:  constants.OrderStatusConfirmed,
				Notes:   "Bukti pembayaran diverifikasi dan disetujui Admin. Pesanan siap diproses gudang.",
			})
		}

		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID,
			Action:    "APPROVE_PAYMENT",
			OldStatus: oldStatus,
			NewStatus: "PAID",
			Notes:     fmt.Sprintf("Verifikasi disetujui oleh Admin: %s", notes),
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

func (s *paymentService) RejectPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	oldStatus := payment.Status
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = "FAILED"
		payment.Notes = notes

		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}

		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID,
			Action:    "REJECT_PAYMENT",
			OldStatus: oldStatus,
			NewStatus: "FAILED",
			Notes:     fmt.Sprintf("Bukti pembayaran ditolak Admin: %s", notes),
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

func (s *paymentService) CancelPayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, paymentID uuid.UUID) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	oldStatus := payment.Status
	payment.Status = "CANCELLED"

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}

		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID,
			Action:    "CANCEL_PAYMENT",
			OldStatus: oldStatus,
			NewStatus: "CANCELLED",
			Notes:     "Pembayaran dibatalkan",
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, userID, userRole, paymentID)
}

func (s *paymentService) RefundPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, req dto.RefundDTO) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, ErrPaymentNotFound
	}

	if payment.Status != "PAID" {
		return nil, ErrOnlyPaidCanRefund
	}

	if req.RefundAmount > payment.GrandTotal {
		return nil, errors.New("refund amount cannot exceed grand total paid")
	}

	oldStatus := payment.Status
	newStatus := "REFUNDED"
	if req.RefundAmount < payment.GrandTotal {
		newStatus = "PARTIALLY_REFUNDED"
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = newStatus
		payment.Notes = fmt.Sprintf("Refund %.2f: %s", req.RefundAmount, req.Reason)

		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}

		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID,
			Action:    "REFUND_PAYMENT",
			OldStatus: oldStatus,
			NewStatus: newStatus,
			Notes:     fmt.Sprintf("Refund sebesar %.2f diproses Admin. Alasan: %s", req.RefundAmount, req.Reason),
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

func mapPaymentToDTO(p *models.Payment) dto.PaymentResponseDTO {
	methodName := "Transfer Manual"
	methodType := "MANUAL_TRANSFER"
	if p.PaymentMethod != nil {
		methodName = p.PaymentMethod.Name
		methodType = p.PaymentMethod.Type
	}

	logs := make([]dto.PaymentLogResponseDTO, len(p.Logs))
	for i, l := range p.Logs {
		logs[i] = dto.PaymentLogResponseDTO{
			ID:        l.ID,
			Action:    l.Action,
			OldStatus: l.OldStatus,
			NewStatus: l.NewStatus,
			Notes:     l.Notes,
			CreatedAt: l.CreatedAt,
		}
	}

	return dto.PaymentResponseDTO{
		ID:            p.ID,
		PaymentNumber: p.PaymentNumber,
		OrderID:       p.OrderID,
		InvoiceID:     p.InvoiceID,
		UserID:        p.UserID,
		MethodName:    methodName,
		MethodType:    methodType,
		Amount:        p.Amount,
		AdminFee:      p.AdminFee,
		Discount:      p.Discount,
		GrandTotal:    p.GrandTotal,
		Status:        p.Status,
		ProofImageURL: p.ProofImageURL,
		PaidAt:        p.PaidAt,
		ExpiredAt:     p.ExpiredAt,
		Notes:         p.Notes,
		Logs:          logs,
		CreatedAt:     p.CreatedAt,
	}
}
