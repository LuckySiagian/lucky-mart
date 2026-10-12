package models

import (
	"github.com/google/uuid"
)

type OrderAddress struct {
	BaseModel
	OrderID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	RecipientName string    `gorm:"type:varchar(100);not null" json:"recipient_name"`
	Phone         string    `gorm:"type:varchar(20);not null" json:"phone"`
	FullAddress   string    `gorm:"type:text;not null" json:"full_address"`
	City          string    `gorm:"type:varchar(100);not null" json:"city"`
	PostalCode    string    `gorm:"type:varchar(10);not null;default:''" json:"postal_code"`
	Notes         string    `gorm:"type:text" json:"notes,omitempty"`
}

func (OrderAddress) TableName() string {
	return "order_addresses"
}
