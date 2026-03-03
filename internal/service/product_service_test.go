package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_Create(t *testing.T) {
	repo := new(MockProductRepository)
	cache := new(MockCacheManager)
	svc := NewProductService(repo, cache)

	req := dto.CreateProductRequest{
		SKU:   "SKU-123",
		Name:  "Test Product",
		Price: 99.99,
	}

	repo.On("FindBySKU", mock.Anything, req.SKU).Return((*models.Product)(nil), apperrors.ErrNotFound)
	repo.On("FindBySlug", mock.Anything, "test-product").Return((*models.Product)(nil), apperrors.ErrNotFound)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Product) bool {
		return p.Name == req.Name && p.Slug == "test-product" && p.SKU == req.SKU && p.Price == req.Price
	})).Return(nil)

	cache.On("DeletePrefix", mock.Anything, "product:").Return(nil)

	product, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "test-product", product.Slug)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestProductService_FindByID_CacheHit(t *testing.T) {
	repo := new(MockProductRepository)
	cache := new(MockCacheManager)
	svc := NewProductService(repo, cache)

	id := uuid.New()
	expected := &models.Product{ID: id, Name: "Cached Product", SKU: "SKU-CACHE"}

	cache.On("Get", mock.Anything, "product:id:"+id.String(), mock.AnythingOfType("*models.Product")).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*models.Product)
			*dest = *expected
		}).Return(nil)

	product, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, product.Name)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestProductService_FindByID_CacheMiss(t *testing.T) {
	repo := new(MockProductRepository)
	cache := new(MockCacheManager)
	svc := NewProductService(repo, cache)

	id := uuid.New()
	expected := &models.Product{ID: id, Name: "DB Product", SKU: "SKU-DB"}

	cache.On("Get", mock.Anything, "product:id:"+id.String(), mock.Anything).Return(errors.New("not found"))
	repo.On("FindByID", mock.Anything, id).Return(expected, nil)
	cache.On("Set", mock.Anything, "product:id:"+id.String(), expected, 72*time.Hour).Return(nil)

	product, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, product.Name)
}

func TestProductService_Update(t *testing.T) {
	repo := new(MockProductRepository)
	cache := new(MockCacheManager)
	svc := NewProductService(repo, cache)

	id := uuid.New()
	existing := &models.Product{ID: id, Name: "Old Name", SKU: "SKU-OLD", Price: 10.0}
	newPrice := 20.0
	req := dto.UpdateProductRequest{
		Name:  "New Name",
		SKU:   "SKU-NEW",
		Price: &newPrice,
	}

	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("FindBySKU", mock.Anything, "SKU-NEW").Return((*models.Product)(nil), apperrors.ErrNotFound)
	repo.On("FindBySlug", mock.Anything, "new-name").Return((*models.Product)(nil), apperrors.ErrNotFound)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Product) bool {
		return p.Name == "New Name" && p.Slug == "new-name" && p.SKU == "SKU-NEW" && p.Price == 20.0
	})).Return(nil)
	cache.On("DeletePrefix", mock.Anything, "product:").Return(nil)

	product, err := svc.Update(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, "new-name", product.Slug)
	assert.Equal(t, "SKU-NEW", product.SKU)
}

func TestProductService_Delete_NotFound(t *testing.T) {
	repo := new(MockProductRepository)
	cache := new(MockCacheManager)
	svc := NewProductService(repo, cache)

	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(apperrors.ErrNotFound)

	err := svc.Delete(context.Background(), id)

	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	cache.AssertNotCalled(t, "DeletePrefix", mock.Anything, mock.Anything)
}
