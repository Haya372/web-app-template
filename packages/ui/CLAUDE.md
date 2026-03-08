# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from `packages/ui/`:

```bash
# Lint (read-only)
pnpm lint         # biome lint

# Format and auto-fix
pnpm check        # biome check --write

# Storybook (component explorer)
pnpm storybook    # dev server at http://localhost:6006
pnpm build-storybook  # static build
```

## Architecture

`packages/ui` is the shared component library consumed by all apps in this monorepo.

```
src/
  components/   # One file per component (PascalCase)
                # Co-located story: ComponentName.stories.tsx
  index.ts      # Public API — all exports go through here
```

- Package name: `@repo/ui`
- Apps import from `@repo/ui` (e.g. `import { Button } from "@repo/ui"`)
- Components are built on **Radix UI** primitives + **Tailwind CSS v4** + **class-variance-authority (CVA)**

## Component Design Rules

### No `className` prop

Components must **not** accept a `className` prop. Design variations are expressed exclusively through `variant`, `size`, or other semantic props defined in the CVA config.

```tsx
// good
type ButtonProps = Omit<ComponentProps<"button">, "className"> &
  VariantProps<typeof buttonVariants> & { asChild?: boolean }

// bad — exposes className and allows ad-hoc overrides
type ButtonProps = ComponentProps<"button"> & VariantProps<typeof buttonVariants>
```

If a new visual pattern is needed in an app, add a new `variant` to the component here rather than overriding styles from the outside.

### CVA for variants

Use `cva` from `class-variance-authority` for all variant definitions. Define `defaultVariants` explicitly.

```tsx
const buttonVariants = cva("/* base classes */", {
  variants: {
    variant: { default: "...", destructive: "...", outline: "..." },
    size: { default: "...", sm: "...", lg: "..." },
  },
  defaultVariants: { variant: "default", size: "default" },
})
```

### Radix UI for interactive primitives

Use Radix UI primitives (Dialog, Select, Tabs, Checkbox, etc.) for interactive and accessible components. Do not implement ARIA roles or keyboard interactions manually when a Radix primitive exists.

### `asChild` pattern

For components that need to render as a different element, support the `asChild` prop via `@radix-ui/react-slot`.

```tsx
function Button({ asChild = false, ...props }: ButtonProps) {
  const Comp = asChild ? Slot : "button"
  return <Comp className={buttonVariants({ variant, size })} {...props} />
}
```

## Storybook

Every component must have a co-located `.stories.tsx` file. Stories must cover:

- All `variant` values
- All `size` values (if applicable)
- Disabled / loading / error states
- An `AllVariants` story that renders every variant side-by-side

Use `satisfies Meta<typeof Component>` for type-safe story metadata.

## Coding Style

- Biome for formatting (tabs, double quotes) and linting; always pass `pnpm lint`
- Named exports only (no default exports except story `default meta`)
- `interface` for props that extend HTML element props; `type` for union / CVA variant intersection types
- `any` is forbidden; `as` casts are strongly discouraged — `as any` is unconditionally forbidden
