package models

import (
	"github.com/google/uuid"
)

type ShippingMethod struct {
	BaseModel
	CourierID     uuid.UUID `gorm:"type:uuid;not null;index" json:"courier_id"`
	Code          string    `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"` // SAME_DAY, REGULAR, CARGO
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`
	EstimatedDays string    `gorm:"type:varchar(20);not null" json:"estimated_days"` // 1-2 Hari
	BaseRate      float64   `gorm:"type:numeric(12,2);not null;default:0.00" json:"base_rate"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`

	// Relationship
	Courier *Courier `gorm:"foreignKey:CourierID" json:"courier,omitempty"`
}

func (ShippingMethod) TableName() string {
	return "shipping_methods"
}
