# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from `apps/react-frontend/`:

```bash
# Type checking
pnpm typecheck

# Generate route tree (run after adding/removing route files)
pnpm generate:route

# Format and lint
pnpm fmt          # biome lint --write + knip --fix

# Tests
pnpm test:agent   # vitest run (single pass, no DB needed)

# Build
pnpm build        # prebuild (generate + typecheck) then vite build
```

## Design Principles

- Keep components small and focused on a single responsibility
- Prefer composition over inheritance; share logic via custom hooks
- Follow TDD where applicable: write a failing test, implement the minimum to pass, then refactor
- Keep PRs small with clear motivation, scope, and test evidence

## Architecture

The app uses TanStack Router (CSR React framework):

```
src/
  routes/          # File-based routing — one file per route
    __root.tsx     # Root layout: HTML shell, Header, Footer
    index.tsx      # Home page (path: /)
    about.tsx      # About page (path: /about)
  features/        # Feature-based modules (co-locate components/hooks/api/types/pages per feature)
    <feature>/
      components/  # Components used only by this feature
      hooks/       # Hooks used only by this feature
      types/       # Domain and response types for this feature
      pages/       # Full-page components rendered by route files
  components/      # Shared UI components used across multiple features (PascalCase)
  hooks/           # Shared hooks used across multiple features
  utils/           # Pure utility functions (no side effects)
  styles.css       # Global styles and CSS custom properties (Tailwind v4)
  router.tsx       # Router instantiation and type registration
  routeTree.gen.ts # Auto-generated route tree — do not edit manually
```

**Key constraints:**

See [frontend-coding-guideline.md](../../docs/guidelines/frontend-coding-guideline.md) for full constraint details.


## Routing Rules (TanStack Router)

See [frontend-coding-guideline.md](../../docs/guidelines/frontend-coding-guideline.md) for TanStack Router rules.

## Coding Style

See [frontend-coding-guideline.md](../../docs/guidelines/frontend-coding-guideline.md) for coding style rules.

## Testing

- Use **Vitest** with **jsdom** environment for unit and component tests.
- Place test files adjacent to the source file: `Header.test.tsx` next to `Header.tsx`.
- Use table-driven tests for utility functions with boundary and error cases.
- There is no integration test setup; focus unit tests on pure logic and component rendering.
- Aim for meaningful coverage of business logic; UI snapshot tests are discouraged.

### Running tests

Run from `apps/react-frontend/`:

```bash
# All tests (agent-friendly reporter)
pnpm test:agent

# Specific file
pnpm test:agent src/features/auth/pages/LoginPage.test.tsx

# Specific test or describe block (substring match on name)
pnpm test:agent -t "LoginPage — rendering"

# Combine: file + name filter
pnpm test:agent src/features/auth/pages/LoginPage.test.tsx -t "renders an email input"
```

The `-t` flag matches against the full test name (`describe` block + `it` label concatenated). Use a unique substring to target a single case.

## Authentication State Management

Global auth state is managed by `AuthProvider` in `src/features/auth/contexts/AuthContext.tsx`. Key rules:

- Access auth state via `useAuth()` — never consume `AuthContext` directly with `useContext`.
- Token persistence is handled by `src/features/auth/utils/tokenStorage.ts` (pure localStorage API). Do **not** call `saveToken`/`removeToken` directly from components or hooks — always go through `login()`/`logout()` from `useAuth()`.
- Protected routes are guarded in `src/routes/_authenticated.tsx` via `beforeLoad`. Since `beforeLoad` runs outside React, it calls `getToken()` directly rather than `useAuth()`.
- `AuthProvider` is placed inside `RootLayout` in `src/routes/__root.tsx`, making `useAuth()` available in all routes.

See [ADR-0011](../../docs/decisions/ADR-0011-FRONTEND-AUTH-STATE-MANAGEMENT.md) for design rationale.

## Security

See [frontend-security-guideline.md](../../docs/guidelines/frontend-security-guideline.md).
