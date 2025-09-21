/*
Package helpers provides utility functions for HTTP response handling.

This package contains helper functions for standardizing HTTP responses
across the forum application. It provides:

- JSON response formatting with consistent structure
- Error response handling with proper status codes
- Pagination metadata support for list endpoints
- Response wrapper for uniform API responses
- Proper content-type headers and error handling

The response helpers ensure consistent API responses throughout
the application while handling edge cases and error conditions.
*/
package helpers

import (
	"encoding/json"
	"net/http"
)

// ResponseWrapper provides a standard structure for API responses
// This wrapper ensures consistent response format across all endpoints
// and supports optional pagination metadata for list operations
type ResponseWrapper struct {
	Info *Info `json:"info,omitzero"` // Optional pagination/metadata information
	Data any   `json:"data,omitzero"` // Actual response data (can be any type)
}

// Info contains metadata for paginated responses and list operations
// This structure provides clients with pagination information needed
// to navigate through large datasets efficiently
type Info struct {
	TotalRecords int `json:"totalRecords"`           // Total number of records available
	CurrentPage  int `json:"currentPage,omitzero"`   // Current page number (1-based)
	PageSize     int `json:"pageSize,omitzero"`      // Number of records per page
	TotalPages   int `json:"totalPages,omitzero"`    // Total number of pages available
	NextPage     int `json:"nextPage,omitzero"`      // Next page number (0 if no next page)
	PrevPage     int `json:"prevPage,omitzero"`      // Previous page number (0 if no previous page)
}

// RespondWithError sends a standardized error response to the client
// This function provides a consistent way to return error messages
// across all HTTP handlers in the application
//
// Parameters:
//   - w: HTTP response writer for sending response to client
//   - code: HTTP status code indicating the type of error
//   - msg: Human-readable error message explaining what went wrong
//
// Response Format:
//   {
//     "error": "Error message here"
//   }
//
// Usage:
//   RespondWithError(w, 400, "Invalid request parameters")
//   RespondWithError(w, 404, "User not found")
//   RespondWithError(w, 500, "Internal server error occurred")
//
// The function automatically sets proper content-type headers
// and handles JSON marshaling errors gracefully
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, nil, map[string]string{"error": msg})
}

// RespondWithJSON sends a JSON response with optional pagination metadata
// This function provides the primary way to send structured data responses
// throughout the forum application with consistent formatting
//
// Parameters:
//   - w: HTTP response writer for sending response to client
//   - code: HTTP status code (200, 201, 400, 404, 500, etc.)
//   - info: Optional pagination metadata (nil for single records or errors)
//   - payload: Data to be sent in the response (any JSON-serializable type)
//
// Response Formats:
//
// For successful responses with pagination:
//   {
//     "info": {
//       "totalRecords": 150,
//       "currentPage": 2,
//       "pageSize": 20,
//       "totalPages": 8,
//       "nextPage": 3,
//       "prevPage": 1
//     },
//     "data": [...response data...]
//   }
//
// For successful responses without pagination:
//   {
//     "data": {...response data...}
//   }
//
// For error responses (status >= 400):
//   {...raw payload without wrapper...}
//
// Error Handling:
//   - JSON marshaling errors result in 500 Internal Server Error
//   - Write errors are handled gracefully with fallback error response
//   - Empty responses default to empty JSON object "{}"
//   - Proper content-type headers are always set
func RespondWithJSON(w http.ResponseWriter, code int, info *Info, payload any) {
	var jsonData []byte
	var err error

	// Handle different response types based on status code and metadata
	switch {
	// Case 1: Error responses (>=400) without metadata - send raw payload
	// This provides direct error messages without wrapper structure
	case code >= http.StatusBadRequest && info == nil:
		jsonData, err = json.Marshal(payload)
		if err != nil {
			// Fallback to built-in error handler if JSON marshaling fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	// Case 2: All other responses - use standard wrapper format
	// This includes successful responses and any responses with metadata
	default:
		response := ResponseWrapper{
			Info: info,    // Pagination metadata (may be nil)
			Data: payload, // Actual response data
		}
		jsonData, err = json.Marshal(response)
		if err != nil {
			// Fallback to built-in error handler if JSON marshaling fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Ensure we always have valid JSON to send
	// This prevents sending empty responses that could confuse clients
	if jsonData == nil {
		jsonData = []byte("{}")
	}

	// Set proper content-type header for JSON responses
	// This ensures clients can properly parse the response
	w.Header().Set("Content-Type", "application/json")

	// Write HTTP status code to response
	w.WriteHeader(code)

	// Write JSON data to response body
	_, err = w.Write(jsonData)
	if err != nil {
		// Handle write errors (client disconnection, network issues, etc.)
		// Use built-in error handler as last resort
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
