package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	// From internal/middleware/rbac_test.go, we need to go up 3 levels:
	// internal/middleware/rbac_test.go -> internal/middleware -> internal -> root
	return filepath.Join(filepath.Dir(b), "../..")
}

func setupCasbinRouter(enforcer *casbin.Enforcer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.TraceIDMiddleware())

	// Mock the effect of AuthMiddleware by setting the role manually in a middleware
	r.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Test-Role")
		if role != "" {
			c.Set("userRole", role)
		}
		c.Next()
	})

	r.Use(middleware.CasbinMiddleware(enforcer))

	r.GET("/api/v1/roles", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.PATCH("/api/v1/users/change-password", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.DELETE("/api/v1/users/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}

func TestCasbinMiddleware(t *testing.T) {
	rootPath := getRootPath()
	modelPath := filepath.Join(rootPath, "model.conf")
	policyPath := filepath.Join(rootPath, "policy.csv")

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	require.NoError(t, err)

	router := setupCasbinRouter(enforcer)

	tests := []struct {
		name           string
		role           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Admin can access everything",
			role:           "admin",
			method:         http.MethodDelete,
			path:           "/api/v1/users/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Manager can access roles",
			role:           "manager",
			method:         http.MethodGet,
			path:           "/api/v1/roles",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Manager can change password",
			role:           "manager",
			method:         http.MethodPatch,
			path:           "/api/v1/users/change-password",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Manager cannot delete users",
			role:           "manager",
			method:         http.MethodDelete,
			path:           "/api/v1/users/123",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Unauthorized role has no access",
			role:           "guest",
			method:         http.MethodGet,
			path:           "/api/v1/roles",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing role is forbidden",
			role:           "",
			method:         http.MethodGet,
			path:           "/api/v1/roles",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			if tt.role != "" {
				req.Header.Set("X-Test-Role", tt.role)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if resp.Code == http.StatusForbidden {
				var pd problemDetails
				err := json.Unmarshal(resp.Body.Bytes(), &pd)
				assert.NoError(t, err)
				assert.Equal(t, "Forbidden", pd.Title)
			}
		})
	}
}

func BenchmarkCasbinMiddleware_Authorized(b *testing.B) {
	rootPath := getRootPath()
	modelPath := filepath.Join(rootPath, "model.conf")
	policyPath := filepath.Join(rootPath, "policy.csv")

	enforcer, _ := casbin.NewEnforcer(modelPath, policyPath)
	router := setupCasbinRouter(enforcer)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	req.Header.Set("X-Test-Role", "admin")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

func BenchmarkCasbinMiddleware_Forbidden(b *testing.B) {
	rootPath := getRootPath()
	modelPath := filepath.Join(rootPath, "model.conf")
	policyPath := filepath.Join(rootPath, "policy.csv")

	enforcer, _ := casbin.NewEnforcer(modelPath, policyPath)
	router := setupCasbinRouter(enforcer)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/123", nil)
	req.Header.Set("X-Test-Role", "manager")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}
