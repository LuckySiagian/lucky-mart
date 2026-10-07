package models

import (
	"github.com/google/uuid"
)

type Cart struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	TotalItems int        `gorm:"not null;default:0" json:"total_items"`
	TotalAmount float64   `gorm:"type:numeric(15,2);not null;default:0.00" json:"total_amount"`

	// Relationships
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (Cart) TableName() string {
	return "carts"
}
