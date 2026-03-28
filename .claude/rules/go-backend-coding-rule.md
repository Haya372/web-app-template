---
paths:
  - "go-backend/**"
---

# Go-Backend Coding Rules

## Language

- Write all source code comments in English.
- Write all test case names (e.g., function names, description strings) in English.

## Naming Conventions

- Follow Go's initialisms convention: acronyms and initialisms in identifiers must be written in consistent case (all-caps or all-lowercase).
  - `ID` instead of `Id`
  - `URL` instead of `Url`
  - `HTTP` instead of `Http`
  - `JSON` instead of `Json`
  - `API` instead of `Api`
  - This applies to struct fields, function names, variable names, and type names.
  - Rationale: Effective Go and the Go Code Review Comments guide require this style. `golangci-lint` (`revive`) enforces it, and `sqlc`-generated code uses it — keeping consistent avoids conversion boilerplate at layer boundaries.
