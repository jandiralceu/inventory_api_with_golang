package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockWarehouseService is a mock implementation of service.WarehouseService.
type MockWarehouseService struct {
	mock.Mock
}

func (m *MockWarehouseService) Create(ctx context.Context, req dto.CreateWarehouseRequest) (*models.Warehouse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateWarehouseRequest) (*models.Warehouse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWarehouseService) FindByID(ctx context.Context, id uuid.UUID) (*models.Warehouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) FindBySlug(ctx context.Context, slug string) (*models.Warehouse, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) FindAll(ctx context.Context, req dto.GetWarehouseListRequest) (dto.PaginatedResponse[models.Warehouse], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.PaginatedResponse[models.Warehouse]), args.Error(1)
}

func TestWarehouseHandler_CreateWarehouse(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.POST("/warehouses", handler.CreateWarehouse)

		req := dto.CreateWarehouseRequest{
			Name: "Main Warehouse",
			Code: "WH01",
			Address: dto.WarehouseAddress{
				Street:  "Main Street",
				Number:  "123",
				City:    "New York",
				State:   "NY",
				Country: "USA",
				ZipCode: "10001",
			},
		}
		expected := &models.Warehouse{ID: uuid.New(), Name: req.Name, Code: req.Code}

		mockService.On("Create", mock.Anything, req).Return(expected, nil)

		w := performRequest(router, "POST", "/warehouses", req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Warehouse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, expected.ID, resp.ID)
		assert.Equal(t, expected.Name, resp.Name)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.POST("/warehouses", handler.CreateWarehouse)

		req := dto.CreateWarehouseRequest{
			Name: "AB", // Too short
		}

		w := performRequest(router, "POST", "/warehouses", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWarehouseHandler_UpdateWarehouse(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.PUT("/warehouses/:id", handler.UpdateWarehouse)

		id := uuid.New()
		isActive := true
		req := dto.UpdateWarehouseRequest{
			Name: "Updated Warehouse",
			Code: "WH01-UPD",
			Address: dto.WarehouseAddress{
				Street:  "Updated St",
				Number:  "456",
				City:    "Boston",
				State:   "MA",
				Country: "USA",
				ZipCode: "02108",
			},
			IsActive: &isActive,
		}
		expected := &models.Warehouse{ID: id, Name: req.Name, Code: req.Code, IsActive: isActive}

		mockService.On("Update", mock.Anything, id, req).Return(expected, nil)

		w := performRequest(router, "PUT", fmt.Sprintf("/warehouses/%s", id), req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Warehouse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, expected.Name, resp.Name)
		assert.Equal(t, true, resp.IsActive)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.PUT("/warehouses/:id", handler.UpdateWarehouse)

		id := uuid.New()
		isActive := true
		req := dto.UpdateWarehouseRequest{
			Name: "Updated Warehouse",
			Code: "WH01-UPD",
			Address: dto.WarehouseAddress{
				Street:  "Valid Street",
				Number:  "1",
				City:    "NY",
				State:   "NY",
				Country: "US",
				ZipCode: "10001",
			},
			IsActive: &isActive,
		}

		mockService.On("Update", mock.Anything, id, req).Return(nil, apperrors.ErrNotFound)

		w := performRequest(router, "PUT", fmt.Sprintf("/warehouses/%s", id), req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWarehouseHandler_FindWarehouseByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.GET("/warehouses/:id", handler.FindWarehouseByID)

		id := uuid.New()
		expected := &models.Warehouse{ID: id, Name: "Test Warehouse"}

		mockService.On("FindByID", mock.Anything, id).Return(expected, nil)

		w := performRequest(router, "GET", fmt.Sprintf("/warehouses/%s", id), nil)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Warehouse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, id, resp.ID)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.GET("/warehouses/:id", handler.FindWarehouseByID)

		w := performRequest(router, "GET", "/warehouses/invalid-uuid", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWarehouseHandler_DeleteWarehouse(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.DELETE("/warehouses/:id", handler.DeleteWarehouse)

		id := uuid.New()
		mockService.On("Delete", mock.Anything, id).Return(nil)

		w := performRequest(router, "DELETE", fmt.Sprintf("/warehouses/%s", id), nil)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockWarehouseService)
		handler := NewWarehouseHandler(mockService)
		router := setupRouter()
		router.DELETE("/warehouses/:id", handler.DeleteWarehouse)

		id := uuid.New()
		mockService.On("Delete", mock.Anything, id).Return(apperrors.ErrNotFound)

		w := performRequest(router, "DELETE", fmt.Sprintf("/warehouses/%s", id), nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWarehouseHandler_FindAllWarehouses(t *testing.T) {
	mockService := new(MockWarehouseService)
	handler := NewWarehouseHandler(mockService)
	router := setupRouter()
	router.GET("/warehouses", handler.FindAllWarehouses)

	expected := dto.PaginatedResponse[models.Warehouse]{
		Data:  []models.Warehouse{{ID: uuid.New(), Name: "Warehouse 1"}},
		Total: 1,
		Page:  1,
		Limit: 10,
	}

	mockService.On("FindAll", mock.Anything, mock.AnythingOfType("dto.GetWarehouseListRequest")).Return(expected, nil)

	w := performRequest(router, "GET", "/warehouses?page=1&limit=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse[models.Warehouse]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, int64(1), resp.Total)
}
