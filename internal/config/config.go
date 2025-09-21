/*
Package config provides application configuration management for the forum.

This package handles:
- Loading configuration from environment variables and .env files
- Setting up default values for all configuration options
- Validating configuration parameters
- Providing typed configuration structures for all components

Configuration is organized into logical groups:
- Server configuration (host, port, timeouts)
- Database configuration (driver, path, migrations)
- Session management configuration (cookies, security)
- Handler timeout configuration

The configuration system supports environment variable overrides and
provides sensible defaults for development environments.
*/
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/path"
)

// Default configuration values in seconds
const (
	readTimeout         = 5      // HTTP read timeout in seconds
	writeTimeout        = 10     // HTTP write timeout in seconds
	idleTimeout         = 15     // HTTP idle timeout in seconds
	configParts         = 2      // Expected parts in key=value configuration
	defaultExpiry       = 86400  // Session expiry: 24 hours in seconds
	cleanupInternal     = 3600   // Session cleanup interval: 1 hour in seconds
	maxSessionsPerUser  = 5      // Maximum concurrent sessions per user
	sessionIDLenght     = 32     // Length of generated session IDs
	userRegisterTimeout = 15     // User registration handler timeout in seconds
	refreshTokenExpiry  = 30     // Refresh token expiry in days
	userLoginTimeout    = 15     // User login handler timeout in seconds
)

// Configuration validation errors
var (
	ErrMissingServerHost    = errors.New("missing SERVER_HOST in config")
	ErrServerPortNotInteger = errors.New("invalid SERVER_PORT: must be integer")
)

// ServerConfig contains all configuration for the forum server
// This is the main configuration structure that aggregates all other configs
type ServerConfig struct {
	Host           string                // Server host address (e.g., "localhost", "0.0.0.0")
	Port           string                // Server port (e.g., "8080")
	Environment    string                // Environment name (development, production, testing)
	APIContext     string                // Base path for API endpoints (e.g., "/api/v1")
	Database       DatabaseConfig        // Database configuration settings
	SessionManager SessionManagerConfig  // Session and cookie management settings
	Timeouts       TimeoutsConfig        // Various timeout configurations
	ReadTimeout    time.Duration         // HTTP server read timeout
	WriteTimeout   time.Duration         // HTTP server write timeout
	IdleTimeout    time.Duration         // HTTP server idle connection timeout
}

// DatabaseConfig contains all database-related configuration
// Supports SQLite with configurable migrations and seeding
type DatabaseConfig struct {
	Driver         string // Database driver name (e.g., "sqlite3")
	Path           string // Database file path (e.g., "data/forum.db")
	Pragma         string // SQLite pragma settings for performance and integrity
	MigrateOnStart bool   // Whether to run database migrations on startup
	SeedOnStart    bool   // Whether to seed database with initial data on startup
	OpenConn       int    // Maximum number of open database connections
}

// SessionManagerConfig contains all session and cookie management settings
// Handles user authentication state and security policies
type SessionManagerConfig struct {
	CookieName         string        // Name of the session cookie
	CookiePath         string        // Cookie path scope (usually "/")
	CookieDomain       string        // Cookie domain scope (empty for current domain)
	SameSite           string        // SameSite cookie attribute for CSRF protection
	DefaultExpiry      time.Duration // How long sessions last before expiring
	CleanupInterval    time.Duration // How often to clean up expired sessions
	MaxSessionsPerUser int           // Maximum concurrent sessions per user
	SessionIDLength    int           // Length of generated session identifiers
	SecureCookie       bool          // Whether to set Secure flag (HTTPS only)
	HTTPOnlyCookie     bool          // Whether to set HttpOnly flag (no JS access)
	EnablePersistence  bool          // Whether sessions persist across server restarts
	LogSessions        bool          // Whether to log session operations for debugging
	RefreshTokenExpiry time.Duration // How long refresh tokens remain valid
}

// TimeoutsConfig groups various timeout configurations
// Helps prevent slow operations from blocking the server
type TimeoutsConfig struct {
	HandlerTimeouts  HandlerTimeoutsConfig  // HTTP handler timeout settings
	UseCasesTimeouts UseCasesTimeoutsConfig // Business logic timeout settings
}

// HandlerTimeoutsConfig contains timeout settings for HTTP handlers
// Prevents slow HTTP requests from blocking other requests
type HandlerTimeoutsConfig struct {
	UserRegister time.Duration // Timeout for user registration requests
	UserLogin    time.Duration // Timeout for user login requests
}

// UseCasesTimeoutsConfig contains timeout settings for business logic operations
// Currently minimal but extensible for future use cases
type UseCasesTimeoutsConfig struct {
	UserRegister time.Duration // Timeout for user registration business logic
}

func LoadConfig() (*ServerConfig, error) {
	resolver := path.NewResolver()
	envFile, _ := os.ReadFile(resolver.GetPath(".env"))
	envMap := helpers.ParseEnv(string(envFile))

	cfg := &ServerConfig{
		Host:         helpers.GetEnv("SERVER_HOST", envMap, "localhost"),
		Port:         helpers.GetEnv("SERVER_PORT", envMap, "8080"),
		Environment:  helpers.GetEnv("SERVER_ENVIRONMENT", envMap, "development"),
		APIContext:   helpers.GetEnv("API_CONTEXT", envMap, "/api/v1"),
		ReadTimeout:  helpers.GetEnvDuration("SERVER_READ_TIMEOUT", envMap, readTimeout),
		WriteTimeout: helpers.GetEnvDuration("SERVER_WRITE_TIMEOUT", envMap, writeTimeout),
		IdleTimeout:  helpers.GetEnvDuration("SERVER_IDLE_TIMEOUT", envMap, idleTimeout),
		Database: DatabaseConfig{
			Driver:         helpers.GetEnv("DB_DRIVER", envMap, "sqlite3"),
			Path:           resolver.GetPath(helpers.GetEnv("DB_PATH", envMap, "data/forum.db")),
			MigrateOnStart: helpers.GetEnvBool("DB_MIGRATE_ON_START", envMap, true),
			SeedOnStart:    helpers.GetEnvBool("DB_SEED_ON_START", envMap, true),
			Pragma:         helpers.GetEnv("DB_PRAGMA", envMap, "_foreign_keys=on&_journal_mode=WAL"),
			OpenConn:       helpers.GetEnvInt("DB_OPEN_CONN", envMap, 1),
		},
		SessionManager: SessionManagerConfig{
			DefaultExpiry:      helpers.GetEnvDuration("SESSION_DEFAULT_EXPIRY", envMap, defaultExpiry),
			SecureCookie:       helpers.GetEnvBool("SESSION_SECURE_COOKIE", envMap, false),
			CookieName:         helpers.GetEnv("SESSION_COOKIE_NAME", envMap, "session_id"),
			CookiePath:         helpers.GetEnv("SESSION_COOKIE_PATH", envMap, "/"),
			CookieDomain:       helpers.GetEnv("SESSION_COOKIE_DOMAIN", envMap, ""),
			HTTPOnlyCookie:     helpers.GetEnvBool("SESSION_HTTPONLY_COOKIE", envMap, true),
			SameSite:           helpers.GetEnv("SESSION_SAMESITE", envMap, "Lax"),
			CleanupInterval:    helpers.GetEnvDuration("SESSION_CLEANUP_INTERVAL", envMap, cleanupInternal),
			MaxSessionsPerUser: helpers.GetEnvInt("SESSION_MAX_SESSIONS_PER_USER", envMap, maxSessionsPerUser),
			SessionIDLength:    helpers.GetEnvInt("SESSION_ID_LENGTH", envMap, sessionIDLenght),
			EnablePersistence:  helpers.GetEnvBool("SESSION_ENABLE_PERSISTENCE", envMap, true),
			LogSessions:        helpers.GetEnvBool("SESSION_LOG_SESSIONS", envMap, false),
			RefreshTokenExpiry: helpers.GetEnvDuration("SESSION_REFRESH_TOKEN_EXPIRY", envMap, refreshTokenExpiry),
		},
		Timeouts: TimeoutsConfig{
			HandlerTimeouts: HandlerTimeoutsConfig{
				UserRegister: helpers.GetEnvDuration("HANDLER_TIMEOUT_REGISTER", envMap, userRegisterTimeout),
				UserLogin:    helpers.GetEnvDuration("HANDLER_TIMEOUT_LOGIN", envMap, userLoginTimeout),
			},
		},
	}

	if cfg.Host == "" {
		return nil, ErrMissingServerHost
	}
	_, err := strconv.Atoi(strings.TrimPrefix(cfg.Port, ":"))
	if err != nil {
		return nil, ErrServerPortNotInteger
	}

	return cfg, nil
}
