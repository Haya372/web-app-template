---
name: create-ticket
description: "Create a well-structured GitHub Issue following the project ticket template in docs/guidlines/task-progress-guide.md. Use when the user wants to file a bug report, feature request, or task."
argument-hint: "<description of the ticket to create>"
user-invocable: true
---

# Create Ticket Skill

## Workflow

### Step 1: Understand the request

`$ARGUMENTS` contains a description of what to turn into a ticket. Extract:
- What the problem or improvement is
- Which layers are affected (Backend, Frontend, DB, Docs, etc.)
- Any known constraints or risks

If information is insufficient, ask the user before proceeding.

### Step 2: Read the template

Read the "チケットテンプレート" section of `docs/guidlines/task-progress-guide.md` to understand the required structure.

### Step 3: Draft the ticket body

Fill in all sections of the template:

- **サマリ**: one-line intent
- **背景 / 課題**: fact-based description of the problem
- **目的 / 成功基準**: declare the desired end state
- **スコープ**: what is and is not in scope
- **要件 / 設計**: API spec, UI flow, DB changes — as appropriate per layer
- **やること**: checklist ordered by dependency
- **完了条件**: objectively verifiable acceptance criteria
- **テスト (受け入れ観点)**: concrete normal and error scenarios
- **リスク / 相談事項**: unknowns and blockers
- **参考**: related tickets, design docs, URLs

### Step 4: Confirm with the user

Present the drafted ticket and ask for approval. Revise as needed.

### Step 5: Create the GitHub Issue

Once approved, create the issue:

```bash
gh issue create --title "<summary>" --body "$(cat <<'EOF'
<ticket body>
EOF
)"
```

Return the created Issue URL to the user.
