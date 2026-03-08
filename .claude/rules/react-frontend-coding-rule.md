---
paths:
  - "apps/react-frontend/**"
---

# React Frontend Coding Rules

## Language

- Write all source code comments in English.
- Write all test case names and test description strings in English.

## Naming Conventions

- **Components:** PascalCase for both the filename and the function name.
  - `Header.tsx` → `function Header()` ✓
  - `header.tsx` → `function header()` ✗
- **Hooks:** camelCase prefixed with `use`.
  - `useTheme.ts` → `function useTheme()` ✓
- **Utilities / non-component files:** camelCase.
  - `apiClient.ts`, `formatDate.ts` ✓
- **Route files:** follow TanStack Router file-based conventions.
  - `src/routes/index.tsx` (path `/`)
  - `src/routes/posts/$id.tsx` (path `/posts/:id`)
  - `src/routes/__root.tsx` (root layout)

## Import Alias

Always use the `@/*` alias for imports across feature boundaries or from shared modules. Relative paths are allowed only within the same feature directory.

```tsx
// good — cross-feature or shared module
import Header from "@/components/Header"
import { useTheme } from "@/hooks/useTheme"
import { fetchUser } from "@/features/user/api/users"

// good — within the same feature (relative is fine)
import { UserCard } from "./components/UserCard"

// bad — deep relative path crossing feature boundaries
import Header from "../../components/Header"
```

## Data Fetching

- Always use **TanStack Query** (`useQuery` / `useMutation`) for server data fetching. `useEffect` + `useState` manual fetching is forbidden.
- Place plain fetch functions in `features/<feature>/api/` and wrap them with TanStack Query in `features/<feature>/hooks/`.

```tsx
// good
export function useUser(id: string) {
  return useQuery({ queryKey: ["user", id], queryFn: () => fetchUser(id) })
}

// bad
useEffect(() => {
  fetch(`/api/users/${id}`).then(r => r.json()).then(setUser)
}, [id])
```

## Tailwind CSS

- Use Tailwind utility classes for all styling; avoid inline `style` props except for dynamic values that cannot be expressed as utilities.
- Define reusable design tokens (colors, spacing, fonts) as CSS custom properties in `src/styles.css`, then reference them via `var(--token-name)` in Tailwind classes.
- Do not add arbitrary pixel values when a Tailwind scale value exists (e.g., prefer `p-4` over `p-[16px]`).
- Responsive variants follow mobile-first order: base → `sm:` → `md:` → `lg:`.

## Component Design

- One component per file; keep components under ~150 lines.
- Extract repeated JSX patterns into named sub-components or map over data arrays.
- Avoid prop drilling beyond two levels — lift state or use context/store.
- Encapsulate side effects (`useEffect`, subscriptions, timers, etc.) in custom hooks. Components should only consume the return values, not contain side effect logic directly.
- Do not use default exports for anything other than route components and layout components (use named exports elsewhere).

## TypeScript

- Enable strict mode; `any` is forbidden — use `unknown` and narrow explicitly.
- `as` casts are strongly discouraged. Use type guard functions (`is` / `asserts`) or validation libraries (e.g. zod) to narrow types safely. If an `as` cast is truly unavoidable, add a comment explaining why. `as any` is unconditionally forbidden.
- Prefer `interface` over `type` for component props and object shapes. Use `type` only when `interface` cannot express it (e.g. union types, mapped types).
- Use `satisfies` or explicit return types on exported functions; avoid implicit `any` returns.

## Directory Structure

Follow feature-based directory structure. Co-locate components, hooks, API clients, and types under the same feature directory. Only put code in `src/components/` or `src/hooks/` when it is shared across multiple features.

```
src/features/<feature-name>/
  components/   # components used only by this feature
  hooks/        # hooks used only by this feature
  api/          # API client functions for this feature
  types/        # domain and response types for this feature
```

- Import files directly by path — both within a feature and across features.
- `index.ts` barrel re-exports are forbidden; they add unnecessary files and obscure dependency paths.

## packages/ui Usage

App-level components (`apps/`) must be built by composing `packages/ui` components.

- Never pass `className` to a `packages/ui` component — doing so breaks design consistency. Use `variant` / `size` or other semantic props instead.
- For layout adjustments (margin, width, alignment), wrap the component and apply Tailwind classes to the wrapper.

```tsx
// good — layout via wrapper only
<div className="mt-4 w-full">
  <Button variant="primary">Submit</Button>
</div>

// bad — overriding design with className on a UI component
<Button className="mt-4 bg-red-500">Submit</Button>
```

- If a new visual pattern is needed, add a new variant to `packages/ui` rather than patching it with `className`.
