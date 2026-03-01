//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/database"
	"github.com/jandiralceu/inventory_api_with_golang/internal/handlers"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
	"github.com/jandiralceu/inventory_api_with_golang/internal/routes"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
	"github.com/jandiralceu/inventory_api_with_golang/tests/integration/testhelpers"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// generateIntegrationRSAKeys creates a temporary RSA key pair for JWT signing during tests.
func generateIntegrationRSAKeys() (privateKeyPEM string, publicKeyPEM string, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}
	privBlock := &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}
	privPEM := pem.EncodeToMemory(privBlock)

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	pubPEM := pem.EncodeToMemory(pubBlock)

	return string(privPEM), string(pubPEM), nil
}

// setupApp boots the entire application stack using real Postgres + Redis containers.
func setupApp(t *testing.T) (*httptest.Server, *gorm.DB, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Resolve migrations path relative to the tests/integration dir.
	migrationsDir := "../../migrations"

	// Start Postgres container with migrations.
	pgContainer, err := testhelpers.SetupPostgres(ctx, migrationsDir)
	require.NoError(t, err, "Failed to start Postgres container")

	// Start Redis container.
	redisContainer, err := testhelpers.SetupRedis(ctx)
	require.NoError(t, err, "Failed to start Redis container")

	// Build config from container info.
	cfg := &config.Config{
		DBHost:       pgContainer.Host,
		DBPort:       pgContainer.Port,
		DBUser:       pgContainer.User,
		DBPassword:   pgContainer.Password,
		DBName:       pgContainer.DBName,
		AppPort:      "0",
		Env:          "test",
		AppName:      "inventory-api-test",
		RedisHost:    redisContainer.Host,
		RedisPort:    redisContainer.Port,
		OTLPEndpoint: "localhost:4317", // Won't be used in test but needs to be string
	}

	// Initialize database (run migrations if needed, but testhelpers already does it)
	db, err := database.Init(ctx, cfg)
	require.NoError(t, err, "Failed to initialize database")

	// Seed roles for RBAC testing
	err = database.SeedRoles(ctx, db)
	require.NoError(t, err, "Failed to seed roles")

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Cache.
	cacheManager := pkg.NewRedisCacheManager(cfg)

	// Generate RSA keys for JWT.
	privPEM, pubPEM, err := generateIntegrationRSAKeys()
	require.NoError(t, err, "Failed to generate RSA keys")

	jwtManager, err := pkg.NewJWTManager(privPEM, pubPEM)
	require.NoError(t, err, "Failed to initialize JWT manager")

	hasher := pkg.NewArgon2PasswordHasher()

	// Initializing repos, services and handlers
	roleRepository := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepository, cacheManager)
	roleHandler := handlers.NewRoleHandler(roleService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, hasher)
	authHandler := handlers.NewAuthHandler(userService, jwtManager, cacheManager, hasher)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Casbin Enforcer for RBAC.
	// We need to point to the files in the project root.
	enforcer, err := casbin.NewEnforcer("../../model.conf", "../../policy.csv")
	require.NoError(t, err, "Failed to initialize Casbin enforcer")

	routeConfig := &routes.RouteConfig{
		AuthHandler: authHandler,
		RoleHandler: roleHandler,
		UserHandler: userHandler,
	}

	// Suppress Gin debug output during tests.
	gin.SetMode(gin.TestMode)
	r := routes.Setup(routeConfig, cfg, jwtManager, enforcer)

	// Health check (same as main.go).
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ts := httptest.NewServer(r)

	cleanup := func() {
		ts.Close()
		cacheManager.Close()
		sqlDB.Close()
		// Termination might take time, use a separate context
		termCtx := context.Background()
		pgContainer.Terminate(termCtx)
		redisContainer.Terminate(termCtx)
	}

	return ts, db, cleanup
}

// signUpUser registers a user using the API.
func signUpUser(t *testing.T, baseURL, name, email, password, roleID string) {
	t.Helper()

	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s","roleId":"%s"}`, name, email, password, roleID)
	resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode, "Register should return 204")
}

// signInUser authenticates a user and returns markers.
func signInUser(t *testing.T, baseURL, email, password string) (string, string) {
	t.Helper()

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	resp, err := http.Post(baseURL+"/api/v1/auth/signin", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)
		t.Fatalf("SignIn failed with status %d: %v", resp.StatusCode, errResp)
	}

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	accessToken, ok1 := result["accessToken"].(string)
	refreshToken, ok2 := result["refreshToken"].(string)
	require.True(t, ok1 && ok2, "Response must contain tokens")

	return accessToken, refreshToken
}

// authedRequest sends an HTTP request with an Authorization Bearer token.
func authedRequest(t *testing.T, method, url, token string, bodyData any) *http.Response {
	t.Helper()

	var reqBody *strings.Reader
	if bodyData != nil {
		jsonBytes, _ := json.Marshal(bodyData)
		reqBody = strings.NewReader(string(jsonBytes))
	} else {
		reqBody = strings.NewReader("")
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// seedInitialData seeds the roles into the test database.
func seedInitialData(t *testing.T, baseURL string) {
	// In the real app we use a seeder script.
	// For integration tests, we can use the repository directly if we had access to it,
	// but setupApp hides it. However, we can call the Role creation endpoint if we had an admin.
	// But wait, the seeder script uses GORM directly.
	// For simplicity in integration tests, we'll assume the DB is migration-ready.
	// If we want the roles, we should seed them.
}
