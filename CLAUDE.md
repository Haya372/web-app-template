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
- `mise.toml` — pins toolchain: Go 1.25.6, Node 24.14.0, pnpm 10.30.3, golangci-lint 2.8.0, psqldef 3.9.7

Run `mise install` once per machine before developing.

See each workspace's `CLAUDE.md` for workspace-specific commands and architecture:
- `go-backend/CLAUDE.md`
- `apps/react-frontend/CLAUDE.md`
- `packages/ui/CLAUDE.md`

## Commit & PR conventions

Format: `<type>(optional-scope): summary (#issue)` — e.g., `feat: add telemetry (#12)`

PR descriptions must include: motivation, test evidence (`make` targets run), linked issues, and any ADR/docs updates. Breaking API changes require a note in `docs/operations`.
