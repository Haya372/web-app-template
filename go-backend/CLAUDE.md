# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from `go-backend/`:

```bash
# Code generation (sqlc + wire + go generate)
make generate

# Format and lint
make fmt
make lint

# Tests
make test-unit          # go test ./... (no DB needed)
make test-integration   # requires Postgres container running
make test-coverage      # integration tests + coverage profile
```

### make target equivalents (use when `make` is unavailable)

| make target | Direct command |
|---|---|
| `make fmt` | `golangci-lint run -c .golangci.toml --fix` |
| `make lint` | `golangci-lint run -c .golangci.toml` |
| `make test-unit` | `go test ./...` |
| `make test-integration` | `go test ./... -tags=integration` |
| `make test-coverage` | `go test ./... -tags=integration -coverprofile=coverage.out && grep -F -v -f .coverageignore coverage.out > coverage.tmp && mv coverage.tmp coverage.out` |
| `make generate` | See `go-backend/Makefile` for full sequence (sqlc → buf → oapi-codegen → mockgen → wire) |

## Design Principles

- Follow DDD: use ubiquitous language and bounded contexts for naming and structure
- Follow TDD: write a failing test first, then implement the minimum to make it pass, then refactor
- Apply SRP/SOLID and avoid over-abstraction — introduce abstractions only when multiple usages exist
- Keep PRs small with clear motivation, scope, and test evidence

## Architecture

The service follows Clean Architecture with strict inward dependency flow:

```
infrastructure (adapters) → usecase → domain
```

**Layer responsibilities:**

- `cmd/api` — composition root: wire DI, config init, graceful shutdown
- `internal/domain` — business rules only, zero I/O; contains:
  - `entity/` — persistent domain objects (immutable: state-change methods return new instances)
  - `vo/` — value objects (identifier-free, immutable, validated at construction)
  - `aggregate/` — multi-entity aggregates used across use cases; fetched via `Repository`, never mutated directly
  - `entity/repository/` and `aggregate/repository/` — port interfaces (implemented in infrastructure)
- `internal/usecase` — application flow, CQRS (Command/Query separation); owns transaction boundaries via `TransactionManager`
- `internal/infrastructure` — adapter implementations (`*_impl`), DB access via sqlc/pgx, OTel tracing/logging
- `internal/common` — cross-cutting: logging, config, shared error handling (no domain logic)

**Key constraints (enforced by ADRs):**

- Echo types must not leak into `usecase` or `domain` layers (ADR-0004)
- Public API routes: `/v{major}/resource-names` with lowercase kebab-case plural nouns (ADR-0005)
- HTTP error responses: `application/problem+json` with stable `type`/`title`/`status` fields; no internal diagnostics in public payloads (ADR-0006)
- Transactions start in `usecase` via `TransactionManager.Do`, propagate via `context`; nested `TransactionManager.Do` calls are forbidden (ADR-0007)
- `wire` is used only at the composition root (`cmd/`)
- **depguard enforces import boundaries at lint time** — violations in `internal/domain` or `internal/usecase` cause `make lint` to fail (see `.golangci.toml` for the full deny list)

**Code generation:**

- `db/schema/schema.sql` — DDL managed by psqldef
- `db/query/` — sqlc query inputs; generated client lives under `internal/` (git-ignored)
- Run `make generate` after schema or query changes; commit schema/query files together with generated output

## Implementation Rules

See [@docs/guidelines/backend-coding-guideline.md](../docs/guidelines/backend-coding-guideline.md) for detailed patterns with code examples (immutability, transaction boundaries, port interfaces, error handling).

## Coding Style

- Apply `go fmt` / `goimports` (optionally `gofumpt`) on save; always pass `golangci-lint`
- Wrap errors with `%w`; use `errors.Is` / `errors.As` for inspection; map error codes to HTTP responses in one place
- Use Value Objects to keep input validation and type safety inside the domain layer
- Read config and secrets from environment variables or a Secrets Manager — never hardcode
- Write all source comments, test case names, and test messages in English

## Testing

- Domain / UseCase: unit tests required; use table-driven tests covering boundary and error cases
- Infrastructure: integration tests with real connections via Testcontainers or local mock servers
- Target branch coverage ≥ 80% overall; critical use cases ≥ 90%
- Bug fixes must include a regression test that fails before the fix

## Observability

- Every server must emit structured logs with a trace ID for per-request tracing
- Use OpenTelemetry SDK; initialise Tracer / Logger / Metrics at the composition root (`cmd/`)
- Expose key metrics: latency, error rate, throughput, DB query count; define SLO/SLA

## Security

- Validate all input before passing to the domain layer; guard against SQL/command injection and CSRF
- Use prepared statements or official SDKs for DB and external API access; store credentials in a Secret Store
- Monitor dependency CVEs regularly; automate updates with Renovate/Dependabot
