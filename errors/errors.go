package errors

import (
	"context"
	"encoding/json"
	"net/http"
)

// ValidationError represents a validation error that should return 400 Bad Request
type ValidationError struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

// NewFieldValidationError creates a new validation error for a specific field
func NewFieldValidationError(field, message string) ValidationError {
	return ValidationError{Message: message, Field: field}
}

// NotFoundError represents a resource not found error that should return 404 Not Found
type NotFoundError struct {
	Message string `json:"message"`
}

func (e NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// HTTPErrorEncoder encodes errors to HTTP responses with appropriate status codes
func HTTPErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var errorResponse ErrorResponse

	switch e := err.(type) {
	case ValidationError:
		statusCode = http.StatusBadRequest
		errorResponse = ErrorResponse{
			Error: ErrorDetail{
				Message: e.Message,
				Field:   e.Field,
			},
		}
	case NotFoundError:
		statusCode = http.StatusNotFound
		errorResponse = ErrorResponse{
			Error: ErrorDetail{
				Message: e.Message,
			},
		}
	default:
		// Default to internal server error for unknown errors
		statusCode = http.StatusInternalServerError
		errorResponse = ErrorResponse{
			Error: ErrorDetail{
				Message: "Internal server error",
			},
		}
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
