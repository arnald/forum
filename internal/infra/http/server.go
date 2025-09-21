/*
Package http provides the HTTP server implementation for the forum application.

This package contains:
- HTTP server setup and configuration
- Route definitions and middleware integration
- Request routing and handler coordination
- Static file serving for frontend assets
- Authentication middleware application

The server follows a clean architecture pattern where HTTP handlers
are kept separate from business logic, and all requests are routed
through appropriate middleware for authentication, CORS, etc.
*/
package http

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/http/category"
	"github.com/arnald/forum/internal/infra/http/comment"
	"github.com/arnald/forum/internal/infra/http/health"
	"github.com/arnald/forum/internal/infra/http/post"
	userLogin "github.com/arnald/forum/internal/infra/http/user/login"
	userLogout "github.com/arnald/forum/internal/infra/http/user/logout"
	userRegister "github.com/arnald/forum/internal/infra/http/user/register"
	"github.com/arnald/forum/internal/infra/http/vote"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/session"
)

// Constants for API configuration
const (
	apiContext   = "/api/v1"                // Base path for all API endpoints
	readTimeout  = 5 * time.Second          // Timeout for reading request
	writeTimeout = 10 * time.Second         // Timeout for writing response
	idleTimeout  = 15 * time.Second         // Timeout for idle connections
)

// Server represents the HTTP server with all its dependencies
// It encapsulates the HTTP router, application services, configuration,
// session management, database connection, and logging functionality
type Server struct {
	appServices    app.Services           // Business logic services
	config         *config.ServerConfig   // Server configuration
	router         *http.ServeMux         // HTTP request router
	sessionManager user.SessionManager    // User session management
	db             *sql.DB                // Database connection
	logger         logger.Logger          // Application logger
}

// NewServer creates and configures a new HTTP server instance
// It initializes all dependencies, sets up routing, and prepares the server for handling requests
//
// Parameters:
//   - cfg: Server configuration containing host, port, timeouts, etc.
//   - db: Database connection for data persistence
//   - logger: Logger instance for application logging
//   - appServices: Business logic services for handling application operations
//
// Returns:
//   - *Server: Configured HTTP server ready to handle requests
func NewServer(cfg *config.ServerConfig, db *sql.DB, logger logger.Logger, appServices app.Services) *Server {
	httpServer := &Server{
		router:      http.NewServeMux(),  // Initialize HTTP router
		appServices: appServices,         // Inject business logic services
		config:      cfg,                 // Store server configuration
		db:          db,                  // Store database connection
		logger:      logger,              // Store logger instance
	}
	// Initialize session management
	httpServer.initSessionManager()
	// Set up all HTTP routes and handlers
	httpServer.AddHTTPRoutes()
	return httpServer
}

// AddHTTPRoutes configures all HTTP routes for the application
// This includes static file serving, frontend page routes, and API endpoints
// Routes are organized into sections: static files, frontend pages, and API endpoints
// Authentication middleware is applied to protected endpoints
func (server *Server) AddHTTPRoutes() {
	// =====================================
	// STATIC FILE SERVING
	// =====================================
	// Serve CSS, JavaScript, and image files from frontend/static directory
	server.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static/"))))

	// =====================================
	// FRONTEND PAGE ROUTES
	// =====================================
	// Serve HTML pages to users - these provide the user interface
	server.router.HandleFunc("/", server.servePage("home.html"))                        // Home page with forum overview
	server.router.HandleFunc("/login", server.servePage("login.html"))                  // Login form for authentication
	server.router.HandleFunc("/register", server.servePage("register.html"))            // Registration form for new users
	server.router.HandleFunc("/posts", server.servePage("posts.html"))                  // Posts listing page with filtering
	server.router.HandleFunc("/create-post", server.servePage("create-post.html"))      // Post creation form for authenticated users
	server.router.HandleFunc("/post/", server.servePage("post-detail.html"))            // Individual post view with comments

	// =====================================
	// API ENDPOINTS
	// =====================================

	// Health endpoint for monitoring server status
	server.router.HandleFunc(apiContext+"/health", health.NewHandler(server.logger).HealthCheck)

	// =====================================
	// USER AUTHENTICATION ENDPOINTS
	// =====================================
	// These endpoints handle user registration, login, and logout
	// They create and manage user sessions for maintaining authentication state

	// Login endpoint for username-based authentication
	server.router.HandleFunc(
		apiContext+"/login/username",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginUsername,
	)
	// Login endpoint for email-based authentication
	server.router.HandleFunc(
		apiContext+"/login/email",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginEmail,
	)
	// User registration endpoint for creating new accounts
	server.router.HandleFunc(
		apiContext+"/register",
		userRegister.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserRegister,
	)
	// Logout endpoint for terminating user sessions
	server.router.HandleFunc(
		apiContext+"/logout",
		userLogout.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLogout,
	)

	// =====================================
	// POST ENDPOINTS
	// =====================================
	// Initialize post handlers with required application services
	createPostHandler := post.NewCreatePostHandler(server.appServices.PostServices.Queries.CreatePost)
	getPostsHandler := post.NewGetPostsHandler(
		server.appServices.PostServices.Queries.GetAllPosts,      // Service for retrieving all posts
		server.appServices.PostServices.Queries.GetPostByID,      // Service for retrieving specific posts
		server.appServices.PostServices.Queries.GetPostsByCategory, // Service for category-filtered posts
	)
	updatePostHandler := post.NewUpdatePostHandler(server.appServices.PostServices.Queries.UpdatePost)

	// Protected post endpoints (require user authentication)
	server.router.Handle(apiContext+"/posts", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(createPostHandler.CreatePost)))
	server.router.Handle(apiContext+"/posts/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(updatePostHandler.UpdatePost)))

	// Public post endpoints (no authentication required)
	server.router.HandleFunc(apiContext+"/posts/all", getPostsHandler.GetAllPosts)   // Get all posts for public viewing
	server.router.HandleFunc(apiContext+"/posts/get", getPostsHandler.GetPostByID)   // Get specific post by ID

	// =====================================
	// COMMENT ENDPOINTS
	// =====================================
	// Initialize comment handlers with required application services
	createCommentHandler := comment.NewCreateCommentHandler(server.appServices.CommentServices.Queries.CreateComment)
	getCommentsHandler := comment.NewGetCommentsHandler(
		server.appServices.CommentServices.Queries.GetCommentsByPost, // Service for retrieving comments on posts
		server.appServices.CommentServices.Queries.GetCommentTree,    // Service for hierarchical comment structure
	)
	updateCommentHandler := comment.NewUpdateCommentHandler(server.appServices.CommentServices.Queries.UpdateComment)

	// Protected comment endpoints (require user authentication)
	server.router.Handle(apiContext+"/comments", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(createCommentHandler.CreateComment)))
	server.router.Handle(apiContext+"/comments/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(updateCommentHandler.UpdateComment)))

	// Public comment endpoints (no authentication required)
	server.router.HandleFunc(apiContext+"/comments/post", getCommentsHandler.GetCommentsByPost) // Get comments for a specific post
	server.router.HandleFunc(apiContext+"/comments/tree", getCommentsHandler.GetCommentTree)   // Get nested comment structure

	// =====================================
	// VOTE ENDPOINTS
	// =====================================
	// Initialize vote handler with all voting-related services
	voteHandler := vote.NewVoteHandler(
		server.appServices.VoteServices.Queries.CastVote,     // Service for casting likes/dislikes
		server.appServices.VoteServices.Queries.GetVoteStatus, // Service for checking vote status
		server.appServices.VoteServices.Queries.GetUserVotes,  // Service for retrieving user's voting history
	)

	// Protected vote endpoints (require user authentication)
	server.router.Handle(apiContext+"/votes/cast", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(voteHandler.CastVote)))     // Cast vote on content
	server.router.Handle(apiContext+"/votes/user", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(voteHandler.GetUserVotes))) // Get user's votes

	// Public vote endpoints (no authentication required)
	server.router.HandleFunc(apiContext+"/votes/status", voteHandler.GetVoteStatus) // Get vote counts and status

	// =====================================
	// CATEGORY ENDPOINTS
	// =====================================
	// Initialize category handler with all category management services
	categoryHandler := category.NewCategoryHandler(
		server.appServices.CategoryServices.Queries.CreateCategory,      // Service for creating new categories
		server.appServices.CategoryServices.Queries.GetAllCategories,    // Service for retrieving all categories
		server.appServices.CategoryServices.Queries.GetCategoryByID,     // Service for retrieving specific categories
		server.appServices.CategoryServices.Queries.UpdateCategory,      // Service for updating category information
		server.appServices.CategoryServices.Queries.DeleteCategory,      // Service for deleting categories
		server.appServices.CategoryServices.Queries.GetCategoryWithPosts, // Service for categories with their posts
	)

	// Protected category endpoints (require user authentication - typically admin only)
	server.router.Handle(apiContext+"/categories", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.CreateCategory)))
	server.router.Handle(apiContext+"/categories/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.UpdateCategory)))
	server.router.Handle(apiContext+"/categories/delete", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.DeleteCategory)))

	// Public category endpoints (no authentication required)
	server.router.HandleFunc(apiContext+"/categories/all", categoryHandler.GetAllCategories)     // Get all categories for browsing
	server.router.HandleFunc(apiContext+"/categories/get", categoryHandler.GetCategoryByID)      // Get specific category details
	server.router.HandleFunc(apiContext+"/categories/posts", categoryHandler.GetCategoryWithPosts) // Get category with associated posts
}

// ListenAndServe starts the HTTP server and begins listening for requests
// This method configures the server with CORS middleware and proper timeouts,
// then starts listening on the configured host and port
//
// The method:
// 1. Wraps the router with CORS middleware for cross-origin requests
// 2. Configures HTTP server timeouts from application configuration
// 3. Logs server startup information
// 4. Starts the server and handles shutdown gracefully
func (server *Server) ListenAndServe() {
	// Wrap router with CORS middleware to handle cross-origin requests
	corsWrappedRouter := middleware.NewCorsMiddleware(server.router)

	// Configure HTTP server with timeouts and address
	srv := &http.Server{
		Addr:         server.config.Host + ":" + server.config.Port, // Server address from configuration
		Handler:      corsWrappedRouter,                             // Router with CORS middleware
		ReadTimeout:  server.config.ReadTimeout,                    // Timeout for reading requests
		WriteTimeout: server.config.WriteTimeout,                   // Timeout for writing responses
		IdleTimeout:  server.config.IdleTimeout,                    // Timeout for idle connections
	}

	// Log server startup information for monitoring and debugging
	server.logger.PrintInfo("Starting server", map[string]string{
		"host":        server.config.Host,        // Server host address
		"port":        server.config.Port,        // Server port number
		"environment": server.config.Environment, // Current environment (dev/prod/test)
	})

	// Start the HTTP server and handle any startup errors
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		// Log fatal error if server fails to start (excluding graceful shutdown)
		server.logger.PrintFatal(err, nil)
	}
}

// initSessionManager initializes the session manager with database and configuration
// This sets up user session management for authentication and authorization
// Session manager handles:
// - Session creation and validation
// - Cookie management for browser sessions
// - Session expiration and cleanup
// - Security settings (HttpOnly, Secure, SameSite)
func (server *Server) initSessionManager() {
	server.sessionManager = session.NewSessionManager(server.db, server.config.SessionManager)
}

// servePage returns an HTTP handler function for serving HTML pages
// This helper function creates handlers that serve static HTML files from the frontend directory
//
// Parameters:
//   - filename: Name of the HTML file to serve (e.g., "home.html", "login.html")
//
// Returns:
//   - http.HandlerFunc: Handler function that serves the specified HTML file
//
// The handler serves files from the "frontend/html/pages/" directory
// and handles proper content-type headers and caching automatically
func (server *Server) servePage(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve the HTML file from the frontend pages directory
		http.ServeFile(w, r, "frontend/html/pages/"+filename)
	}
}
