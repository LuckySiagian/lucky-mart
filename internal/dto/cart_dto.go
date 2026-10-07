package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddToCartDTO struct {
	ProductID     uuid.UUID `json:"product_id" binding:"required"`
	ProductUnitID uuid.UUID `json:"product_unit_id" binding:"required"`
	Quantity      int       `json:"quantity" binding:"required,gte=1"`
	Notes         string    `json:"notes" binding:"omitempty,max=255"`
}

type UpdateCartItemDTO struct {
	Quantity int    `json:"quantity" binding:"required,gte=1"`
	Notes    string `json:"notes" binding:"omitempty,max=255"`
}

type CartItemResponseDTO struct {
	ID                 uuid.UUID `json:"id"`
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	SKU                string    `json:"sku"`
	ProductUnitID      uuid.UUID `json:"product_unit_id"`
	UnitName           string    `json:"unit_name"`
	Quantity           int       `json:"quantity"`
	UnitPrice          float64   `json:"unit_price"`
	Subtotal           float64   `json:"subtotal"`
	Notes              string    `json:"notes,omitempty"`
	AvailableStock     int       `json:"available_stock"`
	IsStockSufficient  bool      `json:"is_stock_sufficient"`
	IsProductActive    bool      `json:"is_product_active"`
}

type CartResponseDTO struct {
	ID          uuid.UUID             `json:"id"`
	UserID      uuid.UUID             `json:"user_id"`
	TotalItems  int                   `json:"total_items"`
	TotalAmount float64               `json:"total_amount"`
	Items       []CartItemResponseDTO `json:"items"`
	CreatedAt   time.Time             `json:"created_at"`
}

type CheckoutSummaryDTO struct {
	CartID               uuid.UUID             `json:"cart_id"`
	UserID               uuid.UUID             `json:"user_id"`
	TotalItems           int                   `json:"total_items"`
	SubtotalAmount       float64               `json:"subtotal_amount"`
	EstimatedShippingFee float64               `json:"estimated_shipping_fee"`
	GrandTotalAmount     float64               `json:"grand_total_amount"`
	Items                []CartItemResponseDTO `json:"items"`
	CanProceedToOrder    bool                  `json:"can_proceed_to_order"`
	ValidationErrors     []string              `json:"validation_errors,omitempty"`
}
