package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock RoleService ---

type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleService) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleService) FindAll(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

// =====================
// CreateRole Tests
// =====================

func TestCreateRoleSuccess(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Role")).
		Return(&models.Role{
			ID:          roleID,
			Name:        "admin",
			Description: "Administrator role",
		}, nil)

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"name":        "Admin",
		"description": "Administrator role",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Role
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, roleID, response.ID)
	assert.Equal(t, "admin", response.Name)
	mockService.AssertExpectations(t)
}

func TestCreateRoleBadRequestMissingName(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"description": "Some description",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Validation Failed", resp.Title)
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateRoleBadRequestMissingDescription(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"name": "Admin",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateRoleBadRequestNameTooShort(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"name":        "ab",
		"description": "Some description",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateRoleConflict(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Role")).
		Return(nil, fmt.Errorf("%w: role name already exists", apperrors.ErrConflict))

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"name":        "Admin",
		"description": "Administrator role",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Conflict", resp.Title)
	mockService.AssertExpectations(t)
}

func TestCreateRoleInternalServerError(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Role")).
		Return(nil, errors.New("database error"))

	router := setupRouter()
	router.POST("/roles", handler.CreateRole)

	body := map[string]any{
		"name":        "Admin",
		"description": "Administrator role",
	}

	w := performRequest(router, "POST", "/roles", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Internal Server Error", resp.Title)
	assert.Equal(t, "An unexpected error occurred. Please try again later.", resp.Detail)
	mockService.AssertExpectations(t)
}

// =====================
// DeleteRole Tests
// =====================

func TestDeleteRoleSuccess(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("Delete", mock.Anything, roleID).Return(nil)

	router := setupRouter()
	router.DELETE("/roles/:id", handler.DeleteRole)

	w := performRequest(router, "DELETE", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteRoleInvalidUUID(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	router := setupRouter()
	router.DELETE("/roles/:id", handler.DeleteRole)

	w := performRequest(router, "DELETE", "/roles/not-a-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Bad Request", resp.Title)
	mockService.AssertNotCalled(t, "Delete")
}

func TestDeleteRoleNotFound(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("Delete", mock.Anything, roleID).
		Return(apperrors.ErrNotFound)

	router := setupRouter()
	router.DELETE("/roles/:id", handler.DeleteRole)

	w := performRequest(router, "DELETE", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteRoleServiceError(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("Delete", mock.Anything, roleID).
		Return(errors.New("delete failed"))

	router := setupRouter()
	router.DELETE("/roles/:id", handler.DeleteRole)

	w := performRequest(router, "DELETE", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// =====================
// FindRoleByID Tests
// =====================

func TestFindRoleByIDSuccess(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("FindByID", mock.Anything, roleID).
		Return(&models.Role{ID: roleID, Name: "admin", Description: "Administrator role"}, nil)

	router := setupRouter()
	router.GET("/roles/:id", handler.FindRoleByID)

	w := performRequest(router, "GET", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Role
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, roleID, response.ID)
	assert.Equal(t, "admin", response.Name)
	mockService.AssertExpectations(t)
}

func TestFindRoleByIDInvalidUUID(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	router := setupRouter()
	router.GET("/roles/:id", handler.FindRoleByID)

	w := performRequest(router, "GET", "/roles/not-a-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Bad Request", resp.Title)
	mockService.AssertNotCalled(t, "FindByID")
}

func TestFindRoleByIDNotFound(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("FindByID", mock.Anything, roleID).
		Return(nil, apperrors.ErrNotFound)

	router := setupRouter()
	router.GET("/roles/:id", handler.FindRoleByID)

	w := performRequest(router, "GET", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Resource Not Found", resp.Title)
	mockService.AssertExpectations(t)
}

func TestFindRoleByIDInternalServerError(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roleID := uuid.New()
	mockService.On("FindByID", mock.Anything, roleID).
		Return(nil, errors.New("database error"))

	router := setupRouter()
	router.GET("/roles/:id", handler.FindRoleByID)

	w := performRequest(router, "GET", fmt.Sprintf("/roles/%s", roleID), nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// =====================
// FindAllRoles Tests
// =====================

func TestFindAllRolesSuccess(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	roles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator role"},
		{ID: uuid.New(), Name: "editor", Description: "Editor role"},
	}

	mockService.On("FindAll", mock.Anything).Return(roles, nil)

	router := setupRouter()
	router.GET("/roles", handler.FindAllRoles)

	w := performRequest(router, "GET", "/roles", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Role
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "admin", response[0].Name)
	assert.Equal(t, "editor", response[1].Name)
	mockService.AssertExpectations(t)
}

func TestFindAllRolesEmptyList(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	mockService.On("FindAll", mock.Anything).Return([]models.Role{}, nil)

	router := setupRouter()
	router.GET("/roles", handler.FindAllRoles)

	w := performRequest(router, "GET", "/roles", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Role
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
	mockService.AssertExpectations(t)
}

func TestFindAllRolesInternalServerError(t *testing.T) {
	mockService := new(MockRoleService)
	handler := NewRoleHandler(mockService)

	mockService.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	router := setupRouter()
	router.GET("/roles", handler.FindAllRoles)

	w := performRequest(router, "GET", "/roles", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Internal Server Error", resp.Title)
	assert.Equal(t, "An unexpected error occurred. Please try again later.", resp.Detail)
	mockService.AssertExpectations(t)
}
