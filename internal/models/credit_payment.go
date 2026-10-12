package models

import (
	"time"

	"github.com/google/uuid"
)

// CreditPayment is one entry in a customer's kasbon history: either credit
// used by an order (CREDIT_USAGE), a repayment received by the store
// (PAYMENT_SETTLEMENT), or a reversal when a kasbon order is cancelled.
type CreditPayment struct {
	BaseModel
	LedgerID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"ledger_id"`
	PaymentID        *uuid.UUID `gorm:"type:uuid;index" json:"payment_id,omitempty"`
	OrderID          *uuid.UUID `gorm:"type:uuid;index" json:"order_id,omitempty"`
	Type             string     `gorm:"type:varchar(30);not null;default:'CREDIT_USAGE';index" json:"type"`
	Amount           float64    `gorm:"type:numeric(15,2);not null;check:amount > 0" json:"amount"`
	PaymentDate      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"payment_date"`
	RemainingBalance float64    `gorm:"type:numeric(15,2);not null" json:"remaining_balance"`
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
}

func (CreditPayment) TableName() string {
	return "credit_payments"
}
