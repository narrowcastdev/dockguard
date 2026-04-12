package rule

import (
	"fmt"
	"strings"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/image"
)

func init() {
	Register(PrivilegedMode{})
	Register(NoCapDrop{})
	Register(NoNewPrivileges{})
}

// PrivilegedMode flags services running with privileged: true.
type PrivilegedMode struct{}

func (r PrivilegedMode) ID() string { return "privileged-mode" }

func (r PrivilegedMode) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.Privileged {
			findings = append(findings, Finding{
				Rule:       r.ID(),
				Service:    name,
				Severity:   Critical,
				Message:    "privileged mode enabled",
				Suggestion: "Remove privileged: true and add only required capabilities",
			})
		}
	}
	return findings
}

func (r PrivilegedMode) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if !svc.Privileged {
			continue
		}
		caps := image.CapabilitiesFor(svc.Image)
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.Privileged = false
				s.CapDrop = []string{"ALL"}
				if len(caps) > 0 {
					s.CapAdd = caps
				}
			},
		})
	}
	return fixes
}

// NoCapDrop flags services that don't drop all capabilities.
type NoCapDrop struct{}

func (r NoCapDrop) ID() string { return "no-cap-drop" }

func (r NoCapDrop) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.Privileged {
			continue // handled by PrivilegedMode rule
		}
		if hasCapDropAll(svc) {
			continue
		}
		caps := image.CapabilitiesFor(svc.Image)
		suggestion := "Add cap_drop: [ALL]"
		if len(caps) > 0 {
			suggestion += fmt.Sprintf(" and cap_add: [%s]", strings.Join(caps, ", "))
		}
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "runs with full Linux capabilities",
			Suggestion: suggestion,
		})
	}
	return findings
}

func (r NoCapDrop) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if svc.Privileged || hasCapDropAll(svc) {
			continue
		}
		caps := image.CapabilitiesFor(svc.Image)
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.CapDrop = []string{"ALL"}
				if len(caps) > 0 {
					s.CapAdd = caps
				}
			},
		})
	}
	return fixes
}

// NoNewPrivileges flags services missing security_opt: no-new-privileges.
type NoNewPrivileges struct{}

func (r NoNewPrivileges) ID() string { return "no-new-privileges" }

func (r NoNewPrivileges) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if hasNoNewPrivileges(svc) {
			continue
		}
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "missing security_opt: no-new-privileges",
			Suggestion: "Add security_opt: [no-new-privileges] to prevent privilege escalation",
		})
	}
	return findings
}

func (r NoNewPrivileges) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if hasNoNewPrivileges(svc) {
			continue
		}
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.SecurityOpt = append(s.SecurityOpt, "no-new-privileges")
			},
		})
	}
	return fixes
}

func hasCapDropAll(svc *compose.Service) bool {
	for _, cap := range svc.CapDrop {
		if strings.EqualFold(cap, "ALL") {
			return true
		}
	}
	return false
}

func hasNoNewPrivileges(svc *compose.Service) bool {
	for _, opt := range svc.SecurityOpt {
		if opt == "no-new-privileges" || opt == "no-new-privileges:true" {
			return true
		}
	}
	return false
}
