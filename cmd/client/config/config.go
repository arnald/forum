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

	client := &Client{
		Host:        helpers.GetEnv("CLIENT_HOST", envMap, "localhost"),
		Port:        helpers.GetEnv("CLIENT_PORT", envMap, "3001"),
		Environment: helpers.GetEnv("CLIENT_ENVIRONMENT", envMap, "development"),
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
