# Backend (Go) Context

## Architecture reference

See `go-backend/CLAUDE.md`.

## Implement prompt (Green)

```
Agent: code-implementer
Prompt: "Implement the minimum code to make the failing tests pass.
Task: <task subject and description>. Issue context: <issue body excerpt>.
Follow the repository architecture in go-backend/CLAUDE.md. Tests are
already written — make them green."
```

Run tests:
```bash
make test-unit
```

If integration paths were changed:
```bash
make test-integration
```

## Review prompt (Refactor)

```
Agent: code-reviewer
Prompt: "Review the implementation for this task. Check: correctness,
security (OWASP top 10), performance (N+1, timeouts), readability, adherence
to repository guidelines in docs/guidelines/backend-coding-guideline.md, and
test coverage. Task: <task description>."
```

## Quality gate

```bash
make fmt && make lint && make test-unit
```

If integration paths were changed:
```bash
make test-integration
```

## Coverage check

```bash
make test-coverage
```

If any package touched by this PR shows decreased coverage compared to `main`, add tests before proceeding.
