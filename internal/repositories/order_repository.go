package repositories

import (
	"context"
	"errors"
	"strings"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	return db.WithContext(ctx).Omit(clause.Associations).Save(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("OrderAddress").
		Preload("OrderItems.Product", withDeleted).
		Preload("OrderItems.ProductUnit", withDeleted).
		Preload("Invoice").
		Preload("Histories", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
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

func (r *orderRepository) list(ctx context.Context, req dto.PaginationRequest, status string, scope func(*gorm.DB) *gorm.DB) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Order{}).
		Preload("User").
		Preload("OrderAddress").
		Preload("OrderItems.Product", withDeleted).
		Preload("OrderItems.ProductUnit", withDeleted).
		Preload("Invoice")

	if scope != nil {
		db = scope(db)
	}
	// status may list several values, e.g. "CONFIRMED,PROCESSING_PACKING".
	if status != "" {
		db = db.Where("status IN ?", strings.Split(status, ","))
	}
	if req.Search != "" {
		term := utils.SearchTerm(req.Search)
		db = db.Where("order_number ILIKE ? OR user_id IN (SELECT id FROM users WHERE full_name ILIKE ? OR phone_number ILIKE ?)", term, term, term)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req, "order_date", "grand_total", "status")
	if err := db.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest, status string) ([]models.Order, int64, error) {
	return r.list(ctx, req, status, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

func (r *orderRepository) FindAll(ctx context.Context, req dto.PaginationRequest, status string) ([]models.Order, int64, error) {
	return r.list(ctx, req, status, nil)
}

func (r *orderRepository) AddHistory(ctx context.Context, tx *gorm.DB, history *models.OrderHistory) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(history).Error
}

// withDeleted lets order history show products that were later removed
// from the catalog.
func withDeleted(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
