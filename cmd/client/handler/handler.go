package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/arnald/forum/cmd/helpers/validation"
	"github.com/arnald/forum/internal/pkg/uuid"
)

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
)

// Post represents a created post (file-backed prototype).
// type Post struct {
// 	ID           int       `json:"id"`
// 	Title        string    `json:"title"`
// 	Body         string    `json:"body"`
// 	CategorySlug string    `json:"category_slug"`
// 	Author       string    `json:"author"`
// 	CreatedAt    time.Time `json:"created_at"`
// // }

// type Topic struct {
//     Title       string    `json:"title"`
//     ID          int       `json:"id"`
//     Content     string    `json:"content"`
//     Author      string    `json:"author"`
//     CreatedAt   time.Time `json:"created_at"`
//     Likes       int       `json:"likes"`
//     Dislikes    int       `json:"dislikes"`
//     Views       int       `json:"views"`
//     Comments    int       `json:"comments"`
//     ImageURL    string    `json:"image_url"`
//     LinkURL     string    `json:"link_url"`
// }

// type Comment struct {
//     ID        int
//     Author    string
//     AvatarURL string
//     Date      string
//     Text      string
//     ImageURL  string
//     Likes     int
//     Dislikes  int
// }

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

// Data for Homepage
type HomePageData struct {
	Categories []Category
	ActivePage string
}

// Data for Single Category
type CategoryPageData struct {
	Categories []Category // for category_details partial
	Category   Category   // current category
	ActivePage string
}

// Data for Single Topic
type TopicPageData struct {
	Category   Category
	Topic      Topic // current topic
	ActivePage string
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

// Create-post form errors
type PostFormErrors struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Categories  string `json:"categories,omitempty"`
	Image       string `json:"image,omitempty"`
	Link        string `json:"link,omitempty"`
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
	// handle / and /categories
	if r.URL.Path != "/" && r.URL.Path != "/categories" {
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

	pageData := HomePageData{
		Categories: data.Data.Categories,
		ActivePage: r.URL.Path,
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
	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Single Category Handler
func CategoryPage(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/category/")
	if slug == "" {
		notFoundHandler(w, r, "Category not found", http.StatusNotFound)
		return
	}

	basePath, _ := os.Getwd()
	jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
	jsonPath = filepath.Clean(jsonPath)

	file, err := os.Open(jsonPath)
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var catData CategoryData
	err = json.NewDecoder(file).Decode(&catData)
	if err != nil {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return

	}

	// Find the selected category by slug
	var selected Category
	found := false
	for _, cat := range catData.Data.Categories {
		if cat.Slug == slug {
			selected = cat
			found = true
			break
		}
	}

	if !found {
		notFoundHandler(w, r, "Category not found", http.StatusNotFound)
		return
	}

	pageData := CategoryPageData{
		Categories: catData.Data.Categories,
		Category:   selected,
		ActivePage: "category",
	}

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/category.html", // render this as content
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/category_details.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
		log.Println("Render error:", err)
	}
}

// Single Topic Handler
func TopicPage(w http.ResponseWriter, r *http.Request) {
	topicIDStr := strings.TrimPrefix(r.URL.Path, "/topic/")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	basePath, _ := os.Getwd()
	jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
	jsonPath = filepath.Clean(jsonPath)

	file, err := os.Open(jsonPath)
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var catData CategoryData
	err = json.NewDecoder(file).Decode(&catData)
	if err != nil {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	// Find the topic and its category
	var topic Topic
	var category Category
	found := false
	for _, cat := range catData.Data.Categories {
		for _, t := range cat.Topics {
			if t.ID == topicID {
				topic = t
				category = cat
				found = true
				break
			}
		}
	}

	if !found {
		notFoundHandler(w, r, "Topic not found", http.StatusNotFound)
		return
	}

	pageData := TopicPageData{
		Category:   category,
		Topic:      topic,
		ActivePage: "topic",
	}

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/topic.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
		log.Println("Render error:", err)
	}
}

// Add Comment Handler
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 20MB file + form data)
	err := r.ParseMultipartForm(20 << 10)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("topic_id")
	content := r.FormValue(("add-comment"))

	// Handle optional image
	var imagePath string
	file, header, err := r.FormFile("comment-image-upload")
	if err == nil {
		defer file.Close()

		// Ensure upload directory exists
		basePath, _ := os.Getwd()
		uploadDir := filepath.Join(basePath, "frontend", "static", "images", "uploads", "comments")
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create upload dir", http.StatusInternalServerError)
			return
		}

		// Generate unique filename
		uuidProvider := uuid.NewProvider()
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s%s", uuidProvider.NewUUID().String(), ext)
		dstPath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		// Public path for frontend + backend
		imagePath = "/static/images/uploads/comments/" + filename
	}

	// Build payload for backend
	payload := map[string]interface{}{
		"topic_id": topicID,
		"author":   "TODO: get from session", // You’d replace this later
		"content":  content,
		"image":    imagePath,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode payload", http.StatusInternalServerError)
		return
	}

	// Send payload to backend API
	resp, err := http.Post("http://localhost:8080/api/comments", "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		http.Error(w, "Failed to reach backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Forward backend response back to JS
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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

// CreatePostPage renders the standalone create_post.html page
func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
	jsonPath = filepath.Clean(jsonPath)

	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("Error opening categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var catData CategoryData
	if err := json.NewDecoder(file).Decode(&catData); err != nil {
		log.Println("Error decoding categories.json:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// View model passed to template: categories + blank form + errors
	view := struct {
		Categories []Category
		Form       struct {
			Title       string
			Description string
			Link        string
			Selected    []string
		}
		Errors PostFormErrors
	}{
		Categories: catData.Data.Categories,
	}

	// render pages/create_post.html
	renderTemplate(w, "create_post", view)
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

	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := LoginFormErrors{
		Identifier: identifier, //preserve the value
	}

	if identifier == "" {
		data.IdentifierError = "Username or Email is required."
	} else if !validation.IsValidEmail(identifier) && len(identifier) < 3 {
		// Could be an invalid email *or* too-short username
		data.IdentifierError = "Invalid username or email."
	}
	data.Password = validation.ValidatePassword(password)

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

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data := RegisterFormErrors{
		Username: username,
		Email:    email,
	}

	data.UsernameError = validation.ValidateUsername(username)
	data.EmailError = validation.ValidateEmail(email)
	data.Password = validation.ValidatePassword(password)

	// If errors, re-render register page with errors
	if data.UsernameError != "" || data.EmailError != "" || data.Password != "" {
		renderTemplate(w, "register", data)
		return
	}

	// TODO: Register user in database
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreatePostPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (adjust maxMemory as needed; 10MB here)
	const maxMemory = 10 << 20 // 10 MB
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get values
	// checkbox values are in r.MultipartForm.Value["categories"]
	var selected []string
	if r.MultipartForm != nil && r.MultipartForm.Value != nil {
		if vals, ok := r.MultipartForm.Value["categories"]; ok {
			selected = vals
		}
	}
	// fallback: also try single value
	if len(selected) == 0 {
		if v := r.FormValue("category"); v != "" {
			selected = []string{v}
		}
	}

	title := strings.TrimSpace(r.FormValue("title"))
	desc := strings.TrimSpace(r.FormValue("description"))
	link := strings.TrimSpace(r.FormValue("link"))

	// Validate
	errors := PostFormErrors{}
	if title == "" {
		errors.Title = "Title is required."
	}
	if desc == "" {
		errors.Description = "Description is required."
	}
	if len(selected) == 0 {
		errors.Categories = "Please select at least one category."
	}
	// Optional: validate link using net/url or allow empty
	if link != "" {
		if _, err := url.ParseRequestURI(link); err != nil {
			errors.Link = "Invalid URL format."
		}
	}

	// If validation errors -> re-render the form with previous values & errors
	if errors.Title != "" || errors.Description != "" || errors.Categories != "" || errors.Image != "" || errors.Link != "" {
		// reload categories
		basePath, err := os.Getwd()
		if err != nil {
			log.Println("Error getting working directory:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		jsonPath := filepath.Join(basePath, "cmd", "client", "data", "categories.json")
		jsonPath = filepath.Clean(jsonPath)

		f, err := os.Open(jsonPath)
		if err != nil {
			log.Println("Error opening categories.json:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var catData CategoryData
		if err := json.NewDecoder(f).Decode(&catData); err != nil {
			log.Println("Error decoding categories.json:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		view := struct {
			Categories []Category
			Form       struct {
				Title       string
				Description string
				Link        string
				Selected    []string
			}
			Errors PostFormErrors
		}{
			Categories: catData.Data.Categories,
			Errors:     errors,
		}
		view.Form.Title = title
		view.Form.Description = desc
		view.Form.Link = link
		view.Form.Selected = selected

		renderTemplate(w, "create_post", view)
		return
	}

	// TODO: persist the post (DB/API) and store file if provided.
	// For now: redirect to first selected category if present, otherwise to "/"
	firstCategory := ""
	if len(selected) > 0 {
		firstCategory = strings.TrimSpace(selected[0])
	}

	if firstCategory != "" {
		http.Redirect(w, r, "/category/"+firstCategory, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreatePostPost handles submitted create post form and redirects to / or to the first chosen category

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
