# Frontend (React) Context

## Architecture reference

See `apps/react-frontend/CLAUDE.md` and `docs/guidelines/frontend-coding-guideline.md`.

## Implement prompt (Green)

```
Agent: code-implementer
Model: sonnet
Prompt: "Implement the minimum code to make the failing tests pass.
Task: <task subject and description>. Issue context: <issue body excerpt>.
Follow the repository architecture in apps/react-frontend/CLAUDE.md and
docs/guidelines/frontend-coding-guideline.md. Tests are already written —
make them green."
```

Run tests (from `apps/react-frontend/`):
```bash
# All tests
pnpm test:agent

# Specific file
pnpm test:agent src/features/auth/pages/LoginPage.test.tsx

# Specific test/describe name (substring match)
pnpm test:agent -t "LoginPage — rendering"

# Combine file and name filter
pnpm test:agent src/features/auth/pages/LoginPage.test.tsx -t "renders an email input"
```

## Review prompt (Refactor)

```
Agent: code-reviewer
Model: sonnet
Prompt: "Review the implementation for this task. Check: correctness,
security (XSS, CSRF), accessibility, performance, readability, adherence
to repository guidelines in docs/guidelines/frontend-coding-guideline.md, and
test coverage.
Task: <task name and description>
Acceptance criteria: <acceptance criteria bullet points from the issue>
Design constraints: <ADR references and architectural constraints>"
```

## Quality gate

Run from `apps/react-frontend/`:
```bash
pnpm lint && pnpm test:agent
```
