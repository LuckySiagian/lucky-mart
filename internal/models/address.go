package models

import (
	"github.com/google/uuid"
)

type Address struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title         string    `gorm:"type:varchar(50);not null" json:"title"` // Rumah, Toko, Gudang
	RecipientName string    `gorm:"type:varchar(100);not null" json:"recipient_name"`
	Phone         string    `gorm:"type:varchar(20);not null" json:"phone"`
	FullAddress   string    `gorm:"type:text;not null" json:"full_address"`
	City          string    `gorm:"type:varchar(100);not null" json:"city"`
	Province      string    `gorm:"type:varchar(100);not null" json:"province"`
	PostalCode    string    `gorm:"type:varchar(10);not null" json:"postal_code"`
	IsDefault     bool      `gorm:"not null;default:false" json:"is_default"`

	// Relationship
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Address) TableName() string {
	return "addresses"
}
