package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers/validation"
)

// LoginPage handles GET requests to /login.
func (cs *ClientServer) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "login", domain.LoginFormErrors{})
}

// LoginPost handles POST requests to /login.
func (cs *ClientServer) LoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := domain.LoginFormErrors{
		Identifier: identifier,
	}

	// FRONTEND VALIDATION - Quick feedback for user
	data.IdentifierError = validateIdentifier(identifier)
	data.PasswordError = validation.ValidatePassword(password)

	// If frontend validation fails, re-render login page with errors
	if data.IdentifierError != "" || data.PasswordError != "" {
		renderTemplate(w, "login", data)
		return
	}

	// BACKEND LOGIN - Send validated data to backend
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	backendResp, backendErr := cs.loginWithBackend(ctx, identifier, password)
	if backendErr != nil {
		// Backend validation/login failed
		data.IdentifierError = ""
		data.PasswordError = ""

		// Parse backend error and show appropriate message
		errorMsg := backendErr.Error()

		// Show error in password field
		data.PasswordError = errorMsg

		renderTemplate(w, "login", data)
		return
	}

	// SUCCESS - User logged in, redirect to homepage
	log.Printf("User logged in successfully: %s (ID: %s)", backendResp.Username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================================================
// API REQUEST METHODS (Methods on ClientServer)
// ============================================================================

// loginWithBackend determines if identifier is email or username and calls appropriate backend endpoint.
func (cs *ClientServer) loginWithBackend(ctx context.Context, identifier string, password string) (*domain.BackendLoginResponse, error) {
	var backendURL string
	var req domain.BackendLoginRequest

	// Determine if identifier is email or username
	if validation.IsValidEmailFormat(identifier) {
		backendURL = backendLoginEmailURL
		req = domain.BackendLoginRequest{
			Email:    identifier,
			Password: password,
		}
	} else {
		backendURL = backendLoginUsernameURL
		req = domain.BackendLoginRequest{
			Username: identifier,
			Password: password,
		}
	}

	return cs.sendLoginRequest(ctx, backendURL, req)
}

// sendLoginRequest sends the login request to the backend API.
func (cs *ClientServer) sendLoginRequest(ctx context.Context, backendURL string, req domain.BackendLoginRequest) (*domain.BackendLoginResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	// Create HTTP request to backend with context
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, backendURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request using the client's HTTP client (with cookie jar)
	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Login request failed: " + err.Error())
	}
	defer resp.Body.Close()

	// Handle different response statuses
	if resp.StatusCode != http.StatusOK {
		// Backend returned an error
		var errResp domain.BackendErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, backendError("Login failed. Please try again.")
		}

		if errResp.Message != "" {
			return nil, backendError(errResp.Message)
		}
		return nil, backendError("Login failed. Please try again.")
	}

	// Success response
	target := domain.BackendLoginResponse{}
	err = DecodeBackendResponse(resp, &target)
	if err != nil {
		return nil, backendError("Failed to decode response: " + err.Error())
	}

	return &target, nil
}

// ============================================================================
// VALIDATION HELPERS
// ============================================================================

// validateIdentifier checks if the identifier (username or email) is valid.
func validateIdentifier(identifier string) string {
	if identifier == "" {
		return "Username or Email is required."
	}

	// If it looks like an email, validate as email
	if validation.IsValidEmailFormat(identifier) {
		return validation.ValidateEmail(identifier)
	}

	// Otherwise validate as username
	if len(identifier) < 3 {
		return "Username must be at least 3 characters long."
	}

	if len(identifier) > 50 {
		return "Username must not exceed 50 characters."
	}

	return ""
}
