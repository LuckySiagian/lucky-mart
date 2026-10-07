package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationLog struct {
	BaseModel
	NotificationID uuid.UUID  `gorm:"type:uuid;not null;index" json:"notification_id"`
	Channel        string     `gorm:"type:varchar(30);not null" json:"channel"`
	Status         string     `gorm:"type:varchar(30);not null" json:"status"` // SENT, FAILED
	ErrorMessage   string     `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount     int        `gorm:"not null;default:0" json:"retry_count"`
	SentAt         *time.Time `gorm:"nullable" json:"sent_at,omitempty"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}
