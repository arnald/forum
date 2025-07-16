package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/arnald/forum/internal/pkg/helpers"
)

var (
	errMissingClientHost    = errors.New("missing CLIENT_HOST in config")
	errClientPortNotInteger = errors.New("invalid CLIENT_PORT: must be integer")
)

type Client struct {
	Host        string
	Port        string
	Environment string
}

func LoadClientConfig() (*Client, error) {
	envFile, _ := os.ReadFile("../../.env")
	envMap := helpers.ParseEnv(string(envFile))

	client := &Client{
		Host:        helpers.GetEnv("CLIENT_HOST", envMap, "localhost"),
		Port:        helpers.GetEnv("CLIENT_PORT", envMap, "3000"),
		Environment: helpers.GetEnv("CLIENT_ENVIRONMENT", envMap, "development"),
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
