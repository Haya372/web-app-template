# Repository Guidelines

This module contains the Go backend that powers the template. Keep business logic, integrations, and database assets inside this directory to avoid leaking module-specific behavior to the workspace root.

## Module Structure & Ownership
- `cmd/` defines binaries (e.g., `cmd/api`); wire dependencies here and keep each entry point minimal.
- `internal/` hosts domain packages (transport, service, repository, etc.). Enforce acyclic dependencies and favor interfaces near the consumer.
- `db/schema` manages DDL, `db/query` manages sqlc inputs, and generated code lives under `internal` (ignored from Git; regenerate via `make generate`).
- `test/` is for integration helpers and fixtures. Colocated unit tests (`*_test.go`) sit beside their source files.

## Development Commands
- `make setup` installs sqlc, mockgen, and wire; rerun after changing Go versions.
- `make generate` (or the specific `generate-db-client` / `generate-di-container` targets) refreshes sqlc and wire output plus `go generate ./...` hooks.
- `make migrate-local` applies `db/schema/schema.sql` to the Postgres container using psqldef.
- `docker compose -f docker-compose.yml up -d db` (run from this directory) spins up the service database; tear it down when finished.

## Coding Style & Naming
- Always run `gofmt` (tabs, newline at EOF) and keep packages lowercase and concise.
- `make fmt` and `make lint` wrap golangci-lint with `.golangci.toml`; both must succeed before commits.
- Exported identifiers use `CamelCase`, private helpers use `camelCase`, and SQL artifacts use `snake_case` to ensure sqlc produces predictable structs.

## Testing Strategy
- `make test-unit` executes `go test ./...` for fast feedback. Keep unit tests hermetic and avoid DB calls.
- `make test-integration` and `make test-coverage` include the `integration` tag (`-tags=integration -p 1`); ensure the Postgres container is running first.
- Update `coverage.out` only when intentionally refreshing coverage data; avoid committing stale profiles.

## Review Checklist
- Commits should stay focused (schema/query + handwritten code changes together, with generation commands run locally). Reference root `AGENTS.md` for repo-wide expectations.
- PR descriptions must state the user problem, summarize changes, and list `make` targets executed. Attach logs or screenshots when behavior is visible externally.
