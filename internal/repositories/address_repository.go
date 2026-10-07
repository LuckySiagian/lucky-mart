package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, addr *models.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *addressRepository) Update(ctx context.Context, addr *models.Address) error {
	return r.db.WithContext(ctx).Save(addr).Error
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, id).Error
}

func (r *addressRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	var addr models.Address
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&addr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addr, nil
}

func (r *addressRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) UnsetDefault(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}
