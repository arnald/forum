// Package handlers - oauth.go implements OAuth authentication handlers
// Handles login/register via Google, Facebook, and GitHub
package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/internal/oauth"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// OAuthConfig holds OAuth configuration for all providers (Google, GitHub, Facebook)
// This global variable is initialized once at startup and used by all OAuth handlers
var OAuthConfig *oauth.Config

// InitOAuth initializes OAuth configuration from environment variables
// Should be called once during application startup before handling any OAuth requests
func InitOAuth() {
	OAuthConfig = oauth.NewConfig()
}

// GoogleLogin redirects to Google OAuth login
// Generates a state token for CSRF protection and redirects user to Google's login page
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if Google OAuth is configured
	if OAuthConfig.Google.ClientID == "" {
		http.Error(w, "Google OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate random state token for CSRF protection
	state := generateStateToken()

	// Store state token in cookie to validate on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute), // Short expiration for security
		HttpOnly: true,                              // Prevent JavaScript access
		Secure:   false,                             // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,              // CSRF protection
	})

	// Generate Google OAuth URL with state token
	url := OAuthConfig.Google.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Redirect user to Google login page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles Google OAuth callback
// Validates state token, exchanges code for token, retrieves user info, and creates/logs in user
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state token to prevent CSRF attacks
	if err := validateStateToken(r); err != nil {
		h.BadRequest(w, r, "Invalid state token")
		return
	}

	// Extract authorization code from query parameters
	code := r.URL.Query().Get("code")

	// Exchange authorization code for access token
	token, err := OAuthConfig.Google.Exchange(r.Context(), code)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Retrieve user information from Google using access token
	userInfo, err := oauth.GetGoogleUserInfo(r.Context(), token, OAuthConfig.Google)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Process OAuth user (create account if new, or log in existing user)
	h.handleOAuthUser(w, r, userInfo)
}

// GitHubLogin redirects to GitHub OAuth login
// Generates a state token for CSRF protection and redirects user to GitHub's authorization page
func (h *Handler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Check if GitHub OAuth is configured
	if OAuthConfig.GitHub.ClientID == "" {
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate random state token for CSRF protection
	state := generateStateToken()

	// Store state token in cookie to validate on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute), // Short expiration for security
		HttpOnly: true,                              // Prevent JavaScript access
		Secure:   false,                             // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,              // CSRF protection
	})

	// Generate GitHub OAuth URL with state token
	url := OAuthConfig.GitHub.AuthCodeURL(state)

	// Redirect user to GitHub authorization page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitHubCallback handles GitHub OAuth callback
// Validates state token, exchanges code for token, retrieves user info, and creates/logs in user
func (h *Handler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state token to prevent CSRF attacks
	if err := validateStateToken(r); err != nil {
		h.BadRequest(w, r, "Invalid state token")
		return
	}

	// Extract authorization code from query parameters
	code := r.URL.Query().Get("code")

	// Exchange authorization code for access token
	token, err := OAuthConfig.GitHub.Exchange(r.Context(), code)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Retrieve user information from GitHub using access token
	userInfo, err := oauth.GetGitHubUserInfo(r.Context(), token, OAuthConfig.GitHub)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Process OAuth user (create account if new, or log in existing user)
	h.handleOAuthUser(w, r, userInfo)
}

// FacebookLogin redirects to Facebook OAuth login
// Generates a state token for CSRF protection and redirects user to Facebook's login dialog
func (h *Handler) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	// Check if Facebook OAuth is configured
	if OAuthConfig.Facebook.ClientID == "" {
		http.Error(w, "Facebook OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate random state token for CSRF protection
	state := generateStateToken()

	// Store state token in cookie to validate on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute), // Short expiration for security
		HttpOnly: true,                              // Prevent JavaScript access
		Secure:   false,                             // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,              // CSRF protection
	})

	// Generate Facebook OAuth URL with state token
	url := OAuthConfig.Facebook.AuthCodeURL(state)

	// Redirect user to Facebook login dialog
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// FacebookCallback handles Facebook OAuth callback
// Validates state token, exchanges code for token, retrieves user info, and creates/logs in user
func (h *Handler) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state token to prevent CSRF attacks
	if err := validateStateToken(r); err != nil {
		h.BadRequest(w, r, "Invalid state token")
		return
	}

	// Extract authorization code from query parameters
	code := r.URL.Query().Get("code")

	// Exchange authorization code for access token
	token, err := OAuthConfig.Facebook.Exchange(r.Context(), code)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Retrieve user information from Facebook using access token
	userInfo, err := oauth.GetFacebookUserInfo(r.Context(), token, OAuthConfig.Facebook)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Process OAuth user (create account if new, or log in existing user)
	h.handleOAuthUser(w, r, userInfo)
}

// handleOAuthUser processes OAuth user login/registration
// This function handles both new OAuth users (creates account) and existing OAuth users (logs them in)
func (h *Handler) handleOAuthUser(w http.ResponseWriter, r *http.Request, userInfo *oauth.UserInfo) {
	// Try to find existing user with this OAuth provider and ID
	user, err := h.auth.FindUserByProvider(userInfo.Provider, userInfo.ID)

	if err != nil {
		// User doesn't exist with this provider, create new account
		username := userInfo.Username
		email := userInfo.Email

		// Check if username or email already exists (from local or other OAuth accounts)
		if h.auth.UserExists(email, username) {
			// Add provider suffix to make username unique (e.g., "john_google")
			username = fmt.Sprintf("%s_%s", username, userInfo.Provider)
		}

		// Create new OAuth user account in database
		user, err = h.auth.CreateOAuthUser(email, username, userInfo.Provider, userInfo.ID)
		if err != nil {
			h.InternalServerError(w, r, err)
			return
		}
	}

	// Create new session for the user (both new and existing users)
	session, err := h.auth.CreateSession(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Set session cookie to maintain login state
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,       // 24 hour expiration
		HttpOnly: true,                     // Prevent JavaScript access for security
		Secure:   false,                    // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,     // CSRF protection
		Path:     "/",                      // Cookie available for entire site
	})

	// Redirect user to home page after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// generateStateToken creates a random state token for OAuth
// The state token is used for CSRF protection during the OAuth flow
// Returns a base64-encoded random 32-byte string
func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b) // Generate cryptographically secure random bytes
	return base64.URLEncoding.EncodeToString(b)
}

// validateStateToken validates OAuth state token for CSRF protection
// Compares the state parameter from OAuth callback with the state cookie we set earlier
// Returns error if state cookie is missing or doesn't match the callback parameter
func validateStateToken(r *http.Request) error {
	// Get state token from cookie (set during initial OAuth redirect)
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		return fmt.Errorf("no state cookie")
	}

	// Get state parameter from OAuth callback URL
	stateParam := r.URL.Query().Get("state")

	// Verify they match to prevent CSRF attacks
	if stateParam != stateCookie.Value {
		return fmt.Errorf("state mismatch")
	}

	return nil
}
