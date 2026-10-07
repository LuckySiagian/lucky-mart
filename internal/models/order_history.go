package models

import (
	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type OrderHistory struct {
	BaseModel
	OrderID uuid.UUID             `gorm:"type:uuid;not null;index" json:"order_id"`
	Status  constants.OrderStatus `gorm:"type:varchar(30);not null" json:"status"`
	Notes   string                `gorm:"type:text" json:"notes,omitempty"`
}

func (OrderHistory) TableName() string {
	return "order_histories"
}
