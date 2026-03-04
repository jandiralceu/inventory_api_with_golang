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

func TestRoleManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL
	adminEmail := "roleadmin@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "Role Manager", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	t.Run("Admin can create a new role", func(t *testing.T) {
		req := dto.CreateRoleRequest{
			Name:        "auditor",
			Description: "View only access for auditing",
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/roles", adminToken, req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Role
		require.NoError(t, decodeResponse(resp, &created))
		assert.Equal(t, "auditor", created.Name)
		assert.Equal(t, "View only access for auditing", created.Description)
		assert.NotEqual(t, [16]byte{}, created.ID)
	})

	t.Run("Admin can list roles", func(t *testing.T) {
		resp := authedRequest(t, "GET", baseURL+"/api/v1/roles", adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var roles []models.Role
		require.NoError(t, decodeResponse(resp, &roles))

		// Should have at least the 3 default seeded roles + the auditor created above
		assert.GreaterOrEqual(t, len(roles), 4)

		found := false
		for _, r := range roles {
			if r.Name == "admin" {
				found = true
				break
			}
		}
		assert.True(t, found, "Admin role should be in the list")
	})

	t.Run("Admin can find role by ID", func(t *testing.T) {
		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, operatorRole.ID), adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var found models.Role
		require.NoError(t, decodeResponse(resp, &found))
		assert.Equal(t, operatorRole.ID, found.ID)
		assert.Equal(t, "operator", found.Name)
	})

	t.Run("Admin can delete a role", func(t *testing.T) {
		// First create a role to delete
		req := dto.CreateRoleRequest{
			Name:        "temp_role",
			Description: "To be deleted",
		}
		respCreate := authedRequest(t, "POST", baseURL+"/api/v1/roles", adminToken, req)
		var tempRole models.Role
		require.NoError(t, decodeResponse(respCreate, &tempRole))

		// Delete it
		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, tempRole.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify it's gone
		respCheck := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, tempRole.ID), adminToken, nil)
		assert.Equal(t, http.StatusNotFound, respCheck.StatusCode)
	})

	t.Run("RBAC: Operator cannot create roles", func(t *testing.T) {
		opEmail := "op_role_maker@example.com"
		signUpUser(t, baseURL, "Op User", opEmail, "Pass123!", operatorRole.ID.String())
		opToken, _ := signInUser(t, baseURL, opEmail, "Pass123!")

		req := dto.CreateRoleRequest{
			Name:        "hacker_role",
			Description: "Unauthorized",
		}

		resp := authedRequest(t, "POST", baseURL+"/api/v1/roles", opToken, req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("RBAC: Operator cannot delete roles", func(t *testing.T) {
		opEmail := "op_role_remover@example.com"
		signUpUser(t, baseURL, "Op User", opEmail, "Pass123!", operatorRole.ID.String())
		opToken, _ := signInUser(t, baseURL, opEmail, "Pass123!")

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, adminRole.ID), opToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Unauthorized access is blocked", func(t *testing.T) {
		resp := authedRequest(t, "GET", fmt.Sprintf("%s/api/v1/roles/%s", baseURL, adminRole.ID), "", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
