---
name: design-plan
description: "GitHub Issueから実装設計書を作成し、design-reviewerエージェントのレビューを経てIssueにコメントするスキル。Issueを起点に実装計画を立てたいとき、設計ドキュメントを作りたいとき、実装前に設計レビューを受けたいときは必ずこのスキルを使うこと。「設計して」「設計計画を立てて」「設計書を作って」「design plan」「issue から設計」などのフレーズが含まれる場合にも積極的に使用する。"
argument-hint: "<issue-number>"
model: opus
user-invocable: true
---

# Design Plan Skill

GitHub IssueをもとにAPIや機能の実装設計書を作成し、`design-reviewer`エージェントによるレビューループを経て、
承認された設計をIssueにコメントとして投稿するスキル。

## Context

- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -5`

---

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` にGitHub Issue番号（例: `42` または `#42`）が含まれているはず。
数値部分を抽出して `ISSUE_NUMBER` に代入する。
番号が見つからない場合は止まってユーザーに確認する: "どのIssue番号を設計しますか?"

---

### Step 1: Issueの要件を把握する

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees,comments
```

以下の点を把握する:
- **目的**: このIssueが解決しようとしていること
- **背景**: なぜ必要か
- **スコープ**: 対象・対象外
- **受け入れ条件**: 完了の定義
- **タスクリスト**: 具体的なやること

Issueの本文が空またはあいまいすぎる場合は止まってユーザーに確認する。

---

### Step 2: 既存実装の調査

**必ずAgentを使って調査すること。** 以下を並行して調査する:

```
Agent: codebase-explorer (subagent_type=codebase-explorer)
Prompt: "以下の目的のために既存のコードベースを調査してください。
目的: <Issueの概要>

調査対象:
1. 関連するドメインエンティティ・値オブジェクト (go-backend/internal/domain/)
2. 関連するユースケース (go-backend/internal/usecase/)
3. 関連するHTTPハンドラ・ルーティング (go-backend/internal/interface/)
4. 関連するリポジトリ実装 (go-backend/internal/infrastructure/)
5. 関連するDBスキーマ・マイグレーション (go-backend/db/)
6. 関連するADR (docs/decisions/)
7. 関連するコーディングガイドライン (docs/guidelines/)

出力: 各ファイルの概要と、設計判断に影響する既存の制約・パターン"
```

調査結果をもとに、以下を整理する:
- 変更が必要なファイル・層の一覧
- 既存の命名規則・パターン
- 関連するADRの制約
- 影響を受ける既存API・インターフェース

---

### Step 3: 実装計画（設計書）を作成する

`docs/guidelines/design-plan-template.md` を読み込み、そのテンプレートに従って日本語で実装設計書を作成する。
Issueの内容・対象の層に応じて不要なセクションは省略してよい。

---

### Step 4: design-reviewer エージェントにレビューを依頼する

**必ずAgentを使ってレビューを依頼すること。** design-reviewer エージェントを呼び出す:

```
Agent: design-reviewer
Prompt: "以下の実装設計書をレビューしてください。

Issue: #<ISSUE_NUMBER> — <Issueタイトル>

<Step 3で作成した設計書の全文>

---
参照してほしいファイル:
- docs/decisions/ (ADRs)
- docs/guidelines/
- <調査で特定した関連ファイルのパス>

---
レビュー結果は修正点を以下の重要度で分類してください:
- 🔴 Critical: 実装前に必ず対応が必要な致命的な問題（セキュリティ脆弱性、ADR違反、設計の根本的な欠陥）
- 🟠 High: 実装前に対応が必要な重要な問題（受け入れ条件の欠如、重大な設計上の懸念）
- 🟡 Moderate: 修正が推奨される問題（パフォーマンス最適化、テストカバレッジなど）
- 🟢 Low: 品質向上のための任意の提案

特に以下の観点を必ず確認してください:
- 迷いなく実装に移れるか（実装者が判断に迷う箇所がないか）
- ADRの決定に背くことはないか（docs/decisions/ を必ず確認）
- セキュリティの懸念がないか（OWASP Top 10、認証・認可設計など）
- パフォーマンスの懸念がないか（N+1、インデックス、キャッシュ戦略など）
- 命名や実装の場所が適切か（Clean Architectureの層との対応、既存命名規則との一貫性）
- その他一般的なベストプラクティス（エラーハンドリング、ロギング、テスタビリティなど）"
```

---

### Step 5: レビュー結果に基づいて設計書を修正し、再レビューを依頼する

design-reviewer の総合評価が **LGTM** になるまで以下を繰り返す:

1. レビュー結果の `🔴 Critical` および `🟠 High` の項目を確認する
2. Critical または High の指摘がある場合: 設計書を修正して Step 4 に戻る
3. `🟡 Moderate` も可能な限り対応する
4. `🟢 Low` はユーザー判断に委ねる
5. Critical/High がすべて解消されたら次へ進む

> ループを抜ける条件: 総合評価が「LGTM」または「Moderate/Low のみ残存」であること。
> Critical/High が残存している場合は必ず修正してから再レビューすること。

---

### Step 6: ユーザーにレビューを依頼する

design-reviewer が承認した設計書をユーザーに提示し、フィードバックを求める:

```
【設計書レビューのお願い】

設計書が design-reviewer の承認を得ました。

<最終的な設計書の全文>

---
ご確認をお願いします。指摘事項があればお知らせください。
問題なければ「OK」と入力していただければ、Issueにコメントします。
```

ユーザーから指摘があった場合:
- 指摘事項を設計書に反映する
- 必要であれば Step 4 に戻って再レビューを依頼する
- 再度ユーザーに提示する

ユーザーが承認したら Step 7 へ進む。

---

### Step 7: Issueに実装設計書をコメントする

承認された設計書をIssueにコメントとして投稿する:

```bash
gh issue comment $ISSUE_NUMBER --body "$(cat <<'EOF'
## 実装設計書

> このコメントは `/design-plan` スキルによって自動生成されました。

<最終的な設計書の全文>

---

**レビュー状況:** design-reviewer ✅ → ユーザー ✅
EOF
)"
```

投稿後、コメントのURLをユーザーに伝える。

---

## Safety rules

- Step 2の既存調査は必ずExploreエージェントを使うこと（直接ファイルを読むだけでは不十分）
- Step 4のレビューは必ずdesign-reviewerエージェントを呼び出すこと（自己レビューは不可）
- design-reviewerが「要修正」「要再設計」と判定した場合は、ユーザーに報告せずに修正→再レビューを繰り返すこと
- ユーザーへの提示は設計書が承認された後にのみ行う
- Issueへの投稿はユーザーの明示的な承認後にのみ行う
- ADRに違反する設計は採用しない（design-reviewerが必ず指摘するが、設計時点でも意識すること）
