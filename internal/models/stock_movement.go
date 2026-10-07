package models

import (
	"github.com/google/uuid"
)

type StockMovement struct {
	BaseModel
	InventoryID     uuid.UUID `gorm:"type:uuid;not null;index" json:"inventory_id"`
	ProductUnitID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_unit_id"`
	MovementType    string    `gorm:"type:varchar(30);not null;index" json:"movement_type"` // INITIAL_STOCK, STOCK_IN, STOCK_OUT, ADJUSTMENT
	ReferenceNumber string    `gorm:"type:varchar(60);not null;index" json:"reference_number"`
	QuantityBefore  int       `gorm:"not null" json:"quantity_before"`
	QuantityChanged int       `gorm:"not null" json:"quantity_changed"`
	QuantityAfter   int       `gorm:"not null" json:"quantity_after"`
	Notes           string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Inventory   *Inventory   `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

func (StockMovement) TableName() string {
	return "stock_movements"
}
