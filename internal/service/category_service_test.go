package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	req := dto.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Devices and gadgets",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Category")).Return(nil)
	mockCache.On("DeletePrefix", ctx, "category:").Return(nil)

	category, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Electronics", category.Name)
	assert.Equal(t, "electronics", category.Slug)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCategoryService_Create_WithParent(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	parentID := uuid.New()
	req := dto.CreateCategoryRequest{
		Name:     "Smartphones",
		ParentID: &parentID,
	}

	// Mock parent exists
	mockRepo.On("FindByID", ctx, parentID).Return(&models.Category{ID: parentID}, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Category")).Return(nil)
	mockCache.On("DeletePrefix", ctx, "category:").Return(nil)

	category, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, &parentID, category.ParentID)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create_ParentNotFound(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	parentID := uuid.New()
	req := dto.CreateCategoryRequest{
		Name:     "Smartphones",
		ParentID: &parentID,
	}

	mockRepo.On("FindByID", ctx, parentID).Return(nil, apperrors.ErrNotFound)

	_, err := svc.Create(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestCategoryService_FindByID_CacheHit(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	id := uuid.New()
	expectedCategory := &models.Category{ID: id, Name: "Electronics"}

	mockCache.On("Get", ctx, "category:id:"+id.String(), mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(2).(*models.Category)
		*dest = *expectedCategory
	}).Return(nil)

	category, err := svc.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.Name, category.Name)
	mockRepo.AssertNotCalled(t, "FindByID", ctx, id)
}

func TestCategoryService_FindByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	id := uuid.New()
	expectedCategory := &models.Category{ID: id, Name: "Electronics"}

	mockCache.On("Get", ctx, "category:id:"+id.String(), mock.Anything).Return(errors.New("cache miss"))
	mockRepo.On("FindByID", ctx, id).Return(expectedCategory, nil)
	mockCache.On("Set", ctx, "category:id:"+id.String(), expectedCategory, mock.Anything).Return(nil)

	category, err := svc.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.Name, category.Name)
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("Delete", ctx, id).Return(nil)
	mockCache.On("DeletePrefix", ctx, "category:").Return(nil)

	err := svc.Delete(ctx, id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCategoryService_FindAll(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockCache := new(MockCacheManager)
	svc := NewCategoryService(mockRepo, mockCache)

	ctx := context.Background()
	req := dto.GetCategoryListRequest{
		Name: "",
	}
	req.Page = 1
	req.Limit = 10
	expectedCategories := []models.Category{{Name: "C1"}, {Name: "C2"}}

	mockRepo.On("FindAll", ctx, mock.MatchedBy(func(f repository.CategoryListFilter) bool {
		return f.Pagination.Limit == 10 && f.Pagination.Page == 1
	})).Return(expectedCategories, int64(2), nil)

	res, err := svc.FindAll(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), res.Total)
	assert.Len(t, res.Data, 2)
}
