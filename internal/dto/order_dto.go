package dto

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

// CreateOrderDTO places an order from the customer's cart. When CartItemIDs
// is empty every item in the cart is ordered. For delivery the address comes
// from a saved address (AddressID) or from the inline fields; a PICKUP order
// needs no address and has no shipping fee.
type CreateOrderDTO struct {
	PaymentMethod    constants.PaymentMethod `json:"payment_method" binding:"required"`
	DeliveryType     string                  `json:"delivery_type" binding:"omitempty,oneof=DELIVERY PICKUP"`
	CartItemIDs      []uuid.UUID             `json:"cart_item_ids" binding:"omitempty"`
	ShippingMethodID *uuid.UUID              `json:"shipping_method_id" binding:"omitempty"`
	AddressID        *uuid.UUID              `json:"address_id" binding:"omitempty"`
	RecipientName    string                  `json:"recipient_name" binding:"omitempty,max=100"`
	Phone            string                  `json:"phone" binding:"omitempty,max=20"`
	FullAddress      string                  `json:"full_address" binding:"omitempty"`
	City             string                  `json:"city" binding:"omitempty,max=100"`
	PostalCode       string                  `json:"postal_code" binding:"omitempty,max=10"`
	Notes            string                  `json:"notes" binding:"omitempty,max=500"`
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
	ImageURL      string    `json:"image_url,omitempty"`
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
	Notes         string `json:"notes,omitempty"`
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

type OrderCustomerDTO struct {
	ID          uuid.UUID          `json:"id"`
	FullName    string             `json:"full_name"`
	PhoneNumber string             `json:"phone_number"`
	Role        constants.UserRole `json:"role"`
}

type OrderResponseDTO struct {
	ID            uuid.UUID                 `json:"id"`
	OrderNumber   string                    `json:"order_number"`
	UserID        uuid.UUID                 `json:"user_id"`
	Customer      *OrderCustomerDTO         `json:"customer,omitempty"`
	Status        constants.OrderStatus     `json:"status"`
	NextStatuses  []constants.OrderStatus   `json:"next_statuses"`
	CanCancel     bool                      `json:"can_cancel"`
	PaymentMethod constants.PaymentMethod   `json:"payment_method"`
	DeliveryType  string                    `json:"delivery_type"`
	TotalAmount   float64                   `json:"total_amount"`
	ShippingFee   float64                   `json:"shipping_fee"`
	GrandTotal    float64                   `json:"grand_total"`
	TotalItems    int                       `json:"total_items"`
	DueDate       *time.Time                `json:"due_date,omitempty"`
	OrderDate     time.Time                 `json:"order_date"`
	Items         []OrderItemResponseDTO    `json:"items"`
	Address       *OrderAddressResponseDTO  `json:"address,omitempty"`
	Invoice       *InvoiceResponseDTO       `json:"invoice,omitempty"`
	Payments      []PaymentResponseDTO      `json:"payments,omitempty"`
	Histories     []OrderHistoryResponseDTO `json:"histories"`
}
