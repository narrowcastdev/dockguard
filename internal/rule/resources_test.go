package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestNoMemoryLimit(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db": {Image: "postgres"},
			"cache": {Image: "redis", Deploy: &compose.Deploy{
				Resources: &compose.Resources{
					Limits: &compose.ResourceLimit{Memory: "256m"},
				},
			}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-memory-limit")
	if len(found) != 1 {
		t.Fatalf("expected 1 no-memory-limit finding, got %d", len(found))
	}
	if found[0].Service != "db" {
		t.Errorf("expected service db, got %s", found[0].Service)
	}
}

func TestNoCPULimit(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db":  {Image: "postgres"},
			"web": {Image: "nginx"},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "no-cpu-limit")
	if len(found) != 2 {
		t.Fatalf("expected 2 no-cpu-limit findings, got %d", len(found))
	}
}
