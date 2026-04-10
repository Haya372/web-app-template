---
name: review-ticket
description: "Review a GitHub Issue or Sub-Issue for quality: clear acceptance criteria, appropriate scope, no ambiguity. Reports deficiencies and suggests improvements. Use when you want to validate a ticket before starting implementation, or when refining requirements."
argument-hint: "<issue-number>"
model: haiku
user-invokable: true
---

# Review Ticket Skill (Refinement Review)

This skill reviews the quality of a GitHub Issue and reports deficiencies and ambiguities.

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` must contain a GitHub Issue number (e.g. `42` or `#42`).
Extract the numeric part and assign it to `ISSUE_NUMBER`.
If no number is found, stop and ask the user: "Which issue number should I review?"

### Step 1: Fetch the issue

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees,comments
```

Read all sections of the issue carefully.

### Step 2: Review against quality checklist

Evaluate the issue against the following criteria. Rate each item Pass / Fail / Warning.

#### 2.1 Clarity of summary and purpose
- [ ] Does the summary describe "what to do" in one sentence?
- [ ] Does the goal/success criteria clearly state "why it is needed"?
- [ ] Is business or user value demonstrated?

#### 2.2 Appropriateness of scope
- [ ] Are "in scope" and "out of scope" explicitly stated?
- [ ] Is the scope small enough to complete in a single iteration? (guideline: ≤ 7 to-do items)
- [ ] Are out-of-scope concerns noted in the risks/discussion section?

#### 2.3 Verifiability of acceptance criteria
- [ ] Can all acceptance criteria be verified objectively (no subjective expressions like "improves")?
- [ ] Do the acceptance criteria cover all major functional requirements?
- [ ] Are error cases and boundary values included in the criteria?

#### 2.4 Specificity of the to-do list
- [ ] Is each task granular enough that an implementer can start without hesitation?
- [ ] Are dependencies between tasks made explicit (when order matters)?
- [ ] Are technical implementation details (API endpoints, DB schema, etc.) documented appropriately?

#### 2.5 Completeness of test perspective
- [ ] Are normal-flow scenarios documented?
- [ ] Are error/edge-case scenarios documented?
- [ ] Are acceptance tests written in an executable form?

#### 2.6 Risks and dependencies
- [ ] Are known risks or unresolved questions documented?
- [ ] Are dependencies on other teams or issues made explicit?

### Step 3: Produce the review report

Output the review result in Japanese using the following format:

```
## チケットレビュー: #<ISSUE_NUMBER> — <タイトル>

### 総合評価
<実装開始可能 / 要修正 / 大幅な見直しが必要>

<総合評価の理由を2〜4文で記述>

### 🚨 必須修正 (Blocking)
実装開始前に必ず対応が必要な問題:

- **[問題のタイトル]**
  - **問題**: <何が問題か、なぜ実装に支障をきたすか>
  - **修正案**: <具体的な修正方法>

### ⚠️ 推奨修正 (Non-blocking)
対応が望ましいが実装開始を妨げない問題:

- **[問題のタイトル]**: <問題の説明と改善案>

### 💡 提案 (Optional)
品質向上のための任意の提案:

- **[提案のタイトル]**: <提案内容>

### ✅ 良い点
チケットの品質として評価できる点 (1〜3項目):
- <良い点>

### 品質スコア
| 観点 | 評価 |
|------|------|
| サマリ・目的 | ✅ Pass / ⚠️ Warning / ❌ Fail |
| スコープ | ✅ Pass / ⚠️ Warning / ❌ Fail |
| 完了条件 | ✅ Pass / ⚠️ Warning / ❌ Fail |
| タスクリスト | ✅ Pass / ⚠️ Warning / ❌ Fail |
| テスト観点 | ✅ Pass / ⚠️ Warning / ❌ Fail |
| リスク・依存関係 | ✅ Pass / ⚠️ Warning / ❌ Fail |
```

### Step 4: Post as comment (optional)

If the user instructs, post the review result as a comment on the issue:

```bash
gh issue comment $ISSUE_NUMBER --body "$(cat <<'EOF'
<review report>
EOF
)"
```

---

## Behavioral guidelines

- **Be specific**: don't just say "ambiguous" — show what is ambiguous and how to fix it
- **Be constructive**: the goal is to make the ticket implementable, not to find faults
- **Don't over-demand**: don't require a perfect ticket; it is enough that an implementer can start confidently
- **Respect context**: review with awareness of project conventions and constraints
