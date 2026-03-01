package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

// CategoryService defines the business logic contract for category management and caching.
type CategoryService interface {
	// Create registers a new category, generating its slug automatically.
	Create(ctx context.Context, req dto.CreateCategoryRequest) (*models.Category, error)
	// Update modifies an existing category, handling slug regeneration and validation.
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*models.Category, error)
	// Delete removes a category and invalidates its cache.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a category by its ID, checking the cache first.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	// FindBySlug retrieves a category by its slug, checking the cache first.
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	// FindAll retrieves all categories filtered by request criteria.
	FindAll(ctx context.Context, req dto.GetCategoryListRequest) (dto.PaginatedResponse[models.Category], error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	cache        pkg.CacheManager
}

var _ CategoryService = (*categoryService)(nil)

const (
	_categoryCachePrefix = "category:"
	_categoryCacheTTL    = 72 * time.Hour
)

// NewCategoryService initializes a CategoryService with repository and cache dependencies.
func NewCategoryService(categoryRepo repository.CategoryRepository, cache pkg.CacheManager) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		cache:        cache,
	}
}

// Create generates a slug from the name and persists the new category.
func (s *categoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*models.Category, error) {
	if req.ParentID != nil {
		if _, err := s.categoryRepo.FindByID(ctx, *req.ParentID); err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        pkg.Slugify(req.Name),
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	s.cache.DeletePrefix(ctx, _categoryCachePrefix)
	return category, nil
}

// Update modifies the category, regenerates the slug, and handles parent validation.
func (s *categoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, fmt.Errorf("%w: category cannot be its own parent", apperrors.ErrInvalidInput)
		}
		if _, err := s.categoryRepo.FindByID(ctx, *req.ParentID); err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
	}

	category.Name = req.Name
	category.Slug = pkg.Slugify(req.Name)
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.IsActive = *req.IsActive

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	s.cache.DeletePrefix(ctx, _categoryCachePrefix)
	return category, nil
}

// Delete removes the category and purges the category cache.
func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.cache.DeletePrefix(ctx, _categoryCachePrefix)
	return nil
}

// FindByID retrieves the category, prioritizing the cache.
func (s *categoryService) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	cacheKey := fmt.Sprintf("%sid:%s", _categoryCachePrefix, id)
	var cached models.Category
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	category, err := s.categoryRepo.FindByID(ctx, id)
	if err == nil {
		s.cache.Set(ctx, cacheKey, category, _categoryCacheTTL)
	}

	return category, err
}

// FindBySlug retrieves the category by slug, prioritizing the cache.
func (s *categoryService) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	cacheKey := fmt.Sprintf("%sslug:%s", _categoryCachePrefix, slug)
	var cached models.Category
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err == nil {
		s.cache.Set(ctx, cacheKey, category, _categoryCacheTTL)
	}

	return category, err
}

// FindAll coordinates data retrieval for paginated results with filtering.
func (s *categoryService) FindAll(ctx context.Context, req dto.GetCategoryListRequest) (dto.PaginatedResponse[models.Category], error) {
	filter := repository.CategoryListFilter{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		IsActive: req.IsActive,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "name", "slug"),
			Order: req.GetOrder(),
		},
	}

	categories, total, err := s.categoryRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[models.Category]{}, err
	}

	return dto.NewPaginatedResponse(categories, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}
