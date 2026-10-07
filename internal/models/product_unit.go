package models

import (
	"github.com/google/uuid"
)

type ProductUnit struct {
	BaseModel
	ProductID        uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	UnitID           uuid.UUID `gorm:"type:uuid;not null;index" json:"unit_id"`
	UnitName         string    `gorm:"type:varchar(50);not null" json:"unit_name"`
	ConversionFactor int       `gorm:"not null;default:1" json:"conversion_factor"`
	RetailPrice      float64   `gorm:"type:numeric(12,2);not null;check:retail_price > 0" json:"retail_price"`
	WholesalePrice   float64   `gorm:"type:numeric(12,2);not null;check:wholesale_price > 0" json:"wholesale_price"`
	MinWholesaleQty  int       `gorm:"not null;default:5" json:"min_wholesale_qty"`
	StockQuantity    int       `gorm:"not null;default:0;check:stock_quantity >= 0" json:"stock_quantity"`

	// Relationships
	Unit *Unit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

func (ProductUnit) TableName() string {
	return "product_units"
}
