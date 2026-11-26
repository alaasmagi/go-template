package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		slug    string
		want    Config
		wantErr string
	}{
		{name: "valid", port: "8080", slug: "my-app", want: Config{AppPort: 8080, AppSlug: "my-app"}},
		{name: "trims whitespace", port: " 8080 ", slug: " my-app ", want: Config{AppPort: 8080, AppSlug: "my-app"}},
		{name: "missing port", port: "", slug: "my-app", wantErr: "APP_PORT"},
		{name: "non-numeric port", port: "abc", slug: "my-app", wantErr: "APP_PORT"},
		{name: "zero port", port: "0", slug: "my-app", wantErr: "APP_PORT"},
		{name: "negative port", port: "-1", slug: "my-app", wantErr: "APP_PORT"},
		{name: "missing slug", port: "8080", slug: "", wantErr: "APP_SLUG"},
		{name: "blank slug", port: "8080", slug: "   ", wantErr: "APP_SLUG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_PORT", tt.port)
			t.Setenv("APP_SLUG", tt.slug)

			got, err := Load()

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Load() error = %v, want error naming %s", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
