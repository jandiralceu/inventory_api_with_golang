package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
)

// ProblemDetails implements the RFC 7807 specification for HTTP API error responses.
// It provides machine-readable details to describe error conditions.
type ProblemDetails struct {
	Type          string                   `json:"type"`
	Title         string                   `json:"title"`
	Detail        string                   `json:"detail"`
	TraceID       string                   `json:"trace_id,omitempty"`
	InvalidParams []apperrors.InvalidParam `json:"invalid_params,omitempty"`
}

// RespondWithError translates application-level errors into appropriate HTTP responses
// using the [ProblemDetails] structure. It maps sentinel errors from [apperrors]
// to their corresponding HTTP status codes.
func RespondWithError(c *gin.Context, err error) {
	var statusCode int
	var title string
	var errorType string
	var detail string
	var invalidParams []apperrors.InvalidParam

	// Extract Trace ID from context (provided by telemetry/request ID middleware)
	traceID := c.GetString("trace_id")

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		statusCode = http.StatusNotFound
		title = "Resource Not Found"
		detail = err.Error()
		errorType = "https://api.example.com/errors/not-found"
	case errors.Is(err, apperrors.ErrConflict):
		statusCode = http.StatusConflict
		title = "Conflict"
		detail = err.Error()
		errorType = "https://api.example.com/errors/conflict"
	case errors.Is(err, apperrors.ErrInvalidInput), errors.Is(err, apperrors.ErrInvalidID):
		statusCode = http.StatusBadRequest
		title = "Bad Request"
		detail = err.Error()
		errorType = "https://api.example.com/errors/bad-request"
	case errors.Is(err, apperrors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		title = "Unauthorized"
		detail = err.Error()
		errorType = "https://api.example.com/errors/unauthorized"
	case errors.Is(err, apperrors.ErrForbidden):
		statusCode = http.StatusForbidden
		title = "Forbidden"
		detail = err.Error()
		errorType = "https://api.example.com/errors/forbidden"
	default:
		// Check for structured validation errors
		var vErr *apperrors.ValidationErrors
		if errors.As(err, &vErr) {
			statusCode = http.StatusBadRequest
			title = "Validation Failed"
			detail = vErr.Error()
			invalidParams = vErr.Params
			errorType = "https://api.example.com/errors/validation-failed"
		} else {
			statusCode = http.StatusInternalServerError
			title = "Internal Server Error"
			// Prevent leaking technical details like SQL syntax errors or connection issues
			detail = "An unexpected error occurred. Please try again later."
			errorType = "https://api.example.com/errors/internal-server-error"

			// Log the actual internal error for debugging via Gin's error collection
			_ = c.Error(err)
		}
	}

	c.JSON(statusCode, ProblemDetails{
		Type:          errorType,
		Title:         title,
		Detail:        detail,
		TraceID:       traceID,
		InvalidParams: invalidParams,
	})
}

// ParseValidationError converts Gin's [validator.ValidationErrors] into the application's
// structured [apperrors.ValidationErrors] type, extracting field names and failure tags.
func ParseValidationError(err error) error {
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		valErr := apperrors.NewValidationErrors("One or more fields failed validation")
		for _, f := range vErrs {
			reason := f.Tag()
			if f.Param() != "" {
				reason = reason + "=" + f.Param()
			}
			valErr.Add(f.Field(), reason)
		}
		return valErr
	}
	return err
}
