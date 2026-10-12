package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSKUAlreadyExists      = errors.New("SKU sudah dipakai produk lain")
	ErrCategoryNotFound      = errors.New("kategori tidak ditemukan")
	ErrCategoryAlreadyExists = errors.New("nama kategori sudah ada")
	ErrCategoryInUse         = errors.New("kategori masih dipakai produk; pindahkan atau hapus produknya dulu")
	ErrProductNotFound       = errors.New("produk tidak ditemukan")
)

type ProductService interface {
	GetProducts(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID, includeInactive bool) (dto.PaginatedResponse, error)
	GetProductByID(ctx context.Context, id uuid.UUID, includeInactive bool) (*dto.ProductResponseDTO, error)
	CreateProduct(ctx context.Context, req dto.CreateProductDTO) (*dto.ProductResponseDTO, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductDTO) (*dto.ProductResponseDTO, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetCategories(ctx context.Context) ([]dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, name, description string) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name, description string) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	db            *gorm.DB
	productRepo   repositories.ProductRepository
	categoryRepo  repositories.CategoryRepository
	inventoryRepo repositories.InventoryRepository
	stockMoveRepo repositories.StockMovementRepository
}

func NewProductService(
	db *gorm.DB,
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	inventoryRepo repositories.InventoryRepository,
	stockMoveRepo repositories.StockMovementRepository,
) ProductService {
	return &productService{
		db:            db,
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		inventoryRepo: inventoryRepo,
		stockMoveRepo: stockMoveRepo,
	}
}

func (s *productService) GetProducts(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID, includeInactive bool) (dto.PaginatedResponse, error) {
	products, total, err := s.productRepo.FindAll(ctx, req, categoryID, includeInactive)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.ProductResponseDTO, len(products))
	for i := range products {
		dtos[i] = mapProductToDTO(&products[i])
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID, includeInactive bool) (*dto.ProductResponseDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil || (!product.IsActive && !includeInactive) {
		return nil, ErrProductNotFound
	}
	res := mapProductToDTO(product)
	return &res, nil
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductDTO) (*dto.ProductResponseDTO, error) {
	cat, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	sku := strings.ToUpper(strings.TrimSpace(req.SKU))
	if sku == "" {
		sku = newReference("SKU")
	}
	if existing, err := s.productRepo.FindBySKU(ctx, sku); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, ErrSKUAlreadyExists
	}

	var productID uuid.UUID
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := &models.Product{
			CategoryID:  req.CategoryID,
			BrandID:     req.BrandID,
			SupplierID:  req.SupplierID,
			SKU:         sku,
			Barcode:     req.Barcode,
			Name:        strings.TrimSpace(req.Name),
			Description: strings.TrimSpace(req.Description),
			IsActive:    true,
		}
		if err := s.productRepo.Create(ctx, tx, product); err != nil {
			return fmt.Errorf("gagal menyimpan produk: %w", err)
		}
		productID = product.ID

		warehouse, err := defaultWarehouse(tx)
		if err != nil {
			return err
		}

		minStock := 10
		if req.MinStockAlert != nil {
			minStock = *req.MinStockAlert
		}

		for _, u := range req.Units {
			unitID, err := resolveUnitID(tx, u.UnitID, u.UnitName)
			if err != nil {
				return err
			}

			conversion := u.ConversionFactor
			if conversion < 1 {
				conversion = 1
			}
			wholesale := u.WholesalePrice
			if wholesale <= 0 {
				wholesale = u.RetailPrice
			}
			minWholesale := u.MinWholesaleQty
			if minWholesale < 1 {
				minWholesale = 5
			}

			unit := &models.ProductUnit{
				ProductID:        product.ID,
				UnitID:           unitID,
				UnitName:         strings.TrimSpace(u.UnitName),
				ConversionFactor: conversion,
				RetailPrice:      u.RetailPrice,
				WholesalePrice:   wholesale,
				MinWholesaleQty:  minWholesale,
				StockQuantity:    0,
			}
			if err := tx.Create(unit).Error; err != nil {
				return fmt.Errorf("gagal menyimpan satuan produk: %w", err)
			}

			inv := &models.Inventory{
				WarehouseID:   warehouse.ID,
				ProductID:     product.ID,
				ProductUnitID: unit.ID,
				CurrentStock:  0,
				MinStockAlert: minStock,
			}
			if err := tx.Create(inv).Error; err != nil {
				return err
			}

			if u.InitialStock > 0 {
				if _, err := applyStockChange(tx, unit.ID, u.InitialStock, constants.MovementInitialStock, sku, "Stok awal produk baru"); err != nil {
					return err
				}
			}
		}

		if url := strings.TrimSpace(req.ImageURL); url != "" {
			if !IsAllowedImageURL(url) {
				return errInvalidProductImage
			}
			if err := tx.Create(&models.ProductImage{ProductID: product.ID, ImageURL: url, IsPrimary: true}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetProductByID(ctx, productID, true)
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductDTO) (*dto.ProductResponseDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	if req.CategoryID != nil && *req.CategoryID != product.CategoryID {
		cat, err := s.categoryRepo.FindByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, ErrCategoryNotFound
		}
		product.CategoryID = *req.CategoryID
	}
	if req.SKU != nil {
		sku := strings.ToUpper(strings.TrimSpace(*req.SKU))
		if sku != product.SKU {
			if existing, err := s.productRepo.FindBySKU(ctx, sku); err != nil {
				return nil, err
			} else if existing != nil && existing.ID != product.ID {
				return nil, ErrSKUAlreadyExists
			}
			product.SKU = sku
		}
	}
	if req.Name != nil {
		product.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		product.Description = strings.TrimSpace(*req.Description)
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Flat price/stock fields target the first (smallest) unit.
	unitUpdates := req.Units
	if (req.RetailPrice != nil || req.WholesalePrice != nil || req.StockQuantity != nil) && len(product.ProductUnits) > 0 {
		unitUpdates = append(unitUpdates, dto.UpdateProductUnitDTO{
			ID:             product.ProductUnits[0].ID,
			RetailPrice:    req.RetailPrice,
			WholesalePrice: req.WholesalePrice,
			StockQuantity:  req.StockQuantity,
		})
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.productRepo.Update(ctx, tx, product); err != nil {
			return fmt.Errorf("gagal memperbarui produk: %w", err)
		}

		for _, u := range unitUpdates {
			var unit *models.ProductUnit
			for i := range product.ProductUnits {
				if product.ProductUnits[i].ID == u.ID {
					unit = &product.ProductUnits[i]
					break
				}
			}
			if unit == nil {
				return ErrProductUnitNotFound
			}

			updates := map[string]interface{}{}
			if u.UnitName != nil && strings.TrimSpace(*u.UnitName) != "" {
				updates["unit_name"] = strings.TrimSpace(*u.UnitName)
			}
			if u.RetailPrice != nil {
				updates["retail_price"] = *u.RetailPrice
			}
			if u.WholesalePrice != nil {
				wholesale := *u.WholesalePrice
				if wholesale <= 0 {
					wholesale = unit.RetailPrice
					if u.RetailPrice != nil {
						wholesale = *u.RetailPrice
					}
				}
				updates["wholesale_price"] = wholesale
			}
			if u.MinWholesaleQty != nil {
				updates["min_wholesale_qty"] = *u.MinWholesaleQty
			}
			if len(updates) > 0 {
				if err := tx.Model(&models.ProductUnit{}).Where("id = ?", unit.ID).Updates(updates).Error; err != nil {
					return err
				}
			}

			if u.StockQuantity != nil && *u.StockQuantity != unit.StockQuantity {
				var current models.ProductUnit
				if err := tx.Where("id = ?", unit.ID).First(&current).Error; err != nil {
					return err
				}
				delta := *u.StockQuantity - current.StockQuantity
				if delta != 0 {
					if _, err := applyStockChange(tx, unit.ID, delta, constants.MovementAdjustment, newReference("ADJ"), "Koreksi stok dari form edit produk"); err != nil {
						return err
					}
				}
			}
		}

		if req.ImageURL != nil {
			url := strings.TrimSpace(*req.ImageURL)
			if url != "" && !IsAllowedImageURL(url) {
				return errInvalidProductImage
			}
			if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
				return err
			}
			if url != "" {
				if err := tx.Create(&models.ProductImage{ProductID: product.ID, ImageURL: url, IsPrimary: true}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetProductByID(ctx, id, true)
}

// DeleteProduct soft-deletes the product. Past orders keep their item
// snapshots; the product disappears from the catalog and carts.
func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Product{}, id).Error
	})
}

// resolveUnitID finds the master unit by id or name, creating it when new.
func resolveUnitID(tx *gorm.DB, unitID *uuid.UUID, unitName string) (uuid.UUID, error) {
	if unitID != nil && *unitID != uuid.Nil {
		var unit models.Unit
		if err := tx.Where("id = ?", *unitID).First(&unit).Error; err == nil {
			return unit.ID, nil
		}
	}

	name := strings.TrimSpace(unitName)
	var unit models.Unit
	err := tx.Where("name ILIKE ? OR short_name ILIKE ?", name, name).First(&unit).Error
	if err == nil {
		return unit.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	short := strings.ToLower(name)
	if len(short) > 15 {
		short = short[:15]
	}
	unit = models.Unit{Name: name, ShortName: short}
	if err := tx.Create(&unit).Error; err != nil {
		return uuid.Nil, fmt.Errorf("gagal membuat satuan '%s': %w", name, err)
	}
	return unit.ID, nil
}

func mapProductToDTO(p *models.Product) dto.ProductResponseDTO {
	res := dto.ProductResponseDTO{
		ID:          p.ID,
		SKU:         p.SKU,
		Barcode:     p.Barcode,
		Name:        p.Name,
		Description: p.Description,
		IsActive:    p.IsActive,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
	}

	if p.Category != nil {
		res.Category = dto.CategoryResponse{
			ID:          p.Category.ID,
			Name:        p.Category.Name,
			Slug:        p.Category.Slug,
			Description: p.Category.Description,
		}
	}
	if p.Brand != nil {
		res.Brand = &dto.BrandResponse{ID: p.Brand.ID, Name: p.Brand.Name}
	}
	if p.Supplier != nil {
		res.Supplier = &dto.SupplierResponse{ID: p.Supplier.ID, Name: p.Supplier.Name}
	}

	res.Units = make([]dto.ProductUnitResponse, len(p.ProductUnits))
	for i, u := range p.ProductUnits {
		res.Units[i] = dto.ProductUnitResponse{
			ID:               u.ID,
			UnitID:           u.UnitID,
			UnitName:         u.UnitName,
			ConversionFactor: u.ConversionFactor,
			RetailPrice:      u.RetailPrice,
			WholesalePrice:   u.WholesalePrice,
			MinWholesaleQty:  u.MinWholesaleQty,
			StockQuantity:    u.StockQuantity,
		}
		res.TotalStock += u.StockQuantity
	}

	res.Images = make([]dto.ProductImageResponse, len(p.Images))
	for i, img := range p.Images {
		res.Images[i] = dto.ProductImageResponse{ID: img.ID, ImageURL: img.ImageURL, IsPrimary: img.IsPrimary}
		if img.IsPrimary || res.ImageURL == "" {
			res.ImageURL = img.ImageURL
		}
	}

	return res
}

func (s *productService) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		dtos[i] = dto.CategoryResponse{ID: c.ID, Name: c.Name, Slug: c.Slug, Description: c.Description}
	}
	return dtos, nil
}

func (s *productService) CreateCategory(ctx context.Context, name, description string) (*dto.CategoryResponse, error) {
	name = strings.TrimSpace(name)
	slug := utils.GenerateSlug(name)

	existing, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var cat *models.Category
	if existing != nil {
		if !existing.DeletedAt.Valid && existing.IsActive {
			return nil, ErrCategoryAlreadyExists
		}
		// Restore a previously deleted category with the same name, since
		// the unique index on slug still holds the deleted row.
		existing.DeletedAt = gorm.DeletedAt{}
		existing.IsActive = true
		existing.Name = name
		existing.Description = strings.TrimSpace(description)
		if err := s.db.WithContext(ctx).Unscoped().Omit("Products").Save(existing).Error; err != nil {
			return nil, err
		}
		cat = existing
	} else {
		cat = &models.Category{Name: name, Slug: slug, Description: strings.TrimSpace(description), IsActive: true}
		if err := s.categoryRepo.Create(ctx, cat); err != nil {
			return nil, fmt.Errorf("gagal menyimpan kategori: %w", err)
		}
	}

	return &dto.CategoryResponse{ID: cat.ID, Name: cat.Name, Slug: cat.Slug, Description: cat.Description}, nil
}

func (s *productService) UpdateCategory(ctx context.Context, id uuid.UUID, name, description string) (*dto.CategoryResponse, error) {
	cat, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	name = strings.TrimSpace(name)
	slug := utils.GenerateSlug(name)
	if slug != cat.Slug {
		if other, err := s.categoryRepo.FindBySlug(ctx, slug); err != nil {
			return nil, err
		} else if other != nil && other.ID != cat.ID {
			return nil, ErrCategoryAlreadyExists
		}
	}

	cat.Name = name
	cat.Slug = slug
	cat.Description = strings.TrimSpace(description)
	if err := s.categoryRepo.Update(ctx, cat); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{ID: cat.ID, Name: cat.Name, Slug: cat.Slug, Description: cat.Description}, nil
}

func (s *productService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	cat, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrCategoryNotFound
	}
	count, err := s.categoryRepo.CountProducts(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryInUse
	}
	return s.categoryRepo.Delete(ctx, id)
}
