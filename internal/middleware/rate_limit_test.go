package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCacheManager for rate limiting tests
type MockCacheManager struct {
	mock.Mock
}

func (m *MockCacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheManager) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheManager) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheManager) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}

func (m *MockCacheManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCacheManager) GetClient() *redis.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.Client)
}

func TestRateLimiter(t *testing.T) {
	// Start miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	// Create a real redis client connecting to miniredis
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() { _ = client.Close() }()

	// Create our mock CacheManager to return the real redis client
	mockCache := new(MockCacheManager)
	mockCache.On("GetClient").Return(client)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simulate getting auth context - optional
	router.Use(func(c *gin.Context) {
		// Just to test the fallback key getter mechanism.
		// Use a specific test IP so we can track rate limit correctly.
		c.Request.RemoteAddr = "192.168.1.1:1234"
		c.Next()
	})

	// Add the rate limiter: limit to 2 requests per minute for tests
	router.Use(middleware.RateLimiter(mockCache, "test", "2-M"))

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// First request: should pass
	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusOK, resp1.Code)

	// Second request: should pass (reached the limit 2 of 2)
	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusOK, resp2.Code)

	// Third request: should fail (Too Many Requests)
	req3, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	assert.Equal(t, http.StatusTooManyRequests, resp3.Code)

	var pd problemDetails
	err = json.Unmarshal(resp3.Body.Bytes(), &pd)
	assert.NoError(t, err)
	assert.Equal(t, "More Requests Than Allowed", pd.Title)
}

func TestRateLimiter_WithUserID(t *testing.T) {
	// Start miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer func() { _ = client.Close() }()

	mockCache := new(MockCacheManager)
	mockCache.On("GetClient").Return(client)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	testUserID := uuid.New()

	// Simulate auth context with UserID
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, testUserID)
		c.Next()
	})

	// Allow 1 request per minute
	router.Use(middleware.RateLimiter(mockCache, "test_user", "1-M"))

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// First request: should pass
	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusOK, resp1.Code)

	// Second request: should fail
	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusTooManyRequests, resp2.Code)
}
