package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository defines the persistence contract for category-related data operations.
type CategoryRepository interface {
	// Create persists a new category record.
	Create(ctx context.Context, category *models.Category) error
	// Update modifies an existing category record.
	Update(ctx context.Context, category *models.Category) error
	// Delete removes a category record by its unique ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single category by its unique ID.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	// FindBySlug retrieves a single category by its unique URL-friendly slug.
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	// FindAll retrieves a list of categories based on filters and pagination parameters.
	FindAll(ctx context.Context, filter CategoryListFilter) (categories []models.Category, total int64, err error)
}

type categoryRepository struct {
	db *gorm.DB
}

var _ CategoryRepository = (*categoryRepository)(nil)

// NewCategoryRepository initializes a GORM-based implementation of CategoryRepository.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// CategoryListFilter encapsulates search and pagination parameters for category listing operations.
type CategoryListFilter struct {
	Name       string
	Slug       string
	ParentID   *uuid.UUID
	IsActive   *bool
	Pagination PaginationParams
}

// Create inserts a new category record into the database.
func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Update saves changes to an existing category record.
func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Save(category)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// Delete removes a category record and returns [apperrors.ErrNotFound] if no record was deleted.
func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// FindByID retrieves a category by its ID.
func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &category, nil
}

// FindBySlug retrieves a category by its slug.
func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &category, nil
}

// FindAll executes a listing query with dynamic filtering and pagination.
func (r *categoryRepository) FindAll(ctx context.Context, filter CategoryListFilter) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Category{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+sanitizeLike(filter.Name)+"%")
	}
	if filter.Slug != "" {
		query = query.Where("slug = ?", filter.Slug)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
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
		Find(&categories).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return categories, total, nil
}
