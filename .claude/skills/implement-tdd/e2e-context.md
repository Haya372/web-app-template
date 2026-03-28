# E2E Context

## Scope

Apply to tasks that add or change UI or user-facing API endpoints.

## Architecture reference

See `docs/guidelines/e2e-testing.md`.

## Test-writer prompt

```
Agent: test-writer
Prompt: "Write an E2E test spec for the following feature using Playwright.
Place it in e2e/tests/<feature>.spec.ts. Follow the conventions in
docs/guidelines/e2e-testing.md:
- Use getByRole / getByLabel selectors
- Set up test data via the backend REST API in beforeEach
- If the UI is not yet implemented, use test.fixme() with a comment
  describing the missing prerequisite
Task: <task subject and description>. Issue context: <issue body excerpt>."
```

## Test execution

From the monorepo root:
```bash
pnpm test:e2e
```

During development (when the stack is already running):
```bash
pnpm --filter e2e docker:up
pnpm --filter e2e test
pnpm --filter e2e docker:down
```

## fixme rules

- If the frontend UI is not yet implemented, mark with `test.fixme()` and include a comment with the missing prerequisite and the related issue number.
- Once the UI is complete, remove the fixme and promote the test to a passing spec.
