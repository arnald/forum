package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/helpers/validation"
)

const (
	accessTokenMaxAge  = 7 // minutes
	refreshTokenMaxAge = 7 // days
)

// LoginPage handles GET requests to /login.
func (cs *ClientServer) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// if !cookies {
	// 	templates.RenderTemplate(w, "login", domain.LoginFormErrors{})
	// }else {
	// 	verifiyCookies with backend
	// /me handler, json response with verified user
	// }
	templates.RenderTemplate(w, "login", domain.LoginFormErrors{})
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

	loginType := strings.TrimSpace(r.FormValue("loginType"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := domain.LoginFormErrors{
		Password: password,
	}

	data.PasswordError = validation.ValidatePassword(password)

	// Route based on login type
	if loginType == "email" {
		cs.handleEmailLogin(w, r, data)
	} else {
		cs.handleUsernameLogin(w, r, data)
	}
}

// handleEmailLogin processes email-based login.
func (cs *ClientServer) handleEmailLogin(w http.ResponseWriter, r *http.Request, data domain.LoginFormErrors) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Email = email

	// FRONTEND VALIDATION
	data.EmailError = validation.ValidateEmail(email)
	data.PasswordError = validation.ValidatePassword(password)

	// If frontend validation fails, re-render login page with errors
	if data.EmailError != "" || data.PasswordError != "" {
		templates.RenderTemplate(w, "login", data)
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
		templates.RenderTemplate(w, "login", data)
		return
	}

	cs.setSessionCookies(w, backendResp.AccessToken, backendResp.RefreshToken)

	// SUCCESS - User logged in, redirect to homepage
	log.Printf("User logged in successfully with email: %s (ID: %s)", backendResp.Username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleUsernameLogin processes username-based login.
func (cs *ClientServer) handleUsernameLogin(w http.ResponseWriter, r *http.Request, data domain.LoginFormErrors) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Username = username

	// FRONTEND VALIDATION
	data.UsernameError = validation.ValidateUsername(username)
	data.PasswordError = validation.ValidatePassword(password)

	// If frontend validation fails, re-render login page with errors
	if data.UsernameError != "" || data.PasswordError != "" {
		templates.RenderTemplate(w, "login", data)
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
		templates.RenderTemplate(w, "login", data)
		return
	}

	// Set cookies for session persistence
	cs.setSessionCookies(w, backendResp.AccessToken, backendResp.RefreshToken)

	// SUCCESS - User logged in, redirect to homepage
	log.Printf("User logged in successfully with username: %s (ID: %s)", backendResp.Username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// loginWithBackendEmail sends login request to backend email endpoint.
func (cs *ClientServer) loginWithBackendEmail(ctx context.Context, email string, password string) (*domain.BackendLoginResponse, error) {
	req := domain.BackendLoginRequest{
		Email:    email,
		Password: password,
	}
	return cs.sendLoginRequest(ctx, backendLoginEmailURL, req)
}

// loginWithBackendUsername sends login request to backend username endpoint.
func (cs *ClientServer) loginWithBackendUsername(ctx context.Context, username string, password string) (*domain.BackendLoginResponse, error) {
	req := domain.BackendLoginRequest{
		Username: username,
		Password: password,
	}
	return cs.sendLoginRequest(ctx, backendLoginUsernameURL, req)
}

// sendLoginRequest sends the login request to the backend API.
func (cs *ClientServer) sendLoginRequest(ctx context.Context, backendURL string, req domain.BackendLoginRequest) (*domain.BackendLoginResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, backendURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Login request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

// setSessionCookies sets the access and refresh tokens as cookies.
func (cs *ClientServer) setSessionCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(float64(accessTokenMaxAge) * time.Minute.Seconds()),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(float64(refreshTokenMaxAge) * time.Hour.Seconds()),
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}
