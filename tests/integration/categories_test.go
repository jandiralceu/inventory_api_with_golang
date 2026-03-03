//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole, managerRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "manager").First(&managerRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL
	adminEmail := "catadmin@example.com"
	managerEmail := "catmanager@example.com"
	operatorEmail := "catoperator@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "Cat Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	// Create and login as Manager
	signUpUser(t, baseURL, "Cat Manager", managerEmail, password, managerRole.ID.String())
	managerToken, _ := signInUser(t, baseURL, managerEmail, password)

	// Create and login as Operator
	signUpUser(t, baseURL, "Cat Operator", operatorEmail, password, operatorRole.ID.String())
	operatorToken, _ := signInUser(t, baseURL, operatorEmail, password)

	t.Run("Manager can create and update a category", func(t *testing.T) {
		// 1. Create
		createReq := dto.CreateCategoryRequest{
			Name:        "Electronics",
			Description: "Devices and gadgets",
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/categories", managerToken, createReq)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Category
		require.NoError(t, decodeResponse(resp, &created))
		assert.Equal(t, "Electronics", created.Name)
		assert.Equal(t, "electronics", created.Slug)

		// 2. Update
		isActive := true
		updateReq := dto.UpdateCategoryRequest{
			Name:        "Electronics & Gadgets",
			Description: "Updated description",
			IsActive:    &isActive,
		}

		respUpdate := authedRequest(t, "PUT", fmt.Sprintf("%s/api/v1/categories/%s", baseURL, created.ID), managerToken, updateReq)
		assert.Equal(t, http.StatusOK, respUpdate.StatusCode)

		var updated models.Category
		require.NoError(t, decodeResponse(respUpdate, &updated))
		assert.Equal(t, updateReq.Name, updated.Name)
		assert.Equal(t, "electronics-gadgets", updated.Slug)
	})

	t.Run("Manager can create subcategory", func(t *testing.T) {
		// Parent
		pReq := dto.CreateCategoryRequest{Name: "Computers"}
		pResp := authedRequest(t, "POST", baseURL+"/api/v1/categories", managerToken, pReq)
		var parent models.Category
		require.NoError(t, decodeResponse(pResp, &parent))

		// Subcategory
		subReq := dto.CreateCategoryRequest{
			Name:     "Laptops",
			ParentID: &parent.ID,
		}
		subResp := authedRequest(t, "POST", baseURL+"/api/v1/categories", managerToken, subReq)
		assert.Equal(t, http.StatusCreated, subResp.StatusCode)

		var sub models.Category
		require.NoError(t, decodeResponse(subResp, &sub))
		assert.Equal(t, &parent.ID, sub.ParentID)
	})

	t.Run("Operator can list and filter categories", func(t *testing.T) {
		// Note: Operator should have GET permission according to policy.csv
		resp := authedRequest(t, "GET", baseURL+"/api/v1/categories?limit=5", operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var list dto.PaginatedResponse[models.Category]
		require.NoError(t, decodeResponse(resp, &list))
		assert.GreaterOrEqual(t, len(list.Data), 2)
	})

	t.Run("Operator can find by ID", func(t *testing.T) {
		// First get an ID (from a category created by manager)
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/categories?limit=1", operatorToken, nil)
		var list dto.PaginatedResponse[models.Category]
		require.NoError(t, decodeResponse(listResp, &list))
		id := list.Data[0].ID

		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/categories/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("RBAC: Operator cannot delete categories", func(t *testing.T) {
		// Get an ID
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/categories?limit=1", operatorToken, nil)
		var list dto.PaginatedResponse[models.Category]
		require.NoError(t, decodeResponse(listResp, &list))
		id := list.Data[0].ID

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/categories/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Admin can delete category", func(t *testing.T) {
		// Create one to delete
		cReq := dto.CreateCategoryRequest{Name: "To Delete"}
		cResp := authedRequest(t, "POST", baseURL+"/api/v1/categories", managerToken, cReq)
		var toDel models.Category
		require.NoError(t, decodeResponse(cResp, &toDel))

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/categories/%s", baseURL, toDel.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Unauthorized access is blocked", func(t *testing.T) {
		resp := authedRequest(t, "GET", baseURL+"/api/v1/categories", "", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
