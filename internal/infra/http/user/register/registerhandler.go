/*
Package userregister provides HTTP handlers for user registration functionality.

This package contains HTTP handlers that manage user account creation
for the forum application. Registration is a critical user journey that includes:

- Input validation and sanitization
- Password security and encryption
- Email format validation and normalization
- Duplicate username/email prevention
- Account creation with proper error handling
- Security logging for audit trails

The registration process ensures data integrity and user security
while providing clear feedback for successful or failed registrations.
*/
package userregister

import (
	"context"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

// RegisterUserReguestModel represents the data required for user registration
// This model defines the structure for incoming registration requests from clients
type RegisterUserReguestModel struct {
	Username string `json:"username"` // Desired username (must be unique across the system)
	Password string `json:"password"` // Plain text password (will be encrypted before storage)
	Email    string `json:"email"`    // Email address (must be unique and valid format)
}

// RegisterUserResponse represents the response sent after successful registration
// This model provides confirmation details and user identification for the client
type RegisterUserResponse struct {
	UserID  string `json:"userdId"` // Unique identifier for the newly created user account
	Message string `json:"message"` // Success message confirming account creation
}

// Handler encapsulates dependencies for user registration HTTP operations
// It coordinates between HTTP layer, application services, and infrastructure
type Handler struct {
	UserServices   app.Services           // Application layer services for business logic
	SessionManager user.SessionManager    // Session management for user authentication
	Config         *config.ServerConfig   // Server configuration including timeouts
	Logger         logger.Logger          // Logger for request tracking and error reporting
}

// NewHandler creates a new user registration handler with required dependencies
// This factory function initializes the handler with all necessary services
//
// Parameters:
//   - config: Server configuration containing timeout and security settings
//   - app: Application services containing user registration business logic
//   - sm: Session manager for handling user authentication sessions
//   - logger: Logger instance for request tracking and error reporting
//
// Returns:
//   - *Handler: Configured registration handler ready to process requests
func NewHandler(config *config.ServerConfig, app app.Services, sm user.SessionManager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,     // Store application services
		SessionManager: sm,      // Store session manager
		Config:         config,  // Store server configuration
		Logger:         logger,  // Store logger instance
	}
}

// UserRegister handles HTTP user registration requests
// This endpoint processes new user account creation with comprehensive validation
//
// HTTP Method: POST only
// Endpoint: /api/v1/register
// Authentication: None required (public endpoint)
// Content-Type: application/json
//
// Request Body:
//   {
//     "username": "john_doe",
//     "password": "securePassword123",
//     "email": "john@example.com"
//   }
//
// Response Format (Success):
//   {
//     "userdId": "uuid-string",
//     "message": "user registered successfully"
//   }
//
// Parameters:
//   - w: HTTP response writer for sending response to client
//   - r: HTTP request containing registration data
//
// Response Codes:
//   - 201 Created: User account successfully created
//   - 400 Bad Request: Invalid input data or validation errors
//   - 405 Method Not Allowed: Non-POST request received
//   - 500 Internal Server Error: Server or database errors
//
// Security Features:
//   - Password encryption using bcrypt
//   - Email normalization (lowercase)
//   - Input validation and sanitization
//   - Request timeout protection
//   - Comprehensive error logging
func (h Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - registration should only use POST
	// This prevents accidental GET requests that might expose sensitive data
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Create timeout context to prevent long-running requests
	// This protects against slow or malicious requests that could affect server performance
	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel() // Ensure context is always cancelled to prevent resource leaks

	// Parse and validate request body into registration model
	var userToRegister RegisterUserReguestModel
	userAny, err := helpers.ParseBodyRequest(r, &userToRegister)
	if err != nil {
		// Log and respond with parsing error
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	defer r.Body.Close() // Ensure request body is always closed

	// Perform comprehensive input validation
	// This validates username format, password strength, and email format
	v := validator.New()
	validator.ValidateUserRegistration(v, userAny)
	if !v.Valid() {
		// Log validation errors and respond with detailed error messages
		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		helpers.RespondWithError(w, http.StatusBadRequest, v.ToStringErrors())
		return
	}

	// Execute user registration through application layer
	// This handles business logic including encryption, uniqueness checks, and persistence
	user, err := h.UserServices.UserServices.Queries.UserRegister.Handle(ctx, queries.UserRegisterRequest{
		Name:     userToRegister.Username,              // Username as provided
		Password: userToRegister.Password,              // Plain text password (will be encrypted)
		Email:    strings.ToLower(userToRegister.Email), // Normalize email to lowercase
	})
	if err != nil {
		// Log registration failure and respond with error
		// Error could be duplicate username/email, database error, etc.
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create success response with user information
	userResponse := RegisterUserResponse{
		UserID:  user.ID,                        // UUID of newly created user
		Message: "user registered successfully", // Confirmation message
	}

	// Send successful registration response
	// HTTP 201 Created indicates resource (user account) was successfully created
	helpers.RespondWithJSON(w, http.StatusCreated, nil, userResponse)

	// Log successful registration for audit trail and monitoring
	// This helps track user growth and identify any registration patterns
	h.Logger.PrintInfo("User registered successfully", map[string]string{
		"userId": user.ID,       // User identifier for tracking
		"email":  user.Email,    // Email for audit purposes
		"name":   user.Username, // Username for audit purposes
	})
}
