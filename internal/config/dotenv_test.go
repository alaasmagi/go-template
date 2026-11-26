package config

import (
	"os"
	"path/filepath"
	"testing"
)

// unsetEnv clears key for the test and restores its original state afterwards.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}

func writeDotEnv(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeDotEnv(t, dir, `# comment
DOTENV_TEST_PLAIN=plain

  DOTENV_TEST_SPACED = spaced value
DOTENV_TEST_DOUBLE="double quoted"
DOTENV_TEST_SINGLE='single quoted'
DOTENV_TEST_MISMATCHED="mismatched'
DOTENV_TEST_EMPTY=
DOTENV_TEST_EQUALS=a=b
DOTENV_TEST_EXISTING=from-file
not a variable
=no-key
`)

	keys := []string{
		"DOTENV_TEST_PLAIN", "DOTENV_TEST_SPACED", "DOTENV_TEST_DOUBLE", "DOTENV_TEST_SINGLE",
		"DOTENV_TEST_MISMATCHED", "DOTENV_TEST_EMPTY", "DOTENV_TEST_EQUALS",
	}
	for _, key := range keys {
		unsetEnv(t, key)
	}
	t.Setenv("DOTENV_TEST_EXISTING", "from-process")

	if err := loadDotEnvFile(filepath.Join(dir, ".env")); err != nil {
		t.Fatalf("loadDotEnvFile() error: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{key: "DOTENV_TEST_PLAIN", want: "plain"},
		{key: "DOTENV_TEST_SPACED", want: "spaced value"},
		{key: "DOTENV_TEST_DOUBLE", want: "double quoted"},
		{key: "DOTENV_TEST_SINGLE", want: "single quoted"},
		{key: "DOTENV_TEST_MISMATCHED", want: `"mismatched'`},
		{key: "DOTENV_TEST_EMPTY", want: ""},
		{key: "DOTENV_TEST_EQUALS", want: "a=b"},
		{key: "DOTENV_TEST_EXISTING", want: "from-process"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := os.LookupEnv(tt.key)
			if !ok {
				t.Fatalf("%s not set", tt.key)
			}
			if got != tt.want {
				t.Errorf("%s = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestLoadDotEnvDoesNotOverrideEmptyProcessValue(t *testing.T) {
	dir := t.TempDir()
	writeDotEnv(t, dir, "DOTENV_TEST_SET_EMPTY=from-file\n")
	t.Setenv("DOTENV_TEST_SET_EMPTY", "")

	if err := loadDotEnvFile(filepath.Join(dir, ".env")); err != nil {
		t.Fatalf("loadDotEnvFile() error: %v", err)
	}

	if got := os.Getenv("DOTENV_TEST_SET_EMPTY"); got != "" {
		t.Errorf("DOTENV_TEST_SET_EMPTY = %q, want empty", got)
	}
}

func TestLoadDotEnvSearchesUpward(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	writeDotEnv(t, root, "DOTENV_TEST_UPWARD=found\n")
	unsetEnv(t, "DOTENV_TEST_UPWARD")
	t.Chdir(nested)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error: %v", err)
	}

	if got := os.Getenv("DOTENV_TEST_UPWARD"); got != "found" {
		t.Errorf("DOTENV_TEST_UPWARD = %q, want %q", got, "found")
	}
}

func TestFindDotEnv(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		envDir   string
		startDir string
		want     string
	}{
		{name: "in start directory", envDir: nested, startDir: nested, want: filepath.Join(nested, ".env")},
		{name: "in ancestor", envDir: root, startDir: nested, want: filepath.Join(root, ".env")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeDotEnv(t, tt.envDir, "")
			t.Cleanup(func() { os.Remove(filepath.Join(tt.envDir, ".env")) })

			got, ok := findDotEnv(tt.startDir)
			if !ok || got != tt.want {
				t.Errorf("findDotEnv() = %q, %v, want %q, true", got, ok, tt.want)
			}
		})
	}
}
