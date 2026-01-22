package server

import (
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
	"github.com/arnald/forum/cmd/client/middleware"
)

const (
	accessTokenMaxAge  = 7 // minutes
	refreshTokenMaxAge = 7 // days
)

type LoginFormErrors struct {
	User          *domain.LoggedInUser `json:"-"`
	Username      string               `json:"-"`
	Email         string               `json:"-"`
	Password      string               `json:"-"`
	UsernameError string               `json:"username,omitempty"`
	EmailError    string               `json:"email,omitempty"`
	PasswordError string               `json:"password,omitempty"`
}

// BackendLoginRequest - sent to backend.
type BackendLoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

// BackendLoginResponse - response from backend.
type BackendLoginResponse struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginPage handles GET requests to /login.
func (cs *ClientServer) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If there is a user redirect him to homepage
	user := middleware.GetUserFromContext(r.Context())
	if user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.RenderTemplate(w, "login", LoginFormErrors{})
}

// LoginPost handles POST requests to /login.
func (cs *ClientServer) LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	loginType := strings.TrimSpace(r.FormValue("loginType"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := LoginFormErrors{
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
func (cs *ClientServer) handleEmailLogin(w http.ResponseWriter, r *http.Request, data LoginFormErrors) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Email = email

	// BACKEND LOGIN
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
		return
	}

	backendResp, backendErr := cs.loginWithBackendEmail(ctx, email, password, ip)
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
func (cs *ClientServer) handleUsernameLogin(w http.ResponseWriter, r *http.Request, data LoginFormErrors) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	data.Username = username

	// BACKEND LOGIN
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	backendResp, backendErr := cs.loginWithBackendUsername(ctx, username, password, ip)
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
func (cs *ClientServer) loginWithBackendEmail(ctx context.Context, email string, password string, ip string) (*BackendLoginResponse, error) {
	req := BackendLoginRequest{
		Email:    email,
		Password: password,
	}
	return cs.sendLoginRequest(ctx, cs.BackendURLs.LoginEmailURL(), req, ip)
}

// loginWithBackendUsername sends login request to backend username endpoint.
func (cs *ClientServer) loginWithBackendUsername(ctx context.Context, username string, password string, ip string) (*BackendLoginResponse, error) {
	req := BackendLoginRequest{
		Username: username,
		Password: password,
	}

	return cs.sendLoginRequest(ctx, cs.BackendURLs.LoginUsernameURL(), req, ip)
}

// sendLoginRequest sends the login request to the backend API.
func (cs *ClientServer) sendLoginRequest(ctx context.Context, backendURL string, req BackendLoginRequest, ip string) (*BackendLoginResponse, error) {
	resp, err := cs.newRequest(
		ctx,
		http.MethodPost,
		backendURL,
		req,
		ip,
	)
	if err != nil {
		defer resp.Body.Close()
		return nil, err
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
	target := BackendLoginResponse{}
	err = helpers.DecodeBackendResponse(resp, &target)
	if err != nil {
		return nil, backendError("Failed to decode response: " + err.Error())
	}

	return &target, nil
}

// setSessionCookies sets the access and refresh tokens as cookies.
func (cs *ClientServer) setSessionCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	// Use secure cookies when in production or when using HTTPS
	isSecure := cs.Config.Environment == "production" || cs.Config.TLSCertFile != ""

	log.Printf("Setting session cookies - isSecure: %v, TLSCertFile: %v", isSecure, cs.Config.TLSCertFile)

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(float64(accessTokenMaxAge) * time.Minute.Seconds()),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(float64(refreshTokenMaxAge) * time.Hour.Seconds()),
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	log.Printf("Cookies set successfully - access_token MaxAge: %d, refresh_token MaxAge: %d", accessCookie.MaxAge, refreshCookie.MaxAge)
}
