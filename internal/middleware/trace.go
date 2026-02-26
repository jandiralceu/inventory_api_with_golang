package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceIDKey is the key used to store the request's Trace ID in the context.
const TraceIDKey = "trace_id"

// TraceIDMiddleware ensures every request has a unique Trace ID for observability.
// It checks for an existing "X-Trace-ID" header, or generates a new one if missing.
// The Trace ID is added both to the Gin context and to the response headers.
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has a Trace ID (from a proxy/gateway)
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Set in Gin context so handlers can use it
		c.Set(TraceIDKey, traceID)

		// Set in response header so the client knows their trace ID
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
