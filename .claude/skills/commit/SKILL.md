---
name: commit
description: Analyze git changes, split them into logical commits at appropriate granularity, and execute them following this repository's conventions. Use this skill whenever the user says "commit", "コミット", "変更をコミット", "commit my changes", "save my work to git", or asks to create a git commit — even if they don't say "commit" explicitly but clearly want to save their changes to git history.
---

# Git Commit Skill

This repository uses the convention: `<type>(<scope>): <summary> (#<issue>)`

## Workflow

### Step 1: Understand all changes

Run in parallel:
- `git status` — full list of staged and unstaged files
- `git diff --staged` — staged diff
- `git diff` — unstaged diff
- `git log --oneline -5` — recent commit style reference

### Step 2: Plan commit groupings

Examine the changes and propose a split into **logical, self-contained commits**. Each commit should represent one coherent intent — something that could be reviewed and reverted independently.

**Good reasons to split into separate commits:**
- Different `type` (e.g. `feat` changes alongside `chore` dependency bumps)
- Different layers/components that change for different reasons (e.g. domain logic vs. infrastructure wiring vs. tests)
- A bug fix that is independent of a new feature added in the same session
- Config/tooling changes that are unrelated to business logic

**Keep in one commit when:**
- Changes are tightly coupled (e.g. an entity, its repository implementation, and its unit test all added together)
- Splitting would leave the repo in a broken/inconsistent state

Present the proposed plan to the user as a numbered list before doing anything:

```
I found changes across 8 files. I suggest 3 commits:

1. feat(go-backend): add user entity and repository interface
   → internal/domain/entity/user.go, internal/domain/entity/repository/user.go

2. feat(go-backend): implement user repository with sqlc
   → internal/infrastructure/user_repository_impl.go, db/query/user.sql

3. test(go-backend): add unit tests for user entity
   → internal/domain/entity/user_test.go
```

Ask: "Does this split look right, or would you like to adjust it?"

### Step 3: Generate commit messages

For each group, apply the format:

```
<type>(<scope>): <summary> (#<issue>)
```

**Type selection:**

| type | when to use |
|------|-------------|
| `feat` | new feature or user-facing behaviour |
| `fix` | bug fix |
| `refactor` | behaviour-preserving restructure |
| `chore` | build, deps, tooling, config |
| `docs` | documentation only |
| `test` | tests added or updated |
| `perf` | performance improvement |
| `ci` | CI/CD pipeline changes |

**Scope:** the affected component — e.g. `go-backend`, `react-frontend`, `auth`, `docker`. Omit when changes span many components.

**Summary:** imperative English verb phrase, ≤ 50 chars.

**Issue number:** use from `$ARGUMENTS` if provided; otherwise infer from branch name (e.g. `feature/42-auth` → `#42`); omit if unknown.

### Step 4: Execute commits in order

After the user approves the plan, execute each commit sequentially:

1. Stage only the files for that commit (use explicit file paths, never `git add -A`)
2. Run the commit
3. Confirm success with `git status` before moving to the next

```bash
git add <file1> <file2> ...
git commit -m "$(cat <<'EOF'
<message>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
EOF
)"
```

After all commits, show a brief summary:
```
Done. 3 commits created:
  abc1234 feat(go-backend): add user entity and repository interface
  def5678 feat(go-backend): implement user repository with sqlc
  ghi9012 test(go-backend): add unit tests for user entity
```

## Safety rules

- Never commit files that may contain secrets (`.env`, `*.key`, credentials)
- Never use `--no-verify` — if a hook fails, investigate and fix the root cause
- Never amend an existing commit unless the user explicitly asks
- Stage files explicitly by path; never use `git add -A` or `git add .`
- If a commit fails, fix the issue and create a NEW commit; do not amend

## Arguments

`$ARGUMENTS` may contain:
- An issue number (e.g. `42` or `#42`) — apply to all commits in the session
- A hint about intent — use it to inform grouping and message generation
- A single commit message — treat as a signal to commit everything in one shot
