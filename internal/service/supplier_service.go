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

// SupplierService defines the business logic contract for supplier management and caching.
type SupplierService interface {
	Create(ctx context.Context, req dto.CreateSupplierRequest) (*models.Supplier, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateSupplierRequest) (*models.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error)
	FindBySlug(ctx context.Context, slug string) (*models.Supplier, error)
	FindAll(ctx context.Context, req dto.GetSupplierListRequest) (dto.PaginatedResponse[models.Supplier], error)
}

type supplierService struct {
	supplierRepo repository.SupplierRepository
	cache        pkg.CacheManager
}

var _ SupplierService = (*supplierService)(nil)

const (
	_supplierCachePrefix = "supplier:"
	_supplierCacheTTL    = 72 * time.Hour
)

// NewSupplierService initializes a SupplierService with repository and cache dependencies.
func NewSupplierService(supplierRepo repository.SupplierRepository, cache pkg.CacheManager) SupplierService {
	return &supplierService{
		supplierRepo: supplierRepo,
		cache:        cache,
	}
}

func (s *supplierService) Create(ctx context.Context, req dto.CreateSupplierRequest) (*models.Supplier, error) {
	supplier := &models.Supplier{
		Name:          req.Name,
		Slug:          pkg.Slugify(req.Name),
		Description:   req.Description,
		TaxID:         req.TaxID,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address.MapToModel(),
		ContactPerson: req.ContactPerson,
		IsActive:      true,
	}

	if err := s.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	s.cache.DeletePrefix(ctx, _supplierCachePrefix)
	return supplier, nil
}

func (s *supplierService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateSupplierRequest) (*models.Supplier, error) {
	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	supplier.Name = req.Name
	supplier.Slug = pkg.Slugify(req.Name)
	supplier.Description = req.Description
	supplier.TaxID = req.TaxID
	supplier.Email = req.Email
	supplier.Phone = req.Phone
	supplier.Address = req.Address.MapToModel()
	supplier.ContactPerson = req.ContactPerson
	supplier.IsActive = *req.IsActive

	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	s.cache.DeletePrefix(ctx, _supplierCachePrefix)
	return supplier, nil
}

func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.supplierRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.cache.DeletePrefix(ctx, _supplierCachePrefix)
	return nil
}

func (s *supplierService) FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error) {
	cacheKey := fmt.Sprintf("%sid:%s", _supplierCachePrefix, id)
	var cached models.Supplier
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err == nil {
		s.cache.Set(ctx, cacheKey, supplier, _supplierCacheTTL)
	}

	return supplier, err
}

func (s *supplierService) FindBySlug(ctx context.Context, slug string) (*models.Supplier, error) {
	cacheKey := fmt.Sprintf("%sslug:%s", _supplierCachePrefix, slug)
	var cached models.Supplier
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	supplier, err := s.supplierRepo.FindBySlug(ctx, slug)
	if err == nil {
		s.cache.Set(ctx, cacheKey, supplier, _supplierCacheTTL)
	}

	return supplier, err
}

func (s *supplierService) FindAll(ctx context.Context, req dto.GetSupplierListRequest) (dto.PaginatedResponse[models.Supplier], error) {
	filter := repository.SupplierListFilter{
		Name:     req.Name,
		TaxID:    req.TaxID,
		Email:    req.Email,
		IsActive: req.IsActive,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "name", "tax_id"),
			Order: req.GetOrder(),
		},
	}

	suppliers, total, err := s.supplierRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[models.Supplier]{}, err
	}

	return dto.NewPaginatedResponse(suppliers, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}
