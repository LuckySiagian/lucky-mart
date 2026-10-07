package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type shipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(ctx context.Context, tx *gorm.DB, shipment *models.Shipment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(shipment).Error
}

func (r *shipmentRepository) Update(ctx context.Context, tx *gorm.DB, shipment *models.Shipment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(shipment).Error
}

func (r *shipmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Courier").
		Preload("ShippingMethod").
		Preload("Trackings").
		Where("id = ?", id).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Courier").
		Preload("ShippingMethod").
		Preload("Trackings").
		Where("order_id = ?", orderID).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) AddTracking(ctx context.Context, tx *gorm.DB, tracking *models.ShipmentTracking) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(tracking).Error
}
