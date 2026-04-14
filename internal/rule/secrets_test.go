package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestPlaintextSecrets(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db": {
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES_PASSWORD": "supersecret",
					"POSTGRES_USER":     "myapp",
					"POSTGRES_DB":       "mydb",
				},
			},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "plaintext-secrets")
	if len(found) != 1 {
		t.Fatalf("expected 1 plaintext-secrets finding, got %d", len(found))
	}
	if found[0].Service != "db" {
		t.Errorf("expected service db, got %s", found[0].Service)
	}
	if found[0].Severity != rule.Critical {
		t.Errorf("expected Critical severity, got %v", found[0].Severity)
	}
}

func TestPlaintextSecretsEnvRef(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db": {
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES_PASSWORD": "${DB_PASSWORD}",
					"API_TOKEN":         "$MY_TOKEN",
				},
			},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "plaintext-secrets")
	if len(found) != 0 {
		t.Fatalf("expected 0 plaintext-secrets findings for env refs, got %d", len(found))
	}
}

func TestPlaintextSecretsEmpty(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"db": {
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES_PASSWORD": "",
				},
			},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "plaintext-secrets")
	if len(found) != 0 {
		t.Fatalf("expected 0 findings for empty value, got %d", len(found))
	}
}
