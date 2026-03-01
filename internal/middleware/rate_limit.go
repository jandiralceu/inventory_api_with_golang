package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiter returns a gin middleware for rate limiting.
// scope: unique identifier for the limiter (e.g. "global", "auth")
// format: "number-period" (e.g. "5-S", "10-M", "1000-H")
func RateLimiter(cache pkg.CacheManager, scope string, formattedRate string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(formattedRate)
	if err != nil {
		panic(fmt.Sprintf("invalid rate limit format: %s", formattedRate))
	}

	// Use our existing Redis client with a scope-specific prefix
	store, err := sredis.NewStoreWithOptions(cache.GetClient(), limiter.StoreOptions{
		Prefix: fmt.Sprintf("rate_limit:%s:", scope),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create redis rate limit store: %v", err))
	}

	// Create a custom key getter to use UserID if available, fallback to IP
	keyGetter := func(c *gin.Context) string {
		// Try to get user ID from context (set by AuthMiddleware)
		if userID, exists := c.Get(UserIDKey); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		// Fallback to client IP
		return c.ClientIP()
	}

	// Create the limiter
	instance := limiter.New(store, rate)

	// Return the gin middleware with custom options
	return mgin.NewMiddleware(instance, mgin.WithKeyGetter(keyGetter), mgin.WithLimitReachedHandler(limitReachedHandler))
}

func limitReachedHandler(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"type":   "https://api.example.com/errors/too-many-requests",
		"title":  "More Requests Than Allowed",
		"status": http.StatusTooManyRequests,
		"detail": "You have exceeded your request limit. Please try again later.",
	})
	c.Abort()
}
