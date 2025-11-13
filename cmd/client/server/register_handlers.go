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

// RegisterPage handles GET requests to /register.
func (cs *ClientServer) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "register", domain.RegisterFormErrors{})
}

// RegisterPost handles POST requests to /register.
func (cs *ClientServer) RegisterPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := domain.RegisterFormErrors{
		Username: username,
		Email:    email,
	}

	// FRONTEND VALIDATION - Quick feedback for user
	data.UsernameError = validation.ValidateUsername(username)
	data.EmailError = validation.ValidateEmail(email)
	data.Password = validation.ValidatePassword(password)

	// If frontend validation fails, re-render register page with errors
	if data.UsernameError != "" || data.EmailError != "" || data.Password != "" {
		renderTemplate(w, "register", data)
		return
	}

	// BACKEND REGISTRATION - Send validated data to backend
	backendReq := domain.BackendRegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	backendResp, backendErr := cs.registerWithBackend(ctx, backendReq)
	if backendErr != nil {
		// Backend validation/registration failed
		data.UsernameError = ""
		data.EmailError = ""
		data.Password = ""

		// Parse backend error and show appropriate message
		errorMsg := backendErr.Error()

		// Try to determine which field the error is about
		switch {
		case strings.Contains(errorMsg, "email"):
			data.EmailError = errorMsg
		case strings.Contains(errorMsg, "username"):
			data.UsernameError = errorMsg
		default:
			data.Password = errorMsg
		}

		renderTemplate(w, "register", data)
		return
	}

	// SUCCESS - User registered, redirect to login or homepage
	log.Printf("User registered successfully: %s (ID: %s)", username, backendResp.UserID)
	// http.Redirect(w, r, "/login", http.StatusSeeOther)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================================================
// API REQUEST METHODS (Methods on ClientServer)
// ============================================================================

// registerWithBackend sends registration request to backend API.
// The HTTP client includes the cookie jar, so cookies will be automatically.
// handled for all subsequent requests.
func (cs *ClientServer) registerWithBackend(ctx context.Context, req domain.BackendRegisterRequest) (*domain.BackendRegisterResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	// Create HTTP request to backend with context
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, backendRegisterURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request using the client's HTTP client (with cookie jar)
	// THIS IS THE KEY: cs.HTTPClient maintains the cookie jar!
	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Registration request failed: " + err.Error())
	}
	defer resp.Body.Close()

	// Handle different response statuses
	if resp.StatusCode != http.StatusCreated {
		// Backend returned an error
		var errResp domain.BackendErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, backendError("Registration failed. Please try again.")
		}

		if errResp.Message != "" {
			return nil, backendError(errResp.Message)
		}
		return nil, backendError("Registration failed. Please try again.")
	}

	// Success response
	target := domain.BackendRegisterResponse{}
	err = DecodeBackendResponse(resp, &target)
	if err != nil {
		return nil, backendError("Failed to decode response " + err.Error())
	}
	return &target, nil
}
