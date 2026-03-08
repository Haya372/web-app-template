---
paths:
  - "packages/ui/**"
---

# Frontend Coding Rules

## packages/ui Component Design

### File naming

Component files in `packages/ui/src/components/` must use **PascalCase** to match the component name.

- `Button.tsx` ✓
- `button.tsx` ✗

### No `className` prop on UI components

Components in `packages/ui` must NOT accept a `className` prop.
Design variations must be expressed exclusively through `variant` and `size` (or other semantic) props using `class-variance-authority`.

**Bad:**
```tsx
function Button({ className, variant, ...props }: ButtonProps) {
  return <button className={cn(buttonVariants({ variant }), className)} {...props} />
}
```

**Good:**
```tsx
function Button({ variant, size, ...props }: ButtonProps) {
  return <button className={buttonVariants({ variant, size })} {...props} />
}
```

- If a new visual pattern is needed, add a new variant instead of passing `className` from the outside.
- The `cn()` utility is still used internally when combining multiple CVA results, but must not be exposed via `className` prop.
