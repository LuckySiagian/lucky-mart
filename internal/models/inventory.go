package models

import (
	"github.com/google/uuid"
)

type Inventory struct {
	BaseModel
	WarehouseID   uuid.UUID `gorm:"type:uuid;not null;index" json:"warehouse_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductUnitID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_unit_id"`
	CurrentStock  int       `gorm:"not null;default:0;check:current_stock >= 0" json:"current_stock"`
	ReservedStock int       `gorm:"not null;default:0;check:reserved_stock >= 0" json:"reserved_stock"`
	MinStockAlert int       `gorm:"not null;default:10" json:"min_stock_alert"`

	// Relationships
	Warehouse   *Warehouse   `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Product     *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

func (Inventory) TableName() string {
	return "inventories"
}

func (i *Inventory) IsLowStock() bool {
	return i.CurrentStock <= i.MinStockAlert
}
