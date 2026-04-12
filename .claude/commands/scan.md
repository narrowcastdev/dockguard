# Command: scan

Scan a Docker Compose file for security misconfigurations and suggest hardened fixes.

## Argument

$ARGUMENTS = path to a docker-compose.yml file.
If blank, look for `docker-compose.yml`, `docker-compose.yaml`, or `compose.yml` in the current directory.

## Your task

Read the compose file. Apply every rule below to every service. Report findings grouped by severity. Then offer to generate a hardened version.

---

## Rules

Apply all 12 rules. For each finding, note the **rule ID**, **service name**, **severity**, and a **plain-English explanation** of why it matters and how to fix it.

### Critical

**privileged-mode** — Service has `privileged: true`. This gives the container full access to the host kernel — effectively root on the host machine. Remove it. If the service needs specific capabilities, use `cap_add` instead.

**exposed-ports** — Port mapping uses `"hostPort:containerPort"` format without binding to `127.0.0.1`. This exposes the port to your entire network (and the internet if port-forwarded). Change to `"127.0.0.1:hostPort:containerPort"` unless you intentionally need LAN/WAN access. For services behind a reverse proxy, always bind to localhost.

**plaintext-secrets** — An environment variable whose name contains PASSWORD, PASSWD, SECRET, TOKEN, API_KEY, APIKEY, DB_PASS, PRIVATE_KEY, or AUTH has a literal string value (not a `${VAR}` reference, not empty). This means the secret is in plaintext in the compose file — anyone with read access sees it. Move to a `.env` file (`${VAR_NAME}`) or Docker secrets.

### Warning

**no-cap-drop** — Service does not have `cap_drop: [ALL]`. By default Docker containers get a broad set of Linux capabilities. Best practice is to drop all and selectively add back only what's needed. Use the image catalog below for which caps each image actually needs.

**no-new-privileges** — Service is missing `security_opt: [no-new-privileges]`. Without this, processes inside the container can gain additional privileges via setuid/setgid binaries. Always add it.

**running-as-root** — Service has no `user:` directive, so it runs as root inside the container. If the container is compromised, the attacker has root. Use the image catalog below for the correct non-root user per image.

**no-read-only** — Service does not have `read_only: true`. A read-only root filesystem prevents an attacker from writing malicious files. Some services need writable dirs — use tmpfs mounts for those.

**default-network** — No custom networks are defined and no service specifies a `networks:` list. All services share the default bridge network and can communicate freely. Create separate networks to isolate services that don't need to talk to each other (e.g., your reverse proxy doesn't need direct access to your database).

**no-memory-limit** — Service has no memory limit (neither `deploy.resources.limits.memory` nor `mem_limit`). A misbehaving container can consume all host memory and crash everything. Use the image catalog for suggested defaults.

**no-cpu-limit** — Service has no CPU limit. Same reasoning as memory — prevent one container from starving others.

### Info

**no-healthcheck** — Service has no `healthcheck` defined. Docker can't tell if the service is actually working. Use the image catalog for service-specific healthchecks.

**restart-no-health** — Service has `restart: always` or `restart: unless-stopped` but no healthcheck. Docker will keep restarting a broken container forever without knowing it's unhealthy. Add a healthcheck first.

---

## Image Catalog

When suggesting fixes, use these service-aware defaults. For images not in this list, use generic defaults: user `1000:1000`, no extra capabilities, memory `512m`, CPU `1.0`, no healthcheck (suggest the user find one).

| Image | User | Capabilities (after cap_drop ALL) | Healthcheck | Memory | CPU |
|-------|------|----------------------------------|-------------|--------|-----|
| postgres | `999:999` | CHOWN, SETUID, SETGID | `pg_isready -U ${POSTGRES_USER:-postgres}` | 512m | 1.0 |
| redis | `999:999` | SETUID, SETGID | `redis-cli ping` | 256m | 0.5 |
| mysql / mariadb | `999:999` | CHOWN, SETUID, SETGID | `mysqladmin ping -h localhost` | 512m | 1.0 |
| mongo | `999:999` | CHOWN, SETUID, SETGID | `mongosh --eval "db.runCommand('ping')"` | 512m | 1.0 |
| nginx | `101:101` | NET_BIND_SERVICE | `curl -f http://localhost/ \|\| exit 1` | 128m | 0.5 |
| traefik | `65534:65534` | NET_BIND_SERVICE | `traefik healthcheck` | 256m | 0.5 |
| caddy | `1000:1000` | NET_BIND_SERVICE | `caddy version` | 128m | 0.5 |
| jellyfin | `1000:1000` | none | `curl -f http://localhost:8096/health` | 1g | 2.0 |
| immich-server | `1000:1000` | none | `curl -f http://localhost:2283/api/server/ping` | 1g | 2.0 |
| nextcloud | `33:33` | none | `curl -f http://localhost/status.php` | 512m | 1.0 |
| vaultwarden | `1000:1000` | none | `curl -f http://localhost/alive` | 256m | 0.5 |
| sonarr | `1000:1000` | none | `curl -f http://localhost:8989/ping` | 512m | 0.5 |
| radarr | `1000:1000` | none | `curl -f http://localhost:7878/ping` | 512m | 0.5 |
| prowlarr | `1000:1000` | none | `curl -f http://localhost:9696/ping` | 256m | 0.5 |
| grafana | `472:472` | none | `curl -f http://localhost:3000/api/health` | 256m | 0.5 |
| prometheus | `65534:65534` | none | `wget --spider http://localhost:9090/-/healthy` | 512m | 0.5 |
| pihole | `999:999` | NET_BIND_SERVICE, CHOWN | `dig +short @127.0.0.1 pi.hole` | 256m | 0.5 |
| homeassistant | `1000:1000` | none | `curl -f http://localhost:8123/api/` | 512m | 1.0 |
| adguardhome | `1000:1000` | NET_BIND_SERVICE | `wget --spider http://localhost:3000/` | 256m | 0.5 |
| plex | `1000:1000` | none | `curl -f http://localhost:32400/identity` | 1g | 2.0 |
| portainer | `1000:1000` | none | `curl -f http://localhost:9000/api/status` | 256m | 0.5 |
| uptime-kuma | `1000:1000` | none | `curl -f http://localhost:3001/` | 256m | 0.5 |

---

## Output format

### 1. Report

Print findings grouped by severity (critical first). Use this format:

```
## Security Scan: {filename}

**Services scanned:** {count}

### 🔴 Critical

**{service}: {message}** `[{rule-id}]`
{1-2 sentence explanation of why this is dangerous and how to fix it}

### 🟡 Warning

(same format)

### 🔵 Info

(same format)

### Summary

{count} critical · {count} warnings · {count} info
```

### 2. Offer to harden

After the report, ask:

> "Want me to generate a hardened version of this compose file? I'll apply fixes for all findings and you can review the diff."

If the user says yes:
- Create a copy of the compose file with all fixable issues resolved
- Use the image catalog for service-aware values
- Add a comment at the top: `# Hardened by dockguard scan — review changes before applying`
- Show the user a summary of what changed per service
- Do NOT overwrite the original file — write to `docker-compose.hardened.yml` or ask the user where

### 3. Explain tradeoffs

Be honest about tradeoffs when they exist:
- Binding ports to `127.0.0.1` breaks LAN access — note this if the service is typically accessed from other devices (Jellyfin, Plex, Home Assistant)
- `read_only: true` breaks some images — note when tmpfs mounts are needed
- Resource limits are suggestions — the user's hardware may need different values
- Some images (especially linuxserver.io images) expect to run as root — note when the user suggestion may cause issues

Don't just apply rules blindly. Use judgment. That's why this is a Claude skill and not just a CLI.
