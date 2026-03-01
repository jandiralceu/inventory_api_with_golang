package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to generate RSA test keys
func generateMemRSAKeys(t *testing.T) (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(privPEM), string(pubPEM)
}

func setupTestRouter(jwtManager *pkg.JWTManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.TraceIDMiddleware())
	r.Use(middleware.AuthMiddleware(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		userID := middleware.GetUserID(c)
		role := middleware.GetUserRole(c)
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID.String(),
			"role":    role,
		})
	})
	return r
}

type problemDetails struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	TraceID string `json:"trace_id"`
}

func TestAuthMiddleware(t *testing.T) {
	privPEM, pubPEM := generateMemRSAKeys(t)
	jwtManager, err := pkg.NewJWTManager(privPEM, pubPEM)
	require.NoError(t, err)

	router := setupTestRouter(jwtManager)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var pd problemDetails
		err := json.Unmarshal(resp.Body.Bytes(), &pd)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", pd.Title)
		assert.Contains(t, pd.Detail, "authorization header is required")
		assert.Equal(t, "https://api.example.com/errors/unauthorized", pd.Type)
		assert.NotEmpty(t, pd.TraceID, "should include trace ID from middleware")
	})

	t.Run("Invalid Header Format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var pd problemDetails
		err := json.Unmarshal(resp.Body.Bytes(), &pd)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", pd.Title)
		assert.Contains(t, pd.Detail, "invalid authorization header format")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var pd problemDetails
		err := json.Unmarshal(resp.Body.Bytes(), &pd)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", pd.Title)
		assert.Contains(t, pd.Detail, "invalid or expired token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		userID := uuid.New()
		role := "admin"
		token, err := jwtManager.GenerateToken(userID, role, 15*time.Minute)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), userID.String())
		assert.Contains(t, resp.Body.String(), role)
	})
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid UserID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedID := uuid.New()
		c.Set(middleware.UserIDKey, expectedID)

		actualID := middleware.GetUserID(c)
		assert.Equal(t, expectedID, actualID)
	})

	t.Run("Missing UserID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		actualID := middleware.GetUserID(c)
		assert.Equal(t, uuid.Nil, actualID)
	})
}

func TestGetUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedRole := "admin"
		c.Set(middleware.UserRoleKey, expectedRole)

		actualRole := middleware.GetUserRole(c)
		assert.Equal(t, expectedRole, actualRole)
	})

	t.Run("Missing Role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		actualRole := middleware.GetUserRole(c)
		assert.Equal(t, "", actualRole)
	})
}
