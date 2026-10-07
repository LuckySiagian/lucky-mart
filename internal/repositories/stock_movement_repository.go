package repositories

import (
	"context"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(ctx context.Context, tx *gorm.DB, movement *models.StockMovement) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(movement).Error
}

func (r *stockMovementRepository) FindByInventoryID(ctx context.Context, inventoryID uuid.UUID, req dto.PaginationRequest) ([]models.StockMovement, int64, error) {
	var movements []models.StockMovement
	var total int64

	db := r.db.WithContext(ctx).
		Model(&models.StockMovement{}).
		Where("inventory_id = ?", inventoryID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&movements).Error; err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}
