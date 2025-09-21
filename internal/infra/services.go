/*
Package infra provides infrastructure layer services and dependency injection.

This package implements the dependency injection pattern for the forum application,
organizing infrastructure components into reusable services. It handles:

- Repository implementations for data access layer
- HTTP server configuration and initialization
- Service aggregation for clean dependency management
- Infrastructure provider factory functions

The infrastructure layer implements interfaces defined in the domain layer,
providing concrete implementations for databases, web servers, external services,
and other infrastructure concerns.
*/
package infra

import (
	"database/sql"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/post"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/http"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

// Services aggregates all infrastructure services for dependency injection
// This struct provides a centralized way to manage infrastructure dependencies
// and makes it easy to pass them to application and domain layers
type Services struct {
	UserRepository     user.Repository     // Repository for user data operations
	PostRepository     post.Repository     // Repository for post data operations
	CommentRepository  comment.Repository  // Repository for comment data operations
	VoteRepository     vote.Repository     // Repository for vote data operations
	CategoryRepository category.Repository // Repository for category data operations
	Server             *http.Server        // HTTP server instance for web requests
}

// NewInfraProviders creates and configures all infrastructure repository services
// This function implements the factory pattern to initialize all data repositories
// using a single SQLite database connection for consistency and performance
//
// Parameters:
//   - db: SQLite database connection for all repository operations
//
// Returns:
//   - Services: Configured infrastructure services with repository implementations
//
// Note: All repositories use the same underlying SQLite implementation (repo)
// which provides a unified data access layer with consistent transaction handling
func NewInfraProviders(db *sql.DB) Services {
	repo := sqlite.NewRepo(db) // Create single SQLite repository instance
	return Services{
		UserRepository:     repo, // User data operations (registration, authentication)
		PostRepository:     repo, // Post data operations (CRUD, voting, categories)
		CommentRepository:  repo, // Comment data operations (CRUD, threading)
		VoteRepository:     repo, // Vote data operations (like/dislike tracking)
		CategoryRepository: repo, // Category data operations (organization)
	}
}

// NewHTTPServer creates and configures the HTTP server with all dependencies
// This function sets up the web server with proper configuration, logging,
// database access, and application services for handling HTTP requests
//
// Parameters:
//   - cfg: Server configuration (host, port, timeouts, security settings)
//   - db: Database connection for session management and data operations
//   - logger: Logger instance for request/error logging throughout the server
//   - appServices: Application layer services containing business logic handlers
//
// Returns:
//   - *http.Server: Configured HTTP server ready to handle forum requests
//
// The server includes:
//   - Middleware for authentication, CORS, and logging
//   - Route handlers for all forum operations (users, posts, comments, votes)
//   - Session management for user authentication
//   - Static file serving for frontend assets
func NewHTTPServer(cfg *config.ServerConfig, db *sql.DB, logger logger.Logger, appServices app.Services) *http.Server {
	return http.NewServer(cfg, db, logger, appServices)
}
