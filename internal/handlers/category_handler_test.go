package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService is a mock implementation of service.CategoryService.
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*models.Category, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*models.Category, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryService) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryService) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryService) FindAll(ctx context.Context, req dto.GetCategoryListRequest) (dto.PaginatedResponse[models.Category], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.PaginatedResponse[models.Category]), args.Error(1)
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupRouter()
	router.POST("/categories", handler.CreateCategory)

	req := dto.CreateCategoryRequest{
		Name: "Electronics",
	}
	expected := &models.Category{ID: uuid.New(), Name: "Electronics", Slug: "electronics"}

	mockService.On("Create", mock.Anything, req).Return(expected, nil)

	w := performRequest(router, "POST", "/categories", req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Category
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestCategoryHandler_FindCategoryByID(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupRouter()
	router.GET("/categories/:id", handler.FindCategoryByID)

	id := uuid.New()
	expected := &models.Category{ID: id, Name: "Electronics"}

	mockService.On("FindByID", mock.Anything, id).Return(expected, nil)

	w := performRequest(router, "GET", fmt.Sprintf("/categories/%s", id), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Category
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, id, resp.ID)
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupRouter()
	router.DELETE("/categories/:id", handler.DeleteCategory)

	id := uuid.New()
	mockService.On("Delete", mock.Anything, id).Return(nil)

	w := performRequest(router, "DELETE", fmt.Sprintf("/categories/%s", id), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCategoryHandler_FindAllCategories(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupRouter()
	router.GET("/categories", handler.FindAllCategories)

	expected := dto.PaginatedResponse[models.Category]{
		Data:  []models.Category{{Name: "Electronics"}},
		Total: 1,
	}

	mockService.On("FindAll", mock.Anything, mock.MatchedBy(func(req dto.GetCategoryListRequest) bool {
		return req.Page == 1 && req.Limit == 10
	})).Return(expected, nil)

	w := performRequest(router, "GET", "/categories?page=1&limit=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse[models.Category]
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, int64(1), resp.Total)
}
