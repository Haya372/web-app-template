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

## Safety rules

- テストなしに実装を始めてはいけない（必ずRedフェーズを先に完了すること）
- レビューはかならず code-reviewer エージェントに依頼すること（自己レビュー不可）
- Blockerが残っているままレビューループを抜けてはいけない
- テストがすべて通過していない状態でレビューを依頼してはいけない
