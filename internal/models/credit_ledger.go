package models

import (
	"github.com/google/uuid"
)

type CreditLedger struct {
	BaseModel
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	CreditLimit        float64   `gorm:"type:numeric(15,2);not null;default:0.00" json:"credit_limit"`
	UsedCredit         float64   `gorm:"type:numeric(15,2);not null;default:0.00;check:used_credit >= 0" json:"used_credit"`
	RemainingCredit    float64   `gorm:"type:numeric(15,2);not null;default:0.00;check:remaining_credit >= 0" json:"remaining_credit"`
	OutstandingBalance float64   `gorm:"type:numeric(15,2);not null;default:0.00;check:outstanding_balance >= 0" json:"outstanding_balance"`

	// Relationships
	User           *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreditPayments []CreditPayment `gorm:"foreignKey:LedgerID" json:"credit_payments,omitempty"`
}

func (CreditLedger) TableName() string {
	return "credit_ledgers"
}
