/*
Package middleware provides HTTP middleware components for cross-cutting concerns.

This package implements middleware for:
- CORS (Cross-Origin Resource Sharing) handling for browser security
- Request/response modification and validation
- Security headers and policy enforcement

Middleware components wrap HTTP handlers to provide functionality that applies
to multiple routes without duplicating code in each handler.
*/
package middleware

import "net/http"

// corsMiddleware implements CORS (Cross-Origin Resource Sharing) functionality
// This middleware allows the forum API to be accessed from different origins (domains)
// which is essential for modern web applications with separate frontend/backend
type corsMiddleware struct {
	handler http.Handler // The next handler in the chain to call after CORS processing
}

// ServeHTTP implements the http.Handler interface for CORS middleware
// This method processes every HTTP request to add CORS headers and handle preflight requests
//
// CORS Process:
// 1. Add necessary CORS headers to allow cross-origin requests
// 2. Handle OPTIONS preflight requests from browsers
// 3. Pass non-preflight requests to the next handler
//
// Parameters:
//   - w: HTTP response writer for sending responses to client
//   - r: HTTP request containing client request data
func (c *corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to enable cross-origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")                                      // Allow requests from any origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")      // Allowed HTTP methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")          // Allowed request headers
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")                       // Headers exposed to JavaScript
	w.Header().Set("Access-Control-Max-Age", "86400")                                       // Cache preflight response for 24 hours

	// Handle preflight OPTIONS requests
	// Browsers send OPTIONS requests before actual requests to check CORS permissions
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK) // Send 200 OK response for preflight
		return                       // Don't continue to next handler for preflight requests
	}

	// For non-preflight requests, continue to the next handler in the chain
	c.handler.ServeHTTP(w, r)
}

// NewCorsMiddleware creates a new CORS middleware wrapper
// This function wraps an existing HTTP handler with CORS functionality
//
// Parameters:
//   - handler: The HTTP handler to wrap with CORS middleware
//
// Returns:
//   - http.Handler: CORS-enabled handler that processes requests with proper headers
//
// Usage:
//   router := http.NewServeMux()
//   corsEnabledRouter := NewCorsMiddleware(router)
//   http.ListenAndServe(":8080", corsEnabledRouter)
func NewCorsMiddleware(handler http.Handler) http.Handler {
	return &corsMiddleware{handler: handler}
}
