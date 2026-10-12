package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	BaseModel
	OrderID       uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductUnitID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_unit_id"`
	ProductName   string    `gorm:"type:varchar(200);not null;default:''" json:"product_name"`
	SKU           string    `gorm:"type:varchar(50);not null;default:''" json:"sku"`
	UnitName      string    `gorm:"type:varchar(50);not null;default:''" json:"unit_name"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	UnitPrice     float64   `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	Subtotal      float64   `gorm:"type:numeric(15,2);not null" json:"subtotal"`
	Notes         string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Product     *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
