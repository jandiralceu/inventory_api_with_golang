//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductLifecycle(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)

	adminEmail := "prodadmin@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, ts.URL, "Prod Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, ts.URL, adminEmail, password)

	// Since product creation requires a valid Category and Supplier (due to foreign keys),
	// we should create those first to have valid UUIDs to link.

	// A. Create a Category
	catReq := dto.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Devices and gadgets",
	}
	catReqBody, _ := json.Marshal(catReq)
	reqCat, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/categories", bytes.NewReader(catReqBody))
	reqCat.Header.Set("Content-Type", "application/json")
	reqCat.Header.Set("Authorization", "Bearer "+adminToken)

	respCat, err := http.DefaultClient.Do(reqCat)
	require.NoError(t, err)
	defer respCat.Body.Close()
	require.Equal(t, http.StatusCreated, respCat.StatusCode)

	var catData map[string]any
	json.NewDecoder(respCat.Body).Decode(&catData)
	catIDStr := catData["id"].(string)
	catID, _ := uuid.Parse(catIDStr)

	// B. Create a Supplier
	supReq := dto.CreateSupplierRequest{
		Name:  "Tech Supplier Inc",
		Email: "sales@techsupplier.com",
		Phone: "1234567890",
		TaxID: "1234567890",
		Address: dto.SupplierAddress{
			Street:  "Supplier St",
			City:    "Tech City",
			State:   "TC",
			Country: "USA",
			ZipCode: "12345",
			Number:  "10",
		},
	}
	supReqBody, _ := json.Marshal(supReq)
	reqSup, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/suppliers", bytes.NewReader(supReqBody))
	reqSup.Header.Set("Content-Type", "application/json")
	reqSup.Header.Set("Authorization", "Bearer "+adminToken)

	respSup, err := http.DefaultClient.Do(reqSup)
	require.NoError(t, err)
	defer respSup.Body.Close()
	require.Equal(t, http.StatusCreated, respSup.StatusCode)

	var supData map[string]any
	json.NewDecoder(respSup.Body).Decode(&supData)
	supIDStr := supData["id"].(string)
	supID, _ := uuid.Parse(supIDStr)

	// 2. Create the Product
	createReq := dto.CreateProductRequest{
		Name:            "Smartphone X",
		SKU:             "SMART-X-001",
		Description:     "A brand new smartphone",
		Price:           999.99,
		CostPrice:       ptrFloat(600.00),
		CategoryID:      &catID,
		SupplierID:      &supID,
		ReorderLevel:    10,
		ReorderQuantity: 50,
		WeightKg:        ptrFloat(0.5),
	}

	createReqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/products", bytes.NewReader(createReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)
		t.Fatalf("Failed to create product: %v", errResp)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdProduct map[string]any
	err = json.NewDecoder(resp.Body).Decode(&createdProduct)
	require.NoError(t, err)
	productID := createdProduct["id"].(string)
	assert.NotEmpty(t, productID)
	assert.Equal(t, "Smartphone X", createdProduct["name"])
	assert.Equal(t, "SMART-X-001", createdProduct["sku"])
	assert.Equal(t, 999.99, createdProduct["price"])

	// 3. Find By ID
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var fetchedProduct map[string]any
	json.NewDecoder(resp.Body).Decode(&fetchedProduct)
	assert.Equal(t, "Smartphone X", fetchedProduct["name"])

	// 4. List and Filter Products
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/api/v1/products?sku=SMART-X-001", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var listResp map[string]any
	json.NewDecoder(resp.Body).Decode(&listResp)

	total := int(listResp["total"].(float64))
	assert.GreaterOrEqual(t, total, 1)

	// 5. Update the Product
	updateReq := dto.UpdateProductRequest{
		Name:  "Smartphone X Updated",
		Price: ptrFloat(1050.50),
	}
	updateReqBody, _ := json.Marshal(updateReq)
	req, _ = http.NewRequest(http.MethodPut, ts.URL+"/api/v1/products/"+productID, bytes.NewReader(updateReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var updatedProduct map[string]any
	json.NewDecoder(resp.Body).Decode(&updatedProduct)
	assert.Equal(t, "Smartphone X Updated", updatedProduct["name"])
	assert.Equal(t, 1050.50, updatedProduct["price"])

	// 6. Delete the Product
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 7. Verify Deletion
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func ptrFloat(f float64) *float64 { return &f }
