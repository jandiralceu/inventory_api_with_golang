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
	"gorm.io/gorm"
)

func TestUserManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	var adminRole, operatorRole models.Role
	require.NoError(t, db.Where("name = ?", "admin").First(&adminRole).Error)
	require.NoError(t, db.Where("name = ?", "operator").First(&operatorRole).Error)

	baseURL := ts.URL
	adminEmail := "superadmin@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "Super Admin", adminEmail, password, adminRole.ID.String())
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	t.Run("Admin can list and search users", func(t *testing.T) {
		// Create a second user to search for
		signUpUser(t, baseURL, "Searchable User", "search@example.com", "Pass123!", adminRole.ID.String())

		resp := authedRequest(t, "GET", baseURL+"/api/v1/users?name=Searchable", adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var listResp dto.UserListResponse
		json.NewDecoder(resp.Body).Decode(&listResp)
		assert.NotEmpty(t, listResp.Data)
		assert.Equal(t, "Searchable User", listResp.Data[0].Name)
	})

	t.Run("User can change their own password", func(t *testing.T) {
		userEmail := "changepass@example.com"
		signUpUser(t, baseURL, "Pass Changer", userEmail, "old-pass-123", adminRole.ID.String())
		token, _ := signInUser(t, baseURL, userEmail, "old-pass-123")

		req := dto.ChangePasswordRequest{
			OldPassword: "old-pass-123",
			NewPassword: "new-pass-456",
		}

		resp := authedRequest(t, "PATCH", baseURL+"/api/v1/users/change-password", token, req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify login with new password
		newToken, _ := signInUser(t, baseURL, userEmail, "new-pass-456")
		assert.NotEmpty(t, newToken)
	})

	t.Run("Admin can delete users", func(t *testing.T) {
		userEmail := "tobedeleted@example.com"
		signUpUser(t, baseURL, "To Be Deleted", userEmail, "Pass123!", adminRole.ID.String())

		var user models.User
		db.Where("email = ?", userEmail).First(&user)

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, user.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify user is gone
		var checkUser models.User
		err := db.Where("id = ?", user.ID).First(&checkUser).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Admin can change user roles", func(t *testing.T) {
		userEmail := "rolechanger@example.com"
		signUpUser(t, baseURL, "Role Changer", userEmail, "Pass123!", adminRole.ID.String())

		var user models.User
		db.Where("email = ?", userEmail).First(&user)

		req := dto.ChangeRoleRequest{
			UserID: user.ID,
			RoleID: operatorRole.ID,
		}

		resp := authedRequest(t, "PATCH", baseURL+"/api/v1/users/change-role", adminToken, req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		var updatedUser models.User
		db.Where("id = ?", user.ID).First(&updatedUser)
		assert.Equal(t, operatorRole.ID, updatedUser.RoleID)
	})

	t.Run("RBAC: Operator cannot delete users", func(t *testing.T) {
		opEmail := "op_delete_bad@example.com"
		signUpUser(t, baseURL, "Op User", opEmail, "Pass123!", operatorRole.ID.String())
		opToken, _ := signInUser(t, baseURL, opEmail, "Pass123!")

		var user models.User
		db.Where("email = ?", adminEmail).First(&user)

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, user.ID), opToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("RBAC: Operator cannot change roles", func(t *testing.T) {
		opEmail := "op_role_bad@example.com"
		signUpUser(t, baseURL, "Op User", opEmail, "Pass123!", operatorRole.ID.String())
		opToken, _ := signInUser(t, baseURL, opEmail, "Pass123!")

		req := dto.ChangeRoleRequest{
			UserID: uuid.New(), // Dummy UUID works here as it should be blocked by middleware
			RoleID: adminRole.ID,
		}

		resp := authedRequest(t, "PATCH", baseURL+"/api/v1/users/change-role", opToken, req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
