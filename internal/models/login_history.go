package models

import (
	"github.com/google/uuid"
)

type LoginHistory struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	IPAddress string    `gorm:"type:varchar(50);not null" json:"ip_address"`
	UserAgent string    `gorm:"type:varchar(500);not null" json:"user_agent"`
	Device    string    `gorm:"type:varchar(100);not null" json:"device"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"` // SUCCESS, FAILED, LOCKED
}

func (LoginHistory) TableName() string {
	return "login_histories"
}
