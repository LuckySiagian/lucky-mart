package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	BaseModel
	UserID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Title   string     `gorm:"type:varchar(200);not null" json:"title"`
	Message string     `gorm:"type:text;not null" json:"message"`
	Type    string     `gorm:"type:varchar(50);not null;index" json:"type"`    // ORDER_CREATED, PAYMENT_PAID, etc.
	Channel string     `gorm:"type:varchar(30);not null;index" json:"channel"` // EMAIL, WHATSAPP, PUSH_NOTIFICATION, IN_APP
	Status  string     `gorm:"type:varchar(30);not null;default:'PENDING';index" json:"status"` // PENDING, SENT, FAILED
	IsRead  bool       `gorm:"not null;default:false;index" json:"is_read"`
	ReadAt  *time.Time `gorm:"nullable" json:"read_at,omitempty"`

	// Relationships
	User *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Logs []NotificationLog `gorm:"foreignKey:NotificationID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
