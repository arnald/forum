package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
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

	// Get the login type from radio button (username or email)
	loginType := strings.TrimSpace(r.FormValue("loginType"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := domain.LoginFormErrors{
		Password: password,
	}

	// Validate password (same for both login types)
	data.PasswordError = validation.ValidatePassword(password)

	// Route based on login type
	if loginType == "email" {
		cs.handleEmailLogin(w, r, data)
	} else {
		cs.handleUsernameLogin(w, r, data)
	}
}

// handleEmailLogin processes email-based login
func (cs *ClientServer) handleEmailLogin(w http.ResponseWriter, r *http.Request, data domain.LoginFormErrors) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Email = email

	// FRONTEND VALIDATION
	data.EmailError = validation.ValidateEmail(email)
	data.PasswordError = validation.ValidatePassword(password)

	// If frontend validation fails, re-render login page with errors
	if data.EmailError != "" || data.PasswordError != "" {
		renderTemplate(w, "login", data)
		return
	}

	// BACKEND LOGIN
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	backendResp, backendErr := cs.loginWithBackendEmail(ctx, email, password)
	if backendErr != nil {
		// Backend validation/login failed
		data.EmailError = ""
		data.PasswordError = backendErr.Error()
		renderTemplate(w, "login", data)
		return
	}

	// SUCCESS - User logged in, redirect to homepage
	log.Printf("User logged in successfully with email: %s (ID: %s)", backendResp.Username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleUsernameLogin processes username-based login
func (cs *ClientServer) handleUsernameLogin(w http.ResponseWriter, r *http.Request, data domain.LoginFormErrors) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Username = username

	// FRONTEND VALIDATION
	data.UsernameError = validation.ValidateUsername(username)
	data.PasswordError = validation.ValidatePassword(password)

	// If frontend validation fails, re-render login page with errors
	if data.UsernameError != "" || data.PasswordError != "" {
		renderTemplate(w, "login", data)
		return
	}

	// BACKEND LOGIN
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	backendResp, backendErr := cs.loginWithBackendUsername(ctx, username, password)
	if backendErr != nil {
		// Backend validation/login failed
		data.UsernameError = ""
		data.PasswordError = backendErr.Error()
		renderTemplate(w, "login", data)
		return
	}

	// SUCCESS - User logged in, redirect to homepage
	log.Printf("User logged in successfully with username: %s (ID: %s)", backendResp.Username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================================================
// API REQUEST METHODS (Methods on ClientServer)
// ============================================================================

// loginWithBackendEmail sends login request to backend email endpoint
func (cs *ClientServer) loginWithBackendEmail(ctx context.Context, email string, password string) (*domain.BackendLoginResponse, error) {
	req := domain.BackendLoginRequest{
		Email:    email,
		Password: password,
	}
	return cs.sendLoginRequest(ctx, backendLoginEmailURL, req)
}

// loginWithBackendUsername sends login request to backend username endpoint
func (cs *ClientServer) loginWithBackendUsername(ctx context.Context, username string, password string) (*domain.BackendLoginResponse, error) {
	req := domain.BackendLoginRequest{
		Username: username,
		Password: password,
	}
	return cs.sendLoginRequest(ctx, backendLoginUsernameURL, req)
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
	// THIS IS THE KEY: cs.HTTPClient maintains the cookie jar!
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
	err = helpers.DecodeBackendResponse(resp, &target)
	if err != nil {
		return nil, backendError("Failed to decode response: " + err.Error())
	}

	return &target, nil
}
