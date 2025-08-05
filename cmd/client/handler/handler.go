package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
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

type LoginFormErrors struct {
	Identifier      string `json:"-"` // the actual value user typed
	Password        string `json:"password,omitempty"`
	IdentifierError string `json:"identifier,omitempty"`
}

type RegisterFormErrors struct {
	Username      string `json:"-"`
	Email         string `json:"-"`
	Password      string `json:"password,omitempty"`
	UsernameError string `json:"username,omitempty"`
	EmailError    string `json:"email,omitempty"`
}

// Validation helpers
func validateEmail(email string) string {
	if email == "" {
		return "Email is required."
	}
	if !isValidEmail(email) {
		return "Invalid email format."
	}
	return ""
}

func validateUsername(username string) string {
	if username == "" {
		return "Username is required."
	}
	if len(username) < 3 {
		return "Username must be at least 3 characters"
	}
	return ""
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hasLower(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

func hasUpper(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func hasSpecial(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func validatePassword(pw string) string {
	if pw == "" {
		return "Password is required."
	}
	if len(pw) < 8 {
		return "Password must be 8+ chars"
	}

	var missing []string
	if !hasLower(pw) {
		missing = append(missing, "lowercase letter")
	}
	if !hasUpper(pw) {
		missing = append(missing, "uppercase letter")
	}
	if !hasDigit(pw) {
		missing = append(missing, "number")
	}
	if !hasSpecial(pw) {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return "Missing: " + strings.Join(missing, ", ")
	}

	return "" // empty string means no error
}

// Helper for rendering different templates (login/register)
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

	// tmpl, err := template.ParseGlob(filepath.Join(basePath, "frontend/html/**/*.html"))
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
	renderTemplate(w, "register", RegisterFormErrors{})
}

// Login Handler GET
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "login", LoginFormErrors{})
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

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	data := LoginFormErrors{
		Identifier: identifier, //preserve the value
	}

	if identifier == "" {
		data.IdentifierError = "Username or Email is required."
	} else if !isValidEmail(identifier) && len(identifier) < 3 {
		// Could be an invalid email *or* too-short username
		data.IdentifierError = "Invalid username or email."
	}
	data.Password = validatePassword(password)

	// If errors, re-render login page with errors
	if data.IdentifierError != "" || data.Password != "" {
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

	data := RegisterFormErrors{
		Username: username,
		Email:    email,
	}

	data.UsernameError = validateUsername(username)
	data.EmailError = validateEmail(email)
	data.Password = validatePassword(password)

	// If errors, re-render register page with errors
	if data.UsernameError != "" || data.EmailError != "" || data.Password != "" {
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
