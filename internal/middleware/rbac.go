package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware returns a Gin middleware that enforces RBAC authorization.
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get user role from context (set by AuthMiddleware)
		role := GetUserRole(c)
		if role == "" {
			respondWithForbidden(c, "user role not found in context")
			return
		}

		// 2. Get path and method from request
		path := c.Request.URL.Path
		method := c.Request.Method

		// 3. Ask Casbin if this role has permission for this path and method
		allowed, err := enforcer.Enforce(role, path, method)
		if err != nil {
			respondWithForbidden(c, "error during authorization check")
			return
		}

		if !allowed {
			respondWithForbidden(c, fmt.Sprintf("no permission for %s on %s", method, path))
			return
		}

		c.Next()
	}
}

// respondWithForbidden writes a consistent ProblemDetails response for 403 errors.
func respondWithForbidden(c *gin.Context, detail string) {
	traceID := c.GetString(TraceIDKey)

	c.AbortWithStatusJSON(http.StatusForbidden, ProblemDetails{
		Type:    "https://api.example.com/errors/forbidden",
		Title:   "Forbidden",
		Detail:  detail,
		TraceID: traceID,
	})
}
