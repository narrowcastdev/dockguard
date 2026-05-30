# dockguard

[![CI](https://github.com/narrowcastdev/dockguard/actions/workflows/ci.yml/badge.svg)](https://github.com/narrowcastdev/dockguard/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/narrowcastdev/dockguard)](https://goreportcard.com/report/github.com/narrowcastdev/dockguard)

**Docker Compose security scanner** — scans your `docker-compose.yml` for misconfigurations (privileged mode, exposed ports, plaintext secrets, missing capability restrictions) and generates a hardened version. Free, open source, single binary.

Built for self-hosters and homelab operators who know their setup is insecure but lack the time to audit and fix it manually.

📖 **[Docker Compose Security: The Complete Guide](https://narrowcast.dev/blog/docker-compose-security-guide/)** — learn the 9 most common misconfigurations and how to fix each one.

## Quick Start

```bash
# Scan a compose file
dockguard docker-compose.yml

# Scan and generate a hardened version
dockguard docker-compose.yml --fix
```

Example output:

```
🔴 CRITICAL
  [postgres] privileged-mode: container runs in privileged mode
  [redis] exposed-ports: port 6379 bound to 0.0.0.0 (all interfaces)

🟡 WARNING
  [nginx] no-cap-drop: no capabilities dropped
  [postgres] running-as-root: no user directive set

Found 8 findings (2 critical, 4 warning, 2 info)
```

## Installation

**Binary** — download from [GitHub Releases](https://github.com/narrowcastdev/dockguard/releases):

```bash
# Linux (amd64)
curl -Lo dockguard https://github.com/narrowcastdev/dockguard/releases/latest/download/dockguard_linux_amd64.tar.gz
tar xzf dockguard_linux_amd64.tar.gz
sudo mv dockguard /usr/local/bin/
```

**Docker:**

```bash
docker run --rm -v $(pwd):/workspace narrowcastdev/dockguard /workspace/docker-compose.yml
```

**Go install:**

```bash
go install github.com/narrowcastdev/dockguard/cmd/dockguard@latest
```

**Build from source:**

```bash
git clone https://github.com/narrowcastdev/dockguard.git
cd dockguard
go build -o dockguard ./cmd/dockguard
```

## What It Checks

| Rule | Severity | What |
|------|----------|------|
| `privileged-mode` | Critical | Privileged containers |
| `exposed-ports` | Critical | Ports bound to 0.0.0.0 |
| `plaintext-secrets` | Critical | Passwords in environment variables |
| `no-cap-drop` | Warning | Missing capability restrictions |
| `no-new-privileges` | Warning | Missing privilege escalation prevention |
| `running-as-root` | Warning | No `user:` directive |
| `no-read-only` | Warning | Writable root filesystem |
| `default-network` | Warning | No network segmentation |
| `no-memory-limit` | Warning | No memory limits |
| `no-cpu-limit` | Warning | No CPU limits |
| `no-healthcheck` | Info | No healthcheck defined |
| `restart-no-health` | Info | Auto-restart without health monitoring |

## Service-Aware Fixes

dockguard knows 20+ common Docker images and generates tailored fixes instead of generic defaults:

| Image | Healthcheck | User | Capabilities | Memory |
|-------|-------------|------|-------------|--------|
| postgres | `pg_isready` | `999:999` | `CHOWN, SETUID, SETGID` | 512m |
| redis | `redis-cli ping` | `999:999` | — | 256m |
| nginx | `curl -f http://localhost/` | `101:101` | `NET_BIND_SERVICE` | 128m |
| mysql | `mysqladmin ping` | `999:999` | `CHOWN, SETUID, SETGID` | 512m |
| traefik | `traefik healthcheck` | `65534:65534` | `NET_BIND_SERVICE` | 256m |

Unknown images get safe generic defaults.

## Usage

```
dockguard [flags] <docker-compose.yml>

Flags:
  --fix              Generate hardened compose file
  -o string          Output path (default "docker-compose.hardened.yml")
  --json             Output findings as JSON
  --severity string  Minimum severity: info, warning, critical (default "info")
```

**Exit codes:** `0` clean, `1` warnings found, `2` critical issues found.

### JSON output

```bash
dockguard docker-compose.yml --json | jq '.[] | select(.severity == "critical")'
```

### CI usage

```bash
# Fail the build on critical issues (exit code 2)
dockguard docker-compose.yml --severity critical
```

## AI Skill

dockguard ships with an interactive skill for AI coding agents in [`skills/scan.md`](skills/scan.md). Instead of just pass/fail output, the skill explains *why* each finding matters and walks you through tradeoffs before generating fixes.

Works with any agent that supports markdown skill files ([Claude Code](https://docs.anthropic.com/en/docs/claude-code), etc). Copy `skills/scan.md` into your project's commands directory, or point your agent at it directly.

```bash
# Claude Code — run from your project directory
cp /path/to/dockguard/skills/scan.md .claude/commands/scan.md
# Then: /project:scan docker-compose.yml
```

Use the CLI for automation (CI, scripts, batch scanning). Use the skill when you want interactive guidance on what to fix and why.

## What This Is NOT

- Not an image vulnerability scanner (use [Trivy](https://trivy.dev/) for that)
- Not a Dockerfile linter (use [Hadolint](https://github.com/hadolint/hadolint))
- Not a container runtime security tool (use [Falco](https://falco.org/))
- No network access, no telemetry, no external API calls — runs entirely offline

## Contributing

Contributions welcome — open an issue, send a PR, or fork it and make it your own.

```bash
git clone https://github.com/narrowcastdev/dockguard.git
cd dockguard
go build ./cmd/dockguard
go test ./...
```

## License

[MIT](LICENSE)
