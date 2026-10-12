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
)

var (
	ErrPaymentNotFound     = errors.New("data pembayaran tidak ditemukan")
	ErrInvoiceAlreadyPaid  = errors.New("tagihan ini sudah lunas")
	ErrOnlyPaidCanRefund   = errors.New("refund hanya bisa untuk pembayaran yang sudah lunas")
	ErrPaymentWrongStatus  = errors.New("status pembayaran tidak memungkinkan aksi ini")
	ErrUnauthorizedPayment = errors.New("anda tidak berhak melihat atau mengubah pembayaran ini")
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreatePaymentDTO) (*dto.PaymentResponseDTO, error)
	GetPaymentByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.PaymentResponseDTO, error)
	GetPayments(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error)
	UploadProof(ctx context.Context, userID uuid.UUID, req dto.UploadProofDTO) (*dto.PaymentResponseDTO, error)
	GetWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	ApprovePayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error)
	RejectPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error)
	CancelPayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, paymentID uuid.UUID) (*dto.PaymentResponseDTO, error)
	RefundPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, req dto.RefundDTO) (*dto.PaymentResponseDTO, error)
}

type paymentService struct {
	db          *gorm.DB
	paymentRepo repositories.PaymentRepository
	methodRepo  repositories.PaymentMethodRepository
	invoiceRepo repositories.InvoiceRepository
	orderRepo   repositories.OrderRepository
	logRepo     repositories.PaymentLogRepository
	bus         eventbus.EventBus
}

func NewPaymentService(
	db *gorm.DB,
	paymentRepo repositories.PaymentRepository,
	methodRepo repositories.PaymentMethodRepository,
	invoiceRepo repositories.InvoiceRepository,
	orderRepo repositories.OrderRepository,
	logRepo repositories.PaymentLogRepository,
	bus eventbus.EventBus,
) PaymentService {
	return &paymentService{
		db:          db,
		paymentRepo: paymentRepo,
		methodRepo:  methodRepo,
		invoiceRepo: invoiceRepo,
		orderRepo:   orderRepo,
		logRepo:     logRepo,
		bus:         bus,
	}
}

// CreatePayment opens a new bill for an unpaid order, e.g. when a customer
// wants to switch from bank transfer to COD after a rejected transfer proof.
// Kasbon is only charged when the order is placed.
func (s *paymentService) CreatePayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.CreatePaymentDTO) (*dto.PaymentResponseDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil || order.UserID != userID {
		return nil, ErrOrderNotFound
	}
	if order.Invoice == nil {
		return nil, errors.New("tagihan pesanan tidak ditemukan")
	}
	if order.Invoice.Status == constants.InvoiceStatusPaid {
		return nil, ErrInvoiceAlreadyPaid
	}
	if order.Status == constants.OrderStatusCancelled || order.Status == constants.OrderStatusCompleted {
		return nil, errors.New("pesanan sudah selesai atau dibatalkan")
	}

	method, err := s.methodRepo.FindByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, err
	}
	if method == nil || !method.IsActive || method.Type == constants.PaymentTypeB2BCredit {
		return nil, ErrInvalidPaymentMethod
	}

	existing, err := s.paymentRepo.FindByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	for _, p := range existing {
		if p.Status == constants.PaymentStatusPending || p.Status == constants.PaymentStatusWaitingConfirmation || p.Status == constants.PaymentStatusPaid {
			return nil, errors.New("pesanan ini masih punya tagihan aktif")
		}
	}

	now := time.Now()
	payment := &models.Payment{
		PaymentNumber:   newReference("PAY"),
		OrderID:         order.ID,
		InvoiceID:       order.Invoice.ID,
		UserID:          userID,
		PaymentMethodID: method.ID,
		Amount:          order.Invoice.TotalAmount,
		GrandTotal:      order.Invoice.TotalAmount,
		Status:          constants.PaymentStatusPending,
		ExpiredAt:       now.Add(24 * time.Hour),
		Notes:           req.Notes,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}
		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("payment_method", method.Code).Error; err != nil {
			return err
		}
		if method.Type == constants.PaymentTypeCOD && order.Status == constants.OrderStatusPendingPayment {
			if err := setOrderStatusTx(tx, order, constants.OrderStatusConfirmed, "Metode pembayaran diganti ke COD. Pesanan dikonfirmasi."); err != nil {
				return err
			}
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "CREATE_PAYMENT", OldStatus: "NONE", NewStatus: payment.Status,
			Notes: "Tagihan baru dibuat dengan metode " + method.Name,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, userID, userRole, payment.ID)
}

func (s *paymentService) GetPaymentByID(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, id uuid.UUID) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil || (!userRole.IsManager() && payment.UserID != userID) {
		return nil, ErrPaymentNotFound
	}

	res := mapPaymentToDTO(payment)
	return &res, nil
}

func (s *paymentService) GetPayments(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, req dto.PaginationRequest, status string) (dto.PaginatedResponse, error) {
	var payments []models.Payment
	var total int64
	var err error

	if userRole.IsManager() {
		payments, total, err = s.paymentRepo.FindAll(ctx, req, status)
	} else {
		payments, total, err = s.paymentRepo.FindByUserID(ctx, userID, req, status)
	}
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.PaymentResponseDTO, len(payments))
	for i := range payments {
		dtos[i] = mapPaymentToDTO(&payments[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *paymentService) UploadProof(ctx context.Context, userID uuid.UUID, req dto.UploadProofDTO) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil || payment.UserID != userID {
		return nil, ErrPaymentNotFound
	}
	if payment.PaymentMethod == nil || payment.PaymentMethod.Type != constants.PaymentTypeManualTransfer {
		return nil, errors.New("bukti transfer hanya untuk pembayaran transfer bank")
	}
	if payment.Status != constants.PaymentStatusPending && payment.Status != constants.PaymentStatusFailed && payment.Status != constants.PaymentStatusWaitingConfirmation {
		return nil, ErrPaymentWrongStatus
	}
	if payment.Order != nil && payment.Order.Status == constants.OrderStatusCancelled {
		return nil, errors.New("pesanan sudah dibatalkan")
	}

	oldStatus := payment.Status
	proof := strings.TrimSpace(req.ProofImageURL)
	payment.ProofImageURL = &proof
	payment.Status = constants.PaymentStatusWaitingConfirmation
	payment.Notes = req.Notes

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "UPLOAD_TRANSFER_PROOF", OldStatus: oldStatus,
			NewStatus: constants.PaymentStatusWaitingConfirmation, Notes: "Pelanggan mengunggah bukti transfer",
		})
	})
	if err != nil {
		return nil, err
	}

	orderNumber := ""
	if payment.Order != nil {
		orderNumber = payment.Order.OrderNumber
	}
	publishEvent(s.bus, "PAYMENT_PROOF_UPLOADED", "PAYMENT", payment.ID, map[string]interface{}{
		"order_number":  orderNumber,
		"notify_staff":  true,
		"staff_title":   "Bukti transfer masuk: " + orderNumber,
		"staff_message": "Pelanggan mengunggah bukti transfer " + formatRupiah(payment.GrandTotal) + ". Mohon diverifikasi.",
	})

	return s.GetPaymentByID(ctx, userID, constants.RoleB2CCustomer, payment.ID)
}

func (s *paymentService) GetWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	payments, total, err := s.paymentRepo.FindWaitingConfirmation(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.PaymentResponseDTO, len(payments))
	for i := range payments {
		dtos[i] = mapPaymentToDTO(&payments[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

// ApprovePayment confirms a bank transfer (with or without an uploaded
// proof, since the store may check its bank statement directly).
func (s *paymentService) ApprovePayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	if payment.Status != constants.PaymentStatusWaitingConfirmation && payment.Status != constants.PaymentStatusPending {
		return nil, ErrPaymentWrongStatus
	}
	order, err := s.orderRepo.FindByID(ctx, payment.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Status == constants.OrderStatusCancelled {
		return nil, errors.New("pesanan sudah dibatalkan; pembayaran tidak bisa disetujui")
	}

	now := time.Now()
	oldStatus := payment.Status

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = constants.PaymentStatusPaid
		payment.PaidAt = &now
		if notes != "" {
			payment.Notes = notes
		}
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}
		if err := markInvoicePaidTx(tx, payment.InvoiceID, now); err != nil {
			return err
		}
		if order.Status == constants.OrderStatusPendingPayment || order.Status == constants.OrderStatusCreated {
			if err := setOrderStatusTx(tx, order, constants.OrderStatusConfirmed, "Pembayaran diverifikasi. Pesanan siap diproses."); err != nil {
				return err
			}
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "APPROVE_PAYMENT", OldStatus: oldStatus, NewStatus: constants.PaymentStatusPaid,
			Notes: strings.TrimSpace("Pembayaran disetujui admin. " + notes),
		})
	})
	if err != nil {
		return nil, err
	}

	publishEvent(s.bus, "PAYMENT_PAID", "PAYMENT", payment.ID, map[string]interface{}{
		"user_id":      payment.UserID.String(),
		"order_number": order.OrderNumber,
		"grand_total":  formatRupiah(payment.GrandTotal),
		"title":        "Pembayaran " + order.OrderNumber + " diterima",
		"message":      "Pembayaran " + formatRupiah(payment.GrandTotal) + " sudah kami terima. Pesanan Anda segera diproses.",
	})

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

func (s *paymentService) RejectPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, notes string) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	if payment.Status != constants.PaymentStatusWaitingConfirmation {
		return nil, ErrPaymentWrongStatus
	}

	oldStatus := payment.Status
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = constants.PaymentStatusFailed
		payment.Notes = notes
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "REJECT_PAYMENT", OldStatus: oldStatus, NewStatus: constants.PaymentStatusFailed,
			Notes: "Bukti transfer ditolak: " + notes,
		})
	})
	if err != nil {
		return nil, err
	}

	orderNumber := ""
	if payment.Order != nil {
		orderNumber = payment.Order.OrderNumber
	}
	publishEvent(s.bus, "PAYMENT_REJECTED", "PAYMENT", payment.ID, map[string]interface{}{
		"user_id":      payment.UserID.String(),
		"order_number": orderNumber,
		"title":        "Bukti transfer " + orderNumber + " ditolak",
		"message":      "Alasan: " + notes + ". Silakan unggah ulang bukti transfer yang benar.",
	})

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

// CancelPayment voids an unpaid bill (staff only). The order itself is
// cancelled through the order flow.
func (s *paymentService) CancelPayment(ctx context.Context, userID uuid.UUID, userRole constants.UserRole, paymentID uuid.UUID) (*dto.PaymentResponseDTO, error) {
	if !userRole.IsManager() {
		return nil, ErrUnauthorizedPayment
	}
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	if payment.Status != constants.PaymentStatusPending && payment.Status != constants.PaymentStatusWaitingConfirmation && payment.Status != constants.PaymentStatusFailed {
		return nil, ErrPaymentWrongStatus
	}

	oldStatus := payment.Status
	payment.Status = constants.PaymentStatusCancelled

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "CANCEL_PAYMENT", OldStatus: oldStatus, NewStatus: constants.PaymentStatusCancelled,
			Notes: "Tagihan dibatalkan admin",
		})
	})
	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, userID, userRole, paymentID)
}

func (s *paymentService) RefundPayment(ctx context.Context, adminID uuid.UUID, paymentID uuid.UUID, req dto.RefundDTO) (*dto.PaymentResponseDTO, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	if payment.Status != constants.PaymentStatusPaid {
		return nil, ErrOnlyPaidCanRefund
	}
	if payment.PaymentMethod != nil && payment.PaymentMethod.Type == constants.PaymentTypeB2BCredit {
		return nil, errors.New("pembayaran kasbon dikembalikan otomatis saat pesanan dibatalkan")
	}
	if req.RefundAmount > payment.GrandTotal {
		return nil, errors.New("jumlah refund melebihi total pembayaran")
	}

	oldStatus := payment.Status
	newStatus := constants.PaymentStatusRefunded
	if req.RefundAmount < payment.GrandTotal {
		newStatus = constants.PaymentStatusPartiallyRefunded
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.Status = newStatus
		payment.Notes = fmt.Sprintf("Refund %s: %s", formatRupiah(req.RefundAmount), req.Reason)
		if err := s.paymentRepo.Update(ctx, tx, payment); err != nil {
			return err
		}
		return s.logRepo.Create(ctx, tx, &models.PaymentLog{
			PaymentID: payment.ID, Action: "REFUND_PAYMENT", OldStatus: oldStatus, NewStatus: newStatus,
			Notes: fmt.Sprintf("Refund %s diproses admin. Alasan: %s", formatRupiah(req.RefundAmount), req.Reason),
		})
	})
	if err != nil {
		return nil, err
	}

	return s.GetPaymentByID(ctx, adminID, constants.RoleSuperAdmin, paymentID)
}

func mapPaymentToDTO(p *models.Payment) dto.PaymentResponseDTO {
	res := dto.PaymentResponseDTO{
		ID:            p.ID,
		PaymentNumber: p.PaymentNumber,
		OrderID:       p.OrderID,
		InvoiceID:     p.InvoiceID,
		UserID:        p.UserID,
		MethodName:    "-",
		Amount:        p.Amount,
		AdminFee:      p.AdminFee,
		Discount:      p.Discount,
		GrandTotal:    p.GrandTotal,
		Status:        p.Status,
		ProofImageURL: p.ProofImageURL,
		PaidAt:        p.PaidAt,
		ExpiredAt:     p.ExpiredAt,
		Notes:         p.Notes,
		CreatedAt:     p.CreatedAt,
	}
	if p.PaymentMethod != nil {
		res.MethodCode = p.PaymentMethod.Code
		res.MethodName = p.PaymentMethod.Name
		res.MethodType = p.PaymentMethod.Type
	}
	if p.Order != nil {
		res.OrderNumber = p.Order.OrderNumber
		res.OrderStatus = string(p.Order.Status)
	}
	if p.User != nil {
		res.CustomerName = p.User.FullName
		res.CustomerPhone = p.User.PhoneNumber
	}

	res.Logs = make([]dto.PaymentLogResponseDTO, len(p.Logs))
	for i, l := range p.Logs {
		res.Logs[i] = dto.PaymentLogResponseDTO{
			ID: l.ID, Action: l.Action, OldStatus: l.OldStatus, NewStatus: l.NewStatus, Notes: l.Notes, CreatedAt: l.CreatedAt,
		}
	}
	return res
}
