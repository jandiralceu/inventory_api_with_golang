package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
)

const (
	// UserIDKey is the key used to store the authenticated user ID in the Gin context.
	UserIDKey = "userID"
	// UserRoleKey is the key used to store the authenticated user role in the Gin context.
	UserRoleKey = "userRole"
)

// ProblemDetails represents a standard error response following RFC 7807.
// This is a local copy to avoid a circular dependency with the handlers package.
type ProblemDetails struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	TraceID string `json:"trace_id,omitempty"`
}

// respondWithUnauthorized writes a consistent ProblemDetails response for 401 errors.
func respondWithUnauthorized(c *gin.Context, detail string) {
	traceID := c.GetString(TraceIDKey)

	c.AbortWithStatusJSON(http.StatusUnauthorized, ProblemDetails{
		Type:    "https://api.example.com/errors/unauthorized",
		Title:   "Unauthorized",
		Detail:  fmt.Sprintf("%s: %s", apperrors.ErrUnauthorized.Error(), detail),
		TraceID: traceID,
	})
}

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// from the Authorization header and injects the user ID into the context.
// It expects the header in the format: "Bearer <token>".
func AuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondWithUnauthorized(c, "authorization header is required")
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondWithUnauthorized(c, "invalid authorization header format")
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			respondWithUnauthorized(c, "invalid or expired token")
			return
		}

		// Store the user ID and role in the context for downstream handlers.
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)

		// Inject into Request context so services can access it without Gin dependency
		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserID extracts the authenticated user ID from the Gin context.
func GetUserID(c *gin.Context) uuid.UUID {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userID
}

// GetUserRole extracts the authenticated user role from the Gin context.
func GetUserRole(c *gin.Context) string {
	val, exists := c.Get(UserRoleKey)
	if !exists {
		return ""
	}

	role, ok := val.(string)
	if !ok {
		return ""
	}

	return role
}
