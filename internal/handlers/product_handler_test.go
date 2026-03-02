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
)

// MockProductService is a mock implementation of service.ProductService.
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(ctx context.Context, req dto.CreateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) FindAll(ctx context.Context, req dto.GetProductListRequest) (*dto.ProductListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductListResponse), args.Error(1)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	req := dto.CreateProductRequest{
		Name:  "Test Product",
		SKU:   "SKU-123",
		Price: 99.99,
	}
	expected := &models.Product{ID: uuid.New(), Name: "Test Product", SKU: "SKU-123"}

	mockService.On("Create", mock.Anything, req).Return(expected, nil)

	w := performRequest(router, "POST", "/products", req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Product
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.SKU, resp.SKU)
}

func TestProductHandler_FindProductByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockProductService)
		handler := NewProductHandler(mockService)
		router := setupRouter()
		router.GET("/products/:id", handler.FindProductByID)

		id := uuid.New()
		expected := &models.Product{ID: id, Name: "Test Product"}

		mockService.On("FindByID", mock.Anything, id).Return(expected, nil)

		w := performRequest(router, "GET", fmt.Sprintf("/products/%s", id), nil)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Product
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockProductService)
		handler := NewProductHandler(mockService)
		router := setupRouter()
		router.GET("/products/:id", handler.FindProductByID)

		id := uuid.New()
		mockService.On("FindByID", mock.Anything, id).Return(nil, apperrors.ErrNotFound)

		w := performRequest(router, "GET", fmt.Sprintf("/products/%s", id), nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.PUT("/products/:id", handler.UpdateProduct)

	id := uuid.New()
	price := 150.0
	req := dto.UpdateProductRequest{
		Name:  "Updated Name",
		Price: &price,
	}
	expected := &models.Product{ID: id, Name: "Updated Name", Price: 150.0}

	mockService.On("Update", mock.Anything, id, req).Return(expected, nil)

	w := performRequest(router, "PUT", fmt.Sprintf("/products/%s", id), req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Product
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expected.Name, resp.Name)
	assert.Equal(t, expected.Price, resp.Price)
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.DELETE("/products/:id", handler.DeleteProduct)

	id := uuid.New()
	mockService.On("Delete", mock.Anything, id).Return(nil)

	w := performRequest(router, "DELETE", fmt.Sprintf("/products/%s", id), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestProductHandler_FindAllProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products", handler.FindAllProducts)

	expected := &dto.ProductListResponse{
		PaginatedResponse: dto.PaginatedResponse[models.Product]{
			Data:  []models.Product{{ID: uuid.New(), Name: "Product 1"}},
			Total: 1,
			Page:  1,
			Limit: 10,
		},
	}

	mockService.On("FindAll", mock.Anything, mock.AnythingOfType("dto.GetProductListRequest")).Return(expected, nil)

	w := performRequest(router, "GET", "/products?page=1&limit=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.ProductListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, int64(1), resp.Total)
}
