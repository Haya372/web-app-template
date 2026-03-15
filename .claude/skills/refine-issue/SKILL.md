---
name: refine-issue
description: "Refine a GitHub Issue into structured requirements (要件定義) and create Sub-Issues as an implementation plan. Use when the user wants to break down an issue into actionable sub-tasks, or when starting work on a complex feature that needs planning before implementation."
argument-hint: "<issue-number>"
user-invokable: true
---

# Refine Issue Skill

This skill takes a GitHub Issue, organizes it as structured requirements (要件定義), and creates Sub-Issues as an implementation plan.

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` must contain a GitHub Issue number (e.g. `42` or `#42`).
Extract the numeric part and assign it to `ISSUE_NUMBER`.
If no number is found, stop and ask the user: "Which issue number should I refine?"

### Step 1: Fetch the issue

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees
```

Understand:
- **Summary**: what needs to be done
- **Background / Problem**: why it is needed
- **Goals / Success criteria**: definition of done
- **Scope**: what is and is not included
- **To-do list**: list of tasks
- **Acceptance criteria**: verifiable conditions

If the issue body is missing or too vague, stop and ask the user before continuing.

### Step 2: Produce structured requirements (要件定義)

Organize the issue content into the following structure and output it in Japanese:

```
## 要件定義: <Issue タイトル>

### 目的
<このIssueが達成すべきことを1〜3文で記述>

### スコープ
**対象:**
- <対象1>
- <対象2>

**対象外:**
- <対象外1>

### 機能要件
| # | 要件 | レイヤー | 優先度 |
|---|------|----------|--------|
| 1 | <要件の説明> | Backend / Frontend / DB / Docs | Must / Should / Could |

### 非機能要件
- <パフォーマンス・セキュリティ・保守性等の要件>

### 完了条件 (Acceptance Criteria)
- [ ] <検証可能な条件1>
- [ ] <検証可能な条件2>

### 依存関係
- <他の Issue や前提条件>
```

Present the requirements to the user and get confirmation before proceeding to Sub-Issue creation.

### Step 3: Split into Sub-Issues (実装計画)

Split into Sub-Issues based on functional requirements and layers.

Splitting rules:
- **1 Sub-Issue = 1 unit of functionality per layer** (e.g. adding a Backend endpoint, implementing a Frontend component)
- If there are dependencies, include `Blocked by #<issue>` in the body
- Each Sub-Issue should be granular enough to be reviewed and reverted independently

Sub-Issue title format: `<type>(<scope>): <requirement summary> (Sub-Issue of #<ISSUE_NUMBER>)`

Sub-Issue body template (write in Japanese):

```
**親 Issue:** #<ISSUE_NUMBER>

**サマリ**
<この Sub-Issue で行うことを1文で記述>

**やること**
- [ ] <具体的なタスク1>
- [ ] <具体的なタスク2>

**完了条件**
- [ ] <検証可能な条件>

**依存関係**
- Blocked by: <なし、または #<issue>>
```

Create each Sub-Issue:

```bash
gh issue create \
  --title "<sub-issue-title>" \
  --body "$(cat <<'EOF'
<sub-issue-body>
EOF
)"
```

### Step 4: Link Sub-Issues to the parent

Add a comment to the parent issue listing all Sub-Issues (write in Japanese):

```bash
gh issue comment $ISSUE_NUMBER --body "$(cat <<'EOF'
## 実装計画 (Sub-Issues)

以下の Sub-Issues に分割しました:

- [ ] #<sub1> — <Sub-Issue タイトル>
- [ ] #<sub2> — <Sub-Issue タイトル>

実装順序: <依存関係に基づいた順序の説明>
EOF
)"
```

### Step 5: Summary

Return a list of created Sub-Issue URLs to the user.

---

## Safety rules

- If the issue body is unclear, do not create Sub-Issues — ask the user first
- Sub-Issues must be more granular than the parent (do not simply duplicate it)
- If a single Sub-Issue is too large (more than 5 to-do items), consider splitting further
