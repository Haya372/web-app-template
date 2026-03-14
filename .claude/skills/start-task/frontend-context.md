# Frontend (React) Context

## Architecture reference

See `apps/react-frontend/CLAUDE.md` and `docs/guidlines/frontend-coding-guidline.md`.

## Implement prompt (Green)

```
Agent: code-implementer
Prompt: "Implement the minimum code to make the failing tests pass.
Task: <task subject and description>. Issue context: <issue body excerpt>.
Follow the repository architecture in apps/react-frontend/CLAUDE.md and
docs/guidlines/frontend-coding-guidline.md. Tests are already written —
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
Prompt: "Review the implementation for this task. Check: correctness,
security (XSS, CSRF), accessibility, performance, readability, adherence
to repository guidelines in docs/guidlines/frontend-coding-guidline.md, and
test coverage. Task: <task description>."
```

## Quality gate

Run from `apps/react-frontend/`:
```bash
pnpm lint && pnpm test:agent
```
