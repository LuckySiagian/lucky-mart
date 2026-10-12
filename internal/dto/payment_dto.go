package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentDTO struct {
	OrderID         uuid.UUID `json:"order_id" binding:"required"`
	PaymentMethodID uuid.UUID `json:"payment_method_id" binding:"required"`
	Notes           string    `json:"notes" binding:"omitempty,max=500"`
}

type UploadProofDTO struct {
	PaymentID     uuid.UUID `json:"payment_id" binding:"required"`
	ProofImageURL string    `json:"proof_image_url" binding:"required,max=500"` // uploaded /uploads/images/... path or an http(s) link
	Notes         string    `json:"notes" binding:"omitempty,max=500"`
}

type ConfirmPaymentDTO struct {
	Notes string `json:"notes" binding:"omitempty,max=500"`
}

type RejectPaymentDTO struct {
	Notes string `json:"notes" binding:"required,min=3,max=500"`
}

type RefundDTO struct {
	RefundAmount float64 `json:"refund_amount" binding:"required,gt=0"`
	Reason       string  `json:"reason" binding:"required,min=3,max=255"`
}

type CreatePaymentMethodDTO struct {
	Code        string `json:"code" binding:"required,min=2,max=30"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Type        string `json:"type" binding:"required,oneof=CASH COD MANUAL_TRANSFER GATEWAY B2B_CREDIT"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdatePaymentMethodDTO struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
}

type CreateCreditLedgerDTO struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	CreditLimit float64   `json:"credit_limit" binding:"gte=0"`
}

type AdjustCreditLimitDTO struct {
	NewCreditLimit float64 `json:"new_credit_limit" binding:"gte=0"`
}

type SettleCreditDTO struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Notes  string  `json:"notes" binding:"omitempty,max=255"`
}

type PaymentMethodDTO struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description,omitempty"`
}

type PaymentLogResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Action    string    `json:"action"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentResponseDTO struct {
	ID            uuid.UUID               `json:"id"`
	PaymentNumber string                  `json:"payment_number"`
	OrderID       uuid.UUID               `json:"order_id"`
	OrderNumber   string                  `json:"order_number,omitempty"`
	OrderStatus   string                  `json:"order_status,omitempty"`
	InvoiceID     uuid.UUID               `json:"invoice_id"`
	UserID        uuid.UUID               `json:"user_id"`
	CustomerName  string                  `json:"customer_name,omitempty"`
	CustomerPhone string                  `json:"customer_phone,omitempty"`
	MethodCode    string                  `json:"method_code"`
	MethodName    string                  `json:"method_name"`
	MethodType    string                  `json:"method_type"`
	Amount        float64                 `json:"amount"`
	AdminFee      float64                 `json:"admin_fee"`
	Discount      float64                 `json:"discount"`
	GrandTotal    float64                 `json:"grand_total"`
	Status        string                  `json:"status"`
	ProofImageURL *string                 `json:"proof_image_url,omitempty"`
	PaidAt        *time.Time              `json:"paid_at,omitempty"`
	ExpiredAt     time.Time               `json:"expired_at"`
	Notes         string                  `json:"notes,omitempty"`
	Logs          []PaymentLogResponseDTO `json:"logs"`
	CreatedAt     time.Time               `json:"created_at"`
}

type CreditLedgerDTO struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	UserName           string    `json:"user_name"`
	UserPhone          string    `json:"user_phone,omitempty"`
	StoreName          string    `json:"store_name,omitempty"`
	CreditLimit        float64   `json:"credit_limit"`
	UsedCredit         float64   `json:"used_credit"`
	RemainingCredit    float64   `json:"remaining_credit"`
	OutstandingBalance float64   `json:"outstanding_balance"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreditPaymentHistoryDTO struct {
	ID               uuid.UUID  `json:"id"`
	LedgerID         uuid.UUID  `json:"ledger_id"`
	PaymentID        *uuid.UUID `json:"payment_id,omitempty"`
	OrderID          *uuid.UUID `json:"order_id,omitempty"`
	Type             string     `json:"type"`
	Amount           float64    `json:"amount"`
	PaymentDate      time.Time  `json:"payment_date"`
	RemainingBalance float64    `json:"remaining_balance"`
	Notes            string     `json:"notes,omitempty"`
}
