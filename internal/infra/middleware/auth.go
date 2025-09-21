/*
Package middleware provides HTTP middleware functions for the forum application.

This package contains middleware for:
- Authentication verification (required and optional)
- Session management and validation
- Request context manipulation
- User identification and authorization

Middleware functions wrap HTTP handlers to provide cross-cutting concerns
like authentication, logging, and request processing.
*/
package middleware

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

// AuthMiddleware provides authentication verification for protected routes
// This middleware requires users to be logged in with a valid session
// It extracts session information from cookies and validates it against the database
//
// Parameters:
//   - sessionManager: Interface for managing user sessions
//   - logger: Logger for recording authentication events and errors
//
// Returns:
//   - Middleware function that wraps HTTP handlers with authentication
//
// Behavior:
//   - Checks for session_token cookie in the request
//   - Validates the session token against the database
//   - Adds user_id to request context for use by handlers
//   - Returns 401 Unauthorized if authentication fails
func AuthMiddleware(sessionManager user.SessionManager, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session token from HTTP cookie
			sessionCookie, err := r.Cookie("session_token")
			if err != nil {
				// No session cookie found - user is not logged in
				helpers.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Validate session token against database
			session, err := sessionManager.GetSession(sessionCookie.Value)
			if err != nil {
				// Session validation failed - log error and deny access
				logger.PrintError(err, nil)
				helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid session")
				return
			}

			if session == nil {
				// Session not found in database - possibly expired or invalid
				helpers.RespondWithError(w, http.StatusUnauthorized, "Session not found")
				return
			}

			// Authentication successful - add user ID to request context
			// This allows handlers to access the authenticated user's ID
			ctx := context.WithValue(r.Context(), "user_id", session.UserID)
			r = r.WithContext(ctx)

			// Continue to the next handler with authenticated context
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuthMiddleware provides optional authentication for public routes
// This middleware checks for authentication but doesn't require it
// It's useful for routes that can show different content to authenticated vs anonymous users
//
// Parameters:
//   - sessionManager: Interface for managing user sessions
//   - logger: Logger for recording authentication events
//
// Returns:
//   - Middleware function that optionally adds authentication context
//
// Behavior:
//   - Attempts to extract and validate session information
//   - Adds user_id to context if authentication is successful
//   - Continues processing regardless of authentication status
//   - Does NOT return errors for missing or invalid authentication
func OptionalAuthMiddleware(sessionManager user.SessionManager, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get session token from cookie (may not exist)
			sessionCookie, err := r.Cookie("session_token")
			if err == nil {
				// Cookie exists - try to validate session
				session, err := sessionManager.GetSession(sessionCookie.Value)
				if err == nil && session != nil {
					// Valid session found - add user ID to context
					ctx := context.WithValue(r.Context(), "user_id", session.UserID)
					r = r.WithContext(ctx)
				}
				// If session validation fails, we silently continue without authentication
			}
			// If no cookie exists, we continue without authentication

			// Always continue to next handler (authentication is optional)
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext extracts the authenticated user ID from request context
// This function is used by handlers to get the current user's ID after authentication
//
// Parameters:
//   - r: HTTP request containing context with user information
//
// Returns:
//   - string: User ID if user is authenticated
//   - bool: true if user ID was found in context, false otherwise
//
// Usage:
//   userID, isAuthenticated := GetUserIDFromContext(r)
//   if !isAuthenticated {
//       // Handle unauthenticated user
//   }
func GetUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value("user_id").(string)
	return userID, ok
}