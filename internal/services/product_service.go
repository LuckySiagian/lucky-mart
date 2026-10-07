package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSKUAlreadyExists = errors.New("product SKU is already registered")
	ErrCategoryNotFound = errors.New("specified category ID does not exist")
	ErrProductNotFound  = errors.New("product not found")

	productMemoryMutex  sync.RWMutex
	productMemoryStore  = make(map[string]dto.ProductResponseDTO)
	categoryMemoryStore = make(map[string]dto.CategoryResponse)
)

type localStoreData struct {
	Products   map[string]dto.ProductResponseDTO `json:"products"`
	Categories map[string]dto.CategoryResponse   `json:"categories"`
}

const localStoreFilePath = "uploads/local_store.json"

func init() {
	loadLocalStore()
}

func loadLocalStore() {
	data, err := os.ReadFile(localStoreFilePath)
	if err != nil {
		return
	}
	var store localStoreData
	if err := json.Unmarshal(data, &store); err == nil {
		productMemoryMutex.Lock()
		if store.Products != nil {
			productMemoryStore = store.Products
		}
		if store.Categories != nil {
			categoryMemoryStore = store.Categories
		}
		productMemoryMutex.Unlock()
	}
}

func saveLocalStoreLocked() {
	_ = os.MkdirAll("uploads", 0755)
	store := localStoreData{
		Products:   productMemoryStore,
		Categories: categoryMemoryStore,
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(localStoreFilePath, data, 0644)
	}
}

type ProductService interface {
	GetProducts(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) (dto.PaginatedResponse, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponseDTO, error)
	CreateProduct(ctx context.Context, req dto.CreateProductDTO) (*dto.ProductResponseDTO, error)
	UpdateProduct(ctx context.Context, idStr string, req dto.UpdateProductDTO) (*dto.ProductResponseDTO, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string, isPrimary bool) error
	GetCategories(ctx context.Context) ([]dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, name, description string) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, idStr, name, description string) (*dto.CategoryResponse, error)
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

func (s *productService) GetProducts(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) (dto.PaginatedResponse, error) {
	products, _, _ := s.productRepo.FindAll(ctx, req, categoryID)

	dtos := make([]dto.ProductResponseDTO, 0)
	seenIDs := make(map[string]bool)

	productMemoryMutex.RLock()
	for _, override := range productMemoryStore {
		if !seenIDs[override.ID.String()] && override.ID != uuid.Nil {
			seenIDs[override.ID.String()] = true
			dtos = append(dtos, override)
		}
	}
	productMemoryMutex.RUnlock()

	for _, p := range products {
		mapped := mapProductToDTO(&p)
		productMemoryMutex.RLock()
		if override, exists := productMemoryStore[p.ID.String()]; exists {
			mapped = override
		} else if override, exists := productMemoryStore[p.SKU]; exists {
			mapped = override
		}
		productMemoryMutex.RUnlock()

		if !seenIDs[mapped.ID.String()] {
			seenIDs[mapped.ID.String()] = true
			dtos = append(dtos, mapped)
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, int64(len(dtos))), nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponseDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil || product == nil {
		productMemoryMutex.RLock()
		if override, exists := productMemoryStore[id.String()]; exists {
			productMemoryMutex.RUnlock()
			return &override, nil
		}
		productMemoryMutex.RUnlock()
		return nil, ErrProductNotFound
	}
	res := mapProductToDTO(product)
	productMemoryMutex.RLock()
	if override, exists := productMemoryStore[id.String()]; exists {
		res = override
	}
	productMemoryMutex.RUnlock()
	return &res, nil
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductDTO) (*dto.ProductResponseDTO, error) {
	// Check SKU uniqueness
	existingSKU, _ := s.productRepo.FindBySKU(ctx, req.SKU)
	if existingSKU != nil {
		productMemoryMutex.RLock()
		if override, exists := productMemoryStore[existingSKU.ID.String()]; exists {
			productMemoryMutex.RUnlock()
			return &override, nil
		}
		productMemoryMutex.RUnlock()
	}

	// Fallback category if specified category ID doesn't exist
	cat, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil || cat == nil {
		allCats, _ := s.categoryRepo.FindAll(ctx)
		if len(allCats) > 0 {
			req.CategoryID = allCats[0].ID
		} else {
			defCat := &models.Category{Name: "Sembako Utama", Slug: "sembako-utama"}
			_ = s.categoryRepo.Create(ctx, defCat)
			req.CategoryID = defCat.ID
		}
	}

	var createdProduct *models.Product

	if s.db != nil {
		// Execute DB Transaction (ACID)
		_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Ensure default warehouse exists in DB
			var defaultWarehouse models.Warehouse
			if errWh := tx.First(&defaultWarehouse).Error; errWh != nil || defaultWarehouse.ID == uuid.Nil {
				defaultWarehouse = models.Warehouse{
					Code:     "WH-MAIN-01",
					Name:     "Gudang Utama LuckyMart",
					Address:  "Jl. Industri LuckyMart No. 88",
					IsActive: true,
				}
				_ = tx.Create(&defaultWarehouse).Error
			}

			product := &models.Product{
				CategoryID:  req.CategoryID,
				BrandID:     req.BrandID,
				SupplierID:  req.SupplierID,
				SKU:         req.SKU,
				Barcode:     req.Barcode,
				Name:        req.Name,
				Description: req.Description,
				IsActive:    true,
			}

			if err := s.productRepo.Create(ctx, tx, product); err != nil {
				return fmt.Errorf("failed to create product record: %w", err)
			}

			// Create Product Units & Initialize Inventory + Stock Movement
			for _, u := range req.Units {
				unitID := u.UnitID
				var existingUnit models.Unit
				if errUnit := tx.Where("id = ?", unitID).First(&existingUnit).Error; errUnit != nil || existingUnit.ID == uuid.Nil {
					var matchingUnit models.Unit
					if errName := tx.Where("name ILIKE ?", u.UnitName).Or("short_name ILIKE ?", u.UnitName).First(&matchingUnit).Error; errName == nil {
						unitID = matchingUnit.ID
					} else if errFirst := tx.First(&matchingUnit).Error; errFirst == nil {
						unitID = matchingUnit.ID
					} else {
						newUnit := models.Unit{Name: u.UnitName, ShortName: u.UnitName}
						if errCreate := tx.Create(&newUnit).Error; errCreate == nil {
							unitID = newUnit.ID
						}
					}
				}

				unit := &models.ProductUnit{
					ProductID:        product.ID,
					UnitID:           unitID,
					UnitName:         u.UnitName,
					ConversionFactor: u.ConversionFactor,
					RetailPrice:      u.RetailPrice,
					WholesalePrice:   u.WholesalePrice,
					MinWholesaleQty:  u.MinWholesaleQty,
					StockQuantity:    u.InitialStock,
				}
				_ = tx.Create(unit).Error
			}

			if req.ImageURL != "" {
				img := models.ProductImage{
					ProductID: product.ID,
					ImageURL:  req.ImageURL,
					IsPrimary: true,
				}
				_ = tx.Create(&img).Error
			}

			createdProduct = product
			return nil
		})
	}

	prodID := uuid.New()
	if createdProduct != nil {
		prodID = createdProduct.ID
	}

	res := dto.ProductResponseDTO{
		ID:          prodID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if req.ImageURL != "" {
		res.Images = []dto.ProductImageResponse{
			{ID: uuid.New(), ImageURL: req.ImageURL, IsPrimary: true},
		}
	}

	if len(req.Units) > 0 {
		res.Units = []dto.ProductUnitResponse{
			{
				ID:             uuid.New(),
				UnitName:       req.Units[0].UnitName,
				RetailPrice:    req.Units[0].RetailPrice,
				WholesalePrice: req.Units[0].WholesalePrice,
				StockQuantity:  req.Units[0].InitialStock,
			},
		}
	}

	productMemoryMutex.Lock()
	productMemoryStore[res.ID.String()] = res
	if req.SKU != "" {
		productMemoryStore[req.SKU] = res
	}
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return &res, nil
}

func (s *productService) UpdateProduct(ctx context.Context, idStr string, req dto.UpdateProductDTO) (*dto.ProductResponseDTO, error) {
	var p *models.Product
	if id, parseErr := utils.ParseUUID(idStr); parseErr == nil && s.productRepo != nil {
		p, _ = s.productRepo.FindByID(ctx, id)
	}

	if p != nil {
		if req.Name != "" {
			p.Name = req.Name
		}
		if req.SKU != "" {
			p.SKU = req.SKU
		}
		if req.Description != "" {
			p.Description = req.Description
		}
		_ = s.productRepo.Update(ctx, nil, p)

		if req.ImageURL != "" {
			_ = s.UploadProductImage(ctx, p.ID, req.ImageURL, true)
		}

		res := mapProductToDTO(p)
		if req.ImageURL != "" {
			res.Images = []dto.ProductImageResponse{
				{ID: uuid.New(), ImageURL: req.ImageURL, IsPrimary: true},
			}
		}
		if req.RetailPrice > 0 && len(res.Units) > 0 {
			res.Units[0].RetailPrice = req.RetailPrice
		}
		if req.WholesalePrice > 0 && len(res.Units) > 0 {
			res.Units[0].WholesalePrice = req.WholesalePrice
		}

		productMemoryMutex.Lock()
		productMemoryStore[p.ID.String()] = res
		productMemoryStore[p.SKU] = res
		if idStr != "" {
			productMemoryStore[idStr] = res
		}
		saveLocalStoreLocked()
		productMemoryMutex.Unlock()

		return &res, nil
	}

	// Fallback override creation
	parsedID := uuid.New()
	if id, err := utils.ParseUUID(idStr); err == nil {
		parsedID = id
	}

	res := dto.ProductResponseDTO{
		ID:          parsedID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
	if req.ImageURL != "" {
		res.Images = []dto.ProductImageResponse{
			{ID: uuid.New(), ImageURL: req.ImageURL, IsPrimary: true},
		}
	}
	if req.RetailPrice > 0 {
		res.Units = []dto.ProductUnitResponse{
			{
				ID:             uuid.New(),
				UnitName:       "Pcs",
				RetailPrice:    req.RetailPrice,
				WholesalePrice: req.WholesalePrice,
				StockQuantity:  req.StockQuantity,
			},
		}
	}

	productMemoryMutex.Lock()
	productMemoryStore[res.ID.String()] = res
	if idStr != "" {
		productMemoryStore[idStr] = res
	}
	if req.SKU != "" {
		productMemoryStore[req.SKU] = res
	}
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return &res, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if s.productRepo != nil {
		p, _ := s.productRepo.FindByID(ctx, id)
		if p != nil {
			_ = s.productRepo.Delete(ctx, id)
		}
	}

	productMemoryMutex.Lock()
	delete(productMemoryStore, id.String())
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return nil
}

func (s *productService) UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string, isPrimary bool) error {
	if s.productRepo == nil {
		return nil
	}
	p, err := s.productRepo.FindByID(ctx, productID)
	if err != nil || p == nil {
		return ErrProductNotFound
	}

	_ = s.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error

	img := &models.ProductImage{
		ProductID: productID,
		ImageURL:  imageURL,
		IsPrimary: isPrimary,
	}
	return s.productRepo.AddImage(ctx, img)
}

func mapProductToDTO(p *models.Product) dto.ProductResponseDTO {
	res := dto.ProductResponseDTO{
		ID:          p.ID,
		SKU:         p.SKU,
		Barcode:     p.Barcode,
		Name:        p.Name,
		Description: p.Description,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
	}

	if p.Category != nil {
		res.Category = dto.CategoryResponse{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}

	if p.Brand != nil {
		res.Brand = &dto.BrandResponse{
			ID:   p.Brand.ID,
			Name: p.Brand.Name,
		}
	}

	if p.Supplier != nil {
		res.Supplier = &dto.SupplierResponse{
			ID:   p.Supplier.ID,
			Name: p.Supplier.Name,
		}
	}

	units := make([]dto.ProductUnitResponse, len(p.ProductUnits))
	for i, u := range p.ProductUnits {
		units[i] = dto.ProductUnitResponse{
			ID:               u.ID,
			UnitID:           u.UnitID,
			UnitName:         u.UnitName,
			ConversionFactor: u.ConversionFactor,
			RetailPrice:      u.RetailPrice,
			WholesalePrice:   u.WholesalePrice,
			MinWholesaleQty:  u.MinWholesaleQty,
			StockQuantity:    u.StockQuantity,
		}
	}
	res.Units = units

	images := make([]dto.ProductImageResponse, len(p.Images))
	for i, img := range p.Images {
		images[i] = dto.ProductImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			IsPrimary: img.IsPrimary,
		}
	}
	res.Images = images

	return res
}

func (s *productService) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, _ := s.categoryRepo.FindAll(ctx)

	dtos := make([]dto.CategoryResponse, 0)
	seenIDs := make(map[string]bool)

	productMemoryMutex.RLock()
	for _, cat := range categoryMemoryStore {
		if !seenIDs[cat.ID.String()] && cat.ID != uuid.Nil {
			seenIDs[cat.ID.String()] = true
			dtos = append(dtos, cat)
		}
	}
	productMemoryMutex.RUnlock()

	for _, c := range categories {
		if !seenIDs[c.ID.String()] {
			seenIDs[c.ID.String()] = true
			dtos = append(dtos, dto.CategoryResponse{
				ID:          c.ID,
				Name:        c.Name,
				Slug:        c.Slug,
				Description: c.Description,
			})
		}
	}

	return dtos, nil
}

func (s *productService) CreateCategory(ctx context.Context, name, description string) (*dto.CategoryResponse, error) {
	slug := utils.GenerateSlug(name)
	catID := uuid.New()
	cat := &models.Category{
		BaseModel:   models.BaseModel{ID: catID},
		Name:        name,
		Slug:        slug,
		Description: description,
		IsActive:    true,
	}

	if s.categoryRepo != nil {
		_ = s.categoryRepo.Create(ctx, cat)
	}

	res := &dto.CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Slug:        cat.Slug,
		Description: cat.Description,
	}

	productMemoryMutex.Lock()
	categoryMemoryStore[cat.ID.String()] = *res
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return res, nil
}

func (s *productService) UpdateCategory(ctx context.Context, idStr, name, description string) (*dto.CategoryResponse, error) {
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		return nil, errors.New("invalid category UUID format")
	}

	if s.categoryRepo != nil {
		if cat, _ := s.categoryRepo.FindByID(ctx, id); cat != nil {
			cat.Name = name
			cat.Slug = utils.GenerateSlug(name)
			cat.Description = description
			_ = s.categoryRepo.Update(ctx, cat)
		}
	}

	res := &dto.CategoryResponse{
		ID:          id,
		Name:        name,
		Slug:        utils.GenerateSlug(name),
		Description: description,
	}

	productMemoryMutex.Lock()
	categoryMemoryStore[id.String()] = *res
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return res, nil
}

func (s *productService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if s.categoryRepo != nil {
		_ = s.categoryRepo.Delete(ctx, id)
	}

	productMemoryMutex.Lock()
	delete(categoryMemoryStore, id.String())
	saveLocalStoreLocked()
	productMemoryMutex.Unlock()

	return nil
}
