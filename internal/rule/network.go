package rule

import (
	"fmt"
	"strings"

	"github.com/narrowcastdev/dockguard/internal/compose"
)

func init() {
	Register(ExposedPorts{})
	Register(DefaultNetwork{})
}

// ExposedPorts flags ports bound to 0.0.0.0 instead of 127.0.0.1.
type ExposedPorts struct{}

func (r ExposedPorts) ID() string { return "exposed-ports" }

func (r ExposedPorts) Check(f *compose.File) []Finding {
	var findings []Finding
	for name, svc := range f.Services {
		for _, port := range svc.Ports {
			if !isExposedPort(port) {
				continue
			}
			hostPort := extractHostPort(port)
			findings = append(findings, Finding{
				Rule:       r.ID(),
				Service:    name,
				Severity:   Critical,
				Message:    fmt.Sprintf("port %s exposed to 0.0.0.0", hostPort),
				Suggestion: fmt.Sprintf("Bind to 127.0.0.1:%s instead", port),
			})
		}
	}
	return findings
}

func (r ExposedPorts) Fix(f *compose.File) []Fix {
	var fixes []Fix
	for name, svc := range f.Services {
		hasExposed := false
		for _, port := range svc.Ports {
			if isExposedPort(port) {
				hasExposed = true
				break
			}
		}
		if !hasExposed {
			continue
		}
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				for i, port := range s.Ports {
					if isExposedPort(port) {
						s.Ports[i] = "127.0.0.1:" + port
					}
				}
			},
		})
	}
	return fixes
}

// isExposedPort returns true if a port mapping is bound to all interfaces.
// Short format: "80:80" (exposed) vs "127.0.0.1:80:80" (bound).
func isExposedPort(port string) bool {
	parts := strings.Split(port, ":")
	switch len(parts) {
	case 2:
		// "hostPort:containerPort" -- bound to 0.0.0.0
		return true
	case 3:
		// "ip:hostPort:containerPort" -- check the IP
		return parts[0] != "127.0.0.1" && parts[0] != "localhost"
	}
	return false
}

func extractHostPort(port string) string {
	parts := strings.Split(port, ":")
	switch len(parts) {
	case 2:
		return parts[0]
	case 3:
		return parts[1]
	}
	return port
}

// DefaultNetwork flags when all services use the default bridge network.
type DefaultNetwork struct{}

func (r DefaultNetwork) ID() string { return "default-network" }

func (r DefaultNetwork) Check(f *compose.File) []Finding {
	if len(f.Networks) > 0 {
		return nil
	}
	for _, svc := range f.Services {
		if len(svc.Networks) > 0 {
			return nil
		}
	}
	if len(f.Services) < 2 {
		return nil
	}
	return []Finding{{
		Rule:       r.ID(),
		Service:    "(all)",
		Severity:   Warning,
		Message:    "all services on default bridge network (no segmentation)",
		Suggestion: "Create separate networks to isolate services that don't need to communicate",
	}}
}

func (r DefaultNetwork) Fix(f *compose.File) []Fix {
	if len(f.Networks) > 0 {
		return nil
	}
	for _, svc := range f.Services {
		if len(svc.Networks) > 0 {
			return nil
		}
	}
	if len(f.Services) < 2 {
		return nil
	}

	// Build dependency groups from depends_on.
	groups := buildNetworkGroups(f)

	var fixes []Fix
	for name := range f.Services {
		network := groups[name]
		svcName := name
		fixes = append(fixes, Fix{
			Service: svcName,
			Apply: func(s *compose.Service) {
				s.Networks = []string{network}
			},
		})
	}
	return fixes
}

// buildNetworkGroups assigns each service to a named network based on depends_on.
func buildNetworkGroups(f *compose.File) map[string]string {
	groups := make(map[string]string)
	groupID := 0

	// First pass: group services connected by depends_on.
	for name, svc := range f.Services {
		if _, ok := groups[name]; ok {
			continue
		}
		deps := extractDependsOn(svc.DependsOn)
		if len(deps) == 0 {
			continue
		}
		netName := fmt.Sprintf("dg-net-%d", groupID)
		groupID++
		groups[name] = netName
		for _, dep := range deps {
			groups[dep] = netName
		}
	}

	// Second pass: isolated services get their own network.
	for name := range f.Services {
		if _, ok := groups[name]; ok {
			continue
		}
		netName := fmt.Sprintf("dg-net-%d", groupID)
		groupID++
		groups[name] = netName
	}

	return groups
}

func extractDependsOn(raw any) []string {
	if raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case []any:
		var deps []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				deps = append(deps, s)
			}
		}
		return deps
	case map[string]any:
		var deps []string
		for k := range v {
			deps = append(deps, k)
		}
		return deps
	}
	return nil
}
