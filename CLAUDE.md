# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Layout

This monorepo pairs root-level documentation and tooling with a Go backend service:

- `go-backend/` — the Go HTTP API service (Echo v5, Clean Architecture)
- `docs/decisions/` — Architecture Decision Records (ADRs) that are mandatory implementation constraints
- `docs/guidelines/` — coding guidelines
- `mise.toml` — pins toolchain: Go 1.25.6, golangci-lint 2.8.0, psqldef 3.9.7

Run `mise install` once per machine before developing.

See `go-backend/CLAUDE.md` for service-specific commands and architecture.

## Commit & PR conventions

Format: `<type>(optional-scope): summary (#issue)` — e.g., `feat: add telemetry (#12)`

PR descriptions must include: motivation, test evidence (`make` targets run), linked issues, and any ADR/docs updates. Breaking API changes require a note in `docs/operations`.
