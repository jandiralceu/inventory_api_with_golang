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

func TestWarehouseService_Create(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	req := dto.CreateWarehouseRequest{
		Name: "Test Warehouse",
		Code: "WH-TEST",
		Address: dto.WarehouseAddress{
			Street: "Street", Number: "1", City: "City", State: "ST", Country: "BR", ZipCode: "123",
		},
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(w *models.Warehouse) bool {
		return w.Name == req.Name && w.Slug == "test-warehouse" && w.Code == req.Code
	})).Return(nil)

	cache.On("DeletePrefix", mock.Anything, "warehouse:").Return(nil)

	warehouse, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, warehouse)
	assert.Equal(t, "test-warehouse", warehouse.Slug)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestWarehouseService_Update(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	existing := &models.Warehouse{ID: id, Name: "Old Name"}
	isActive := true
	req := dto.UpdateWarehouseRequest{
		Name: "New Name",
		Code: "WH-NEW",
		Address: dto.WarehouseAddress{
			Street: "New Street", Number: "2", City: "New City", State: "NS", Country: "BR", ZipCode: "456",
		},
		IsActive: &isActive,
	}

	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(w *models.Warehouse) bool {
		return w.Name == "New Name" && w.Slug == "new-name" && w.IsActive == true
	})).Return(nil)
	cache.On("DeletePrefix", mock.Anything, "warehouse:").Return(nil)

	warehouse, err := svc.Update(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, "new-name", warehouse.Slug)
}

func TestWarehouseService_FindByID_CacheHit(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	id := uuid.New()
	expected := &models.Warehouse{ID: id, Name: "Cached Warehouse"}

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
	expected := &models.Warehouse{ID: id, Name: "DB Warehouse"}

	cache.On("Get", mock.Anything, "warehouse:id:"+id.String(), mock.Anything).Return(errors.New("not found"))
	repo.On("FindByID", mock.Anything, id).Return(expected, nil)
	cache.On("Set", mock.Anything, "warehouse:id:"+id.String(), expected, 72*time.Hour).Return(nil)

	warehouse, err := svc.FindByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, warehouse.Name)
}

func TestWarehouseService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockWarehouseRepository)
		cache := new(MockCacheManager)
		svc := NewWarehouseService(repo, cache)

		id := uuid.New()
		repo.On("Delete", mock.Anything, id).Return(nil)
		cache.On("DeletePrefix", mock.Anything, "warehouse:").Return(nil)

		err := svc.Delete(context.Background(), id)

		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := new(MockWarehouseRepository)
		cache := new(MockCacheManager)
		svc := NewWarehouseService(repo, cache)

		id := uuid.New()
		repo.On("Delete", mock.Anything, id).Return(apperrors.ErrNotFound)

		err := svc.Delete(context.Background(), id)

		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		cache.AssertNotCalled(t, "DeletePrefix", mock.Anything, mock.Anything)
	})
}
func TestWarehouseService_FindAll(t *testing.T) {
	repo := new(MockWarehouseRepository)
	cache := new(MockCacheManager)
	svc := NewWarehouseService(repo, cache)

	req := dto.GetWarehouseListRequest{
		PaginationRequest: dto.PaginationRequest{
			Page:  1,
			Limit: 10,
		},
	}

	expectedWarehouses := []models.Warehouse{{ID: uuid.New(), Name: "Warehouse 1"}}
	var expectedTotal int64 = 1

	repo.On("FindAll", mock.Anything, mock.AnythingOfType("repository.WarehouseListFilter")).
		Return(expectedWarehouses, expectedTotal, nil)

	resp, err := svc.FindAll(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, expectedTotal, resp.Total)
}
