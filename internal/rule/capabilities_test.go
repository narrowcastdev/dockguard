package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestPrivilegedMode(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"nginx": {Image: "nginx", Privileged: true},
			"redis": {Image: "redis"},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "privileged-mode")
	if len(found) != 1 {
		t.Fatalf("expected 1 privileged-mode finding, got %d", len(found))
	}
	if found[0].Service != "nginx" {
		t.Errorf("expected service nginx, got %s", found[0].Service)
	}
	if found[0].Severity != rule.Critical {
		t.Errorf("expected Critical severity, got %v", found[0].Severity)
	}
}

func TestNoCapDrop(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"postgres": {Image: "postgres:16"},
			"redis":    {Image: "redis", CapDrop: []string{"ALL"}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-cap-drop")
	if len(found) != 1 {
		t.Fatalf("expected 1 no-cap-drop finding, got %d", len(found))
	}
	if found[0].Service != "postgres" {
		t.Errorf("expected service postgres, got %s", found[0].Service)
	}
}

func TestNoNewPrivileges(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"app": {Image: "myapp"},
			"db":  {Image: "postgres", SecurityOpt: []string{"no-new-privileges"}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-new-privileges")
	if len(found) != 1 {
		t.Fatalf("expected 1 no-new-privileges finding, got %d", len(found))
	}
	if found[0].Service != "app" {
		t.Errorf("expected service app, got %s", found[0].Service)
	}
}

func findByRule(findings []rule.Finding, ruleID string) []rule.Finding {
	var result []rule.Finding
	for _, f := range findings {
		if f.Rule == ruleID {
			result = append(result, f)
		}
	}
	return result
}
