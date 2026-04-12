package compose_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
)

func TestParseBasicFile(t *testing.T) {
	f, err := compose.ParseFile("../../testdata/basic.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(f.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(f.Services))
	}

	pg := f.Services["postgres"]
	if pg == nil {
		t.Fatal("postgres service not found")
	}
	if pg.Image != "postgres:16" {
		t.Errorf("expected image postgres:16, got %s", pg.Image)
	}
	if pg.Environment["POSTGRES_PASSWORD"] != "supersecret" {
		t.Errorf("expected POSTGRES_PASSWORD=supersecret, got %s", pg.Environment["POSTGRES_PASSWORD"])
	}
	if pg.Restart != "always" {
		t.Errorf("expected restart=always, got %s", pg.Restart)
	}

	nginx := f.Services["nginx"]
	if !nginx.Privileged {
		t.Error("expected nginx to be privileged")
	}
	if len(nginx.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(nginx.Ports))
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := compose.ParseFile("nonexistent.yml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := compose.Parse([]byte("{{invalid"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseNoServices(t *testing.T) {
	_, err := compose.Parse([]byte("version: '3'\n"))
	if err == nil {
		t.Fatal("expected error for compose file with no services")
	}
}

func TestParseEnvironmentListFormat(t *testing.T) {
	input := []byte(`
services:
  app:
    image: myapp
    environment:
      - FOO=bar
      - BAZ=qux
      - DEBUG
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := f.Services["app"]
	if app.Environment["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", app.Environment["FOO"])
	}
	if app.Environment["DEBUG"] != "" {
		t.Errorf("expected DEBUG='', got %s", app.Environment["DEBUG"])
	}
}
