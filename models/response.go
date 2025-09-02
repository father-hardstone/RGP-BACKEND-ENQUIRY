package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response represents a standardized success response structure
// Used across all API endpoints to maintain consistency
type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Timestamp  string      `json:"timestamp"`
}

// ErrorResponse represents a standardized error response structure
// Used across all API endpoints to maintain consistency
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// CreateSuccessResponse creates a new success response with the given parameters
// statusCode: HTTP status code (e.g., 200, 201)
// message: Human-readable success message
// data: Optional data payload to include in the response
func CreateSuccessResponse(statusCode int, message string, data interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Status:     "success",
		Message:    message,
		Data:       data,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// CreateErrorResponse creates a new error response with the given parameters
// statusCode: HTTP status code (e.g., 400, 500)
// message: Human-readable error message
// error: Technical error details (optional)
func CreateErrorResponse(statusCode int, message string, error string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Status:     "error",
		Message:    message,
		Error:      error,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// SendJSONResponse sends a JSON response to the HTTP client
// w: HTTP response writer
// statusCode: HTTP status code to send
// data: Data to encode as JSON
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
