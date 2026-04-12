package rule_test

import (
	"testing"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/rule"
)

func TestExposedPorts(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web": {Image: "nginx", Ports: []string{"80:80", "127.0.0.1:443:443"}},
			"db":  {Image: "postgres", Ports: []string{"5432:5432"}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "exposed-ports")
	// web port 80 and db port 5432 are exposed to 0.0.0.0, web port 443 is not
	if len(found) != 2 {
		t.Fatalf("expected 2 exposed-ports findings, got %d", len(found))
	}
}

func TestExposedPortsAllBound(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web": {Image: "nginx", Ports: []string{"127.0.0.1:80:80"}},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "exposed-ports")
	if len(found) != 0 {
		t.Fatalf("expected 0 exposed-ports findings, got %d", len(found))
	}
}

func TestDefaultNetwork(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web": {Image: "nginx"},
			"db":  {Image: "postgres"},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "default-network")
	if len(found) != 1 {
		t.Fatalf("expected 1 default-network finding, got %d", len(found))
	}
}

func TestDefaultNetworkWithCustomNetworks(t *testing.T) {
	f := &compose.File{
		Services: map[string]*compose.Service{
			"web": {Image: "nginx", Networks: []string{"frontend"}},
			"db":  {Image: "postgres", Networks: []string{"backend"}},
		},
		Networks: map[string]*compose.Network{
			"frontend": {},
			"backend":  {},
		},
	}
	findings := rule.Run(f)
	found := findByRule(findings, "default-network")
	if len(found) != 0 {
		t.Fatalf("expected 0 default-network findings, got %d", len(found))
	}
}
