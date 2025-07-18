package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	readTimeout  = 5
	writeTimeout = 10
	idleTimeout  = 15
	configParts  = 2
)

var (
	ErrMissingServerHost    = errors.New("missing SERVER_HOST in config")
	ErrServerPortNotInteger = errors.New("invalid SERVER_PORT: must be integer")
)

type ServerConfig struct {
	Host         string
	Port         string
	Environment  string
	APIContext   string
	Database     DatabaseConfig
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Driver         string
	Path           string
	Pragma         string
	MigrateOnStart bool
	SeedOnStart    bool
	OpenConn       int
}

func LoadConfig() (*ServerConfig, error) {
	envFile, _ := os.ReadFile("../../.env")
	envMap := parseEnv(string(envFile))

	cfg := &ServerConfig{
		Host:         getEnv("SERVER_HOST", envMap, "localhost"),
		Port:         getEnv("SERVER_PORT", envMap, "8080"),
		Environment:  getEnv("SERVER_ENVIRONMENT", envMap, "development"),
		APIContext:   getEnv("API_CONTEXT", envMap, "/api/v1"),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", envMap, readTimeout),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", envMap, writeTimeout),
		IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", envMap, idleTimeout),
		Database: DatabaseConfig{
			Driver:         getEnv("DB_DRIVER", envMap, "sqlite3"),
			Path:           getEnv("DB_PATH", envMap, "./data/forum.db"),
			MigrateOnStart: getEnvBool("DB_MIGRATE_ON_START", envMap, true),
			SeedOnStart:    getEnvBool("DB_SEED_ON_START", envMap, true),
			Pragma:         getEnv("DB_PRAGMA", envMap, "_foreign_keys=on&_journal_mode=WAL"),
			OpenConn:       getEnvInt("DB_OPEN_CONN", envMap, 1),
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

func parseEnv(content string) map[string]string {
	env := make(map[string]string)
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", configParts)
		if len(parts) == configParts {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			env[key] = value
		}
	}
	return env
}

// Check OS environment -> .env file -> default values.
func getEnv(key string, envMap map[string]string, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	if val, exists := envMap[key]; exists {
		return val
	}

	return defaultValue
}

func getEnvDuration(key string, envMap map[string]string, defaultSeconds int) time.Duration {
	strValue := getEnv(key, envMap, "")
	if strValue == "" {
		return time.Duration(defaultSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(strValue)
	if err != nil {
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func getEnvBool(key string, envMap map[string]string, defaultValue bool) bool {
	strVal := getEnv(key, envMap, "")
	if strVal == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(strVal)
	if err != nil {
		return defaultValue
	}
	return b
}

func getEnvInt(s string, envMap map[string]string, defaultValue int) int {
	strVal := getEnv(s, envMap, "")
	if strVal == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(strVal)
	if err != nil {
		return defaultValue
	}
	return i
}
