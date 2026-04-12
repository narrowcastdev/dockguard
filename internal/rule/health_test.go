package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestNoHealthcheck(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db": {Image: "postgres"},
			"cache": {Image: "redis", Healthcheck: &compose.Healthcheck{
				Test: []string{"CMD", "redis-cli", "ping"},
			}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-healthcheck")
	if len(found) != 1 {
		t.Fatalf("expected 1 no-healthcheck finding, got %d", len(found))
	}
	if found[0].Service != "db" {
		t.Errorf("expected service db, got %s", found[0].Service)
	}
	if found[0].Severity != rule.Info {
		t.Errorf("expected Info severity, got %v", found[0].Severity)
	}
}

func TestRestartNoHealth(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web": {Image: "nginx", Restart: "always"},
			"db": {Image: "postgres", Restart: "unless-stopped", Healthcheck: &compose.Healthcheck{
				Test: []string{"CMD-SHELL", "pg_isready"},
			}},
			"worker": {Image: "myapp", Restart: "on-failure"},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "restart-no-health")
	if len(found) != 1 {
		t.Fatalf("expected 1 restart-no-health finding, got %d", len(found))
	}
	if found[0].Service != "web" {
		t.Errorf("expected service web, got %s", found[0].Service)
	}
}
