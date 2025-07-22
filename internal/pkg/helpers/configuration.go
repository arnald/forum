package helpers

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const configParts = 2

func ParseEnv(content string) map[string]string {
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

// GetEnv Check OS environment -> .env file -> default values.
func GetEnv(key string, envMap map[string]string, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	if val, exists := envMap[key]; exists {
		return val
	}

	return defaultValue
}

func GetEnvDuration(key string, envMap map[string]string, defaultSeconds int) time.Duration {
	strValue := GetEnv(key, envMap, "")
	if strValue == "" {
		return time.Duration(defaultSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(strValue)
	if err != nil {
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func GetEnvBool(key string, envMap map[string]string, defaultValue bool) bool {
	strVal := GetEnv(key, envMap, "")
	if strVal == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(strVal)
	if err != nil {
		return defaultValue
	}
	return b
}

func GetEnvInt(s string, envMap map[string]string, defaultValue int) int {
	strVal := GetEnv(s, envMap, "")
	if strVal == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(strVal)
	if err != nil {
		return defaultValue
	}
	return i
}
