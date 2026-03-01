package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_Create(t *testing.T) {
	repo := new(MockSupplierRepository)
	cache := new(MockCacheManager)
	svc := NewSupplierService(repo, cache)

	req := dto.CreateSupplierRequest{
		Name:  "New Supplier",
		TaxID: "999888777",
		Address: dto.SupplierAddress{
			Street:  "Tech St",
			City:    "New York",
			State:   "NY",
			Country: "USA",
			ZipCode: "10001",
			Number:  "123",
		},
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
		return s.Name == req.Name && s.Slug == "new-supplier"
	})).Return(nil)

	cache.On("DeletePrefix", mock.Anything, "supplier:").Return(nil)

	supplier, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, supplier)
	assert.Equal(t, "new-supplier", supplier.Slug)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestSupplierService_FindByID_CacheHit(t *testing.T) {
	repo := new(MockSupplierRepository)
	cache := new(MockCacheManager)
	svc := NewSupplierService(repo, cache)

	id := uuid.New()
	expected := &models.Supplier{ID: id, Name: "Cached Supplier"}

	cache.On("Get", mock.Anything, "supplier:id:"+id.String(), mock.AnythingOfType("*models.Supplier")).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*models.Supplier)
			*dest = *expected
		}).Return(nil)

	supplier, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, supplier.Name)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestSupplierService_FindByID_CacheMiss(t *testing.T) {
	repo := new(MockSupplierRepository)
	cache := new(MockCacheManager)
	svc := NewSupplierService(repo, cache)

	id := uuid.New()
	expected := &models.Supplier{ID: id, Name: "DB Supplier"}

	cache.On("Get", mock.Anything, "supplier:id:"+id.String(), mock.Anything).Return(errors.New("not found"))
	repo.On("FindByID", mock.Anything, id).Return(expected, nil)
	cache.On("Set", mock.Anything, "supplier:id:"+id.String(), expected, mock.Anything).Return(nil)

	supplier, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, supplier.Name)
}

func TestSupplierService_Update(t *testing.T) {
	repo := new(MockSupplierRepository)
	cache := new(MockCacheManager)
	svc := NewSupplierService(repo, cache)

	id := uuid.New()
	existing := &models.Supplier{ID: id, Name: "Old Name"}
	req := dto.UpdateSupplierRequest{
		Name: "New Name",
		Address: dto.SupplierAddress{
			Street:  "Tech St",
			City:    "New York",
			State:   "NY",
			Country: "USA",
			ZipCode: "10001",
			Number:  "123",
		},
		IsActive: ptrBool(true),
	}

	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
		return s.Name == "New Name" && s.Slug == "new-name"
	})).Return(nil)
	cache.On("DeletePrefix", mock.Anything, "supplier:").Return(nil)

	supplier, err := svc.Update(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, "new-name", supplier.Slug)
}

func TestSupplierService_Delete_NotFound(t *testing.T) {
	repo := new(MockSupplierRepository)
	cache := new(MockCacheManager)
	svc := NewSupplierService(repo, cache)

	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(apperrors.ErrNotFound)

	err := svc.Delete(context.Background(), id)

	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	cache.AssertNotCalled(t, "DeletePrefix", mock.Anything, mock.Anything)
}

func ptrBool(b bool) *bool { return &b }
