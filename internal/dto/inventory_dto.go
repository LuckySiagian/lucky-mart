package dto

import (
	"time"

	"github.com/google/uuid"
)

type StockAdjustmentDTO struct {
	InventoryID     uuid.UUID `json:"inventory_id" binding:"required"`
	ProductUnitID   uuid.UUID `json:"product_unit_id" binding:"required"`
	AdjustmentType  string    `json:"adjustment_type" binding:"required,oneof=STOCK_IN STOCK_OUT ADJUSTMENT"`
	QuantityChanged int       `json:"quantity_changed" binding:"required"`
	ReferenceNumber string    `json:"reference_number" binding:"required"`
	Notes           string    `json:"notes" binding:"required"`
}

type StockMovementResponseDTO struct {
	ID              uuid.UUID `json:"id"`
	InventoryID     uuid.UUID `json:"inventory_id"`
	MovementType    string    `json:"movement_type"`
	ReferenceNumber string    `json:"reference_number"`
	QuantityBefore  int       `json:"quantity_before"`
	QuantityChanged int       `json:"quantity_changed"`
	QuantityAfter   int       `json:"quantity_after"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type InventoryStatusResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	WarehouseName string    `json:"warehouse_name"`
	ProductName   string    `json:"product_name"`
	SKU           string    `json:"sku"`
	UnitName      string    `json:"unit_name"`
	CurrentStock  int       `json:"current_stock"`
	ReservedStock int       `json:"reserved_stock"`
	MinStockAlert int       `json:"min_stock_alert"`
	IsLowStock    bool      `json:"is_low_stock"`
}
