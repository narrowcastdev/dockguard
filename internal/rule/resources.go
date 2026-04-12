package rule

import (
	"fmt"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/image"
)

func init() {
	Register(NoMemoryLimit{})
	Register(NoCPULimit{})
}

// NoMemoryLimit flags services with no memory limit.
type NoMemoryLimit struct{}

func (r NoMemoryLimit) ID() string { return "no-memory-limit" }

func (r NoMemoryLimit) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if hasMemoryLimit(svc) {
			continue
		}
		suggested := image.MemoryLimitFor(svc.Image)
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "no memory limit set",
			Suggestion: fmt.Sprintf("Add deploy.resources.limits.memory: %s", suggested),
		})
	}
	return findings
}

func (r NoMemoryLimit) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if hasMemoryLimit(svc) {
			continue
		}
		limit := image.MemoryLimitFor(svc.Image)
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				ensureDeploy(s)
				s.Deploy.Resources.Limits.Memory = limit
			},
		})
	}
	return fixes
}

// NoCPULimit flags services with no CPU limit.
type NoCPULimit struct{}

func (r NoCPULimit) ID() string { return "no-cpu-limit" }

func (r NoCPULimit) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		if hasCPULimit(svc) {
			continue
		}
		suggested := image.CPULimitFor(svc.Image)
		findings = append(findings, Finding{
			Rule:       r.ID(),
			Service:    name,
			Severity:   Warning,
			Message:    "no CPU limit set",
			Suggestion: fmt.Sprintf("Add deploy.resources.limits.cpus: %s", suggested),
		})
	}
	return findings
}

func (r NoCPULimit) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		if hasCPULimit(svc) {
			continue
		}
		limit := image.CPULimitFor(svc.Image)
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				ensureDeploy(s)
				s.Deploy.Resources.Limits.CPUs = limit
			},
		})
	}
	return fixes
}

func hasMemoryLimit(svc *compose.Service) bool {
	if svc.Deploy == nil || svc.Deploy.Resources == nil || svc.Deploy.Resources.Limits == nil {
		return false
	}
	return svc.Deploy.Resources.Limits.Memory != ""
}

func hasCPULimit(svc *compose.Service) bool {
	if svc.Deploy == nil || svc.Deploy.Resources == nil || svc.Deploy.Resources.Limits == nil {
		return false
	}
	return svc.Deploy.Resources.Limits.CPUs != ""
}

func ensureDeploy(svc *compose.Service) {
	if svc.Deploy == nil {
		svc.Deploy = &compose.Deploy{}
	}
	if svc.Deploy.Resources == nil {
		svc.Deploy.Resources = &compose.Resources{}
	}
	if svc.Deploy.Resources.Limits == nil {
		svc.Deploy.Resources.Limits = &compose.ResourceLimit{}
	}
}
