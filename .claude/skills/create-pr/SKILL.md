---
name: create-pr
description: This skill should be used when the user wants to create a pull request for the current branch. It gathers context, validates the branch state, and opens a PR following this repository's conventions.
disable-model-invocation: true
allowed-tools: Bash, Read, Glob, Grep
---

# Create Pull Request

## Context

- Current branch: !`git branch --show-current`
- Base branch: !`git remote show origin | grep 'HEAD branch' | awk '{print $NF}'`
- Commits to be merged: !`git log main..HEAD --oneline`
- Staged/unstaged changes: !`git status --short`
- Full diff from main: !`git diff main...HEAD --stat`

## PR Creation Process

### 1. Validate branch state

Ensure the current branch is not `main`. If on `main`, stop and ask the user to switch to a feature branch.

Verify there are commits ahead of `main` to merge. If there are no commits, stop and inform the user.

### 2. Determine PR metadata

Infer the following from the commit history and diff:

- **Title** — Follow `<type>(optional-scope): summary` format (no issue number in title). Types: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`, `perf`. Keep under 72 characters. **Write in English.**
- **Linked issues** — Identify issue numbers from commit messages (e.g., `#12`). If none are found, ask the user.
- **Breaking changes** — Check if any API routes, response shapes, or public interfaces changed. Breaking changes require a note in `docs/operations/`.
- **ADR/docs updates** — Check whether `docs/decisions/` or `docs/guidelines/` files were added or modified.

### 3. Check test evidence

Identify which `make` targets are relevant to the changes:

| Changed area | Required targets |
|---|---|
| Go source code | `make fmt && make lint && make test-unit` |
| Integration paths | `make test-integration` |
| DB schema / queries | `make generate && make migrate-local` |

Report which targets were run and whether they passed. If tests have not been run, prompt the user to run them before creating the PR.

### 4. Handle breaking API changes

If the diff includes changes to public API routes (`/v{n}/...`), HTTP response schemas, or exported Go interfaces:

1. Check whether `docs/operations/` contains a file documenting the breaking change.
2. If not, create or update `docs/operations/<topic>.md` with a migration note before opening the PR.

### 5. Create the PR

Run `gh pr create` with the following body structure. **Write all body content in Japanese.**

```
## 背景・目的

<なぜこの変更が必要か。関連 issue を参照すること。>

## 変更内容

<変更点の箇条書き。>

## テスト実施結果

<実行した `make` ターゲットと結果の一覧。>

## 関連 issue

Closes #<issue-number>

## ADR / ドキュメント更新

<追加・更新した ADR やガイドラインへの参照、または「なし」。>

## 破壊的変更

<破壊的な API 変更の説明と docs/operations エントリへのリンク、または「なし」。>
```

Use `--base main` unless the user specifies a different base branch. Push the branch first if it has no upstream:

```bash
git push -u origin HEAD
gh pr create --base main --title "<title>" --body "<body>"
```
