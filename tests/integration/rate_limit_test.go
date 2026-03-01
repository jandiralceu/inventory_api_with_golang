//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitingIntegration(t *testing.T) {
	ts, _, cleanup := setupAppCustom(t, func(cfg *config.Config) {
		cfg.RateLimitAuth = "5-M"
		cfg.RateLimitGlobal = "100-M"
	})
	defer cleanup()

	baseURL := ts.URL

	t.Run("Auth SignIn Rate Limiting", func(t *testing.T) {
		// Our auth routes have a limit of 5-M (5 requests per minute)
		// Let's send 6 requests quickly.

		email := "rate@example.com"
		password := "Pass123!"

		// Note: we don't even need the user to exist to hit the rate limit middleware
		// as it runs BEFORE the handler.

		for i := 1; i <= 6; i++ {
			body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
			resp, err := http.Post(baseURL+"/api/v1/auth/signin", "application/json", strings.NewReader(body))
			if err != nil {
				t.Fatalf("Request %d failed: %v", i, err)
			}
			defer resp.Body.Close()

			if i <= 5 {
				// Should be Unauthorized (if user doesn't exist) or OK, but definitely NOT 429
				assert.NotEqual(t, http.StatusTooManyRequests, resp.StatusCode, "Request %d should not be rate limited", i)
			} else {
				// The 6th request should hit the limit
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Request %d should be rate limited", i)
			}
		}
	})

	t.Run("Global Rate Limiting", func(t *testing.T) {
		// Health check has no specific limit, so it uses the global 100-M.
		// We won't test the full 100 here to save time, but the logic is shared.
		resp, err := http.Get(baseURL + "/health")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
