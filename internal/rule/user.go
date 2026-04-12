package rule

import (
	"fmt"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/image"
)

func init() {
	Register(RunningAsRoot{})
	Register(NoReadOnly{})
}

// RunningAsRoot flags services with no user: directive (defaulting to root).
type RunningAsRoot struct{}

func (r RunningAsRoot) ID() string { return "running-as-root" }

func (r RunningAsRoot) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.User != "" {
			continue
		}
		suggested := image.UserFor(svc.Image)
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "running as root (no user: directive)",
			Suggestion: fmt.Sprintf("Add user: %q", suggested),
		})
	}
	return findings
}

func (r RunningAsRoot) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if svc.User != "" {
			continue
		}
		user := image.UserFor(svc.Image)
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.User = user
			},
		})
	}
	return fixes
}

// NoReadOnly flags services without read_only: true.
type NoReadOnly struct{}

func (r NoReadOnly) ID() string { return "no-read-only" }

func (r NoReadOnly) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.ReadOnly {
			continue
		}
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "root filesystem is writable",
			Suggestion: "Add read_only: true (may require tmpfs mounts for writable dirs)",
		})
	}
	return findings
}

func (r NoReadOnly) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if svc.ReadOnly {
			continue
		}
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.ReadOnly = true
			},
		})
	}
	return fixes
}
