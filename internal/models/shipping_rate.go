package models

import (
	"github.com/google/uuid"
)

type ShippingRate struct {
	BaseModel
	ShippingMethodID uuid.UUID `gorm:"type:uuid;not null;index" json:"shipping_method_id"`
	City             string    `gorm:"type:varchar(100);not null;index" json:"city"`
	RatePerKg        float64   `gorm:"type:numeric(12,2);not null;default:0.00" json:"rate_per_kg"`
	FlatRate         float64   `gorm:"type:numeric(12,2);not null;default:15000.00" json:"flat_rate"`

	// Relationship
	ShippingMethod *ShippingMethod `gorm:"foreignKey:ShippingMethodID" json:"shipping_method,omitempty"`
}

func (ShippingRate) TableName() string {
	return "shipping_rates"
}
