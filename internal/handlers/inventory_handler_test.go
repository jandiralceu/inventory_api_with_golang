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

// MockInventoryService is a mock implementation of service.InventoryService.
type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) Create(ctx context.Context, req dto.CreateInventoryRequest) (*models.Inventory, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (*models.Inventory, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInventoryService) FindAll(ctx context.Context, req dto.GetInventoryListRequest) (*dto.InventoryListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InventoryListResponse), args.Error(1)
}

func (m *MockInventoryService) FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryService) AddStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) RemoveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) GetTransactionHistory(ctx context.Context, req dto.TransactionListRequest) (dto.TransactionListResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.TransactionListResponse), args.Error(1)
}

func TestInventoryHandler_CreateInventory(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.POST("/inventory", handler.CreateInventory)

	req := dto.CreateInventoryRequest{
		ProductID:    uuid.New(),
		WarehouseID:  uuid.New(),
		Quantity:     10,
		LocationCode: "A1",
	}
	expected := &models.Inventory{ID: uuid.New(), ProductID: req.ProductID, WarehouseID: req.WarehouseID}

	mockService.On("Create", mock.Anything, req).Return(expected, nil)

	w := performRequest(router, "POST", "/inventory", req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Inventory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, expected.ID, resp.ID)
}

func TestInventoryHandler_FindInventoryByID(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.GET("/inventory/:id", handler.FindInventoryByID)

	id := uuid.New()
	expected := &models.Inventory{ID: id, Quantity: 50}

	mockService.On("FindByID", mock.Anything, id).Return(expected, nil)

	w := performRequest(router, "GET", fmt.Sprintf("/inventory/%s", id), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Inventory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, id, resp.ID)
}

func TestInventoryHandler_AddStock(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.POST("/inventory/:id/add", handler.AddStock)

	id := uuid.New()
	req := dto.StockOperationRequest{Quantity: 10}

	mockService.On("AddStock", mock.Anything, id, 10).Return(nil)

	w := performRequest(router, "POST", fmt.Sprintf("/inventory/%s/add", id), req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInventoryHandler_GetTransactionHistory(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.GET("/inventory/transactions", handler.GetTransactionHistory)

	expected := dto.TransactionListResponse{
		PaginatedResponse: dto.PaginatedResponse[dto.TransactionResponse]{
			Data:  []dto.TransactionResponse{{ID: uuid.New(), TransactionType: "ADJUSTMENT"}},
			Total: 1,
		},
	}

	mockService.On("GetTransactionHistory", mock.Anything, mock.MatchedBy(func(r dto.TransactionListRequest) bool {
		return true
	})).Return(expected, nil)

	w := performRequest(router, "GET", "/inventory/transactions?page=1&limit=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInventoryHandler_UpdateInventory_Success(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.PUT("/inventory/:id", handler.UpdateInventory)

	id := uuid.New()
	qty := 200
	req := dto.UpdateInventoryRequest{
		Quantity: &qty,
	}
	expected := &models.Inventory{ID: id, Quantity: 200}

	mockService.On("Update", mock.Anything, id, req).Return(expected, nil)

	w := performRequest(router, "PUT", fmt.Sprintf("/inventory/%s", id), req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Inventory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 200, resp.Quantity)
}

func TestInventoryHandler_DeleteInventory_Success(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.DELETE("/inventory/:id", handler.DeleteInventory)

	id := uuid.New()
	mockService.On("Delete", mock.Anything, id).Return(nil)

	w := performRequest(router, "DELETE", fmt.Sprintf("/inventory/%s", id), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestInventoryHandler_StockOperations_Errors(t *testing.T) {
	mockService := new(MockInventoryService)
	handler := NewInventoryHandler(mockService)
	router := setupRouter()
	router.POST("/inventory/:id/remove", handler.RemoveStock)
	router.POST("/inventory/:id/reserve", handler.ReserveStock)
	router.POST("/inventory/:id/release", handler.ReleaseStock)

	id := uuid.New()
	req := dto.StockOperationRequest{Quantity: 100}

	t.Run("RemoveStock_Insufficient", func(t *testing.T) {
		mockService.On("RemoveStock", mock.Anything, id, 100).Return(apperrors.ErrInvalidInput).Once()
		w := performRequest(router, "POST", fmt.Sprintf("/inventory/%s/remove", id), req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ReserveStock_NotFound", func(t *testing.T) {
		mockService.On("ReserveStock", mock.Anything, id, 100).Return(apperrors.ErrNotFound).Once()
		w := performRequest(router, "POST", fmt.Sprintf("/inventory/%s/reserve", id), req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ReleaseStock_Success", func(t *testing.T) {
		mockService.On("ReleaseStock", mock.Anything, id, 100).Return(nil).Once()
		w := performRequest(router, "POST", fmt.Sprintf("/inventory/%s/release", id), req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
