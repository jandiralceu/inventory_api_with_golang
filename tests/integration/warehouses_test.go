//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouseManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole, managerRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "manager").First(&managerRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL
	adminEmail := "whadmin@example.com"
	managerEmail := "whmanager@example.com"
	operatorEmail := "whoperator@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "WH Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	// Create and login as Manager
	signUpUser(t, baseURL, "WH Manager", managerEmail, password, managerRole.ID.String())
	managerToken, _ := signInUser(t, baseURL, managerEmail, password)

	// Create and login as Operator
	signUpUser(t, baseURL, "WH Operator", operatorEmail, password, operatorRole.ID.String())
	operatorToken, _ := signInUser(t, baseURL, operatorEmail, password)

	t.Run("Manager can create and update a warehouse", func(t *testing.T) {
		// 1. Create
		createReq := dto.CreateWarehouseRequest{
			Name: "Main Logistics Hub",
			Code: "WH-MA-01",
			Address: dto.WarehouseAddress{
				Street:  "Industrial Blvd",
				Number:  "1000",
				City:    "Chicago",
				State:   "IL",
				Country: "USA",
				ZipCode: "60601",
			},
			ManagerName: "Jack Bauer",
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/warehouses", managerToken, createReq)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Warehouse
		json.NewDecoder(resp.Body).Decode(&created)
		assert.Equal(t, createReq.Name, created.Name)
		assert.Equal(t, "main-logistics-hub", created.Slug)
		assert.Equal(t, createReq.Code, created.Code)

		// 2. Update
		active := true
		updateReq := dto.UpdateWarehouseRequest{
			Name: "Main Logistics Global",
			Code: "WH-MA-01-G",
			Address: dto.WarehouseAddress{
				Street:  "Industrial Blvd",
				Number:  "1001",
				City:    "Chicago",
				State:   "IL",
				Country: "USA",
				ZipCode: "60601",
			},
			IsActive: &active,
		}

		respUpdate := authedRequest(t, "PUT", fmt.Sprintf("%s/api/v1/warehouses/%s", baseURL, created.ID), managerToken, updateReq)
		assert.Equal(t, http.StatusOK, respUpdate.StatusCode)

		var updated models.Warehouse
		json.NewDecoder(respUpdate.Body).Decode(&updated)
		assert.Equal(t, "Main Logistics Global", updated.Name)
		assert.Equal(t, "WH-MA-01-G", updated.Code)
		assert.Equal(t, "1001", updated.Address.Number)
	})

	t.Run("Operator can list and filter warehouses", func(t *testing.T) {
		resp := authedRequest(t, "GET", baseURL+"/api/v1/warehouses?name=Logistics", operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var list dto.WarehouseListResponse
		json.NewDecoder(resp.Body).Decode(&list)
		assert.GreaterOrEqual(t, len(list.Data), 1)
		assert.Contains(t, list.Data[0].Name, "Logistics")
	})

	t.Run("Operator can find by ID", func(t *testing.T) {
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/warehouses?limit=1", operatorToken, nil)
		var list dto.WarehouseListResponse
		json.NewDecoder(listResp.Body).Decode(&list)
		id := list.Data[0].ID

		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/warehouses/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var found models.Warehouse
		json.NewDecoder(resp.Body).Decode(&found)
		assert.Equal(t, id, found.ID)
	})

	t.Run("RBAC: Operator cannot delete warehouses", func(t *testing.T) {
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/warehouses?limit=1", operatorToken, nil)
		var list dto.WarehouseListResponse
		json.NewDecoder(listResp.Body).Decode(&list)
		id := list.Data[0].ID

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/warehouses/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Admin can delete warehouse", func(t *testing.T) {
		cReq := dto.CreateWarehouseRequest{
			Name: "Delete Me WH",
			Code: "WH-DEL-01",
			Address: dto.WarehouseAddress{
				Street: "Main Street", Number: "0", City: "New York", State: "NY", Country: "USA", ZipCode: "00000",
			},
		}
		cResp := authedRequest(t, "POST", baseURL+"/api/v1/warehouses", managerToken, cReq)
		var toDel models.Warehouse
		json.NewDecoder(cResp.Body).Decode(&toDel)

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/warehouses/%s", baseURL, toDel.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		vResp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/warehouses/%s", baseURL, toDel.ID), adminToken, nil)
		assert.Equal(t, http.StatusNotFound, vResp.StatusCode)
	})
}
