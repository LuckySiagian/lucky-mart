package models

import (
	"github.com/google/uuid"
)

type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ImageURL  string    `gorm:"type:varchar(500);not null" json:"image_url"`
	IsPrimary bool      `gorm:"not null;default:false" json:"is_primary"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
