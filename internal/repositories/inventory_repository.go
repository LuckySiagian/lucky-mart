package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) FindByProductUnitID(ctx context.Context, unitID uuid.UUID) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("Product").
		Preload("ProductUnit").
		Where("product_unit_id = ?", unitID).
		First(&inv).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("Product").
		Preload("ProductUnit").
		Where("id = ?", id).
		First(&inv).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) FindLowStock(ctx context.Context) ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Preload("Product").
		Preload("ProductUnit").
		Joins("JOIN products ON products.id = inventories.product_id AND products.deleted_at IS NULL").
		Where("inventories.current_stock <= inventories.min_stock_alert").
		Find(&inventories).Error

	return inventories, err
}

func (r *inventoryRepository) Create(ctx context.Context, tx *gorm.DB, inv *models.Inventory) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(inv).Error
}

func (r *inventoryRepository) UpdateStock(ctx context.Context, tx *gorm.DB, inv *models.Inventory) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Omit(clause.Associations).Save(inv).Error
}

func (r *inventoryRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Inventory, int64, error) {
	var inventories []models.Inventory
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Inventory{}).
		Preload("Warehouse").
		Preload("Product.Category").
		Preload("ProductUnit").
		Joins("JOIN products ON products.id = inventories.product_id AND products.deleted_at IS NULL")

	if req.Search != "" {
		term := utils.SearchTerm(req.Search)
		db = db.Where("products.name ILIKE ? OR products.sku ILIKE ?", term, term)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	req.SetDefaults()
	sort := "products.name ASC"
	if req.SortBy == "current_stock" {
		sort = "inventories.current_stock " + map[bool]string{true: "ASC", false: "DESC"}[req.SortDir == "ASC" || req.SortDir == "asc"]
	}
	if err := db.Order(sort).Offset(req.GetOffset()).Limit(req.Limit).Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}
