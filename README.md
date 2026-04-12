# dockguard

Docker Compose security scanner with fix generation.

Scans your `docker-compose.yml` for security misconfigurations and generates a hardened version with service-aware fixes.

## Quick start

```bash
# Scan a compose file
dockguard docker-compose.yml

# Scan and generate hardened file
dockguard docker-compose.yml --fix
```

Or with Docker:

```bash
docker run --rm -v $(pwd):/workspace narrowcastdev/dockguard /workspace/docker-compose.yml
```

## What it checks

| Rule | Severity | What |
|------|----------|------|
| `privileged-mode` | Critical | Privileged containers |
| `exposed-ports` | Critical | Ports bound to 0.0.0.0 |
| `plaintext-secrets` | Critical | Passwords in environment variables |
| `no-cap-drop` | Warning | Missing capability restrictions |
| `no-new-privileges` | Warning | Missing privilege escalation prevention |
| `running-as-root` | Warning | No user: directive |
| `no-read-only` | Warning | Writable root filesystem |
| `default-network` | Warning | No network segmentation |
| `no-memory-limit` | Warning | No memory limits |
| `no-cpu-limit` | Warning | No CPU limits |
| `no-healthcheck` | Info | No healthcheck defined |
| `restart-no-health` | Info | Auto-restart without health monitoring |

## Service-aware fixes

dockguard knows 20+ common Docker images and generates tailored fixes:

- **postgres** gets `pg_isready` healthcheck, user `999:999`, caps `CHOWN, SETUID, SETGID`
- **redis** gets `redis-cli ping` healthcheck, 256m memory limit
- **nginx** gets `NET_BIND_SERVICE` capability, ports bound to 127.0.0.1
- ...and more

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

**Exit codes:** 0 (clean), 1 (warnings), 2 (critical issues)

## Install

**Binary:** Download from [GitHub Releases](https://github.com/narrowcastdev/dockguard/releases).

**Docker:**
```bash
docker pull narrowcastdev/dockguard
```

**From source:**
```bash
go install github.com/narrowcastdev/dockguard/cmd/dockguard@latest
```

## What this is NOT

- Not an image vulnerability scanner (use [Trivy](https://trivy.dev/) for that)
- Not a Dockerfile linter (use [Hadolint](https://github.com/hadolint/hadolint))
- No network access, no telemetry, no external API calls

## License

MIT
