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

type creditLedgerRepository struct {
	db *gorm.DB
}

func NewCreditLedgerRepository(db *gorm.DB) CreditLedgerRepository {
	return &creditLedgerRepository{db: db}
}

func (r *creditLedgerRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.CreditLedger, error) {
	var ledger models.CreditLedger
	err := r.db.WithContext(ctx).
		Preload("User.B2BProfile").
		Where("user_id = ?", userID).
		First(&ledger).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ledger, nil
}

func (r *creditLedgerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.CreditLedger, error) {
	var ledger models.CreditLedger
	err := r.db.WithContext(ctx).
		Preload("User.B2BProfile").
		Where("id = ?", id).
		First(&ledger).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ledger, nil
}

func (r *creditLedgerRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.CreditLedger, int64, error) {
	var ledgers []models.CreditLedger
	var total int64

	db := r.db.WithContext(ctx).Model(&models.CreditLedger{}).Preload("User.B2BProfile")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&ledgers).Error; err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *creditLedgerRepository) Create(ctx context.Context, tx *gorm.DB, ledger *models.CreditLedger) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(ledger).Error
}

func (r *creditLedgerRepository) Update(ctx context.Context, tx *gorm.DB, ledger *models.CreditLedger) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Omit(clause.Associations).Save(ledger).Error
}

type creditPaymentRepository struct {
	db *gorm.DB
}

func NewCreditPaymentRepository(db *gorm.DB) CreditPaymentRepository {
	return &creditPaymentRepository{db: db}
}

func (r *creditPaymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *models.CreditPayment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(payment).Error
}

func (r *creditPaymentRepository) FindByLedgerID(ctx context.Context, ledgerID uuid.UUID, req dto.PaginationRequest) ([]models.CreditPayment, int64, error) {
	var payments []models.CreditPayment
	var total int64

	db := r.db.WithContext(ctx).Model(&models.CreditPayment{}).Where("ledger_id = ?", ledgerID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
