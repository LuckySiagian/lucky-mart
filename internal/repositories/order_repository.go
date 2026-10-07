package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) Update(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("OrderAddress").
		Preload("OrderItems.Product").
		Preload("OrderItems.ProductUnit").
		Preload("Invoice").
		Preload("Histories").
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Order{}).
		Preload("OrderAddress").
		Preload("OrderItems").
		Preload("Invoice").
		Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Order{}).
		Preload("User").
		Preload("OrderAddress").
		Preload("OrderItems").
		Preload("Invoice")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) AddHistory(ctx context.Context, tx *gorm.DB, history *models.OrderHistory) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(history).Error
}
