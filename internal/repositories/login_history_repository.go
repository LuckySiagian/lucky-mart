package repositories

import (
	"context"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type loginHistoryRepository struct {
	db *gorm.DB
}

func NewLoginHistoryRepository(db *gorm.DB) LoginHistoryRepository {
	return &loginHistoryRepository{db: db}
}

func (r *loginHistoryRepository) Create(ctx context.Context, history *models.LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *loginHistoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]models.LoginHistory, error) {
	var histories []models.LoginHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
