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

const (
	readTimeout         = 5
	writeTimeout        = 10
	idleTimeout         = 15
	configParts         = 2
	defaultExpiry       = 86400
	cleanupInternal     = 3600
	maxSessionsPerUser  = 5
	sessionIDLenght     = 32
	userRegisterTimeout = 15
	refreshTokenExpiry  = 30
	userLoginTimeout    = 15
)

var (
	ErrMissingServerHost    = errors.New("missing SERVER_HOST in config")
	ErrServerPortNotInteger = errors.New("invalid SERVER_PORT: must be integer")
)

type ServerConfig struct {
	Host           string
	Port           string
	Environment    string
	APIContext     string
	Database       DatabaseConfig
	SessionManager SessionManagerConfig
	Timeouts       TimeoutsConfig
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
}

type DatabaseConfig struct {
	Driver         string
	Path           string
	Pragma         string
	MigrateOnStart bool
	SeedOnStart    bool
	OpenConn       int
}

type SessionManagerConfig struct {
	CookieName         string
	CookiePath         string
	CookieDomain       string
	SameSite           string
	DefaultExpiry      time.Duration
	CleanupInterval    time.Duration
	MaxSessionsPerUser int
	SessionIDLength    int
	SecureCookie       bool
	HTTPOnlyCookie     bool
	EnablePersistence  bool
	LogSessions        bool
	RefreshTokenExpiry time.Duration
}

type TimeoutsConfig struct {
	HandlerTimeouts  HandlerTimeoutsConfig
	UseCasesTimeouts UseCasesTimeoutsConfig
}

type HandlerTimeoutsConfig struct {
	UserRegister time.Duration
	UserLogin    time.Duration
}

type UseCasesTimeoutsConfig struct { // Not implemented yet, but can be used for future use cases
	UserRegister time.Duration
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
