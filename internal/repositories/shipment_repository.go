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
	return db.WithContext(ctx).Omit(clause.Associations).Save(shipment).Error
}

func (r *shipmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.WithContext(ctx).
		Preload("Order.User").
		Preload("Order.OrderAddress").
		Preload("Courier").
		Preload("ShippingMethod").
		Preload("Trackings", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
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
		Preload("Order.User").
		Preload("Order.OrderAddress").
		Preload("Courier").
		Preload("ShippingMethod").
		Preload("Trackings", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
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

func (r *shipmentRepository) FindAll(ctx context.Context, req dto.PaginationRequest, status string) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Shipment{}).
		Preload("Order.User").
		Preload("Order.OrderAddress").
		Preload("Courier").
		Preload("ShippingMethod").
		Preload("Trackings", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") })

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if req.Search != "" {
		db = db.Where("tracking_number ILIKE ?", utils.SearchTerm(req.Search))
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req, "status")
	if err := db.Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}
