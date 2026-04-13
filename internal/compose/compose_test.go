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

func TestParseHealthcheckStringForm(t *testing.T) {
	input := []byte(`
services:
  redis:
    image: redis:7
    healthcheck:
      test: redis-cli ping || exit 1
      interval: 30s
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	redis := f.Services["redis"]
	if redis.Healthcheck == nil {
		t.Fatal("expected healthcheck to be set")
	}
	if len(redis.Healthcheck.Test) != 2 {
		t.Fatalf("expected 2 test elements, got %d: %v", len(redis.Healthcheck.Test), redis.Healthcheck.Test)
	}
	if redis.Healthcheck.Test[0] != "CMD-SHELL" {
		t.Errorf("expected Test[0]=CMD-SHELL, got %s", redis.Healthcheck.Test[0])
	}
	if redis.Healthcheck.Test[1] != "redis-cli ping || exit 1" {
		t.Errorf("expected Test[1]='redis-cli ping || exit 1', got %s", redis.Healthcheck.Test[1])
	}
}

func TestParseHealthcheckListForm(t *testing.T) {
	input := []byte(`
services:
  postgres:
    image: postgres:16
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pg := f.Services["postgres"]
	if pg.Healthcheck == nil {
		t.Fatal("expected healthcheck to be set")
	}
	if len(pg.Healthcheck.Test) != 2 {
		t.Fatalf("expected 2 test elements, got %d", len(pg.Healthcheck.Test))
	}
	if pg.Healthcheck.Test[0] != "CMD-SHELL" {
		t.Errorf("expected Test[0]=CMD-SHELL, got %s", pg.Healthcheck.Test[0])
	}
	if pg.Healthcheck.Test[1] != "pg_isready -U postgres" {
		t.Errorf("expected Test[1]='pg_isready -U postgres', got %s", pg.Healthcheck.Test[1])
	}
}

func TestParseHealthcheckDisableFalse(t *testing.T) {
	input := []byte(`
services:
  app:
    image: myapp
    healthcheck:
      disable: false
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := f.Services["app"]
	if app.Healthcheck == nil {
		t.Fatal("expected healthcheck to be set")
	}
	if len(app.Healthcheck.Test) != 0 {
		t.Errorf("expected empty test, got %v", app.Healthcheck.Test)
	}
}

func TestParseNetworksMapForm(t *testing.T) {
	input := []byte(`
services:
  wireguard:
    image: ghcr.io/wg-easy/wg-easy
    networks:
      wg:
        ipv4_address: 10.42.42.42
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wg := f.Services["wireguard"]
	if len(wg.Networks) != 1 {
		t.Fatalf("expected 1 network, got %d", len(wg.Networks))
	}
	if wg.Networks[0] != "wg" {
		t.Errorf("expected network 'wg', got %s", wg.Networks[0])
	}
}

func TestParseNetworksListForm(t *testing.T) {
	input := []byte(`
services:
  app:
    image: myapp
    networks:
      - frontend
      - backend
`)
	f, err := compose.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := f.Services["app"]
	if len(app.Networks) != 2 {
		t.Fatalf("expected 2 networks, got %d", len(app.Networks))
	}
}

func TestParseAdvancedFixture(t *testing.T) {
	f, err := compose.ParseFile("../../testdata/advanced.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Services) != 5 {
		t.Fatalf("expected 5 services, got %d", len(f.Services))
	}

	// String healthcheck
	redis := f.Services["redis"]
	if redis.Healthcheck == nil || redis.Healthcheck.Test[0] != "CMD-SHELL" {
		t.Error("redis healthcheck string form not normalized")
	}

	// List healthcheck
	pg := f.Services["postgres"]
	if pg.Healthcheck == nil || pg.Healthcheck.Test[0] != "CMD-SHELL" {
		t.Error("postgres healthcheck list form not preserved")
	}

	// disable: false
	immich := f.Services["immich-server"]
	if immich.Healthcheck == nil {
		t.Error("immich-server healthcheck should not be nil")
	}

	// Map networks
	wg := f.Services["wireguard"]
	if len(wg.Networks) != 1 || wg.Networks[0] != "wg" {
		t.Errorf("wireguard networks not normalized: %v", wg.Networks)
	}

	// List networks
	app := f.Services["app"]
	if len(app.Networks) != 2 {
		t.Errorf("app networks not preserved: %v", app.Networks)
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
