package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentTransaction struct {
	BaseModel
	PaymentID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"payment_id"`
	GatewayName     string     `gorm:"type:varchar(50);not null" json:"gateway_name"` // MANUAL, MIDTRANS, XENDIT
	TransactionID   string     `gorm:"type:varchar(100);not null;index" json:"transaction_id"`
	ReferenceNumber string     `gorm:"type:varchar(100);not null" json:"reference_number"`
	RequestPayload  string     `gorm:"type:text" json:"request_payload,omitempty"`
	ResponsePayload string     `gorm:"type:text" json:"response_payload,omitempty"`
	Signature       string     `gorm:"type:varchar(255)" json:"signature,omitempty"`
	Status          string     `gorm:"type:varchar(30);not null" json:"status"`
	GatewayStatus   string     `gorm:"type:varchar(50)" json:"gateway_status,omitempty"`
	PaidAt          *time.Time `gorm:"nullable" json:"paid_at,omitempty"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}
