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
	server.router.HandleFunc("/", server.servePage("home.html"))                        // Home page
	server.router.HandleFunc("/login", server.servePage("login.html"))                  // Login form
	server.router.HandleFunc("/register", server.servePage("register.html"))            // Registration form
	server.router.HandleFunc("/posts", server.servePage("posts.html"))                  // Posts listing page
	server.router.HandleFunc("/create-post", server.servePage("create-post.html"))      // Post creation form
	server.router.HandleFunc("/post/", server.servePage("post-detail.html"))            // Individual post view

	// Health endpoint
	server.router.HandleFunc(apiContext+"/health", health.NewHandler(server.logger).HealthCheck)

	// User authentication endpoints
	server.router.HandleFunc(
		apiContext+"/login/username",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginUsername,
	)
	server.router.HandleFunc(
		apiContext+"/login/email",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginEmail,
	)
	server.router.HandleFunc(
		apiContext+"/register",
		userRegister.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserRegister,
	)
	server.router.HandleFunc(
		apiContext+"/logout",
		userLogout.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLogout,
	)

	// Post endpoints
	createPostHandler := post.NewCreatePostHandler(server.appServices.PostServices.Queries.CreatePost)
	getPostsHandler := post.NewGetPostsHandler(
		server.appServices.PostServices.Queries.GetAllPosts,
		server.appServices.PostServices.Queries.GetPostByID,
		server.appServices.PostServices.Queries.GetPostsByCategory,
	)
	updatePostHandler := post.NewUpdatePostHandler(server.appServices.PostServices.Queries.UpdatePost)

	// Protected post endpoints (require authentication)
	server.router.Handle(apiContext+"/posts", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(createPostHandler.CreatePost)))
	server.router.Handle(apiContext+"/posts/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(updatePostHandler.UpdatePost)))

	// Public post endpoints
	server.router.HandleFunc(apiContext+"/posts/all", getPostsHandler.GetAllPosts)
	server.router.HandleFunc(apiContext+"/posts/get", getPostsHandler.GetPostByID)

	// Comment endpoints
	createCommentHandler := comment.NewCreateCommentHandler(server.appServices.CommentServices.Queries.CreateComment)
	getCommentsHandler := comment.NewGetCommentsHandler(
		server.appServices.CommentServices.Queries.GetCommentsByPost,
		server.appServices.CommentServices.Queries.GetCommentTree,
	)
	updateCommentHandler := comment.NewUpdateCommentHandler(server.appServices.CommentServices.Queries.UpdateComment)

	// Protected comment endpoints (require authentication)
	server.router.Handle(apiContext+"/comments", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(createCommentHandler.CreateComment)))
	server.router.Handle(apiContext+"/comments/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(updateCommentHandler.UpdateComment)))

	// Public comment endpoints
	server.router.HandleFunc(apiContext+"/comments/post", getCommentsHandler.GetCommentsByPost)
	server.router.HandleFunc(apiContext+"/comments/tree", getCommentsHandler.GetCommentTree)

	// Vote endpoints
	voteHandler := vote.NewVoteHandler(
		server.appServices.VoteServices.Queries.CastVote,
		server.appServices.VoteServices.Queries.GetVoteStatus,
		server.appServices.VoteServices.Queries.GetUserVotes,
	)

	// Protected vote endpoints (require authentication)
	server.router.Handle(apiContext+"/votes/cast", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(voteHandler.CastVote)))
	server.router.Handle(apiContext+"/votes/user", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(voteHandler.GetUserVotes)))

	// Public vote endpoints
	server.router.HandleFunc(apiContext+"/votes/status", voteHandler.GetVoteStatus)

	// Category endpoints
	categoryHandler := category.NewCategoryHandler(
		server.appServices.CategoryServices.Queries.CreateCategory,
		server.appServices.CategoryServices.Queries.GetAllCategories,
		server.appServices.CategoryServices.Queries.GetCategoryByID,
		server.appServices.CategoryServices.Queries.UpdateCategory,
		server.appServices.CategoryServices.Queries.DeleteCategory,
		server.appServices.CategoryServices.Queries.GetCategoryWithPosts,
	)

	// Protected category endpoints (require authentication)
	server.router.Handle(apiContext+"/categories", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.CreateCategory)))
	server.router.Handle(apiContext+"/categories/update", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.UpdateCategory)))
	server.router.Handle(apiContext+"/categories/delete", middleware.AuthMiddleware(server.sessionManager, server.logger)(http.HandlerFunc(categoryHandler.DeleteCategory)))

	// Public category endpoints
	server.router.HandleFunc(apiContext+"/categories/all", categoryHandler.GetAllCategories)
	server.router.HandleFunc(apiContext+"/categories/get", categoryHandler.GetCategoryByID)
	server.router.HandleFunc(apiContext+"/categories/posts", categoryHandler.GetCategoryWithPosts)
}

func (server *Server) ListenAndServe() {
	corsWrappedRouter := middleware.NewCorsMiddleware(server.router)

	srv := &http.Server{
		Addr:         server.config.Host + ":" + server.config.Port,
		Handler:      corsWrappedRouter,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
		IdleTimeout:  server.config.IdleTimeout,
	}
	server.logger.PrintInfo("Starting server", map[string]string{
		"host":        server.config.Host,
		"port":        server.config.Port,
		"environment": server.config.Environment,
	})
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.logger.PrintFatal(err, nil)
	}
}

func (server *Server) initSessionManager() {
	server.sessionManager = session.NewSessionManager(server.db, server.config.SessionManager)
}

func (server *Server) servePage(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/html/pages/"+filename)
	}
}
