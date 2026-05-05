# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from `apps/next-frontend/`:

```bash
# Type checking
pnpm typecheck

# Format and lint
pnpm fmt          # biome check --write + eslint --fix + knip --fix
pnpm lint         # biome check + eslint + knip (read-only)

# Tests
pnpm test:agent   # vitest run (single pass, no DB needed)

# Dev server (uses NEXT_PORT from .env.local, default 3001)
pnpm dev

# Build
pnpm build        # next build
```

## Architecture

Next.js 16 App Router with `src/app/` layout:

```
src/
  app/
    layout.tsx       # Root layout: html/body shell, global CSS, metadata
    page.tsx         # Home page (path: /)
    globals.css      # Tailwind v4 entry (@import "tailwindcss")
  test-setup.ts      # Vitest setup: imports @testing-library/jest-dom
```

**Key constraints:**

- Server Components by default; use `"use client"` only when browser APIs or React hooks are needed.
- `@/` path alias maps to `src/` (configured in both `tsconfig.json` and `vitest.config.ts`).
- Tailwind v4 via `@tailwindcss/postcss` (PostCSS plugin) — do NOT use `@tailwindcss/vite`.
- packages/ui integration (`@repo/ui`) is deferred to Issue #106.
- Connect-RPC integration is deferred to Issue #106.

## Testing

- Use **Vitest** with **jsdom** environment for sync Server Component and utility tests.
- Async RSC (Server Components that fetch data) → test via Playwright E2E.
- Place test files adjacent to the source file: `page.test.tsx` next to `page.tsx`.

### Running tests

Run from `apps/next-frontend/`:

```bash
# All tests (agent-friendly reporter)
pnpm test:agent

# Specific file
pnpm test:agent src/app/page.test.tsx

# Specific test or describe block (substring match on name)
pnpm test:agent -t "renders the title heading"
```

## Port

Default port is `3001`. Override via `NEXT_PORT` in `.env.local` (managed by `mise run ports:init`).

## Security

See [frontend-security-guideline.md](../../docs/guidelines/frontend-security-guideline.md).
