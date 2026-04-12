package rule

import (
	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/image"
)

func init() {
	Register(NoHealthcheck{})
	Register(RestartNoHealth{})
}

// NoHealthcheck flags services with no healthcheck defined.
type NoHealthcheck struct{}

func (r NoHealthcheck) ID() string { return "no-healthcheck" }

func (r NoHealthcheck) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.Healthcheck != nil {
			continue
		}
		hc := image.HealthcheckFor(svc.Image)
		suggestion := "Add a healthcheck"
		if hc != nil && len(hc.Test) > 0 {
			suggestion = "Add healthcheck with: " + hc.Test[len(hc.Test)-1]
		}
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Info,
			Message:    "no healthcheck defined",
			Suggestion: suggestion,
		})
	}
	return findings
}

func (r NoHealthcheck) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if svc.Healthcheck != nil {
			continue
		}
		hc := image.HealthcheckFor(svc.Image)
		if hc == nil {
			continue // can't auto-fix without a known healthcheck
		}
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.Healthcheck = hc
			},
		})
	}
	return fixes
}

// RestartNoHealth flags services with restart: always/unless-stopped but no healthcheck.
type RestartNoHealth struct{}

func (r RestartNoHealth) ID() string { return "restart-no-health" }

func (r RestartNoHealth) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if svc.Healthcheck != nil {
			continue
		}
		if svc.Restart != "always" && svc.Restart != "unless-stopped" {
			continue
		}
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Info,
			Message:    "restart: " + svc.Restart + " without healthcheck",
			Suggestion: "Add a healthcheck so Docker can detect when the service is unhealthy",
		})
	}
	return findings
}

func (r RestartNoHealth) Fix(f *compose.File) []Fix {
	// This rule doesn't auto-fix — it just flags. The NoHealthcheck rule
	// handles adding healthchecks when a known image profile exists.
	return nil
}
