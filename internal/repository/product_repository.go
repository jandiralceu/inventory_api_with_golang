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

// ProductRepository defines the persistence contract for product-related data operations.
type ProductRepository interface {
	// Create persists a new product record.
	Create(ctx context.Context, product *models.Product) error
	// Update modifies an existing product record.
	Update(ctx context.Context, product *models.Product) error
	// Delete removes a product record by its unique ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single product by its ID, preloading associations.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	// FindBySlug retrieves a single product by its URL-friendly slug.
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	// FindBySKU retrieves a single product by its unique SKU code.
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	// FindAll retrieves a list of products based on filters and pagination parameters.
	FindAll(ctx context.Context, filter ProductListFilter) ([]models.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository initializes a GORM-based implementation of ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product record into the database.
func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Update saves changes to an existing product record.
func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Delete removes a product record and returns [apperrors.ErrNotFound] if no record was deleted.
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

// FindByID retrieves a specific product using its unique UUID.
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").First(&product, id).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

// FindBySlug retrieves a product by its URL-friendly slug.
func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

// FindBySKU retrieves a product by its SKU code.
func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Supplier").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &product, nil
}

// FindAll executes a listing query with dynamic filtering, preloading, and pagination.
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
