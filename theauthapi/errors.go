// theauthapi/errors.go
package theauthapi

import (
	"fmt"
	"net/http"
)

// APIError represents a generic API error
type APIError struct {
	Code    int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// IsNotFound checks if the error is a 404 Not Found error
func (e *APIError) IsNotFound() bool {
	return e.Code == http.StatusNotFound
}

// IsUnauthorized checks if the error is a 401 Unauthorized error
func (e *APIError) IsUnauthorized() bool {
	return e.Code == http.StatusUnauthorized
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		Code:    statusCode,
		Message: message,
		Err:     err,
	}
}