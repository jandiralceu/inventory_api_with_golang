package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test helpers ---

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performRequest(router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// =====================
// RespondWithError Tests
// =====================

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should include TraceID from context", func(t *testing.T) {
		r := gin.New()
		traceID := "test-trace-id"
		r.GET("/error", func(c *gin.Context) {
			c.Set("trace_id", traceID)
			RespondWithError(c, apperrors.ErrNotFound)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, traceID, resp.TraceID)
	})

	t.Run("should return 404 for ErrNotFound", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrNotFound)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Resource Not Found", resp.Title)
		assert.Equal(t, "https://api.example.com/errors/not-found", resp.Type)
	})

	t.Run("should return 409 for ErrConflict", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrConflict)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Conflict", resp.Title)
	})

	t.Run("should return 400 for ErrInvalidInput", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrInvalidInput)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Bad Request", resp.Title)
	})

	t.Run("should return 400 for ErrInvalidID", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrInvalidID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 401 for ErrUnauthorized", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrUnauthorized)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Unauthorized", resp.Title)
	})

	t.Run("should return 403 for ErrForbidden", func(t *testing.T) {
		r := gin.New()
		r.GET("/error", func(c *gin.Context) {
			RespondWithError(c, apperrors.ErrForbidden)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Forbidden", resp.Title)
	})

	t.Run("should include detailed validation errors", func(t *testing.T) {
		r := gin.New()
		r.GET("/validation", func(c *gin.Context) {
			vErr := apperrors.NewValidationErrors("Validation failed")
			vErr.Add("email", "required")
			vErr.Add("password", "min=8")
			RespondWithError(c, vErr)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/validation", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "Validation Failed", resp.Title)
		assert.Len(t, resp.InvalidParams, 2)
		assert.Equal(t, "email", resp.InvalidParams[0].Field)
		assert.Equal(t, "required", resp.InvalidParams[0].Reason)
		assert.Equal(t, "password", resp.InvalidParams[1].Field)
		assert.Equal(t, "min=8", resp.InvalidParams[1].Reason)
	})

	t.Run("should mask internal errors detail", func(t *testing.T) {
		r := gin.New()
		r.GET("/internal", func(c *gin.Context) {
			RespondWithError(c, errors.New("secret db error"))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/internal", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp ProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "An unexpected error occurred. Please try again later.", resp.Detail)
		assert.Equal(t, "Internal Server Error", resp.Title)
	})
}
