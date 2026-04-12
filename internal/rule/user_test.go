package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestRunningAsRoot(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db":  {Image: "postgres"},
			"app": {Image: "myapp", User: "1000:1000"},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "running-as-root")
	if len(found) != 1 {
		t.Fatalf("expected 1 running-as-root finding, got %d", len(found))
	}
	if found[0].Service != "db" {
		t.Errorf("expected service db, got %s", found[0].Service)
	}
}

func TestNoReadOnly(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web":   {Image: "nginx"},
			"cache": {Image: "redis", ReadOnly: true},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-read-only")
	if len(found) != 1 {
		t.Fatalf("expected 1 no-read-only finding, got %d", len(found))
	}
	if found[0].Service != "web" {
		t.Errorf("expected service web, got %s", found[0].Service)
	}
}
