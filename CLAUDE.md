# dockguard -- Claude Code Context

## Read this first

**Start every session by reading `BUSINESS.md` in the narrowcast repo.** It is the source of truth for the Narrowcast side hustle that this project is part of.

```
read /Volumes/Untitled/narrowcast/BUSINESS.md
```

## What is dockguard?

A Go CLI tool that scans Docker Compose files for security misconfigurations and generates hardened versions with service-aware fixes. It targets self-hosters and homelab operators running 10-40 Docker containers who know their setup is insecure but lack the expertise or time to audit and fix it manually.

**Open source (MIT license).** The priority is distribution and community trust -- self-hosters share free tools, not paid ones.

**Repo:** `github.com/narrowcastdev/dockguard`
**GitHub org:** `narrowcastdev`

## Key documents

- **Design spec:** `/Volumes/Untitled/narrowcast/docs/superpowers/specs/2026-04-11-dockguard-design.md`
- **Implementation plan:** `/Volumes/Untitled/narrowcast/docs/superpowers/plans/2026-04-11-dockguard.md`

## Tech stack

- **Go** (stdlib + `gopkg.in/yaml.v3`) for the CLI
- **GoReleaser** for cross-platform distribution

## Dev environment

This project uses **Nix flakes**. The `flake.nix` provides Go and Go tools. Run `direnv allow` to activate the dev shell automatically.

## Project structure

```
dockguard/
  cmd/dockguard/main.go            -- CLI entry point, flag parsing
  internal/
    compose/compose.go             -- YAML parsing into typed Go model
    rule/
      rule.go                      -- Rule interface, Finding, Fix, Severity, registry
      capabilities.go              -- no-cap-drop, privileged-mode, no-new-privileges
      network.go                   -- default-network, exposed-ports
      secrets.go                   -- plaintext-secrets
      resources.go                 -- no-memory-limit, no-cpu-limit
      health.go                    -- no-healthcheck, restart-no-health
      user.go                      -- running-as-root, no-read-only
    image/catalog.go               -- Service-aware image database (~20 images)
    report/report.go               -- Terminal + JSON output formatting
    hardener/hardener.go           -- Apply fixes to yaml.Node tree, write YAML
```
