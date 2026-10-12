package models

import (
	"github.com/google/uuid"
)

type EventLog struct {
	BaseModel
	EventName   string    `gorm:"type:varchar(100);not null;index" json:"event_name"` // USER_REGISTERED, ORDER_CREATED, PAYMENT_PAID, etc.
	Aggregate   string    `gorm:"type:varchar(50);not null;index" json:"aggregate"`   // USER, ORDER, PAYMENT, SHIPMENT
	AggregateID uuid.UUID `gorm:"type:uuid;not null;index" json:"aggregate_id"`
	Payload     string    `gorm:"type:text;not null" json:"payload"`
	Status      string    `gorm:"type:varchar(30);not null;default:'PUBLISHED';index" json:"status"` // PUBLISHED, PROCESSED, FAILED
}

func (EventLog) TableName() string {
	return "event_logs"
}
