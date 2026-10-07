package models

import (
	"github.com/google/uuid"
)

type Product struct {
	BaseModel
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"category_id"`
	BrandID     *uuid.UUID `gorm:"type:uuid;index" json:"brand_id,omitempty"`
	SupplierID  *uuid.UUID `gorm:"type:uuid;index" json:"supplier_id,omitempty"`
	SKU         string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"sku"`
	Barcode     *string    `gorm:"type:varchar(50);uniqueIndex;nullable" json:"barcode,omitempty"`
	Name        string     `gorm:"type:varchar(200);not null;index" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`

	// Relationships
	Category     *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Brand        *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Supplier     *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	ProductUnits []ProductUnit  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product_units,omitempty"`
	Images       []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
