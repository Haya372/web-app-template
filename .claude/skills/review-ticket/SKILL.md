---
name: review-ticket
description: "Review a GitHub Issue or Sub-Issue for quality: clear acceptance criteria, appropriate scope, no ambiguity. Reports deficiencies and suggests improvements. Use when you want to validate a ticket before starting implementation, or when refining requirements."
argument-hint: "<issue-number>"
user-invokable: true
---

# Review Ticket Skill (Refinement Review)

このスキルは GitHub Issue の品質をレビューし、不備・曖昧さを指摘します。

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` に GitHub Issue 番号が含まれている必要があります（例: `42` または `#42`）。
数値部分を抽出して `ISSUE_NUMBER` に代入します。
番号が見つからない場合は止まり、ユーザーに質問します: "Which issue number should I review?"

### Step 1: Fetch the issue

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees,comments
```

Issue のすべてのセクションを精読します。

### Step 2: Review against quality checklist

以下の観点で Issue を評価します。各項目を Pass / Fail / Warning で判定します。

#### 2.1 サマリ・目的の明確さ
- [ ] サマリが1文で「何をするか」を述べているか
- [ ] 目的・成功基準が「なぜ必要か」を明確に述べているか
- [ ] ビジネス価値またはユーザー価値が示されているか

#### 2.2 スコープの妥当性
- [ ] 「対象」と「対象外」が明記されているか
- [ ] スコープが単一のイテレーションで完遂できる大きさか（目安: やること が 7 項目以内）
- [ ] スコープ外の懸念が「リスク / 相談事項」に記載されているか

#### 2.3 完了条件の検証可能性
- [ ] 完了条件がすべて客観的に検証できるか（「良くなる」などの主観的表現がないか）
- [ ] 完了条件がすべての主要機能要件をカバーしているか
- [ ] エラーケース・境界値を含む完了条件があるか

#### 2.4 やること（タスクリスト）の具体性
- [ ] 各タスクが実装者が迷わず着手できる粒度か
- [ ] タスク間の依存関係が明示されているか（順序が重要な場合）
- [ ] 技術的な実装詳細（APIエンドポイント、DBスキーマ等）が適切に記載されているか

#### 2.5 テスト観点の充実度
- [ ] 正常系シナリオが記載されているか
- [ ] 異常系・エラーシナリオが記載されているか
- [ ] 受け入れテストが実行可能な形で記述されているか

#### 2.6 リスクと依存関係
- [ ] 既知のリスクや未解決の問題が記載されているか
- [ ] 他のチームや Issue への依存関係が明示されているか

### Step 3: Produce the review report

以下の形式でレビュー結果を出力します:

```
## Ticket Review: #<ISSUE_NUMBER> — <タイトル>

### 総合評価
<Ready for implementation / Needs revision / Needs major rework>

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

ユーザーが指示した場合、レビュー結果を Issue にコメントとして投稿します:

```bash
gh issue comment $ISSUE_NUMBER --body "$(cat <<'EOF'
<review report>
EOF
)"
```

---

## Behavioral guidelines

- **具体的に指摘する**: 「曖昧」と言うだけでなく、何が曖昧でどう修正すべきかを示す
- **建設的に**: 問題を見つけることではなく、実装可能な状態にすることが目的
- **過剰要求しない**: 完璧なチケットを要求しない。実装者が安心して着手できれば十分
- **コンテキストを尊重する**: プロジェクトの慣習・制約を踏まえてレビューする
