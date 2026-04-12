package rule

import (
	"fmt"
	"strings"

	"github.com/narrowcastdev/dockguard/internal/compose"
)

func init() {
	Register(PlaintextSecrets{})
}

var secretPatterns = []string{
	"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY", "APIKEY",
	"DB_PASS", "PRIVATE_KEY", "AUTH",
}

// PlaintextSecrets flags environment variables containing secrets as literal values.
type PlaintextSecrets struct{}

func (r PlaintextSecrets) ID() string { return "plaintext-secrets" }

func (r PlaintextSecrets) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		for key, val := range svc.Environment {
			if !looksLikeSecret(key) {
				continue
			}
			if val == "" || isEnvRef(val) {
				continue
			}
			findings = append(findings, Finding{
				Rule:       r.ID(),
				Service:    name,
				Severity:   Critical,
				Message:    fmt.Sprintf("%s contains plaintext secret", key),
				Suggestion: fmt.Sprintf("Use ${%s} with a .env file or Docker secrets", key),
			})
		}
	}
	return findings
}

func (r PlaintextSecrets) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		var secretKeys []string
		for key, val := range svc.Environment {
			if looksLikeSecret(key) && val != "" && !isEnvRef(val) {
				secretKeys = append(secretKeys, key)
			}
		}
		if len(secretKeys) == 0 {
			continue
		}
		svcName := name
		keys := make([]string, len(secretKeys))
		copy(keys, secretKeys)
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				for _, key := range keys {
					s.Environment[key] = fmt.Sprintf("${%s}", key)
				}
			},
		})
	}
	return fixes
}

func looksLikeSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range secretPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

func isEnvRef(val string) bool {
	return strings.HasPrefix(val, "${") || strings.HasPrefix(val, "$")
}
