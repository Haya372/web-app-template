---
name: implement-tdd
description: "1つのタスクをTDD（Red→Green→Refactor）サイクルで実装し、code-reviewerエージェントのレビューがLGTMになるまで繰り返すスキル。start-taskから呼び出されるほか、単独で特定タスクの実装にも使える。「TDDで実装して」「テスト書いてから実装して」「このタスクを実装して」などのフレーズでも積極的に使用する。"
argument-hint: "<layer> <task-description>"
model: sonnet
user-invocable: true
---

# Implement TDD Skill

1タスクをTDD（Red → Green → Refactor）サイクルで実装する。
レビューがLGTMになるまで実装とレビューを繰り返す。

## 引数

`$ARGUMENTS` から以下を読み取る:

- **layer**: `backend` / `frontend` / `e2e` のいずれか
- **task description**: タスクの内容（issue番号・タスク名・説明など）

layer が不明な場合はタスクの内容から推測する。

---

## Workflow

### Step 1: コンテキストファイルを読む

layer に応じたコンテキストファイルを読み込む:

- Backend (Go): `.claude/skills/implement-tdd/backend-context.md`
- Frontend (React): `.claude/skills/implement-tdd/frontend-context.md`
- E2E: `.claude/skills/implement-tdd/e2e-context.md`

---

### Step 2: テストを書く（Red）

**test-writer** エージェントを呼び出す:

```
Agent: test-writer
Model: sonnet
Prompt: コンテキストファイルの「Test-writer prompt」セクションを参照してプロンプトを構築する。
        タスクの内容とissueコンテキストを埋め込む。
        TDDのRedフェーズ — テストはこの時点で必ず失敗すること。
```

テストを書いたら、コンテキストファイルのテスト実行コマンドでテストが**失敗する**ことを確認する。
失敗しない場合は test-writer に再依頼する。

---

### Step 3: 実装する（Green）

**code-implementer** エージェントを呼び出す:

```
Agent: code-implementer
Model: sonnet
Prompt: コンテキストファイルの「Implement prompt」セクションを参照してプロンプトを構築する。
        タスクの内容とissueコンテキストを埋め込む。
        テストはすでに書かれているので、テストをグリーンにすることだけに集中する。
```

実装後、テスト実行コマンドで全テストが**通過する**ことを確認する。
通過しない場合は code-implementer に再依頼する。

**2回以上再依頼してもテストが通過しない場合（実装に行き詰まったとき）**は、
**design-advisor** エージェントに設計相談をしてから再試する:

```
Agent: design-advisor
Model: opus
Prompt: 以下の実装に行き詰まっている。設計上の問題点や解決策を教えてほしい。
        - Task: <タスクの内容>
        - テストの内容: <失敗しているテストの概要>
        - 試みたアプローチ: <code-implementer が試みた実装方針>
        - エラー内容: <テスト失敗・コンパイルエラーの内容>
```

design-advisor の回答を踏まえて、code-implementer に再依頼する。

---

### Step 4: レビューを依頼する（Refactor）

Red → Green サイクルが完了し、全テストが通過していることを確認してから **code-reviewer** エージェントを呼び出す:

```
Agent: code-reviewer
Model: sonnet
Prompt: コンテキストファイルの「Review prompt」セクションを参照してプロンプトを構築する。
        以下のタスクコンテキストを必ず含める:
        - Task: <タスク名と説明>
        - Acceptance criteria: <完了条件の箇条書き>
        - Design constraints: <ADR参照・設計上の制約など>
```

---

### Step 5: レビュー結果を確認してループ

レビュー結果を確認する:

- **Blockerがある場合**: 指摘事項を修正し、テスト実行コマンドで全テストが通過することを確認してから Step 4 に戻る
- **Warningのみの場合**: 可能な限り対応して Step 4 に戻る（任意）
- **LGTMの場合**: Step 6 へ進む

> ループを抜ける条件: code-reviewer が「LGTM」と判定すること。
> Blockerが残っている限り必ず修正→全テスト通過確認→再レビューを繰り返すこと。

---

### Step 6: 品質ゲートを通す

コンテキストファイルの「Quality gate」セクションのコマンドをすべて実行する。
いずれかが失敗した場合は修正してから再実行する。

---

### Step 7: チェックポイントコミット

品質ゲートがすべて通過したら、**commit** スキルを呼び出して中間コミットを作成する。

issue番号は引数またはブランチ名（例: `feature/42-foo` → `#42`）から取得する:

```
Skill: commit
Args: #<issue-number> checkpoint: <task description>
```

**コミットメッセージ識別規約:**
- Args に `checkpoint:` プレフィックスを含めることで、commit スキルが summary に `checkpoint:` を含むメッセージを生成する
- セッション再開時に `git log --oneline --grep="checkpoint"` で完了済みタスクを検出可能になる

このコミットにより:
- コンテキスト圧縮後も `git log` で実装内容を参照可能
- セッション中断時に成果が失われない
- 後続タスクが前タスクの変更内容を git diff/log で確認できる

> 注意: start-task の最終コミット（start-task Step 7）とは別に、タスク単位で中間コミットを作成する。
> commit スキルは `git diff` ベースで差分を検出するため、チェックポイントコミット済みの変更は最終コミットで再度コミットされない。
> 最終的に PR マージ時に squash するかはユーザー判断に委ねる。

---

## Safety rules

- テストなしに実装を始めてはいけない（必ずRedフェーズを先に完了すること）
- テストが書けない場合（純粋な設定・ドキュメント変更など）は理由を記載してRedフェーズをスキップしてよい
- レビューはかならず code-reviewer エージェントに依頼すること（自己レビュー不可）
- Blockerが残っているままレビューループを抜けてはいけない
- テストがすべて通過していない状態でレビューを依頼してはいけない
