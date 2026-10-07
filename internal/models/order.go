package models

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type Order struct {
	BaseModel
	OrderNumber       string                  `gorm:"type:varchar(40);not null;uniqueIndex" json:"order_number"`
	UserID            uuid.UUID               `gorm:"type:uuid;not null;index" json:"user_id"`
	ShippingAddressID uuid.UUID               `gorm:"type:uuid;not null;index" json:"shipping_address_id"`
	Status            constants.OrderStatus   `gorm:"type:varchar(30);not null;default:'CREATED';index" json:"status"`
	PaymentMethod     constants.PaymentMethod `gorm:"type:varchar(30);not null" json:"payment_method"`
	TotalAmount       float64                 `gorm:"type:numeric(15,2);not null" json:"total_amount"`
	ShippingFee       float64                 `gorm:"type:numeric(12,2);not null;default:0.00" json:"shipping_fee"`
	GrandTotal        float64                 `gorm:"type:numeric(15,2);not null" json:"grand_total"`
	DueDate           *time.Time              `gorm:"nullable" json:"due_date,omitempty"`
	OrderDate         time.Time               `gorm:"not null;default:CURRENT_TIMESTAMP" json:"order_date"`

	// Relationships
	User            *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderAddress    *OrderAddress   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_address,omitempty"`
	OrderItems      []OrderItem     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items,omitempty"`
	Invoice         *Invoice        `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"invoice,omitempty"`
	Histories       []OrderHistory  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"histories,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == constants.OrderStatusPendingPayment ||
		o.Status == constants.OrderStatusConfirmed ||
		o.Status == constants.OrderStatusProcessingPacking
}
