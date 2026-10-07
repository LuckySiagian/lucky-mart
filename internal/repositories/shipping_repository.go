package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) CourierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) FindAll(ctx context.Context) ([]models.Courier, error) {
	var couriers []models.Courier
	err := r.db.WithContext(ctx).Where("is_active = true").Find(&couriers).Error
	return couriers, err
}

func (r *courierRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Courier, error) {
	var courier models.Courier
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&courier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &courier, nil
}

type shippingMethodRepository struct {
	db *gorm.DB
}

func NewShippingMethodRepository(db *gorm.DB) ShippingMethodRepository {
	return &shippingMethodRepository{db: db}
}

func (r *shippingMethodRepository) FindAll(ctx context.Context) ([]models.ShippingMethod, error) {
	var methods []models.ShippingMethod
	err := r.db.WithContext(ctx).Preload("Courier").Where("is_active = true").Find(&methods).Error
	return methods, err
}

func (r *shippingMethodRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ShippingMethod, error) {
	var method models.ShippingMethod
	err := r.db.WithContext(ctx).Preload("Courier").Where("id = ?", id).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &method, nil
}

type shippingRateRepository struct {
	db *gorm.DB
}

func NewShippingRateRepository(db *gorm.DB) ShippingRateRepository {
	return &shippingRateRepository{db: db}
}

func (r *shippingRateRepository) FindByMethodAndCity(ctx context.Context, methodID uuid.UUID, city string) (*models.ShippingRate, error) {
	var rate models.ShippingRate
	err := r.db.WithContext(ctx).
		Preload("ShippingMethod.Courier").
		Where("shipping_method_id = ? AND city ILIKE ?", methodID, city).
		First(&rate).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rate, nil
}
