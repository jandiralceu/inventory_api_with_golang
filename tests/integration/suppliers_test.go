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

func TestSupplierManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole, managerRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "manager").First(&managerRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL
	adminEmail := "suppadmin@example.com"
	managerEmail := "suppmanager@example.com"
	operatorEmail := "suppoperator@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "Supp Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	// Create and login as Manager
	signUpUser(t, baseURL, "Supp Manager", managerEmail, password, managerRole.ID.String())
	managerToken, _ := signInUser(t, baseURL, managerEmail, password)

	// Create and login as Operator
	signUpUser(t, baseURL, "Supp Operator", operatorEmail, password, operatorRole.ID.String())
	operatorToken, _ := signInUser(t, baseURL, operatorEmail, password)

	t.Run("Manager can create and update a supplier", func(t *testing.T) {
		// 1. Create
		createReq := dto.CreateSupplierRequest{
			Name:  "Tech Supplies Co",
			TaxID: "123.456.789-01",
			Email: "contact@techsupplies.com",
			Address: dto.SupplierAddress{
				Street:  "Innovation Ave",
				Number:  "500",
				City:    "Silicon Valley",
				State:   "CA",
				Country: "USA",
				ZipCode: "94025",
			},
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/suppliers", managerToken, createReq)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Supplier
		json.NewDecoder(resp.Body).Decode(&created)
		assert.Equal(t, createReq.Name, created.Name)
		assert.Equal(t, "tech-supplies-co", created.Slug)
		assert.Equal(t, "Innovation Ave", created.Address.Street)

		// 2. Update
		active := true
		updateReq := dto.UpdateSupplierRequest{
			Name:  "Tech Supplies Global",
			TaxID: "123.456.789-01",
			Address: dto.SupplierAddress{
				Street:  "Innovation Ave",
				Number:  "501", // changed
				City:    "Silicon Valley",
				State:   "CA",
				Country: "USA",
				ZipCode: "94025",
			},
			IsActive: &active,
		}

		respUpdate := authedRequest(t, "PUT", fmt.Sprintf("%s/api/v1/suppliers/%s", baseURL, created.ID), managerToken, updateReq)
		assert.Equal(t, http.StatusOK, respUpdate.StatusCode)

		var updated models.Supplier
		json.NewDecoder(respUpdate.Body).Decode(&updated)
		assert.Equal(t, "Tech Supplies Global", updated.Name)
		assert.Equal(t, "tech-supplies-global", updated.Slug)
		assert.Equal(t, "501", updated.Address.Number)
	})

	t.Run("Operator can list and filter suppliers", func(t *testing.T) {
		resp := authedRequest(t, "GET", baseURL+"/api/v1/suppliers?name=Tech", operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var list dto.PaginatedResponse[models.Supplier]
		json.NewDecoder(resp.Body).Decode(&list)
		assert.GreaterOrEqual(t, len(list.Data), 1)
		assert.Contains(t, list.Data[0].Name, "Tech")
	})

	t.Run("Operator can find by ID", func(t *testing.T) {
		// Get first supplier
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/suppliers?limit=1", operatorToken, nil)
		var list dto.PaginatedResponse[models.Supplier]
		json.NewDecoder(listResp.Body).Decode(&list)
		id := list.Data[0].ID

		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/suppliers/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var found models.Supplier
		json.NewDecoder(resp.Body).Decode(&found)
		assert.Equal(t, id, found.ID)
	})

	t.Run("RBAC: Operator cannot delete suppliers", func(t *testing.T) {
		listResp := authedRequest(t, "GET", baseURL+"/api/v1/suppliers?limit=1", operatorToken, nil)
		var list dto.PaginatedResponse[models.Supplier]
		json.NewDecoder(listResp.Body).Decode(&list)
		id := list.Data[0].ID

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/suppliers/%s", baseURL, id), operatorToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Admin can delete supplier", func(t *testing.T) {
		// Create as manager
		cReq := dto.CreateSupplierRequest{
			Name:  "Delete Me",
			TaxID: "999.999.999-99",
			Address: dto.SupplierAddress{
				Street: "Main Street", Number: "0", City: "Y", State: "ZZ", Country: "W", ZipCode: "00000",
			},
		}
		cResp := authedRequest(t, "POST", baseURL+"/api/v1/suppliers", managerToken, cReq)
		var toDel models.Supplier
		json.NewDecoder(cResp.Body).Decode(&toDel)
		assert.NotEmpty(t, toDel.ID, "ID should not be empty")

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/suppliers/%s", baseURL, toDel.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify deletion
		vResp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/suppliers/%s", baseURL, toDel.ID), adminToken, nil)
		assert.Equal(t, http.StatusNotFound, vResp.StatusCode)
	})
}
