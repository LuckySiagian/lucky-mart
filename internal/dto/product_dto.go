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
	CategoryID    uuid.UUID              `json:"category_id" binding:"required"`
	BrandID       *uuid.UUID             `json:"brand_id" binding:"omitempty"`
	SupplierID    *uuid.UUID             `json:"supplier_id" binding:"omitempty"`
	SKU           string                 `json:"sku" binding:"omitempty,max=50"`
	Barcode       *string                `json:"barcode" binding:"omitempty"`
	Name          string                 `json:"name" binding:"required,min=2,max=200"`
	Description   string                 `json:"description" binding:"omitempty"`
	ImageURL      string                 `json:"image_url" binding:"omitempty,max=500"`
	MinStockAlert *int                   `json:"min_stock_alert" binding:"omitempty,gte=0"`
	Units         []CreateProductUnitDTO `json:"units" binding:"required,min=1,dive"`
}

type CreateProductUnitDTO struct {
	UnitID           *uuid.UUID `json:"unit_id" binding:"omitempty"`
	UnitName         string     `json:"unit_name" binding:"required,max=50"`
	ConversionFactor int        `json:"conversion_factor" binding:"omitempty,gte=1"`
	RetailPrice      float64    `json:"retail_price" binding:"required,gt=0"`
	WholesalePrice   float64    `json:"wholesale_price" binding:"omitempty,gte=0"`
	MinWholesaleQty  int        `json:"min_wholesale_qty" binding:"omitempty,gte=1"`
	InitialStock     int        `json:"initial_stock" binding:"gte=0"`
}

// UpdateProductDTO updates product details. Price and stock changes may be
// sent per unit in Units, or through the flat retail_price / wholesale_price /
// stock_quantity fields, which apply to the product's first (smallest) unit.
// A stock change is recorded as an ADJUSTMENT stock movement.
type UpdateProductDTO struct {
	CategoryID     *uuid.UUID             `json:"category_id" binding:"omitempty"`
	SKU            *string                `json:"sku" binding:"omitempty,min=2,max=50"`
	Name           *string                `json:"name" binding:"omitempty,min=2,max=200"`
	Description    *string                `json:"description" binding:"omitempty"`
	ImageURL       *string                `json:"image_url" binding:"omitempty,max=500"`
	IsActive       *bool                  `json:"is_active" binding:"omitempty"`
	RetailPrice    *float64               `json:"retail_price" binding:"omitempty,gt=0"`
	WholesalePrice *float64               `json:"wholesale_price" binding:"omitempty,gte=0"`
	StockQuantity  *int                   `json:"stock_quantity" binding:"omitempty,gte=0"`
	Units          []UpdateProductUnitDTO `json:"units" binding:"omitempty,dive"`
}

type UpdateProductUnitDTO struct {
	ID              uuid.UUID `json:"id" binding:"required"`
	UnitName        *string   `json:"unit_name" binding:"omitempty,max=50"`
	RetailPrice     *float64  `json:"retail_price" binding:"omitempty,gt=0"`
	WholesalePrice  *float64  `json:"wholesale_price" binding:"omitempty,gte=0"`
	MinWholesaleQty *int      `json:"min_wholesale_qty" binding:"omitempty,gte=1"`
	StockQuantity   *int      `json:"stock_quantity" binding:"omitempty,gte=0"`
}

type ProductResponseDTO struct {
	ID          uuid.UUID              `json:"id"`
	SKU         string                 `json:"sku"`
	Barcode     *string                `json:"barcode,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	IsActive    bool                   `json:"is_active"`
	CategoryID  uuid.UUID              `json:"category_id"`
	Category    CategoryResponse       `json:"category"`
	Brand       *BrandResponse         `json:"brand,omitempty"`
	Supplier    *SupplierResponse      `json:"supplier,omitempty"`
	Units       []ProductUnitResponse  `json:"units"`
	Images      []ProductImageResponse `json:"images"`
	ImageURL    string                 `json:"image_url,omitempty"`
	TotalStock  int                    `json:"total_stock"`
	CreatedAt   time.Time              `json:"created_at"`
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
