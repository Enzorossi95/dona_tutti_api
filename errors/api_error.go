package errors

// APIError represents a standard error response for the API
type APIError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request format"`
}
