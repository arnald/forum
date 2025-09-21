/*
Package health provides HTTP handlers for health check endpoints.

This package contains HTTP handlers that provide server health monitoring
capabilities for the forum application. Health checks are essential for:

- Load balancer health monitoring
- Container orchestration readiness probes
- Service monitoring and alerting
- Deployment verification and rollback decisions

The health check endpoint provides basic server status information
including uptime, timestamp, and service availability status.
*/
package health

import (
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app/health/queries"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

// Handler encapsulates dependencies for health check HTTP operations
// It provides health monitoring endpoints for external systems
type Handler struct {
	Logger logger.Logger // Logger for health check request logging and error reporting
}

// NewHandler creates a new health check handler with required dependencies
// This factory function initializes the handler with logging capability
//
// Parameters:
//   - logger: Logger instance for request logging and error reporting
//
// Returns:
//   - *Handler: Configured health handler ready to process health check requests
//
// Usage:
//   healthHandler := NewHandler(logger)
//   http.HandleFunc("/health", healthHandler.HealthCheck)
func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		Logger: logger, // Store logger for health check operations
	}
}

// HealthCheck handles HTTP health check requests
// This endpoint provides basic server status information for monitoring systems
//
// HTTP Method: GET only
// Endpoint: /api/v1/health
// Authentication: None required (public endpoint)
//
// Response Format:
//   {
//     "status": "up",
//     "timestamp": "2023-12-07T10:30:00Z"
//   }
//
// Parameters:
//   - w: HTTP response writer for sending response to client
//   - r: HTTP request containing client request information
//
// Response Codes:
//   - 200 OK: Server is healthy and operational
//   - 405 Method Not Allowed: Non-GET request received
//
// Usage by external systems:
//   - Load balancers check this endpoint before routing traffic
//   - Monitoring systems poll this endpoint for availability
//   - Container orchestrators use this for readiness probes
//   - CI/CD pipelines verify deployment health
func (h Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - health checks should only use GET
	// This prevents accidental or malicious non-GET requests
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Create health check response with current status and timestamp
	// Status "up" indicates the server is accepting requests
	// RFC3339 timestamp provides precise timing for monitoring
	response := queries.HealthResponse{
		Status:    queries.StatusUp,                    // Server operational status
		Timestamp: time.Now().Format(time.RFC3339),     // Current server time in standard format
	}

	// Send successful health check response
	// HTTP 200 OK indicates the server is healthy and operational
	// JSON format provides structured data for automated monitoring
	helpers.RespondWithJSON(
		w,              // Response writer
		http.StatusOK,  // 200 OK status code
		nil,            // No additional headers needed
		response,       // Health status data
	)
}
