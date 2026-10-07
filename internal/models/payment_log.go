package models

import (
	"github.com/google/uuid"
)

type PaymentLog struct {
	BaseModel
	PaymentID uuid.UUID `gorm:"type:uuid;not null;index" json:"payment_id"`
	Action    string    `gorm:"type:varchar(50);not null" json:"action"`
	OldStatus string    `gorm:"type:varchar(30);not null" json:"old_status"`
	NewStatus string    `gorm:"type:varchar(30);not null" json:"new_status"`
	Notes     string    `gorm:"type:text" json:"notes,omitempty"`
}

func (PaymentLog) TableName() string {
	return "payment_logs"
}
