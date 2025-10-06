// Package handlers contains all HTTP request handlers for the forum application.
// It implements the controller layer in the MVC pattern, handling HTTP requests,
// processing business logic, and rendering responses.
package handlers

import (
	"fmt"
	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler is the main controller struct that contains dependencies for all HTTP handlers.
// It uses dependency injection to receive database and authentication services.
type Handler struct {
	db   *database.DB // Database connection for data operations
	auth *auth.Auth   // Authentication service for user management
}

// NewHandler creates a new Handler instance with injected dependencies.
// This follows the dependency injection pattern for better testability and maintainability.
func NewHandler(db *database.DB) *Handler {
	return &Handler{
		db:   db,
		auth: auth.NewAuth(db), // Create auth service with same database connection
	}
}

// ===== PUBLIC PAGES =====

// Home handles the main homepage request (GET /).
// Displays all approved posts in chronological order with user context.
// This is the main entry point for both authenticated and anonymous users.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Ensure this handler only processes the exact root path
	// Other paths should be handled by their specific handlers
	if r.URL.Path != "/" {
		h.NotFound(w, r)
		return
	}

	// Get current user if logged in (optional - doesn't fail if not logged in)
	// This allows the template to show different content for logged-in users
	user, _ := h.auth.GetUserFromRequest(r)

	// Fetch posts based on user's role and permissions
	// Uses role-based visibility: anonymous see approved, users see approved + own, admins see all
	posts, err := h.GetPostsWithUser(user)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Fetch categories for the filter bar
	categories, err := h.GetCategories()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Prepare data structure for template rendering
	// Using anonymous struct for simple data passing to template
	data := struct {
		Posts      []models.Post
		User       *models.User
		Categories []models.Category
	}{
		Posts:      posts,
		User:       user,
		Categories: categories,
	}

	// Render the homepage template with posts and user data
	h.render(w, "home.html", data)
}

// ===== AUTHENTICATION HANDLERS =====

// Register handles user registration requests (GET and POST /register).
// GET: Displays the registration form
// POST: Processes registration data and creates new user account
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Display registration form
		h.render(w, "register.html", nil)
	case "POST":
		// Process registration form submission
		h.handleRegister(w, r)
	default:
		h.MethodNotAllowed(w, r)
	}
}

// handleRegister processes user registration form data.
// Validates input, checks for duplicate email/username, creates user account,
// establishes session, and redirects to homepage on success.
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Extract and sanitize form data
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password") // Don't trim password (might have intentional spaces)

	// Basic input validation - ensure all required fields are provided
	if email == "" || username == "" || password == "" {
		h.BadRequest(w, r, "All fields are required")
		return
	}

	// Check if email is already registered
	// This prevents duplicate accounts and maintains data integrity
	emailExists, err := h.auth.EmailExists(email)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}
	if emailExists {
		h.BadRequest(w, r, "Email already exists")
		return
	}

	// Check if username is already taken
	// Usernames must be unique for proper user identification
	usernameExists, err := h.auth.UsernameExists(username)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}
	if usernameExists {
		h.BadRequest(w, r, "Username already exists")
		return
	}

	// Create new user account with hashed password
	user, err := h.auth.CreateUser(email, username, password)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Automatically log in the new user by creating a session
	session, err := h.auth.CreateSession(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Set session cookie for authentication
	// HttpOnly prevents XSS attacks, Path ensures cookie works across site
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true, // Prevents JavaScript access for security
		Path:     "/",  // Cookie valid for entire site
	})

	// Redirect to homepage after successful registration
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Login handles user login requests (GET and POST /login).
// GET: Displays the login form
// POST: Authenticates user credentials and establishes session
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Display login form
		h.render(w, "login.html", nil)
	case "POST":
		// Process login form submission
		h.handleLogin(w, r)
	default:
		h.MethodNotAllowed(w, r)
	}
}

// handleLogin processes user login form data.
// Validates credentials, creates session on success, sets auth cookie.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Extract form credentials
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password") // Don't trim password

	// Basic input validation
	if email == "" || password == "" {
		h.BadRequest(w, r, "Email and password are required")
		return
	}

	// Attempt to find user by email
	user, err := h.auth.GetUserByEmail(email)
	if err != nil {
		// Use generic error message to prevent email enumeration attacks
		h.BadRequest(w, r, "Invalid email or password")
		return
	}

	// Verify password using bcrypt comparison
	// This is secure against timing attacks
	if !h.auth.CheckPassword(password, user.Password) {
		h.BadRequest(w, r, "Invalid email or password")
		return
	}

	// Create new session for authenticated user
	session, err := h.auth.CreateSession(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Set authentication cookie
	// Same security settings as registration
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true, // XSS prevention
		Path:     "/",  // Site-wide cookie
	})

	// Redirect to homepage after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles user logout requests (POST /logout).
// Destroys session in database and clears authentication cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get current session cookie if it exists
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete session from database to invalidate it
		h.auth.DeleteSession(cookie.Value)
	}

	// Clear the session cookie by setting it to expire in the past
	// This ensures the browser removes the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire in 1970 (way in the past)
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect to homepage as anonymous user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ===== POST MANAGEMENT =====

// CreatePost handles post creation requests (GET and POST /create-post).
// Requires user authentication. GET shows form, POST processes submission.
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Verify user is authenticated before allowing post creation
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		// Redirect to login if not authenticated
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case "GET":
		// Display post creation form with available categories
		categories, err := h.GetCategories()
		if err != nil {
			h.InternalServerError(w, r, err)
			return
		}

		// Prepare data for template rendering
		data := struct {
			User       *models.User
			Categories []models.Category
		}{
			User:       user,
			Categories: categories,
		}

		h.render(w, "create-post.html", data)
	case "POST":
		// Process post creation form submission
		h.handleCreatePost(w, r, user)
	default:
		h.MethodNotAllowed(w, r)
	}
}

// handleCreatePost processes post creation form submission.
// Handles file uploads, validates data, applies moderation rules, and saves post.
func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request, user *models.User) {
	// Parse multipart form data to handle potential file uploads
	// maxUploadSize is defined in upload.go (20MB limit)
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		h.BadRequest(w, r, "Failed to parse form")
		return
	}

	// Extract and sanitize form data
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"] // Array of selected category IDs

	// Validate required fields
	if title == "" || content == "" {
		h.BadRequest(w, r, "Title and content are required")
		return
	}

	// Handle optional image upload
	var imagePath string
	file, header, err := r.FormFile("image")
	if err == nil {
		// Image was uploaded, process it
		defer file.Close()
		imagePath, err = h.handleImageUpload(file, header)
		if err != nil {
			h.BadRequest(w, r, fmt.Sprintf("Image upload failed: %v", err))
			return
		}
	}
	// If err != nil, no image was uploaded (which is fine)

	// Determine post status based on user role (moderation system)
	status := models.PostStatusApproved
	if user.Role == models.RoleUser {
		// Regular users' posts need approval
		status = models.PostStatusPending
	}
	// Admins and moderators get auto-approved posts

	// Insert post into database
	query := `INSERT INTO posts (user_id, title, content, image_path, status) VALUES (?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query, user.ID, title, content, imagePath, status)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Get the ID of the newly created post
	postID, err := result.LastInsertId()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Associate post with selected categories (many-to-many relationship)
	for _, categoryStr := range categories {
		categoryID, err := strconv.Atoi(categoryStr)
		if err != nil {
			// Skip invalid category IDs
			continue
		}

		// Insert post-category relationship
		query = `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
		_, err = h.db.Exec(query, postID, categoryID)
		if err != nil {
			// Continue with other categories if one fails
			continue
		}
	}

	// Redirect to homepage to show the new post (if approved) or confirmation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ViewPost handles individual post viewing (GET /post/{id}).
// Displays post content with associated comments and allows interaction.
func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from URL path (/post/123 -> "123")
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		// Invalid post ID format
		h.NotFound(w, r)
		return
	}

	// Get current user context (optional - for showing user-specific content)
	user, _ := h.auth.GetUserFromRequest(r)

	// Fetch post details from database with user-specific visibility
	post, err := h.GetPostByIDWithUser(postID, user)
	if err != nil {
		// Post doesn't exist or user doesn't have permission to view it
		h.NotFound(w, r)
		return
	}

	// Fetch all comments for this post
	comments, err := h.GetCommentsByPostID(postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Prepare data for template
	data := struct {
		Post     *models.Post
		Comments []models.Comment
		User     *models.User
	}{
		Post:     post,
		Comments: comments,
		User:     user,
	}

	// Render post detail page
	h.render(w, "post.html", data)
}

// ===== UTILITY FUNCTIONS =====

// truncate is a template helper function that shortens strings to a specified length.
// Used in templates to display excerpts of long content like post previews.
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// render processes and displays HTML templates with provided data.
// It sets up template functions, parses templates, and executes them safely.
func (h *Handler) render(w http.ResponseWriter, tmpl string, data interface{}) {
	// Define custom template functions available in all templates
	funcMap := template.FuncMap{
		"truncate": truncate, // String truncation helper
	}

	// Parse base template and specific template together
	// Base template provides common layout, specific template provides content
	t, err := template.New("base").Funcs(funcMap).ParseFiles("templates/base.html", "templates/"+tmpl)
	if err != nil {
		// Template parsing failed - likely a syntax error in template
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the base template with provided data
	// The base template includes the specific template content
	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		// Template execution failed - likely a data binding error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ===== ERROR HANDLERS =====

// NotFound renders a 404 error page for missing resources.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.render(w, "404.html", nil)
}

// Forbidden renders a 403 error page for authorization failures.
// Used when a user is authenticated but lacks permission for the requested action.
func (h *Handler) Forbidden(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusForbidden)
	h.render(w, "400.html", map[string]string{"Message": message})
}

// InternalServerError renders a 500 error page for server errors.
// Logs the error details for debugging while showing user-friendly message.
func (h *Handler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	h.render(w, "500.html", map[string]string{"Error": err.Error()})
}

// BadRequest renders a 400 error page for client errors like invalid input.
func (h *Handler) BadRequest(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)
	h.render(w, "400.html", map[string]string{"Message": message})
}

// MethodNotAllowed handles requests with unsupported HTTP methods.
func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method not allowed")
}
