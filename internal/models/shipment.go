package models

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	BaseModel
	TrackingNumber   string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"tracking_number"`
	OrderID          uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	CourierID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"courier_id"`
	ShippingMethodID uuid.UUID  `gorm:"type:uuid;not null;index" json:"shipping_method_id"`
	Status           string     `gorm:"type:varchar(30);not null;default:'PENDING';index" json:"status"` // PENDING, READY_TO_SHIP, PICKED_UP, IN_TRANSIT, OUT_FOR_DELIVERY, DELIVERED, FAILED, RETURNED
	ShippingFee      float64    `gorm:"type:numeric(12,2);not null" json:"shipping_fee"`
	ShippedAt        *time.Time `gorm:"nullable" json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time `gorm:"nullable" json:"delivered_at,omitempty"`

	// Relationships
	Order          *Order             `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Courier        *Courier           `gorm:"foreignKey:CourierID" json:"courier,omitempty"`
	ShippingMethod *ShippingMethod    `gorm:"foreignKey:ShippingMethodID" json:"shipping_method,omitempty"`
	Trackings      []ShipmentTracking `gorm:"foreignKey:ShipmentID;constraint:OnDelete:CASCADE" json:"trackings,omitempty"`
}

func (Shipment) TableName() string {
	return "shipments"
}
