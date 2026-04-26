---
name: lessons-learned
description: Review the current session and record worthy lessons to the Lessons Learned issue. Use this skill when the user says "教訓を記録", "lessons learned", "/lessons-learned", or after a session ends with significant friction.
---

# Lessons Learned Skill

このセッションの会話をレビューし、記録に値する教訓があればGitHub Issueに追記する。

## 記録対象（どれか1つに該当すれば記録）

以下のカテゴリのみを対象とする。コード実装レベルのバグ修正は対象外。

| カテゴリ | 絵文字 | 該当する例 |
|---|---|---|
| 環境・ツール制約 | 🔧 | サンドボックスのブロック、`make` / CLIコマンドの失敗、キャッシュ・パーミッションエラー、環境固有の回避策が必要だった問題 |
| ワークフロー障害 | 🔄 | スキル・hookの設計ミス、エージェント間の引き継ぎロス、コンテキスト上限による会話の途切れ、想定外のツール制限 |
| 繰り返し試行の失敗 | 🔁 | 同じアプローチを2回以上試みて失敗（設計ミス・誤った前提・仕様の読み間違いなど） |

## 記録対象外（以下は除外）

- 通常のTDDサイクル内のコード修正（Redフェーズで失敗するのは正常）
- 1回で解決した実装エラー
- コードレビュー指摘への通常の対応
- 単純なQ&Aセッション
- 特に問題のなかったセッション

## Workflow

### Step 1: 親Issueを特定する

```bash
gh issue list --search 'Lessons Learned まとめ' --state open --json number,title
```

親Issueがない場合は作成する:

```bash
gh issue create --title 'Lessons Learned まとめ' --body "$(cat <<'EOF'
## Claude Code セッション教訓一覧

このIssueはClaude Codeセッションで遭遇した失敗・課題を蓄積します。

| 状態 | 意味 |
|------|------|
| 🔴 `- [ ]` | 未対応 |
| 🟡 `- [ ]` | 対応策発見済み（docs/rules/skillへの反映が未完了） |
| 🟢 `- [x]` | docs・rules・skillに反映して再発予防済み |

---

EOF
)"
```

### Step 2: セッションを分析する

会話全体を振り返り、上記の記録対象カテゴリに該当する問題を特定する。
該当がなければここで終了（Issueを更新しない）。

### Step 3: 教訓をコメントとして追記する

各教訓を以下のフォーマットで1件ずつまとめ、`gh issue comment` でコメントを追加する:

```bash
gh issue comment <number> --body "$(cat <<'EOF'
## セッション教訓追記 (YYYY-MM)

- [ ] 🔴 **[YYYY-MM] 🔧 タイトル** — 何が起きたか（1文）。回避策: 〇〇。(参照: #issue番号)
- [ ] 🔴 **[YYYY-MM] 🔄 タイトル** — ...
EOF
)"
```

**フォーマット規則:**
- 新規追加は必ず `🔴` と未チェック `- [ ]`
- カテゴリ絵文字を必ず含める（🔧 / 🔄 / 🔁）
- 回避策が不明な場合は「回避策: 未発見」と書く
- 関連Issueがある場合は `(参照: #番号)` を末尾に付ける
- 1セッションの教訓をまとめて1コメントに収める

## Safety rules

- 教訓がない場合はIssueを更新しない（空コメントを残さない）
- コード実装の詳細（型の設計・アルゴリズム）は記録しない
- 個人を特定・批判する内容は書かない
