package models

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	BaseModel
	InvoiceNumber string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"invoice_number"`
	OrderID       uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalAmount   float64    `gorm:"type:numeric(15,2);not null" json:"total_amount"`
	Status        string     `gorm:"type:varchar(20);not null;default:'UNPAID';index" json:"status"` // UNPAID, PAID, CANCELLED
	DueDate       time.Time  `gorm:"not null" json:"due_date"`
	IssuedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"issued_at"`
	PaidAt        *time.Time `gorm:"nullable" json:"paid_at,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}
