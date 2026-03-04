//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFlowAndRBAC(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	// 1. Get role IDs for testing
	var adminRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL

	t.Run("Full Auth Flow: Register -> SignIn -> Refresh -> SignOut", func(t *testing.T) {
		email := "tester@example.com"
		password := "Password123!"

		// Register
		signUpUser(t, baseURL, "Tester", email, password, adminRole.ID.String())

		// SignIn
		accessToken, refreshToken := signInUser(t, baseURL, email, password)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Refresh Token
		refreshBody := map[string]string{"refreshToken": refreshToken}
		resp := authedRequest(t, "POST", baseURL+"/api/v1/auth/refresh", "", refreshBody)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshResult map[string]string
		require.NoError(t, decodeResponse(resp, &refreshResult))
		newAccessToken := refreshResult["accessToken"]
		newRefreshToken := refreshResult["refreshToken"]
		assert.NotEmpty(t, newAccessToken)
		assert.NotEqual(t, refreshToken, newRefreshToken)

		// SignOut
		signOutBody := map[string]string{"refreshToken": newRefreshToken}
		resp = authedRequest(t, "POST", baseURL+"/api/v1/auth/signout", "", signOutBody)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("RBAC: Admin can list roles, Operator cannot", func(t *testing.T) {
		adminEmail := "admin_test@example.com"
		opEmail := "op_test@example.com"
		password := "Pass123!"

		// Create Admin
		signUpUser(t, baseURL, "Admin User", adminEmail, password, adminRole.ID.String())
		adminToken, _ := signInUser(t, baseURL, adminEmail, password)

		// Create Operator
		signUpUser(t, baseURL, "Op User", opEmail, password, operatorRole.ID.String())
		opToken, _ := signInUser(t, baseURL, opEmail, password)

		// Admin attempts to find role by ID (protected route)
		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, adminRole.ID), adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin should be able to find role by ID")

		// Operator attempts to create a role (protected route)
		createReq := map[string]string{"name": "hacker", "description": "Should fail"}
		resp = authedRequest(t, "POST", baseURL+"/api/v1/roles", opToken, createReq)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Operator should NOT be able to create roles")
	})

	t.Run("RBAC: Accessing non-existent route returns 404 (with auth)", func(t *testing.T) {
		adminEmail := "admin_404@example.com"
		password := "Pass123!"
		signUpUser(t, baseURL, "Admin 404", adminEmail, password, adminRole.ID.String())
		adminToken, _ := signInUser(t, baseURL, adminEmail, password)

		resp := authedRequest(t, "GET", baseURL+"/api/v1/non-existent", adminToken, nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Auth: Request without token to protected route returns 401", func(t *testing.T) {
		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, adminRole.ID), "", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
