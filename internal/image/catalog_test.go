package image_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/image"
)

func TestNormalizeImage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"postgres:16", "postgres"},
		{"postgres:16-alpine", "postgres"},
		{"library/postgres:latest", "postgres"},
		{"redis:7-alpine", "redis"},
		{"nginx:latest", "nginx"},
		{"lscr.io/linuxserver/sonarr:latest", "sonarr"},
		{"ghcr.io/immich-app/immich-server:release", "immich-server"},
		{"myregistry.com:5000/custom/app:v2", "app"},
		{"unknown-app", "unknown-app"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := image.NormalizeImage(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeImage(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLookupKnownImage(t *testing.T) {
	profile := image.Lookup("postgres:16-alpine")
	if profile == nil {
		t.Fatal("expected profile for postgres")
	}
	if profile.User != "999:999" {
		t.Errorf("expected user 999:999, got %s", profile.User)
	}
	if profile.Healthcheck == nil {
		t.Error("expected healthcheck for postgres")
	}
	if profile.MemoryLimit != "512m" {
		t.Errorf("expected memory 512m, got %s", profile.MemoryLimit)
	}
}

func TestLookupUnknownImage(t *testing.T) {
	profile := image.Lookup("my-custom-app:latest")
	if profile != nil {
		t.Error("expected nil profile for unknown image")
	}
}
