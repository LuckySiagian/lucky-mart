package models

import (
	"github.com/google/uuid"
)

type ShipmentTracking struct {
	BaseModel
	ShipmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"shipment_id"`
	Status     string    `gorm:"type:varchar(30);not null" json:"status"`
	Location   string    `gorm:"type:varchar(150);not null" json:"location"`
	Notes      string    `gorm:"type:text" json:"notes,omitempty"`
}

func (ShipmentTracking) TableName() string {
	return "shipment_trackings"
}
