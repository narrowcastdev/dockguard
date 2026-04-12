package rule

import (
	"github.com/narrowcastdev/dockguard/internal/compose"
)

// Severity levels for findings.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "info"
	case Warning:
		return "warning"
	case Critical:
		return "critical"
	}
	return "unknown"
}

// Finding represents a single security issue found by a rule.
type Finding struct {
	Rule       string
	Service    string
	Severity   Severity
	Message    string
	Suggestion string
}

// Fix represents a mutation to apply to a service in the compose model.
type Fix struct {
	Service string
	Apply   func(svc *compose.Service)
}

// Rule is the interface all security rules implement.
type Rule interface {
	ID() string
	Check(f *compose.File) []Finding
	Fix(f *compose.File) []Fix
}

var registry []Rule

// Register adds a rule to the global registry.
func Register(r Rule) {
	registry = append(registry, r)
}

// All returns all registered rules.
func All() []Rule {
	return registry
}

// Run executes all registered rules against a compose file and returns findings.
func Run(f *compose.File) []Finding {
	var findings []Finding
	for _, r := range registry {
		findings = append(findings, r.Check(f)...)
	}
	return findings
}

// Fixes collects all fixes from registered rules against a compose file.
func Fixes(f *compose.File) []Fix {
	var fixes []Fix
	for _, r := range registry {
		fixes = append(fixes, r.Fix(f)...)
	}
	return fixes
}

// MaxSeverity returns the highest severity among findings.
func MaxSeverity(findings []Finding) Severity {
	maxSev := Info
	for _, f := range findings {
		if f.Severity > maxSev {
			maxSev = f.Severity
		}
	}
	return maxSev
}
