//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorMappingIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	baseURL := ts.URL

	// 1. Setup Admin User
	adminEmail := "admin_error_test@example.com"
	adminPass := "Admin@123"

	var adminRoleID string
	err := db.Raw("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID).Error
	require.NoError(t, err)

	signUpUser(t, baseURL, "Admin User", adminEmail, adminPass, adminRoleID)
	token, _ := signInUser(t, baseURL, adminEmail, adminPass)

	// 2. Create Category and Supplier for Product
	cat := models.Category{ID: uuid.New(), Name: "Integration Cat", Slug: "integration-cat"}
	require.NoError(t, db.Create(&cat).Error)

	sup := models.Supplier{ID: uuid.New(), Name: "Integration Sup", Slug: "integration-sup", Email: "sup@test.com", TaxID: "123456"}
	require.NoError(t, db.Create(&sup).Error)

	// 3. Create Warehouse and Product
	wh := models.Warehouse{ID: uuid.New(), Name: "Integration WH", Slug: "integration-wh", Code: "WH-ERR-01"}
	require.NoError(t, db.Create(&wh).Error)

	prod := models.Product{ID: uuid.New(), Name: "Integration Prod", Slug: "integration-prod", SKU: "PROD-ERR-01", Price: 100.0, CategoryID: &cat.ID, SupplierID: &sup.ID}
	require.NoError(t, db.Create(&prod).Error)

	whID := wh.ID
	prodID := prod.ID

	t.Run("Conflict - Unique Violation", func(t *testing.T) {
		// Try to create inventory for same prod/wh twice
		body := map[string]any{
			"productId":   prodID,
			"warehouseId": whID,
			"quantity":    10,
		}

		// First create (Success)
		resp1 := authedRequest(t, "POST", baseURL+"/api/v1/inventory", token, body)
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Second create (Conflict)
		resp2 := authedRequest(t, "POST", baseURL+"/api/v1/inventory", token, body)
		assert.Equal(t, http.StatusConflict, resp2.StatusCode)
	})

	t.Run("InvalidInput - Check Violation", func(t *testing.T) {
		// Negative quantity
		body := map[string]any{
			"productId":   prodID,
			"warehouseId": uuid.New(), // New WH to avoid unique violation if we could skip FK
			"quantity":    -10,
		}

		// But wait, the WarehouseID MUST exist for FK. Let's create another WH.
		whID2 := uuid.New()
		db.Exec("INSERT INTO warehouses (id, name, slug, code) VALUES (?, ?, ?, ?)", whID2, "Integration WH 2", "integration-wh-2", "WH-ERR-02")
		body["warehouseId"] = whID2

		resp := authedRequest(t, "POST", baseURL+"/api/v1/inventory", token, body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidInput - Not Null Violation", func(t *testing.T) {
		// Missing ProductID (The handler might catch this before DB, but let's see)
		// To trigger DB not null violation, we'd need to bypass DTO validation if any.
		// The CreateInventoryDTO likely has 'required' tags.

		body := map[string]any{
			"warehouseId": whID,
			"quantity":    10,
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/inventory", token, body)
		// If DTO validation catches it, it's 400. If DB catches it, it's also 400 now!
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
