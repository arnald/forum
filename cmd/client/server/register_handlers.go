package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
	val "github.com/arnald/forum/internal/pkg/validator"
)

type RegisterFormErrors struct {
	Username      string `json:"-"`
	Email         string `json:"-"`
	Password      string `json:"-"`
	UsernameError string `json:"username,omitempty"`
	EmailError    string `json:"email,omitempty"`
	PasswordError string `json:"password,omitempty"`
}

// BackendRegisterRequest - matches backend RegisterUserReguestModel.
type BackendRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BackendRegisterResponse - matches backend RegisterUserResponse.
type BackendRegisterResponse struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

// RegisterPage handles GET requests to /register.
func (cs *ClientServer) RegisterPage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.RenderTemplate(w, "register", RegisterFormErrors{})
}

// RegisterPost handles POST requests to /register.
func (cs *ClientServer) RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := RegisterFormErrors{
		Username: username,
		Email:    email,
		Password: password,
	}

	validator := val.New()

	val.ValidateUserRegistration(validator, &data)
	if !validator.Valid() {
		data.UsernameError = validator.Errors["Username"]
		data.EmailError = validator.Errors["Email"]
		data.PasswordError = validator.Errors["Password"]
		templates.RenderTemplate(w, "register", data)
		return
	}

	// BACKEND REGISTRATION - Send validated data to backend
	backendReq := BackendRegisterRequest{
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

		templates.RenderTemplate(w, "register", data)
		return
	}

	// SUCCESS - User registered, redirect to login or homepage
	log.Printf("User registered successfully: %s (ID: %s)", username, backendResp.UserID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// registerWithBackend sends registration request to backend API.
// The HTTP client includes the cookie jar, so cookies will be automatically.
// handled for all subsequent requests.
func (cs *ClientServer) registerWithBackend(ctx context.Context, req BackendRegisterRequest) (*BackendRegisterResponse, error) {
	resp, err := cs.newRequest(
		ctx,
		http.MethodPost,
		backendRegisterURL,
		req,
	)
	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
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
	target := BackendRegisterResponse{}
	err = helpers.DecodeBackendResponse(resp, &target)
	if err != nil {
		return nil, backendError("Failed to decode response " + err.Error())
	}
	return &target, nil
}
