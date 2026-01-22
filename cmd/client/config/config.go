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
	readHeaderTimeout = 5
	readTimeout       = 10
	writeTimeout      = 20
	idleTimeout       = 30
)

var (
	errMissingClientHost    = errors.New("missing CLIENT_HOST in config")
	errClientPortNotInteger = errors.New("invalid CLIENT_PORT: must be integer")
)

type Client struct {
	Host         string
	Port         string
	Environment  string
	BackendURL   string
	TLSCertFile  string
	TLSKeyFile   string
	HTTPTimeouts HTTPTimeouts
}

type HTTPTimeouts struct {
	ReadHeader time.Duration
	Read       time.Duration
	Write      time.Duration
	Idle       time.Duration
}

func LoadClientConfig() (*Client, error) {
	resolver := path.NewResolver()
	envFile, _ := os.ReadFile(resolver.GetPath(".env"))
	envMap := helpers.ParseEnv(string(envFile))

	tlsCertFile := helpers.GetEnv("CLIENT_TLS_CERT_FILE", envMap, "")
	tlsKeyFile := helpers.GetEnv("CLIENT_TLS_KEY_FILE", envMap, "")

	// Determine default backend URL based on whether certs exist
	defaultBackendURL := "http://localhost:8080/api/v1"
	if tlsCertFile != "" && tlsKeyFile != "" {
		// Check if cert files actually exist
		if _, err := os.Stat(resolver.GetPath(tlsCertFile)); err == nil {
			if _, err := os.Stat(resolver.GetPath(tlsKeyFile)); err == nil {
				defaultBackendURL = "https://localhost:8080/api/v1"
			}
		}
	}

	client := &Client{
		Host:        helpers.GetEnv("CLIENT_HOST", envMap, "localhost"),
		Port:        helpers.GetEnv("CLIENT_PORT", envMap, "3001"),
		Environment: helpers.GetEnv("CLIENT_ENVIRONMENT", envMap, "development"),
		BackendURL:  helpers.GetEnv("BACKEND_URL", envMap, defaultBackendURL),
		TLSCertFile: tlsCertFile,
		TLSKeyFile:  tlsKeyFile,
		HTTPTimeouts: HTTPTimeouts{
			ReadHeader: helpers.GetEnvDuration("CLIENT_READ_HEADER_TIMEOUT", envMap, readHeaderTimeout),
			Read:       helpers.GetEnvDuration("CLIENT_READ_TIMEOUT", envMap, readTimeout),
			Write:      helpers.GetEnvDuration("CLIENT_WRITE_TIMEOUT", envMap, writeTimeout),
			Idle:       helpers.GetEnvDuration("CLIENT_IDLE_TIMEOUT", envMap, idleTimeout),
		},
	}

	if client.Host == "" {
		return nil, errMissingClientHost
	}
	_, err := strconv.Atoi(strings.TrimPrefix(client.Port, ":"))
	if err != nil {
		return nil, errClientPortNotInteger
	}

	return client, nil
}
