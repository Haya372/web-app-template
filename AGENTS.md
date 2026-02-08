# Repository Guidelines

This repository is a template that pairs shared documentation and tooling at the root with a Go backend service in `go-backend/`. Treat the root as the coordination layer: keep it clean, reusable, and free of environment-specific hacks so downstream projects can rely on it verbatim.

## Workspace Layout
- `docs/` houses ADRs, runbooks, and operational references (`docs/decisions`, `docs/operations`, etc.). Update the relevant note whenever behavior or architecture shifts.
- `mise.toml` pins the toolchain. Run `mise install` once per machine to install Go 1.25.6, golangci-lint 2.8.0, and psqldef 3.9.7 before touching submodules.
- Service-specific contributor instructions live in `go-backend/AGENTS.md`; cross-link it in downstream READMEs when the Go service is used.

## Root Tooling & Environments
- Use `docker compose -f go-backend/docker-compose.yml up -d db` to boot the shared Postgres dependency; run `... down` when finished to avoid lingering state.
- Keep orchestration logic in versioned scripts or Make targets. If a command is shared by multiple modules, define it once at the root and call it from module-specific makefiles.
- Track lifecycle tasks (backups, resets, seeding) under `docs/operations` rather than private scratchpads.

## Collaboration & Hygiene
- Favor deterministic tooling (`mise`, Make, docker-compose) over manual setup. Every newcomer should be able to run the documented commands without guessing.
- When you add a module, mirror this documentation pattern: `AGENTS.md` at the module root plus a pointer from this file.
- Changes to CI/CD, repo-wide hooks, or shared infrastructure require a note in `docs/` and, if impactful, an ADR that explains the trade-offs.

## Commit & Pull Requests
- Follow the conventional short log already in history: `<type>(optional-scope): summary (#issue)`, e.g., `feat: add telemetry (#12)`.
- Each PR must describe motivation, test evidence, linked issues, and screenshots when user-facing artifacts change. Mention any docs you updated.
- Ensure linters/tests relevant to the touched modules pass locally before requesting review; root-level tweaks should be validated across all modules they affect.
