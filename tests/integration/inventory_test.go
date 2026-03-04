//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryHistoryLifecycle(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)

	adminEmail := "invadmin@example.com"
	password := "SecurePass123!"

	// 1. Setup Admin and Products/Warehouses
	signUpUser(t, ts.URL, "Inv Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, ts.URL, adminEmail, password)

	// Create Category
	catID := createCategory(t, ts.URL, adminToken)
	// Create Supplier
	supID := createSupplier(t, ts.URL, adminToken)
	// Create Product
	prodID := createProduct(t, ts.URL, adminToken, catID, supID)
	// Create Warehouse
	whID := createWarehouse(t, ts.URL, adminToken)

	// 2. Create Initial Inventory
	invReq := dto.CreateInventoryRequest{
		ProductID:    prodID,
		WarehouseID:  whID,
		Quantity:     100,
		LocationCode: "A-01-01",
	}
	resp := authedRequest(t, "POST", ts.URL+"/api/v1/inventories", adminToken, invReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var inv models.Inventory
	require.NoError(t, decodeResponse(resp, &inv))
	invID := inv.ID

	// 3. Add Stock
	stockReq := dto.StockOperationRequest{Quantity: 50}
	resp = authedRequest(t, "POST", fmt.Sprintf("%s/api/v1/inventories/%s/add", ts.URL, invID), adminToken, stockReq)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 4. Remove Stock
	stockReq = dto.StockOperationRequest{Quantity: 30}
	resp = authedRequest(t, "POST", fmt.Sprintf("%s/api/v1/inventories/%s/remove", ts.URL, invID), adminToken, stockReq)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 5. Verify Transaction History
	// This is the "new case" the user is asking about.
	resp = authedRequest(t, "GET", ts.URL+"/api/v1/inventories/transactions?inventory_id="+invID.String(), adminToken, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var txHistory dto.TransactionListResponse
	require.NoError(t, decodeResponse(resp, &txHistory))

	// We expect at least 2 transactions (Add and Remove)
	// Depending on implementation, the initial creation might also be a transaction.
	assert.GreaterOrEqual(t, txHistory.Total, int64(2), "Should have at least 2 transactions logged")

	// Check the latest transaction (Remove Stock)
	if len(txHistory.Data) > 0 {
		latest := txHistory.Data[0] // Assuming desc order by default
		assert.Equal(t, "OUT", latest.TransactionType)
		assert.Equal(t, -30, latest.QuantityChange)
		assert.Equal(t, 120, latest.QuantityBalance) // 100 + 50 - 30 = 120
	}
}

// Helper functions to keep the test clean

func createCategory(t *testing.T, baseURL, token string) uuid.UUID {
	req := dto.CreateCategoryRequest{Name: "Cat " + uuid.New().String()}
	resp := authedRequest(t, "POST", baseURL+"/api/v1/categories", token, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data models.Category
	require.NoError(t, decodeResponse(resp, &data))
	return data.ID
}

func createSupplier(t *testing.T, baseURL, token string) uuid.UUID {
	req := dto.CreateSupplierRequest{
		Name:  "Sup " + uuid.New().String(),
		Email: uuid.New().String() + "@test.com",
		TaxID: "123456789",
		Address: dto.SupplierAddress{
			Street:  "Main St",
			Number:  "1",
			City:    "City",
			State:   "NY",
			Country: "USA",
			ZipCode: "12345",
		},
	}
	resp := authedRequest(t, "POST", baseURL+"/api/v1/suppliers", token, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data models.Supplier
	require.NoError(t, decodeResponse(resp, &data))
	return data.ID
}

func createWarehouse(t *testing.T, baseURL, token string) uuid.UUID {
	req := dto.CreateWarehouseRequest{
		Name: "WH " + uuid.New().String(),
		Code: "WH" + uuid.New().String()[:5],
		Address: dto.WarehouseAddress{
			Street:  "Warehouse St",
			Number:  "200",
			City:    "Logistics City",
			State:   "AL",
			Country: "USA",
			ZipCode: "54321",
		},
	}
	resp := authedRequest(t, "POST", baseURL+"/api/v1/warehouses", token, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data models.Warehouse
	require.NoError(t, decodeResponse(resp, &data))
	return data.ID
}

func createProduct(t *testing.T, baseURL, token string, catID, supID uuid.UUID) uuid.UUID {
	req := dto.CreateProductRequest{
		Name:       "Prod " + uuid.New().String(),
		SKU:        "SKU" + uuid.New().String()[:8],
		CategoryID: &catID,
		SupplierID: &supID,
		Price:      10.0,
	}
	resp := authedRequest(t, "POST", baseURL+"/api/v1/products", token, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var data models.Product
	require.NoError(t, decodeResponse(resp, &data))
	return data.ID
}

func decodeResponse(resp *http.Response, target any) error {
	defer func() { _ = resp.Body.Close() }()
	return json.NewDecoder(resp.Body).Decode(target)
}
