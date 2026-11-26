package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadDotEnv loads the nearest .env file found upward from the working directory.
// Variables already set in the process are never overridden. A missing file is not an error.
func LoadDotEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	path, ok := findDotEnv(dir)
	if !ok {
		return nil
	}

	return loadDotEnvFile(path)
}

func findDotEnv(dir string) (string, bool) {
	for {
		candidate := filepath.Join(dir, ".env")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !found || key == "" {
			continue
		}

		if _, set := os.LookupEnv(key); set {
			continue
		}

		if err := os.Setenv(key, unquote(strings.TrimSpace(value))); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func unquote(value string) string {
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return value[1 : len(value)-1]
		}
	}

	return value
}
