package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

// WarehouseService defines the business logic contract for warehouse management and caching.
type WarehouseService interface {
	// Create registers a new warehouse and invalidates the warehouse cache.
	Create(ctx context.Context, req dto.CreateWarehouseRequest) (*models.Warehouse, error)
	// Update modifies an existing warehouse and invalidates the warehouse cache.
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateWarehouseRequest) (*models.Warehouse, error)
	// Delete removes a warehouse by ID and invalidates the warehouse cache.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single warehouse by ID, using cache when available.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error)
	// FindBySlug retrieves a single warehouse by its URL slug, using cache when available.
	FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error)
	// FindAll returns a paginated list of warehouses with optional filtering.
	FindAll(ctx context.Context, req dto.GetWarehouseListRequest) (dto.PaginatedResponse[models.Warehouse], error)
}

type warehouseService struct {
	warehouseRepo repository.WarehouseRepository
	cache         pkg.CacheManager
}

var _ WarehouseService = (*warehouseService)(nil)

const (
	_warehouseCachePrefix = "warehouse:"
	_warehouseCacheTTL    = 72 * time.Hour
)

// NewWarehouseService initializes a WarehouseService with repository and cache dependencies.
func NewWarehouseService(warehouseRepo repository.WarehouseRepository, cache pkg.CacheManager) WarehouseService {
	return &warehouseService{
		warehouseRepo: warehouseRepo,
		cache:         cache,
	}
}

// Create registers a new warehouse and invalidates the warehouse cache.
func (s *warehouseService) Create(ctx context.Context, req dto.CreateWarehouseRequest) (*models.Warehouse, error) {
	warehouse := &models.Warehouse{
		Name:        req.Name,
		Slug:        pkg.Slugify(req.Name),
		Code:        req.Code,
		Description: req.Description,
		Address:     req.Address.MapToModel(),
		ManagerName: req.ManagerName,
		Phone:       req.Phone,
		Email:       req.Email,
		Notes:       req.Notes,
		IsActive:    true,
	}

	if err := s.warehouseRepo.Create(ctx, warehouse); err != nil {
		return nil, err
	}

	_ = s.cache.DeletePrefix(ctx, _warehouseCachePrefix)
	return warehouse, nil
}

// Update modifies an existing warehouse and invalidates the warehouse cache.
func (s *warehouseService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateWarehouseRequest) (*models.Warehouse, error) {
	warehouse, err := s.warehouseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	warehouse.Name = req.Name
	warehouse.Slug = pkg.Slugify(req.Name)
	warehouse.Code = req.Code
	warehouse.Description = req.Description
	warehouse.Address = req.Address.MapToModel()
	warehouse.ManagerName = req.ManagerName
	warehouse.Phone = req.Phone
	warehouse.Email = req.Email
	warehouse.Notes = req.Notes
	warehouse.IsActive = *req.IsActive

	if err := s.warehouseRepo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	_ = s.cache.DeletePrefix(ctx, _warehouseCachePrefix)
	return warehouse, nil
}

// Delete removes a warehouse by ID and invalidates the warehouse cache.
func (s *warehouseService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.warehouseRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.cache.DeletePrefix(ctx, _warehouseCachePrefix)
	return nil
}

// FindByID retrieves a single warehouse by ID, using cache when available.
func (s *warehouseService) FindByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
	cacheKey := fmt.Sprintf("%sid:%s", _warehouseCachePrefix, id)
	var cached models.Warehouse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	warehouse, err := s.warehouseRepo.FindByID(ctx, id)
	if err == nil {
		_ = s.cache.Set(ctx, cacheKey, warehouse, _warehouseCacheTTL)
	}

	return warehouse, err
}

// FindBySlug retrieves a single warehouse by its URL slug, using cache when available.
func (s *warehouseService) FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error) {
	cacheKey := fmt.Sprintf("%sslug:%s", _warehouseCachePrefix, slug)
	var cached models.Warehouse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	warehouse, err := s.warehouseRepo.FindBySlug(ctx, slug)
	if err == nil {
		_ = s.cache.Set(ctx, cacheKey, warehouse, _warehouseCacheTTL)
	}

	return warehouse, err
}

// FindAll returns a paginated list of warehouses with optional filtering.
func (s *warehouseService) FindAll(ctx context.Context, req dto.GetWarehouseListRequest) (dto.PaginatedResponse[models.Warehouse], error) {
	filter := repository.WarehouseListFilter{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: req.IsActive,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "name", "code"),
			Order: req.GetOrder(),
		},
	}

	warehouses, total, err := s.warehouseRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[models.Warehouse]{}, err
	}

	return dto.NewPaginatedResponse(warehouses, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}
