package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	BaseModel
	PaymentNumber   string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"payment_number"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	InvoiceID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"invoice_id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	PaymentMethodID uuid.UUID  `gorm:"type:uuid;not null;index" json:"payment_method_id"`
	Amount          float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	AdminFee        float64    `gorm:"type:numeric(12,2);not null;default:0.00" json:"admin_fee"`
	Discount        float64    `gorm:"type:numeric(12,2);not null;default:0.00" json:"discount"`
	GrandTotal      float64    `gorm:"type:numeric(15,2);not null" json:"grand_total"`
	Status          string     `gorm:"type:varchar(30);not null;default:'PENDING';index" json:"status"` // PENDING, WAITING_CONFIRMATION, PAID, FAILED, EXPIRED, REFUNDED, PARTIALLY_REFUNDED, CANCELLED
	ProofImageURL   *string    `gorm:"type:varchar(500);nullable" json:"proof_image_url,omitempty"`
	PaidAt          *time.Time `gorm:"nullable" json:"paid_at,omitempty"`
	ExpiredAt       time.Time  `gorm:"not null" json:"expired_at"`
	Notes           string     `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Order         *Order               `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Invoice       *Invoice             `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	User          *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PaymentMethod *PaymentMethod       `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Transactions  []PaymentTransaction `gorm:"foreignKey:PaymentID" json:"transactions,omitempty"`
	Logs          []PaymentLog         `gorm:"foreignKey:PaymentID" json:"logs,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}
