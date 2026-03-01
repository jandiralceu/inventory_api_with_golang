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

func TestWarehouseService_Create(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	req := dto.CreateWarehouseRequest{
		Name: "New Warehouse",
		Code: "WH-01",
		Address: dto.WarehouseAddress{
			Street:  "Main St",
			City:    "New York",
			State:   "NY",
			Country: "USA",
			ZipCode: "10001",
			Number:  "123",
		},
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(w *models.Warehouse) bool {
		return w.Name == req.Name && w.Slug == "new-warehouse" && w.Code == req.Code
	})).Return(nil)

	cache.On("DeletePrefix", mock.Anything, "warehouse:").Return(nil)

	warehouse, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, warehouse)
	assert.Equal(t, "new-warehouse", warehouse.Slug)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestWarehouseService_FindByID_CacheHit(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	expected := &models.Warehouse{ID: id, Name: "Cached Warehouse", Code: "WH-01"}

	cache.On("Get", mock.Anything, "warehouse:id:"+id.String(), mock.AnythingOfType("*models.Warehouse")).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*models.Warehouse)
			*dest = *expected
		}).Return(nil)

	warehouse, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, warehouse.Name)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestWarehouseService_FindByID_CacheMiss(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	expected := &models.Warehouse{ID: id, Name: "DB Warehouse", Code: "WH-01"}

	cache.On("Get", mock.Anything, "warehouse:id:"+id.String(), mock.Anything).Return(errors.New("not found"))
	repo.On("FindByID", mock.Anything, id).Return(expected, nil)
	cache.On("Set", mock.Anything, "warehouse:id:"+id.String(), expected, mock.Anything).Return(nil)

	warehouse, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, warehouse.Name)
}

func TestWarehouseService_Update(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	existing := &models.Warehouse{ID: id, Name: "Old Name", Code: "WH-01"}
	req := dto.UpdateWarehouseRequest{
		Name: "New Name",
		Code: "WH-01-UPD",
		Address: dto.WarehouseAddress{
			Street:  "Main St",
			City:    "New York",
			State:   "NY",
			Country: "USA",
			ZipCode: "10001",
			Number:  "123",
		},
		IsActive: ptrBool(true),
	}

	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(w *models.Warehouse) bool {
		return w.Name == "New Name" && w.Slug == "new-name" && w.Code == "WH-01-UPD"
	})).Return(nil)
	cache.On("DeletePrefix", mock.Anything, "warehouse:").Return(nil)

	warehouse, err := svc.Update(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, "new-name", warehouse.Slug)
}

func TestWarehouseService_Delete_NotFound(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(apperrors.ErrNotFound)

	err := svc.Delete(context.Background(), id)

	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	cache.AssertNotCalled(t, "DeletePrefix", mock.Anything, mock.Anything)
}
