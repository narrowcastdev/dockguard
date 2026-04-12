package compose

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// File represents a parsed Docker Compose file.
type File struct {
	Version  string              `yaml:"version,omitempty"`
	Services map[string]*Service `yaml:"services"`
	Networks map[string]*Network `yaml:"networks,omitempty"`
	Secrets  map[string]*Secret  `yaml:"secrets,omitempty"`
}

// Service represents a single service in a compose file.
type Service struct {
	Image       string            `yaml:"image,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"-"`
	RawEnv      any               `yaml:"environment,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	CapDrop     []string          `yaml:"cap_drop,omitempty"`
	CapAdd      []string          `yaml:"cap_add,omitempty"`
	Privileged  bool              `yaml:"privileged,omitempty"`
	SecurityOpt []string          `yaml:"security_opt,omitempty"`
	User        string            `yaml:"user,omitempty"`
	ReadOnly    bool              `yaml:"read_only,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
	Healthcheck *Healthcheck      `yaml:"healthcheck,omitempty"`
	Deploy      *Deploy           `yaml:"deploy,omitempty"`
	Secrets     []string          `yaml:"secrets,omitempty"`
	DependsOn   any               `yaml:"depends_on,omitempty"`
}

// Healthcheck represents a service healthcheck configuration.
type Healthcheck struct {
	Test     []string `yaml:"test,omitempty"`
	Interval string   `yaml:"interval,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Retries  int      `yaml:"retries,omitempty"`
}

// Deploy represents the deploy configuration for a service.
type Deploy struct {
	Resources *Resources `yaml:"resources,omitempty"`
}

// Resources represents resource constraints.
type Resources struct {
	Limits *ResourceLimit `yaml:"limits,omitempty"`
}

// ResourceLimit represents CPU and memory limits.
type ResourceLimit struct {
	Memory string `yaml:"memory,omitempty"`
	CPUs   string `yaml:"cpus,omitempty"`
}

// Network represents a top-level network definition.
type Network struct {
	Driver string         `yaml:"driver,omitempty"`
	Extra  map[string]any `yaml:",inline"`
}

// Secret represents a top-level secret definition.
type Secret struct {
	File     string         `yaml:"file,omitempty"`
	External bool           `yaml:"external,omitempty"`
	Extra    map[string]any `yaml:",inline"`
}

// ParseFile reads and parses a Docker Compose file from disk.
func ParseFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file: %w", err)
	}
	return Parse(data)
}

// Parse parses raw YAML bytes into a compose File.
// Handles both map and list formats for environment variables.
func Parse(data []byte) (*File, error) {
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	if len(f.Services) == 0 {
		return nil, fmt.Errorf("not a Docker Compose file: no services defined")
	}

	// Normalize environment from RawEnv (handles both map and list formats).
	for _, svc := range f.Services {
		svc.Environment = normalizeEnvironment(svc.RawEnv)
	}

	return &f, nil
}

// normalizeEnvironment converts the raw environment value (map or list) into
// a map[string]string.
func normalizeEnvironment(raw any) map[string]string {
	if raw == nil {
		return nil
	}

	switch v := raw.(type) {
	case map[string]any:
		env := make(map[string]string, len(v))
		for key, val := range v {
			if val == nil {
				env[key] = ""
				continue
			}
			env[key] = fmt.Sprintf("%v", val)
		}
		return env
	case []any:
		env := make(map[string]string, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				continue
			}
			parts := strings.SplitN(s, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			} else {
				env[parts[0]] = ""
			}
		}
		return env
	}
	return nil
}
