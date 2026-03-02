package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// WarehouseRepository defines the persistence contract for warehouse-related data operations.
type WarehouseRepository interface {
	// Create persists a new warehouse record.
	// Returns ErrConflict if a warehouse with the same code or slug already exists.
	Create(ctx context.Context, warehouse *models.Warehouse) error

	// Update modifies an existing warehouse record.
	// Returns ErrNotFound if the record does not exist.
	Update(ctx context.Context, warehouse *models.Warehouse) error

	// Delete removes a warehouse record by its unique ID.
	// Returns ErrNotFound if the record does not exist.
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByID retrieves a warehouse by its unique identifier.
	// Returns ErrNotFound if the record does not exist.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error)

	// FindBySlug retrieves a warehouse by its URL-friendly slug.
	// Returns ErrNotFound if the record does not exist.
	FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error)

	// FindByCode retrieves a warehouse by its unique business code.
	// Returns ErrNotFound if the record does not exist.
	FindByCode(ctx context.Context, code string) (*models.Warehouse, error)

	// FindAll retrieves a paginated list of warehouses based on the provided filters.
	FindAll(ctx context.Context, filter WarehouseListFilter) (warehouses []models.Warehouse, total int64, err error)
}

type warehouseRepository struct {
	db *gorm.DB
}

var _ WarehouseRepository = (*warehouseRepository)(nil)

// NewWarehouseRepository initializes a GORM-based implementation of WarehouseRepository.
func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

// WarehouseListFilter encapsulates search and pagination parameters for warehouse listing operations.
type WarehouseListFilter struct {
	Name       string
	Code       string
	IsActive   *bool
	Pagination PaginationParams
}

// Create persists a new warehouse record.
func (r *warehouseRepository) Create(ctx context.Context, warehouse *models.Warehouse) error {
	if err := r.db.WithContext(ctx).Create(warehouse).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Update modifies an existing warehouse record.
func (r *warehouseRepository) Update(ctx context.Context, warehouse *models.Warehouse) error {
	result := r.db.WithContext(ctx).
		Model(&models.Warehouse{}).
		Where("id = ?", warehouse.ID).
		Updates(warehouse)

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}
	return nil
}

// Delete removes a warehouse record by its unique ID.
func (r *warehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Warehouse{}, "id = ?", id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}
	return nil
}

// FindByID retrieves a warehouse by its unique identifier.
func (r *warehouseRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := r.db.WithContext(ctx).First(&warehouse, "id = ?", id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &warehouse, nil
}

// FindBySlug retrieves a warehouse by its URL-friendly slug.
func (r *warehouseRepository) FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&warehouse).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &warehouse, nil
}

// FindByCode retrieves a warehouse by its unique business code.
func (r *warehouseRepository) FindByCode(ctx context.Context, code string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&warehouse).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &warehouse, nil
}

// FindAll retrieves a paginated list of warehouses based on the provided filters.
func (r *warehouseRepository) FindAll(ctx context.Context, filter WarehouseListFilter) ([]models.Warehouse, int64, error) {
	var warehouses []models.Warehouse
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Warehouse{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+sanitizeLike(filter.Name)+"%")
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	err := query.
		Order(filter.Pagination.GetOrderBy()).
		Offset(filter.Pagination.GetOffset()).
		Limit(filter.Pagination.Limit).
		Find(&warehouses).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return warehouses, total, nil
}
