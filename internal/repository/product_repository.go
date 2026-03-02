package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

type ProductListFilter struct {
	Name       string
	SKU        string
	CategoryID *uuid.UUID
	SupplierID *uuid.UUID
	IsActive   *bool
	MinPrice   *float64
	MaxPrice   *float64
	Page       int
	Limit      int
	Sort       string
	Order      string
}

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	FindAll(ctx context.Context, filter ProductListFilter) ([]models.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Product{}, id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").First(&product, id).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context, filter ProductListFilter) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.SKU != "" {
		query = query.Where("sku ILIKE ?", "%"+filter.SKU+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	if filter.Sort != "" {
		order := "asc"
		if filter.Order == "desc" {
			order = "desc"
		}
		query = query.Order(filter.Sort + " " + order)
	} else {
		query = query.Order("created_at desc")
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Preload("Category").Preload("Supplier").Offset(offset).Limit(filter.Limit).Find(&products).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return products, total, nil
}
