package dto

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type CreateOrderDTO struct {
	PaymentMethod constants.PaymentMethod `json:"payment_method" binding:"required"`
	RecipientName string                  `json:"recipient_name" binding:"required,min=3,max=100"`
	Phone         string                  `json:"phone" binding:"required,phone_id"`
	FullAddress   string                  `json:"full_address" binding:"required"`
	City          string                  `json:"city" binding:"required"`
	PostalCode    string                  `json:"postal_code" binding:"required"`
}

type CancelOrderDTO struct {
	Reason string `json:"reason" binding:"required,min=3,max=255"`
}

type UpdateOrderStatusDTO struct {
	Status constants.OrderStatus `json:"status" binding:"required"`
	Notes  string                `json:"notes,omitempty"`
}


type OrderItemResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	SKU           string    `json:"sku"`
	ProductUnitID uuid.UUID `json:"product_unit_id"`
	UnitName      string    `json:"unit_name"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float64   `json:"unit_price"`
	Subtotal      float64   `json:"subtotal"`
	Notes         string    `json:"notes,omitempty"`
}

type OrderAddressResponseDTO struct {
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	FullAddress   string `json:"full_address"`
	City          string `json:"city"`
	PostalCode    string `json:"postal_code"`
}

type InvoiceResponseDTO struct {
	ID            uuid.UUID  `json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	OrderID       uuid.UUID  `json:"order_id"`
	TotalAmount   float64    `json:"total_amount"`
	Status        string     `json:"status"`
	DueDate       time.Time  `json:"due_date"`
	IssuedAt      time.Time  `json:"issued_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

type OrderHistoryResponseDTO struct {
	ID        uuid.UUID             `json:"id"`
	Status    constants.OrderStatus `json:"status"`
	Notes     string                `json:"notes,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
}

type OrderResponseDTO struct {
	ID            uuid.UUID               `json:"id"`
	OrderNumber   string                  `json:"order_number"`
	UserID        uuid.UUID               `json:"user_id"`
	Status        constants.OrderStatus   `json:"status"`
	PaymentMethod constants.PaymentMethod `json:"payment_method"`
	TotalAmount   float64                 `json:"total_amount"`
	ShippingFee   float64                 `json:"shipping_fee"`
	GrandTotal    float64                 `json:"grand_total"`
	DueDate       *time.Time              `json:"due_date,omitempty"`
	OrderDate     time.Time               `json:"order_date"`
	Items         []OrderItemResponseDTO  `json:"items,omitempty"`
	Address       *OrderAddressResponseDTO `json:"address,omitempty"`
	Invoice       *InvoiceResponseDTO     `json:"invoice,omitempty"`
	Histories     []OrderHistoryResponseDTO `json:"histories,omitempty"`
}
