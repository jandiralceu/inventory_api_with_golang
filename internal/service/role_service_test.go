package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock RoleRepository ---

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

// =====================
// Create Tests
// =====================

func TestRoleServiceCreateSuccess(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	role := &models.Role{Name: "admin", Description: "Administrator role"}
	created := &models.Role{ID: uuid.New(), Name: role.Name, Description: role.Description}

	mockRepo.On("Create", mock.Anything, role).Return(created, nil)
	mockCache.On("DeletePrefix", mock.Anything, "role:").Return(nil)

	result, err := svc.Create(context.Background(), role)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.ID, result.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRoleServiceCreateError(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	role := &models.Role{Name: "admin", Description: "Administrator role"}

	mockRepo.On("Create", mock.Anything, role).Return(nil, errors.New("duplicate key"))

	result, err := svc.Create(context.Background(), role)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertNotCalled(t, "DeletePrefix")
}

// =====================
// FindAll Tests
// =====================

func TestRoleServiceFindAllCacheMiss(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	expected := []models.Role{
		{ID: uuid.New(), Name: "admin"},
		{ID: uuid.New(), Name: "editor"},
	}

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("miss"))
	mockRepo.On("FindAll", mock.Anything).Return(expected, nil)
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)

	result, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "admin", result[0].Name)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRoleServiceFindAllCacheHit(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)

	_, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "FindAll")
}

func TestRoleServiceFindAllRepoError(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("miss"))
	mockRepo.On("FindAll", mock.Anything).Return(nil, errors.New("db error"))

	result, err := svc.FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertNotCalled(t, "Set")
}

// =====================
// FindByID Tests
// =====================

func TestRoleServiceFindByIDCacheMiss(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	roleID := uuid.New()
	expected := &models.Role{ID: roleID, Name: "admin"}

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("miss"))
	mockRepo.On("FindByID", mock.Anything, roleID).Return(expected, nil)
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)

	result, err := svc.FindByID(context.Background(), roleID)

	assert.NoError(t, err)
	assert.Equal(t, roleID, result.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRoleServiceFindByIDCacheHit(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	roleID := uuid.New()

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)

	_, err := svc.FindByID(context.Background(), roleID)

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestRoleServiceFindByIDNotFound(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	roleID := uuid.New()

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(errors.New("miss"))
	mockRepo.On("FindByID", mock.Anything, roleID).Return(nil, errors.New("not found"))

	result, err := svc.FindByID(context.Background(), roleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertNotCalled(t, "Set")
}

// =====================
// Delete Tests
// =====================

func TestRoleServiceDeleteSuccess(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	roleID := uuid.New()

	mockRepo.On("Delete", mock.Anything, roleID).Return(nil)
	mockCache.On("DeletePrefix", mock.Anything, "role:").Return(nil)

	err := svc.Delete(context.Background(), roleID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRoleServiceDeleteError(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	mockCache := new(MockCacheManager)
	svc := NewRoleService(mockRepo, mockCache)

	roleID := uuid.New()

	mockRepo.On("Delete", mock.Anything, roleID).Return(errors.New("not found"))

	err := svc.Delete(context.Background(), roleID)

	assert.Error(t, err)
	mockCache.AssertNotCalled(t, "DeletePrefix")
}
