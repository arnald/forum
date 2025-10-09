// Package main is the entry point for the Forum web application.
// This is a comprehensive forum system built with Go's standard library,
// featuring user authentication, posts, comments, categories, notifications,
// and admin functionality.
package main

import (
	"forum/internal/database"
	"forum/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// main is the application entry point that sets up the database,
// initializes handlers, configures routes, and starts the HTTP server.
func main() {
	// Load environment variables from .env file
	// Ignore error if .env doesn't exist (allow production deployments with system env vars)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database path from environment or use default location
	dbPath := database.GetDBPath()

	// Create data directory if it doesn't exist
	// This directory will store the SQLite database file and uploaded images
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Initialize database connection
	// NewDB creates a connection to SQLite database and verifies connectivity
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close() // Ensure database connection is closed when main exits

	// Initialize database schema
	// This creates all necessary tables if they don't exist
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create main handler with database dependency injection
	// The handler contains all business logic and database operations
	h := handlers.NewHandler(db)

	// Initialize OAuth configuration
	handlers.InitOAuth()

	// === PUBLIC ROUTES ===
	// Home page - displays all approved posts
	http.HandleFunc("/", h.Home)

	// Static file serving for CSS, JS, images, and uploads
	// StripPrefix removes "/static/" from URL before serving files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// === AUTHENTICATION ROUTES ===
	// Rate limiting is applied to prevent brute force attacks
	http.HandleFunc("/register", h.RateLimitedHandler(h.Register)) // User registration
	http.HandleFunc("/login", h.RateLimitedHandler(h.Login))       // User login
	http.HandleFunc("/logout", h.Logout)                          // User logout

	// === OAUTH AUTHENTICATION ROUTES ===
	// OAuth login and callback handlers for third-party authentication
	http.HandleFunc("/auth/google", h.GoogleLogin)           // Redirect to Google OAuth
	http.HandleFunc("/auth/google/callback", h.GoogleCallback) // Handle Google OAuth callback
	http.HandleFunc("/auth/github", h.GitHubLogin)           // Redirect to GitHub OAuth
	http.HandleFunc("/auth/github/callback", h.GitHubCallback) // Handle GitHub OAuth callback
	http.HandleFunc("/auth/facebook", h.FacebookLogin)       // Redirect to Facebook OAuth
	http.HandleFunc("/auth/facebook/callback", h.FacebookCallback) // Handle Facebook OAuth callback

	// === POST MANAGEMENT ROUTES ===
	http.HandleFunc("/create-post", h.RateLimitedHandler(h.CreatePost)) // Create new post (rate limited)
	http.HandleFunc("/post/", h.ViewPost)                               // View individual post with comments
	http.HandleFunc("/edit-post/", h.EditPost)                          // Edit existing post (owner/admin only)
	http.HandleFunc("/delete-post/", h.DeletePost)                      // Delete post (owner/admin only)

	// === COMMENT MANAGEMENT ROUTES ===
	http.HandleFunc("/create-comment", h.RateLimitedHandler(h.CreateCommentHandler)) // Add comment to post
	http.HandleFunc("/edit-comment/", h.EditComment)                                 // Edit comment (owner/admin only)
	http.HandleFunc("/delete-comment/", h.DeleteComment)                             // Delete comment (owner/admin only)

	// === LIKE/DISLIKE SYSTEM ROUTES ===
	// These routes handle the voting system for posts and comments
	http.HandleFunc("/like-post/", h.LikePostHandler)       // Like a post
	http.HandleFunc("/dislike-post/", h.DislikePostHandler) // Dislike a post
	http.HandleFunc("/like-comment/", h.LikeCommentHandler) // Like a comment
	http.HandleFunc("/dislike-comment/", h.DislikeCommentHandler) // Dislike a comment

	// === FILTERING AND SEARCH ROUTES ===
	http.HandleFunc("/filter", h.FilterPosts) // Filter posts by category or user-specific criteria

	// === NOTIFICATION SYSTEM ROUTES ===
	http.HandleFunc("/notifications", h.Notifications)                  // View user notifications
	http.HandleFunc("/mark-read", h.MarkNotificationReadHandler)        // Mark single notification as read
	http.HandleFunc("/mark-all-read", h.MarkAllNotificationsRead)       // Mark all notifications as read

	// === USER ACTIVITY ROUTES ===
	http.HandleFunc("/activity", h.Activity) // View user's activity (posts, comments, likes)

	// === ADMIN PANEL ROUTES ===
	// These routes are protected and only accessible by admins/moderators
	http.HandleFunc("/admin", h.AdminPanel)                                    // Main admin dashboard
	http.HandleFunc("/admin/update-role", h.UpdateUserRole)                    // Change user roles
	http.HandleFunc("/admin/approve-post/", h.ApprovePost)                     // Approve pending posts
	http.HandleFunc("/admin/reject-post/", h.RejectPost)                       // Reject pending posts
	http.HandleFunc("/admin/create-category", h.CreateCategory)                // Create new categories
	http.HandleFunc("/admin/delete-category/", h.DeleteCategory)               // Delete categories
	http.HandleFunc("/admin/approve-moderator-posts", h.ApproveModeratorPosts) // Bulk approve moderator posts
	http.HandleFunc("/admin/bulk-approve-all", h.BulkApproveAllPending)        // Bulk approve all pending posts
	http.HandleFunc("/admin/bulk-reject-all", h.BulkRejectAllPending)          // Bulk reject all pending posts

	// Start HTTP server on port 8080
	// The server will handle all incoming requests and route them to appropriate handlers
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
