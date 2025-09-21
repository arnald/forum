/*
Forum Server Application

This is the main entry point for the forum web application server.
The application follows Clean Architecture principles with the following layers:

1. Domain Layer: Core business entities and interfaces
2. Application Layer: Use cases and business logic
3. Infrastructure Layer: External concerns (HTTP, database, etc.)

The server provides:
- RESTful API endpoints for forum operations
- User authentication and session management
- Post creation, reading, updating with categories
- Comment system with threaded replies
- Voting system for posts and comments
- Static file serving for frontend assets

Technology Stack:
- Go for backend API
- SQLite for data persistence
- Vanilla HTML/CSS/JavaScript for frontend
- bcrypt for password hashing
- UUID for entity identification
- Cookie-based session management

To run the server:
  go run ./cmd/server/main.go

The server will start on the configured host:port (default localhost:8080)
*/
package main

import (
	"log"
	"os"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

// main is the application entry point that orchestrates server startup
// It follows the dependency injection pattern to wire up all components
func main() {
	// Step 1: Load application configuration
	// This includes database settings, server host/port, timeouts, etc.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Step 2: Initialize database connection
	// Sets up SQLite database with schema migrations and seeding
	db, err := sqlite.InitializeDB(*cfg)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close() // Ensure database connection is closed on exit

	// Step 3: Initialize infrastructure layer components
	// Create repository layer that handles data persistence
	userRepo := sqlite.NewRepo(db)

	// Initialize application logger for debugging and monitoring
	logger := logger.New(os.Stdout, logger.LevelInfo)

	// Create infrastructure providers (repositories for all entities)
	infraProviders := infra.NewInfraProviders(userRepo.DB)

	// Step 4: Initialize application services layer
	// Wire up business logic services with repository dependencies
	appServices := app.NewServices(
		infraProviders.UserRepository,
		infraProviders.PostRepository,
		infraProviders.CommentRepository,
		infraProviders.VoteRepository,
		infraProviders.CategoryRepository,
	)

	// Step 5: Initialize HTTP server with all dependencies
	// Create HTTP server with routes, middleware, and handlers
	infraHTTPServer := infra.NewHTTPServer(cfg, db, logger, appServices)

	// Step 6: Start the HTTP server
	// This blocks until the server is stopped or encounters an error
	infraHTTPServer.ListenAndServe()
}
