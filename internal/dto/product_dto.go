package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

type CreateProductDTO struct {
	CategoryID  uuid.UUID              `json:"category_id" binding:"required"`
	BrandID     *uuid.UUID             `json:"brand_id" binding:"omitempty"`
	SupplierID  *uuid.UUID             `json:"supplier_id" binding:"omitempty"`
	SKU         string                 `json:"sku" binding:"required"`
	Barcode     *string                `json:"barcode" binding:"omitempty"`
	Name        string                 `json:"name" binding:"required,min=2,max=200"`
	Description string                 `json:"description" binding:"omitempty"`
	ImageURL    string                 `json:"image_url" binding:"omitempty"`
	Units       []CreateProductUnitDTO `json:"units" binding:"required,min=1"`
}

type UpdateProductDTO struct {
	CategoryID     *uuid.UUID `json:"category_id" binding:"omitempty"`
	SKU            string     `json:"sku" binding:"omitempty"`
	Name           string     `json:"name" binding:"omitempty"`
	Description    string     `json:"description" binding:"omitempty"`
	ImageURL       string     `json:"image_url" binding:"omitempty"`
	RetailPrice    float64    `json:"retail_price" binding:"omitempty"`
	WholesalePrice float64    `json:"wholesale_price" binding:"omitempty"`
	StockQuantity  int        `json:"stock_quantity" binding:"omitempty"`
}

type CreateProductUnitDTO struct {
	UnitID           uuid.UUID `json:"unit_id" binding:"required"`
	UnitName         string    `json:"unit_name" binding:"required"`
	ConversionFactor int       `json:"conversion_factor" binding:"required,gte=1"`
	RetailPrice      float64   `json:"retail_price" binding:"required,gt=0"`
	WholesalePrice   float64   `json:"wholesale_price" binding:"required,gt=0"`
	MinWholesaleQty  int       `json:"min_wholesale_qty" binding:"required,gte=1"`
	InitialStock     int       `json:"initial_stock" binding:"required,gte=0"`
}

type ProductResponseDTO struct {
	ID          uuid.UUID          `json:"id"`
	SKU         string             `json:"sku"`
	Barcode     *string            `json:"barcode,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	IsActive    bool               `json:"is_active"`
	Category    CategoryResponse   `json:"category"`
	Brand       *BrandResponse     `json:"brand,omitempty"`
	Supplier    *SupplierResponse  `json:"supplier,omitempty"`
	Units       []ProductUnitResponse `json:"units"`
	Images      []ProductImageResponse `json:"images,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
}

type BrandResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SupplierResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductUnitResponse struct {
	ID               uuid.UUID `json:"id"`
	UnitID           uuid.UUID `json:"unit_id"`
	UnitName         string    `json:"unit_name"`
	ConversionFactor int       `json:"conversion_factor"`
	RetailPrice      float64   `json:"retail_price"`
	WholesalePrice   float64   `json:"wholesale_price"`
	MinWholesaleQty  int       `json:"min_wholesale_qty"`
	StockQuantity    int       `json:"stock_quantity"`
}

type ProductImageResponse struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
}
