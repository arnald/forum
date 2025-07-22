package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"text/template"
	"unicode"
)

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
)

type Topic struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Category struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Topics      []Topic `json:"topics"`
	Logo        Logo    `json:"logo"`
	ID          int     `json:"id"`
}

type CategoryData struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}

/*
type BasePageData struct {
	Categories []Category
	User       *User // when you add auth
}
// When .User is nil, you show login/register links.
// When .User is populated, you show the user’s name, avatar, logout, etc.
*/

// TemplateData is an empty interface that marks types that can be passed to templates
type TemplateData interface{}

type LoginErrors struct {
	Email    string
	Password string
}

type RegisterErrors struct {
	Username string
	Email    string
	Password string
}

// helper functions
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, c := range pw {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmplPath := filepath.Join(basePath, "frontend", "html", "pages", templateName+".html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing %s: %v", templateName, err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("Error executing %s: %v", templateName, err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}

}

// HomePage Handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r, notFoundMessage, http.StatusNotFound)

		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Load categories.json
	jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
	jsonPath = filepath.Clean(jsonPath)

	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("Error opening categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data CategoryData
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	/*
		data := BasePageData{
			Categories: data.Data.Categories,
			User:       nil, // or retrieve user from session/cookie when you implement auth
		}
		tmpl.ExecuteTemplate(w, "base", data)
	*/

	tmpl, err := template.ParseGlob(filepath.Join(basePath, "frontend/html/**/*.html"))
	if err != nil {
		log.Println("Error loading home.html:", err)
		notFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)

		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data.Data.Categories)
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
	renderTemplate(w, "register", nil)
}

// Login Handler GET
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "login", nil)
}

// Login Handler POST
func LoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	data := LoginErrors{
		Email:    email,
		Password: password,
	}

	// Validate email
	if email == "" {
		data.Email = "Email is required."
	} else if !isValidEmail(email) {
		data.Email = "Invalid email format."
	}

	// Validate password
	if password == "" {
		data.Password = "Password is required."
	} else if !isValidPassword(password) {
		data.Password = "Password must be 8+ chars, including uppercase, lowercase, number, and special char."
	}

	// If errors, re-render login page with errors
	if data.Email != "" || data.Password != "" {
		renderTemplate(w, "login", data)
		return
	}

	// TODO: Authenticate user (check DB)
	// If auth fails, set an error like:
	// data.Errors.Password = "Invalid email or password."
	// renderLoginTemplate(w, data)
	// return

	// On success, redirect to homepage or user dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := RegisterErrors{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Validate username
	if username == "" {
		data.Username = "Username is required."
	}

	// Validate email
	if email == "" {
		data.Email = "Email is required."
	} else if !isValidEmail(email) {
		data.Email = "Invalid email format."
	}

	// Validate password
	if password == "" {
		data.Password = "Password is required."
	} else if !isValidPassword(password) {
		data.Password = "Password must be 8+ chars, including uppercase, lowercase, number, and special char."
	}

	// If errors, re-render register page with errors
	if data.Username != "" || data.Email != "" || data.Password != "" {
		renderTemplate(w, "register", data)
		return
	}

	// TODO: Register user in database
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request, errorMessage string, httpStatus int) {
	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templatePath := filepath.Join(basePath, "frontend", "html", "pages", "not_found.html")
	templatePath = filepath.Clean(templatePath)

	tmpl, err := template.ParseFiles(templatePath)
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
