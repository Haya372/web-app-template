---
name: start-task
description: "Start working on a GitHub Issue end-to-end: fetch the issue, check/create a design doc (design-plan), create a feature branch, implement via TDD (implement-tdd skill per task), commit, and open a PR."
argument-hint: "<issue-number>"
user-invocable: true
---

# Start Task Skill

GitHub IssueをEnd-to-Endで実装する。設計書がなければ自動生成し、TDDで実装してPRを開くまでを一気通貫で行う。

## Context

- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -5`

---

## Workflow

### Step 0: Parse arguments

`$ARGUMENTS` にGitHub Issue番号（例: `42` または `#42`）が含まれているはず。
数値部分を抽出して `ISSUE_NUMBER` に代入する。
番号が見つからない場合は止まってユーザーに確認する: "Which issue number should I start on?"

---

### Step 0.5: セッション再開の検出

Issue番号を取得した後、現在のブランチを確認する:

```bash
git branch --show-current
```

**Case A: 現在のブランチが `feature/<ISSUE_NUMBER>-*` に一致する場合（セッション再開）:**

1. ブランチ作成をスキップ（Step 3 不要）
2. Step 1 + Step 2 は通常通り実行する（Issue取得・設計書読み込み）
3. `git log --oneline --grep="checkpoint"` で完了済みタスクのコミットメッセージを取得する
   - チェックポイントコミットは implement-tdd の Step 7 で `checkpoint: <task description>` 形式で作成される
   - 例: `feat(skills): checkpoint: implement field validation for user profile`
4. Step 2 で取得した設計書の「実装順序」タスクリストと突合する
   - コミットメッセージ中の `checkpoint:` 以降のテキストとタスク説明を照合する
   - 一致するタスクを「完了済み」と判定する（部分一致・意味的な一致で判断する）
5. 設計書のタスクリストのうち、完了済みでないものを未完了タスクとする
6. 未完了タスクのみ `TaskCreate` で登録（Step 4 相当）し、最初の未完了タスクから Step 5 を開始する
7. セッション内タスクカウンタは0から開始する
   （このセッションで implement-tdd を実行したタスク数をカウント。前セッションで完了したタスクはカウントに含めない）

**Case B: 現在のブランチが `main` の場合（通常の新規開始）:**

- 既存の Step 1 〜 Step 8 をそのまま実行する

**Case C: 現在のブランチが `main` でも `feature/<ISSUE_NUMBER>-*` でもない場合:**

- ユーザーに確認する:
  ```
  現在 `<branch-name>` ブランチにいます。mainに切り替えて新規開始しますか？
  ```

---

### Step 1: Issueを取得する

```bash
gh issue view $ISSUE_NUMBER --json number,title,body,labels,assignees,comments
```

以下を把握する:
- **目的・背景**: このIssueが解決しようとしていること
- **スコープ**: 対象レイヤー（Backend / Frontend / DB / E2E）
- **完了条件**: 「完了条件」「Acceptance Criteria」セクション
- **タスクリスト**: 「やること」セクションの具体的なサブタスク

Issueの本文が空またはあいまいすぎる場合は止まってユーザーに確認する。

---

### Step 2: 設計書を確認する

Issueのコメント一覧から `## 実装設計書` を含むコメントを探す。

**設計書コメントがある場合:**
- 設計書の内容を読み込む
- Step 3 へ進む

**設計書コメントがない場合:**
- ユーザーに以下を伝える:
  ```
  設計書が見つかりませんでした。design-plan スキルを実行して設計書を作成します。
  ```
- **design-plan** スキルを実行する:
  ```
  Skill: design-plan
  Args: $ISSUE_NUMBER
  ```
- design-plan が完了したら Step 3 へ進む

---

### Step 3: featureブランチを作成する

現在 `main` ブランチにいることを確認してから分岐する。

ブランチ名フォーマット: `feature/<ISSUE_NUMBER>-<short-slug>`

`<short-slug>` はIssueタイトルを小文字・ハイフン区切り英語に変換したもの（最大5単語）。

```bash
git checkout main && git pull origin main
git checkout -b feature/<ISSUE_NUMBER>-<short-slug>
```

---

### Step 4: タスクをタスクリストに登録する

`TaskCreate` を使って Step 1 で特定したサブタスクをそれぞれ登録する。
`TaskUpdate` で依存関係（`addBlockedBy`）を設定する（例: ドメインロジック → HTTPハンドラ、Backend → Frontend の順）。

最低限以下を作成する:
- タッチするレイヤーごとに1タスク（Backend / Frontend / DB / Docs）
- UIまたはユーザー向けAPIが追加・変更される場合はE2Eタスク（`.claude/skills/start-task/e2e-context.md` 参照）
- すべての実装タスクに blocked by された「commit & PR」タスク

---

### Step 5: 各タスクをTDDで実装する

依存関係の順にタスクを処理する。各タスクの開始前に `TaskUpdate` でステータスを `in_progress` に更新する。

各タスクに対して **implement-tdd** スキルを呼び出す:

```
Skill: implement-tdd
Args: <layer> <task subject and description>
      Issue context: <issue body excerpt and design doc excerpt>
```

- `layer`: タスクが属するレイヤー（`backend` / `frontend` / `e2e`）
- implement-tdd が完了（LGTM）したら `TaskUpdate` でステータスを `completed` に更新する
- すべての実装タスクが `completed` になるまで繰り返す

**セッション内タスク数の管理:**

実装タスクの完了数をカウントし（セッション再開時もカウンタは0から開始）、以下のルールに従う:

- **3タスク完了時**: 次のタスクを開始する前にユーザーへ以下を表示する:
  ```
  ⚠️ セッション内で3タスクを完了しました。
  コンテキスト圧縮による品質低下を防ぐため、新しいセッションで続きを実行することを推奨します。

  続行する場合: 「続行」と入力してください。
  新セッションで再開する場合: 新しい会話を開始し /start-task <ISSUE_NUMBER> を実行してください。
  （チェックポイントコミット済みのため、未完了タスクから自動的に再開されます）
  ```
- ユーザーが「続行」を選んだ場合は次のタスクに進む
- **推奨は3タスク以下**だが、タスクの複雑度が低い場合はユーザーの判断で続行可能

> 根拠: Issue #171 の実例では複数VOのTDDサイクルを繰り返した際にコンテキスト圧縮が発生。
> 1タスクあたり概ね15,000〜25,000トークン消費するため、3タスクで圧縮限界に近づく。

---

### Step 6: 最終品質ゲートを通す

タッチしたレイヤーごとにコンテキストファイルの「Quality gate」コマンドをすべて実行する:
- Backend (Go): `.claude/skills/implement-tdd/backend-context.md`
- Frontend (React): `.claude/skills/implement-tdd/frontend-context.md`

いずれかが失敗またはカバレッジが低下した場合は修正してから再実行する。

---

### Step 7: コミットする

**commit** スキルを呼び出す:

```
Skill: commit
Args: $ISSUE_NUMBER
```

---

### Step 8: PRを作成する

**create-pr** スキルを呼び出す:

```
Skill: create-pr
```

PRの説明文には `Closes #<ISSUE_NUMBER>` を含める。

---

## Safety rules

- `main` ブランチで作業を始めてはいけない。必ずfeatureブランチを先に作成する。
- テストをスキップしてはいけない。テストが書けない場合（純粋な設定変更など）は理由をドキュメントに残す。
- シークレットや認証情報をコミットしてはいけない。
- いずれかのステップが失敗または不明確な結果を返した場合は、継続する前にユーザーに報告する。
- `TaskUpdate` でタスクステータスをワークフロー全体を通して最新に保つ。
