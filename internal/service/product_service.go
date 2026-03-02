package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*models.Product, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*models.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	FindAll(ctx context.Context, req dto.GetProductListRequest) (*dto.ProductListResponse, error)
}

type productService struct {
	repo  repository.ProductRepository
	cache pkg.CacheManager
}

func NewProductService(repo repository.ProductRepository, cache pkg.CacheManager) ProductService {
	return &productService{
		repo:  repo,
		cache: cache,
	}
}

func (s *productService) Create(ctx context.Context, req dto.CreateProductRequest) (*models.Product, error) {
	// Check SKU uniqueness
	if _, err := s.repo.FindBySKU(ctx, req.SKU); err == nil {
		return nil, apperrors.ErrConflict
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	slug := pkg.Slugify(req.Name)
	if existing, _ := s.repo.FindBySlug(ctx, slug); existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	product := &models.Product{
		SKU:             req.SKU,
		Name:            req.Name,
		Slug:            slug,
		Description:     req.Description,
		Price:           req.Price,
		CostPrice:       req.CostPrice,
		CategoryID:      req.CategoryID,
		SupplierID:      req.SupplierID,
		ReorderLevel:    req.ReorderLevel,
		ReorderQuantity: req.ReorderQuantity,
		WeightKg:        req.WeightKg,
		Images:          req.Images,
		Metadata:        req.Metadata,
		IsActive:        true,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	_ = s.cache.DeletePrefix(ctx, "product:")
	return product, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.SKU != "" && req.SKU != product.SKU {
		if _, err := s.repo.FindBySKU(ctx, req.SKU); err == nil {
			return nil, apperrors.ErrConflict // SKU already exists
		} else if !errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		product.SKU = req.SKU
	}

	if req.Name != "" {
		product.Name = req.Name
		newSlug := pkg.Slugify(req.Name)
		if newSlug != product.Slug {
			if existing, _ := s.repo.FindBySlug(ctx, newSlug); existing != nil && existing.ID != id {
				newSlug = newSlug + "-" + uuid.New().String()[:8]
			}
			product.Slug = newSlug
		}
	}

	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CostPrice != nil {
		product.CostPrice = req.CostPrice
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.SupplierID != nil {
		product.SupplierID = req.SupplierID
	}
	if req.ReorderLevel != nil {
		product.ReorderLevel = *req.ReorderLevel
	}
	if req.ReorderQuantity != nil {
		product.ReorderQuantity = *req.ReorderQuantity
	}
	if req.WeightKg != nil {
		product.WeightKg = req.WeightKg
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.Metadata != nil {
		product.Metadata = req.Metadata
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	_ = s.cache.DeletePrefix(ctx, "product:")
	return product, nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.cache.DeletePrefix(ctx, "product:")
	return nil
}

func (s *productService) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:id:%s", id)
	var product models.Product

	if err := s.cache.Get(ctx, cacheKey, &product); err == nil {
		return &product, nil
	}

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, p, 72*time.Hour)
	return p, nil
}

func (s *productService) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:slug:%s", slug)
	var product models.Product

	if err := s.cache.Get(ctx, cacheKey, &product); err == nil {
		return &product, nil
	}

	p, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, p, 72*time.Hour)
	return p, nil
}

func (s *productService) FindAll(ctx context.Context, req dto.GetProductListRequest) (*dto.ProductListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	sort := req.Sort
	if sort == "" {
		sort = "created_at"
	}
	order := req.Order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	filter := repository.ProductListFilter{
		Name:       req.Name,
		SKU:        req.SKU,
		CategoryID: req.CategoryID,
		SupplierID: req.SupplierID,
		IsActive:   req.IsActive,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Page:       page,
		Limit:      limit,
		Sort:       sort,
		Order:      order,
	}

	products, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResponse{
		PaginatedResponse: dto.PaginatedResponse[models.Product]{
			Data:  products,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}, nil
}
