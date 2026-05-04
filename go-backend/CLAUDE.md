## Commands

Run from `go-backend/`:

```bash
# Code generation
# When `make` is unavailable, see `go-backend/Makefile` for full sequence
make generate

golangci-lint run -c .golangci.toml --fix # Format and lint

# Tests
go test ./...          # go test ./... (no DB needed)
go test ./... -tags=integration   # requires Postgres container running
# integration tests + coverage profile
go test ./... -tags=integration -coverprofile=coverage.out && grep -F -v -f .coverageignore coverage.out > coverage.tmp && mv coverage.tmp coverage.out
```

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

- Public API routes: `/v{major}/resource-names` with lowercase kebab-case plural nouns ([ADR-0005](../docs/decisions/ADR-0005-API-VERSIONING-AND-ROUTE-DESIGN-FOR-GO-BACKEND.md))
- HTTP error responses: `application/problem+json` with stable `type`/`title`/`status` fields; no internal diagnostics in public payloads ([ADR-0006](../docs/decisions/ADR-0006-ERROR-CONTRACT-AND-MAPPING-POLICY-FOR-GO-BACKEND.md))
- Transactions start in `usecase` via `TransactionManager.Do`, propagate via `context`; nested `TransactionManager.Do` calls are forbidden ([ADR-0007](../docs/decisions/ADR-0007-TRANSACTION-BOUNDARY-AND-PROPAGATION-FOR-GO-BACKEND.md))

**Code generation:**

- `db/schema/schema.sql` — DDL managed by psqldef
- `db/query/` — sqlc query inputs; generated client lives under `internal/` (git-ignored)
- Run `make generate` after schema or query changes

## Implementation Rules

See [@docs/guidelines/backend-coding-guideline.md](../docs/guidelines/backend-coding-guideline.md) for detailed patterns with code examples.

## Coding Style

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

- Use OpenTelemetry SDK; initialise Tracer / Logger / Metrics at the composition root (`cmd/`)
- Expose key metrics: latency, error rate, throughput, DB query count; define SLO/SLA

## Security

See [backend-security-guideline.md](../docs/guidelines/backend-security-guideline.md).
