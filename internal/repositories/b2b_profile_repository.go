package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type b2bProfileRepository struct {
	db *gorm.DB
}

func NewB2BProfileRepository(db *gorm.DB) B2BProfileRepository {
	return &b2bProfileRepository{db: db}
}

func (r *b2bProfileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.B2BProfile, error) {
	var profile models.B2BProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *b2bProfileRepository) Create(ctx context.Context, profile *models.B2BProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *b2bProfileRepository) Update(ctx context.Context, profile *models.B2BProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}
