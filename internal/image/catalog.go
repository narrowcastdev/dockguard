package image

import (
	"strings"

	"github.com/narrowcastdev/dockguard/internal/compose"
)

// ImageProfile contains service-aware security defaults for a known Docker image.
type ImageProfile struct {
	Capabilities []string
	User         string
	Healthcheck  *compose.Healthcheck
	MemoryLimit  string
	CPULimit     string
}

var catalog = map[string]*ImageProfile{
	"postgres": {
		Capabilities: []string{"CHOWN", "SETUID", "SETGID"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"},
			Interval: "10s", Timeout: "5s", Retries: 5,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"redis": {
		Capabilities: []string{"SETUID", "SETGID"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: "10s", Timeout: "5s", Retries: 5,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"mysql": {
		Capabilities: []string{"CHOWN", "SETUID", "SETGID"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "mysqladmin ping -h localhost"},
			Interval: "10s", Timeout: "5s", Retries: 5,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"mariadb": {
		Capabilities: []string{"CHOWN", "SETUID", "SETGID"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "mysqladmin ping -h localhost"},
			Interval: "10s", Timeout: "5s", Retries: 5,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"mongo": {
		Capabilities: []string{"CHOWN", "SETUID", "SETGID"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", `mongosh --eval "db.runCommand('ping')"`},
			Interval: "10s", Timeout: "5s", Retries: 5,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"nginx": {
		Capabilities: []string{"NET_BIND_SERVICE"},
		User:         "101:101",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost/ || exit 1"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "128m", CPULimit: "0.5",
	},
	"traefik": {
		Capabilities: []string{"NET_BIND_SERVICE"},
		User:         "65534:65534",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD", "traefik", "healthcheck"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"caddy": {
		Capabilities: []string{"NET_BIND_SERVICE"},
		User:         "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD", "caddy", "version"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "128m", CPULimit: "0.5",
	},
	"jellyfin": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:8096/health || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "1g", CPULimit: "2.0",
	},
	"immich-server": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:2283/api/server/ping || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "1g", CPULimit: "2.0",
	},
	"nextcloud": {
		User: "33:33",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost/status.php || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"vaultwarden": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost/alive || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"sonarr": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:8989/ping || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "512m", CPULimit: "0.5",
	},
	"radarr": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:7878/ping || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "512m", CPULimit: "0.5",
	},
	"prowlarr": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:9696/ping || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"grafana": {
		User: "472:472",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:3000/api/health || exit 1"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"prometheus": {
		User: "65534:65534",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "wget --spider -q http://localhost:9090/-/healthy"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "512m", CPULimit: "0.5",
	},
	"pihole": {
		Capabilities: []string{"NET_BIND_SERVICE", "CHOWN"},
		User:         "999:999",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "dig +short @127.0.0.1 pi.hole || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"homeassistant": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:8123/api/ || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "512m", CPULimit: "1.0",
	},
	"adguardhome": {
		Capabilities: []string{"NET_BIND_SERVICE"},
		User:         "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "wget --spider -q http://localhost:3000/"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"plex": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:32400/identity || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "1g", CPULimit: "2.0",
	},
	"portainer": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:9000/api/status || exit 1"},
			Interval: "10s", Timeout: "5s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
	"uptime-kuma": {
		User: "1000:1000",
		Healthcheck: &compose.Healthcheck{
			Test:     []string{"CMD-SHELL", "curl -f http://localhost:3001/ || exit 1"},
			Interval: "30s", Timeout: "10s", Retries: 3,
		},
		MemoryLimit: "256m", CPULimit: "0.5",
	},
}

// NormalizeImage extracts the base image name from a full reference.
// "library/postgres:16-alpine" -> "postgres"
// "lscr.io/linuxserver/sonarr:latest" -> "sonarr"
// "ghcr.io/immich-app/immich-server:release" -> "immich-server"
func NormalizeImage(ref string) string {
	if ref == "" {
		return ""
	}
	// Strip tag (everything after last colon, but not port numbers).
	if idx := strings.LastIndex(ref, ":"); idx != -1 {
		// Check if this looks like a port (digits only after colon before /).
		afterColon := ref[idx+1:]
		if !strings.Contains(afterColon, "/") {
			ref = ref[:idx]
		}
	}
	// Take the last path segment.
	if idx := strings.LastIndex(ref, "/"); idx != -1 {
		ref = ref[idx+1:]
	}
	return ref
}

// Lookup returns the ImageProfile for a given image reference, or nil if unknown.
func Lookup(imageRef string) *ImageProfile {
	name := NormalizeImage(imageRef)
	return catalog[name]
}

// CapabilitiesFor returns the capabilities a known image needs after cap_drop ALL.
func CapabilitiesFor(imageRef string) []string {
	p := Lookup(imageRef)
	if p == nil {
		return nil
	}
	return p.Capabilities
}

// UserFor returns the suggested non-root user for a known image.
func UserFor(imageRef string) string {
	p := Lookup(imageRef)
	if p == nil {
		return "1000:1000"
	}
	if p.User == "" {
		return "1000:1000"
	}
	return p.User
}

// HealthcheckFor returns the suggested healthcheck for a known image.
func HealthcheckFor(imageRef string) *compose.Healthcheck {
	p := Lookup(imageRef)
	if p == nil {
		return nil
	}
	return p.Healthcheck
}

// MemoryLimitFor returns the suggested memory limit for a known image.
func MemoryLimitFor(imageRef string) string {
	p := Lookup(imageRef)
	if p == nil {
		return "512m"
	}
	if p.MemoryLimit == "" {
		return "512m"
	}
	return p.MemoryLimit
}

// CPULimitFor returns the suggested CPU limit for a known image.
func CPULimitFor(imageRef string) string {
	p := Lookup(imageRef)
	if p == nil {
		return "1.0"
	}
	if p.CPULimit == "" {
		return "1.0"
	}
	return p.CPULimit
}
