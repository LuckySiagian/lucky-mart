package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		Where("current_stock <= min_stock_alert").
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
	return db.WithContext(ctx).Save(inv).Error
}
