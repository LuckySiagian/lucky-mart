package dto

import (
	"time"

	"github.com/google/uuid"
)

type StockAdjustmentDTO struct {
	InventoryID     *uuid.UUID `json:"inventory_id" binding:"omitempty"`
	ProductUnitID   uuid.UUID  `json:"product_unit_id" binding:"required"`
	AdjustmentType  string     `json:"adjustment_type" binding:"required,oneof=STOCK_IN STOCK_OUT ADJUSTMENT"`
	QuantityChanged int        `json:"quantity_changed" binding:"required"`
	ReferenceNumber string     `json:"reference_number" binding:"omitempty,max=60"`
	Notes           string     `json:"notes" binding:"required,min=3"`
}

type UpdateMinStockDTO struct {
	MinStockAlert int `json:"min_stock_alert" binding:"gte=0"`
}

type StockMovementResponseDTO struct {
	ID              uuid.UUID `json:"id"`
	InventoryID     uuid.UUID `json:"inventory_id"`
	ProductUnitID   uuid.UUID `json:"product_unit_id"`
	ProductName     string    `json:"product_name,omitempty"`
	SKU             string    `json:"sku,omitempty"`
	UnitName        string    `json:"unit_name,omitempty"`
	MovementType    string    `json:"movement_type"`
	ReferenceNumber string    `json:"reference_number"`
	QuantityBefore  int       `json:"quantity_before"`
	QuantityChanged int       `json:"quantity_changed"`
	QuantityAfter   int       `json:"quantity_after"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type InventoryStatusResponseDTO struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductUnitID   uuid.UUID `json:"product_unit_id"`
	WarehouseName   string    `json:"warehouse_name"`
	ProductName     string    `json:"product_name"`
	SKU             string    `json:"sku"`
	CategoryName    string    `json:"category_name,omitempty"`
	UnitName        string    `json:"unit_name"`
	RetailPrice     float64   `json:"retail_price"`
	CurrentStock    int       `json:"current_stock"`
	ReservedStock   int       `json:"reserved_stock"`
	MinStockAlert   int       `json:"min_stock_alert"`
	IsLowStock      bool      `json:"is_low_stock"`
	IsProductActive bool      `json:"is_product_active"`
	UpdatedAt       time.Time `json:"updated_at"`
}
