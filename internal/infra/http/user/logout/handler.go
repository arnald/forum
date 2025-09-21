/*
Package userlogout handles user logout functionality for the forum application.

This package provides HTTP handlers for:
- User session termination
- Cookie cleanup
- Session database cleanup
- Logout response handling

The logout process involves both server-side session invalidation
and client-side cookie removal to ensure complete logout.
*/
package userlogout

import (
	"net/http"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

// Handler encapsulates dependencies needed for user logout operations
type Handler struct {
	UserServices   app.Services           // Application services for business logic
	SessionManager user.SessionManager    // Session management interface
	Config         *config.ServerConfig   // Server configuration
	Logger         logger.Logger          // Application logger
}

// NewHandler creates a new logout handler with all required dependencies
//
// Parameters:
//   - config: Server configuration for timeouts and settings
//   - app: Application services for business logic operations
//   - sm: Session manager for handling user sessions
//   - logger: Logger for recording logout events and errors
//
// Returns:
//   - *Handler: Configured logout handler ready to process requests
func NewHandler(config *config.ServerConfig, app app.Services, sm user.SessionManager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,     // Store application services
		SessionManager: sm,      // Store session manager
		Config:         config,  // Store server configuration
		Logger:         logger,  // Store logger instance
	}
}

// UserLogout handles HTTP POST requests for user logout
// This endpoint terminates the user's session and clears authentication cookies
//
// Process:
// 1. Validates HTTP method (must be POST)
// 2. Extracts session token from cookies
// 3. Deletes session from database
// 4. Clears session cookie from browser
// 5. Returns success response
//
// Parameters:
//   - w: HTTP response writer for sending response
//   - r: HTTP request containing session cookie
//
// Expected Request:
//   - Method: POST
//   - Cookie: session_token (containing valid session token)
//
// Response:
//   - 200 OK: Logout successful
//   - 401 Unauthorized: No session found
//   - 405 Method Not Allowed: Invalid HTTP method
func (h Handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - logout must be POST for security
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Extract session token from HTTP cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		// No session cookie found - user may not be logged in
		helpers.RespondWithError(w, http.StatusUnauthorized, "No session found")
		return
	}

	// Delete session from database to invalidate it server-side
	err = h.SessionManager.DeleteSession(sessionCookie.Value)
	if err != nil {
		// Log error but continue - we still want to clear client-side cookie
		h.Logger.PrintError(err, nil)
	}

	// Clear the session cookie from user's browser
	// This ensures the user is logged out client-side even if DB deletion failed
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",              // Empty value clears the cookie
		Path:     "/",             // Apply to entire site
		MaxAge:   -1,              // Negative MaxAge deletes the cookie
		HttpOnly: true,            // Prevent JavaScript access for security
		Secure:   false,           // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // CSRF protection
	})

	// Send successful logout response
	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		map[string]interface{}{
			"message": "Logged out successfully",
		},
	)

	// Log successful logout for monitoring and debugging
	h.Logger.PrintInfo(
		"User logout successfully",
		map[string]string{
			"session_token": sessionCookie.Value,
		},
	)
}