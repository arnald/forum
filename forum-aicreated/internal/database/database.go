// Package database handles all database operations for the forum application.
// It uses SQLite as the database engine with proper schema design and relationships.
package database

import (
	"database/sql"
	"os"

	// Import SQLite driver for database/sql
	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the standard sql.DB to provide custom methods for our forum
// It embeds *sql.DB so all standard database methods are available
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection and verifies it works
// dataSourceName is the path to the SQLite database file
func NewDB(dataSourceName string) (*DB, error) {
	// Open connection to SQLite database
	// The sqlite3 driver is registered via the blank import above
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Verify the connection actually works
	// Ping sends a simple query to check connectivity
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Return our custom DB wrapper
	return &DB{db}, nil
}

// Init initializes the database schema by creating all necessary tables
// This is called once when the application starts
func (db *DB) Init() error {
	return db.createTables()
}

// createTables creates all necessary database tables and relationships
// This function defines the entire database schema for the forum application
func (db *DB) createTables() error {
	query := `
	-- ===== USERS TABLE =====
	-- Stores user account information and authentication data
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,    -- Unique user identifier
		email TEXT UNIQUE NOT NULL,             -- User's email (must be unique)
		username TEXT UNIQUE NOT NULL,          -- Display name (must be unique)
		password TEXT NOT NULL,                 -- Bcrypt hashed password
		role TEXT DEFAULT 'user',               -- User role: 'user', 'moderator', 'admin'
		provider TEXT DEFAULT '',               -- OAuth provider ('google', 'github', or empty for local)
		provider_id TEXT DEFAULT '',            -- OAuth provider user ID
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- ===== SESSIONS TABLE =====
	-- Manages user login sessions for authentication
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,                    -- UUID session identifier stored in cookies
		user_id INTEGER NOT NULL,               -- Links to users table
		expires_at DATETIME NOT NULL,           -- When this session expires
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	-- ===== CATEGORIES TABLE =====
	-- Stores post categories for organizing content
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL               -- Category name must be unique
	);

	-- ===== POSTS TABLE =====
	-- Main content storage for forum posts
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,               -- Author of the post
		title TEXT NOT NULL,                    -- Post title
		content TEXT NOT NULL,                  -- Post content/body
		image_path TEXT DEFAULT '',             -- Optional uploaded image path
		status TEXT DEFAULT 'approved',         -- 'pending', 'approved', 'rejected' for moderation
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	-- ===== POST_CATEGORIES TABLE =====
	-- Many-to-many relationship between posts and categories
	-- One post can belong to multiple categories
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),     -- Composite primary key prevents duplicates
		FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
	);

	-- ===== COMMENTS TABLE =====
	-- User comments on posts
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,               -- Which post this comment belongs to
		user_id INTEGER NOT NULL,               -- Who wrote the comment
		content TEXT NOT NULL,                  -- Comment text
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	-- ===== POST_LIKES TABLE =====
	-- Tracks likes and dislikes on posts
	CREATE TABLE IF NOT EXISTS post_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		is_like BOOLEAN NOT NULL,               -- TRUE = like, FALSE = dislike
		UNIQUE(post_id, user_id),              -- One vote per user per post
		FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	-- ===== COMMENT_LIKES TABLE =====
	-- Tracks likes and dislikes on comments
	CREATE TABLE IF NOT EXISTS comment_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		is_like BOOLEAN NOT NULL,               -- TRUE = like, FALSE = dislike
		UNIQUE(comment_id, user_id),           -- One vote per user per comment
		FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);

	-- ===== NOTIFICATIONS TABLE =====
	-- Real-time notification system for user interactions
	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,               -- Who receives the notification
		actor_id INTEGER NOT NULL,              -- Who triggered the notification
		type TEXT NOT NULL,                     -- 'like', 'dislike', 'comment'
		post_id INTEGER,                        -- Optional: related post
		comment_id INTEGER,                     -- Optional: related comment
		message TEXT NOT NULL,                  -- Human-readable notification message
		read BOOLEAN DEFAULT FALSE,             -- Whether user has seen it
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (actor_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
	);

	-- ===== REPORTS TABLE =====
	-- Content reporting system for moderation
	CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		reporter_id INTEGER NOT NULL,           -- Who filed the report
		post_id INTEGER,                        -- Optional: reported post
		comment_id INTEGER,                     -- Optional: reported comment
		reason TEXT NOT NULL,                   -- Reason for reporting
		description TEXT NOT NULL,              -- Detailed description
		status TEXT DEFAULT 'pending',          -- 'pending', 'reviewed', 'resolved'
		admin_id INTEGER,                       -- Admin who handled the report
		admin_notes TEXT DEFAULT '',            -- Admin's notes on the report
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reporter_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
		FOREIGN KEY (admin_id) REFERENCES users (id) ON DELETE SET NULL
	);

	-- ===== ROLE_REQUESTS TABLE =====
	-- System for users to request role upgrades (user -> moderator)
	CREATE TABLE IF NOT EXISTS role_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,               -- Who is requesting the role change
		requested_role TEXT NOT NULL,           -- Which role they want
		reason TEXT NOT NULL,                   -- Why they want the role
		status TEXT DEFAULT 'pending',          -- 'pending', 'approved', 'rejected'
		admin_id INTEGER,                       -- Admin who handled the request
		admin_notes TEXT DEFAULT '',            -- Admin's decision notes
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (admin_id) REFERENCES users (id) ON DELETE SET NULL
	);

	-- ===== RATE_LIMITS TABLE =====
	-- Rate limiting system to prevent spam and abuse
	CREATE TABLE IF NOT EXISTS rate_limits (
		ip TEXT PRIMARY KEY,                    -- Client IP address
		requests INTEGER DEFAULT 0,            -- Number of requests made
		reset_time DATETIME NOT NULL           -- When the rate limit resets
	);

	-- ===== INITIAL DATA =====
	-- Insert default categories (INSERT OR IGNORE prevents duplicates)
	INSERT OR IGNORE INTO categories (name) VALUES
		('General'), ('Technology'), ('Sports'), ('Entertainment'), ('Science'), ('Politics');

	-- Insert default admin user
	-- Password is 'admin123' hashed with bcrypt (cost 10)
	INSERT OR IGNORE INTO users (email, username, password, role) VALUES
		('admin@forum.com', 'admin', '$2a$10$CiTwPi/Ovih4fpVMKaERlu.OKFXiN6OLuYV3ZlnaQY0JxudmnDzQC', 'admin');
	`

	// Execute the entire schema creation as one transaction
	_, err := db.Exec(query)
	return err
}

// GetDBPath returns the database file path from environment variable or default
// This allows for different database locations in development vs production
func GetDBPath() string {
	// Check if DB_PATH environment variable is set
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Use default path if environment variable is not set
		dbPath = "./data/forum.db"
	}
	return dbPath
}
