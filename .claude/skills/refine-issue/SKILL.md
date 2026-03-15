---
name: refine-issue
description: "Refine a GitHub Issue into structured requirements (要件定義) and create Sub-Issues as an implementation plan. Use when the user wants to break down an issue into actionable sub-tasks, or when starting work on a complex feature that needs planning before implementation."
argument-hint: "<issue-number>"
user-invokable: true
---

# Refine Issue Skill

このスキルは GitHub Issue を受け取り、要件定義として整理し、実装計画を Sub-Issue として作成します。

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` に GitHub Issue 番号が含まれている必要があります（例: `42` または `#42`）。
数値部分を抽出して `ISSUE_NUMBER` に代入します。
番号が見つからない場合は止まり、ユーザーに質問します: "Which issue number should I refine?"

### Step 1: Fetch the issue

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees
```

以下を把握します:
- **サマリ**: 何をするのか
- **背景 / 課題**: なぜ必要か
- **目的 / 成功基準**: 完了の定義
- **スコープ**: 対象範囲
- **やること**: 実施事項リスト
- **完了条件**: 検証可能な受け入れ条件

Issue 本文が不足・曖昧な場合は、次の Step に進む前にユーザーに確認します。

### Step 2: Produce structured requirements (要件定義)

Issue の内容を以下の構造で整理します:

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

要件定義の結果をユーザーに提示し、確認を得てから Sub-Issue 作成に進みます。

### Step 3: Split into Sub-Issues (実装計画)

機能要件とレイヤーに基づいて、Sub-Issues に分割します。

分割の基準:
- **1 Sub-Issue = 1 レイヤーの 1 機能単位** (例: Backend のエンドポイント追加、Frontend のコンポーネント実装)
- 依存関係がある場合は `Blocked by #<issue>` を本文に記載
- 各 Sub-Issue は独立してレビュー・リバートできる粒度にする

Sub-Issue のタイトル形式: `<type>(<scope>): <要件の概要> (Sub-Issue of #<ISSUE_NUMBER>)`

Sub-Issue の本文テンプレート:

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

各 Sub-Issue を作成します:

```bash
gh issue create \
  --title "<sub-issue-title>" \
  --body "$(cat <<'EOF'
<sub-issue-body>
EOF
)"
```

### Step 4: Link Sub-Issues to the parent

親 Issue にコメントを追加して Sub-Issue の一覧を記録します:

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

作成した Sub-Issues の URL 一覧をユーザーに返します。

---

## Safety rules

- Issue 本文が不明瞭な場合は Sub-Issue を作成せず、ユーザーに確認する
- Sub-Issue は親 Issue より細かい粒度にする（親をそのまま複製しない）
- 1 つの Sub-Issue が大きすぎる場合（やること が 5 項目超）はさらに分割を検討する
