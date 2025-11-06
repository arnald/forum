package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arnald/forum/cmd/client/domain"
	h "github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/validation"
	"github.com/arnald/forum/internal/pkg/path"
)

const (
	notFoundMessage       = "Oops! The page you're looking for has vanished into the digital void."
	backendURL            = "http://localhost:8080/api/v1/register"
	backendRequestTimeout = 10 * time.Second
)

// Helper for rendering different templates (login/register)
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	resolver := path.NewResolver()

	tmplPath := resolver.GetPath("frontend/html/pages/" + templateName + ".html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing %s: %v", tmplPath, err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("Error executing %s template: %v", templateName, err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// HomePage Handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	// handle / and /categories
	if r.URL.Path != "/" && r.URL.Path != "/categories" {
		notFoundHandler(w, r, notFoundMessage, http.StatusNotFound)
		return
	}

	resolver := path.NewResolver()

	file, err := os.Open(resolver.GetPath("cmd/client/data/categories.json"))
	if err != nil {
		log.Println("Error opening categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var categoryData domain.CategoryData
	err = json.NewDecoder(file).Decode(&categoryData)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.PrepareCategories(categoryData.Data.Categories)

	pageData := domain.HomePageData{
		Categories: categoryData.Data.Categories,
		ActivePage: r.URL.Path,
	}

	/*
		data := BasePageData{
			Categories: categoryData.Data.Categories,
			User:       nil, // or retrieve user from session/cookie when you implement auth
		}
		tmpl.ExecuteTemplate(w, "base", data)
	*/

	// tmpl, err := template.ParseGlob(resolver.GetPath("frontend/html/**/*.html"))
	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/home.html", // render homepage for actual content
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/category_details.html",
		"frontend/html/partials/categories.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)

		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Register Handler GET
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	renderTemplate(w, "register", domain.RegisterFormErrors{})
}

// Register Handler POST
func RegisterPost(w http.ResponseWriter, r *http.Request) {
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

	// FrontEnd Validation - Quick feedback for user
	data.UsernameError = validation.ValidateUsername(username)
	data.EmailError = validation.ValidateEmail(email)
	data.Password = validation.ValidatePassword(password)

	// If errors, re-render register page with errors
	if data.UsernameError != "" || data.EmailError != "" || data.Password != "" {
		renderTemplate(w, "register", data)
		return
	}

	// Sending validated data to backend for backend Registration
	backendReq := domain.BackendRegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	backendResp, backendErr := registerWithBackend(r.Context(), backendReq)
	if backendErr != nil {
		// Backend validation/registration failed
		data.UsernameError = ""
		data.EmailError = ""
		data.Password = ""

		// Parse backend error and show appropriate message
		errorMsg := backendErr.Error()

		// Try to determine which field the error is about
		if strings.Contains(errorMsg, "email") {
			data.EmailError = errorMsg
		} else if strings.Contains(errorMsg, "username") {
			data.UsernameError = errorMsg
		} else {
			data.Password = errorMsg
		}

		renderTemplate(w, "register", data)
		return
	}

	// SUCCESS - User registered, redirect to login or home
	log.Printf("User registered successfully: %s (ID: %s)", username, backendResp.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// registerWithBackend - Makes HTTP request to backend register endpoint
func registerWithBackend(ctx context.Context, req domain.BackendRegisterRequest) (*domain.BackendRegisterResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request to backend
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, backendURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle different response statuses
	if resp.StatusCode != http.StatusCreated {
		// Backend returns an error
		var errResp domain.BackendErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Message != "" {
			return nil, backendError(errResp.Message)
		}
		return nil, backendError("Registration failed. Please try again.")
	}

	// Success response
	var successResp domain.BackendRegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&successResp)
	if err != nil {
		return nil, err
	}

	return &successResp, nil
}

// backendError - Custom error type for backend errors
type backendError string

func (e backendError) Error() string {
	return string(e)
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
	resolver := path.NewResolver()

	tmpl, err := template.ParseFiles(resolver.GetPath("frontend/html/pages/not_found.html"))
	if err != nil {
		http.Error(w, errorMessage, httpStatus)
		log.Println("Error loading not_found_page.html:", err)

		return
	}

	data := struct {
		StatusText   string
		ErrorMessage string
		StatusCode   int
	}{
		StatusText:   http.StatusText(httpStatus),
		ErrorMessage: errorMessage,
		StatusCode:   httpStatus,
	}

	w.WriteHeader(httpStatus)
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, errorMessage, httpStatus)
	}
}
