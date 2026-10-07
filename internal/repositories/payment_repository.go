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

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) Update(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Invoice").
		Preload("User").
		Preload("PaymentMethod").
		Preload("Transactions").
		Preload("Logs").
		Where("id = ?", id).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Invoice").
		Preload("PaymentMethod").
		Where("invoice_id = ?", invoiceID).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Payment{}).
		Preload("Invoice").
		Preload("PaymentMethod").
		Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *paymentRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Payment{}).
		Preload("User").
		Preload("Invoice").
		Preload("PaymentMethod")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *paymentRepository) FindWaitingConfirmation(ctx context.Context, req dto.PaginationRequest) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Payment{}).
		Preload("User").
		Preload("Invoice").
		Preload("PaymentMethod").
		Where("status = 'WAITING_CONFIRMATION'")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) FindAll(ctx context.Context) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	err := r.db.WithContext(ctx).Where("is_active = true").Find(&methods).Error
	return methods, err
}

func (r *paymentMethodRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &method, nil
}

func (r *paymentMethodRepository) FindByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &method, nil
}

func (r *paymentMethodRepository) Create(ctx context.Context, method *models.PaymentMethod) error {
	return r.db.WithContext(ctx).Create(method).Error
}

func (r *paymentMethodRepository) Update(ctx context.Context, method *models.PaymentMethod) error {
	return r.db.WithContext(ctx).Save(method).Error
}

func (r *paymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.PaymentMethod{}, id).Error
}

type paymentTransactionRepository struct {
	db *gorm.DB
}

func NewPaymentTransactionRepository(db *gorm.DB) PaymentTransactionRepository {
	return &paymentTransactionRepository{db: db}
}

func (r *paymentTransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *models.PaymentTransaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(transaction).Error
}

type paymentLogRepository struct {
	db *gorm.DB
}

func NewPaymentLogRepository(db *gorm.DB) PaymentLogRepository {
	return &paymentLogRepository{db: db}
}

func (r *paymentLogRepository) Create(ctx context.Context, tx *gorm.DB, log *models.PaymentLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(log).Error
}
