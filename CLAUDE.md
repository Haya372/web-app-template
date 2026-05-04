# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Layout

This monorepo contains a Go backend, a React frontend, a shared UI library, and E2E tests:

- `go-backend/` — the Go HTTP API service (Echo v5, Clean Architecture)
- `apps/react-frontend/` — the React frontend (TanStack Router, Vite, Tailwind v4)
- `packages/ui/` — shared UI component library (`@repo/ui`, Radix UI + CVA)
- `e2e/` — Playwright E2E tests (runs against containerised stack)
- `openapi/` — OpenAPI schema (single source of truth for API contract)
- `docs/decisions/` — Architecture Decision Records (ADRs) that are mandatory implementation constraints
- `docs/guidelines/` — coding guidelines
- `mise.toml` — pins toolchain versions (see file for details)

See each workspace's `CLAUDE.md` for workspace-specific commands and architecture:
- `go-backend/CLAUDE.md`
- `apps/react-frontend/CLAUDE.md`
- `packages/ui/CLAUDE.md`

## Commit & PR conventions

Format: `<type>(optional-scope): summary (#issue)` — e.g., `feat: add telemetry (#12)`

PR descriptions must include: motivation, test evidence (make targets or equivalent commands run), linked issues, and any ADR/docs updates. Breaking API changes require a note in `docs/operations`.

## Sandbox environment notes

Claude Code runs in a sandboxed environment. Use these workarounds for known constraints:

| Constraint | Workaround |
|-----------|-----------|
| Edit tool cannot modify `.github/workflows/` files | Use Bash heredoc: `cat > file << 'EOF' ... EOF` |
| Semgrep blocked from writing to HOME | Env vars set in `.claude/settings.json` redirect logs/cache to `/tmp/claude/` |
| `pnpm add` / `pnpm install` are denied | Edit `package.json` manually, ask user to run install locally. `pnpm run`/`lint`/`test` work normally |
| `make` may fail to find Makefile when CWD ≠ repo root | Use `cd go-backend && make ...` or run direct commands below |
| Go test cache write denied | `GOCACHE=/tmp/claude/gocache` is set in `.claude/settings.json` |

See `go-backend/CLAUDE.md` for make target equivalents when `make` is unavailable.
