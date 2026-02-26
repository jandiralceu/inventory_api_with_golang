package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should generate a new trace ID if none provided", func(t *testing.T) {
		r := gin.New()
		r.Use(TraceIDMiddleware())
		r.GET("/test", func(c *gin.Context) {
			traceID, exists := c.Get(TraceIDKey)
			assert.True(t, exists)
			assert.NotEmpty(t, traceID)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))
	})

	t.Run("should reuse existing trace ID from header", func(t *testing.T) {
		existingID := "existing-trace-id"
		r := gin.New()
		r.Use(TraceIDMiddleware())
		r.GET("/test", func(c *gin.Context) {
			traceID, _ := c.Get(TraceIDKey)
			assert.Equal(t, existingID, traceID)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Trace-ID", existingID)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, existingID, w.Header().Get("X-Trace-ID"))
	})
}
