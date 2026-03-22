# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from `go-backend/`:

```bash
# One-time setup: installs sqlc, mockgen, wire
make setup

# Code generation (sqlc + wire + go generate)
make generate

# Format and lint
make fmt
make lint

# Tests
make test-unit          # go test ./... (no DB needed)
make test-integration   # requires Postgres container running
make test-coverage      # integration tests + coverage profile

# Database
docker compose -f docker-compose.yml up -d db   # start Postgres on :55432
make migrate-local                               # apply schema.sql via psqldef + seed data
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

- Echo types must not leak into `usecase` or `domain` layers (ADR-0004)
- Public API routes: `/v{major}/resource-names` with lowercase kebab-case plural nouns (ADR-0005)
- HTTP error responses: `application/problem+json` with stable `type`/`title`/`status` fields; no internal diagnostics in public payloads (ADR-0006)
- Transactions start in `usecase` via `TransactionManager.Do`, propagate via `context`; nested `TransactionManager.Do` calls are forbidden (ADR-0007)
- `wire` is used only at the composition root (`cmd/`)
- **depguard enforces import boundaries at lint time** — violations in `internal/domain` or `internal/usecase` cause `make lint` to fail:
  - `internal/domain/**` must not import `echo`, `wire`, `pgx`/`pgconn`/`pgtype`, `database/sql`, `internal/usecase`, or `internal/infrastructure`
  - `internal/usecase/**` must not import `echo`, `wire`, `pgx`/`pgconn`/`pgtype`, `database/sql`, or `internal/infrastructure`

**Code generation:**

- `db/schema/schema.sql` — DDL managed by psqldef
- `db/query/` — sqlc query inputs; generated client lives under `internal/` (git-ignored)
- Run `make generate` after schema or query changes; commit schema/query files together with generated output

## Implementation Rules

**Domain layer — immutability is mandatory:**

```go
// good: state-change returns a new instance
func (u User) UpdateStatus(s Status) (User, error) { ... }

// good: VO validates at construction, no external dependencies
func NewPassword(raw string) (*Password, error) {
    if len(raw) < 8 {
        return nil, ErrTooShort
    }
    pwd := Password(raw)
    return &pwd, nil
}

// bad: domain object issues SQL or knows HTTP types
func (u *User) Save(ctx context.Context, db *sql.DB) error { ... }
```

**UseCase layer — orchestrate via interfaces, own transaction boundary:**

```go
// good
func (uc *signupUseCaseImpl) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    return uc.txManager.Do(ctx, func(ctx context.Context) error {
        user, err := entity.NewUser(...)
        if err != nil { return err }
        _, err = uc.userRepository.Create(ctx, user)
        return err
    })
}

// bad: HTTP types, raw SQL, or multiple responsibilities inside a use case
func (uc *signupUseCaseImpl) Execute(ctx context.Context, req *echo.Context) error { ... }
```

**Infrastructure layer — implement ports, contain all side effects:**

```go
// good: implements the repository port, tracing/logging here
func (r *userRepositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
    return runInTx(ctx, func(q sqlc.Queries) error {
        return q.CreateUser(ctx, mapToParams(user))
    })
}

// bad: global variable, no port interface
var globalDB *sql.DB
func SaveUser(ctx context.Context, u *entity.User) error { ... }
```

Even when schema constraints exist, keep defensive validation when converting DB rows to Value Objects for early detection of unexpected data.

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
