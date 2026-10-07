package repositories

import (
	"context"
	"errors"

	"backend-sembako/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, tx *gorm.DB, invoice *models.Invoice) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(invoice).Error
}

func (r *invoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *invoiceRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, invoiceID uuid.UUID, status string) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).
		Model(&models.Invoice{}).
		Where("id = ?", invoiceID).
		Update("status", status).Error
}
