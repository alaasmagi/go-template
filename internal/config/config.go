// Package config loads the application configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all settings the application needs, loaded once at startup.
type Config struct {
	AppPort int
	AppSlug string
}

// Load reads and validates the configuration from the process environment.
func Load() (Config, error) {
	port, err := requiredPositiveInt("APP_PORT")
	if err != nil {
		return Config{}, err
	}

	slug, err := required("APP_SLUG")
	if err != nil {
		return Config{}, err
	}

	return Config{AppPort: port, AppSlug: slug}, nil
}

func required(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("missing required configuration: set %s", key)
	}

	return value, nil
}

func requiredPositiveInt(key string) (int, error) {
	value, err := required(key)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}

	return parsed, nil
}
