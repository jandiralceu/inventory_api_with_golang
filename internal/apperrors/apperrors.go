// Package apperrors defines sentinel errors and validation types used across the
// application to classify failures into standard HTTP problem categories.
//
// Sentinel errors map directly to HTTP status codes in the response layer:
//
//	ErrNotFound     → 404    ErrConflict     → 409
//	ErrInvalidInput → 400    ErrUnauthorized → 401
//	ErrForbidden    → 403    ErrInternal     → 500
//
// Wrap them with [fmt.Errorf] to add context:
//
//	fmt.Errorf("%w: email already in use", apperrors.ErrConflict)
package apperrors

import "errors"

// Sentinel errors used to classify application-level failures.
var (
	ErrNotFound      = errors.New("resource not found")
	ErrConflict      = errors.New("resource conflict")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInternal      = errors.New("internal server error")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidID     = errors.New("invalid identifier")
	ErrForbidden     = errors.New("forbidden")
)

// InvalidParam represents a single field-level validation failure,
// following the RFC 9457 "invalid-params" extension.
type InvalidParam struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// ValidationErrors aggregates one or more [InvalidParam] into an error
// that the response layer serialises as a Problem Details "invalid-params" member.
type ValidationErrors struct {
	Message string         `json:"message"`
	Params  []InvalidParam `json:"params"`
}

func (v *ValidationErrors) Error() string {
	if v.Message != "" {
		return v.Message
	}
	return "validation failed"
}

// NewValidationErrors creates a [ValidationErrors] with the given summary message.
func NewValidationErrors(message string) *ValidationErrors {
	return &ValidationErrors{
		Message: message,
		Params:  []InvalidParam{},
	}
}

// Add appends a field-level validation failure.
func (v *ValidationErrors) Add(field, reason string) {
	v.Params = append(v.Params, InvalidParam{
		Field:  field,
		Reason: reason,
	})
}
