package models

import (
	"github.com/google/uuid"
)

type CartItem struct {
	BaseModel
	CartID        uuid.UUID `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductUnitID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_unit_id"`
	Quantity      int       `gorm:"not null;default:1;check:quantity >= 1" json:"quantity"`
	UnitPrice     float64   `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	Subtotal      float64   `gorm:"type:numeric(15,2);not null" json:"subtotal"`
	Notes         string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Product     *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductUnit *ProductUnit `gorm:"foreignKey:ProductUnitID" json:"product_unit,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
