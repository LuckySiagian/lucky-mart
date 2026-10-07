package models

import (
	"time"

	"github.com/google/uuid"
)

type CreditPayment struct {
	BaseModel
	LedgerID         uuid.UUID `gorm:"type:uuid;not null;index" json:"ledger_id"`
	PaymentID        uuid.UUID `gorm:"type:uuid;not null;index" json:"payment_id"`
	Amount           float64   `gorm:"type:numeric(15,2);not null;check:amount > 0" json:"amount"`
	PaymentDate      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"payment_date"`
	RemainingBalance float64   `gorm:"type:numeric(15,2);not null" json:"remaining_balance"`
}

func (CreditPayment) TableName() string {
	return "credit_payments"
}
