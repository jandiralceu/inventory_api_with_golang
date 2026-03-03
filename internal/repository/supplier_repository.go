package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// SupplierRepository defines the persistence contract for supplier-related data operations.
type SupplierRepository interface {
	// Create persists a new supplier record.
	Create(ctx context.Context, supplier *models.Supplier) error
	// Update modifies an existing supplier record.
	Update(ctx context.Context, supplier *models.Supplier) error
	// Delete removes a supplier record by its unique ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single supplier by its unique identifier.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error)
	// FindBySlug retrieves a single supplier by its URL-friendly slug.
	FindBySlug(ctx context.Context, slug string) (*models.Supplier, error)
	// FindAll retrieves a paginated list of suppliers based on the provided filters.
	FindAll(ctx context.Context, filter SupplierListFilter) (suppliers []models.Supplier, total int64, err error)
}

type supplierRepository struct {
	db *gorm.DB
}

var _ SupplierRepository = (*supplierRepository)(nil)

// NewSupplierRepository initializes a GORM-based implementation of SupplierRepository.
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

// SupplierListFilter encapsulates search and pagination parameters for supplier listing operations.
type SupplierListFilter struct {
	Name       string
	TaxID      string
	Email      string
	IsActive   *bool
	Pagination PaginationParams
}

// Create inserts a new supplier record into the database.
func (r *supplierRepository) Create(ctx context.Context, supplier *models.Supplier) error {
	if err := r.db.WithContext(ctx).Create(supplier).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Update saves changes to an existing supplier record.
func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	result := r.db.WithContext(ctx).Save(supplier)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}
	return nil
}

// Delete removes a supplier record by its unique identifier.
func (r *supplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Supplier{}, "id = ?", id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}
	return nil
}

// FindByID retrieves a supplier by its unique identifier.
func (r *supplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := r.db.WithContext(ctx).First(&supplier, "id = ?", id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &supplier, nil
}

// FindBySlug retrieves a supplier by its URL-friendly slug.
func (r *supplierRepository) FindBySlug(ctx context.Context, slug string) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&supplier).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &supplier, nil
}

// FindAll executes a listing query with dynamic filtering and pagination for suppliers.
func (r *supplierRepository) FindAll(ctx context.Context, filter SupplierListFilter) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Supplier{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+sanitizeLike(filter.Name)+"%")
	}
	if filter.TaxID != "" {
		query = query.Where("tax_id = ?", filter.TaxID)
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+sanitizeLike(filter.Email)+"%")
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
		Find(&suppliers).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return suppliers, total, nil
}
