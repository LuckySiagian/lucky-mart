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

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Product{}).
		Preload("Category").
		Preload("Brand").
		Preload("Supplier").
		Preload("ProductUnits").
		Preload("Images").
		Where("is_active = true")

	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Supplier").
		Preload("ProductUnits").
		Preload("Images").
		Where("id = ?", id).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, tx *gorm.DB, product *models.Product) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, tx *gorm.DB, product *models.Product) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) AddImage(ctx context.Context, image *models.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}
