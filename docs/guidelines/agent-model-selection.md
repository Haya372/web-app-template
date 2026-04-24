# Agent Model Selection Guidelines

スキルファイル内のエージェント呼び出し時に使用するモデルの選択基準を定める。

## 基本方針

タスクの性質に応じて最適なモデルを使用し、コストと品質のバランスをとる。
エージェント呼び出しブロックには必ず `Model:` フィールドを明示すること。

```
Agent: <agent-name>
Model: sonnet  # または opus / haiku
Prompt: "..."
```

## モデル選択基準

| エージェント | 推奨モデル | 理由 |
| --- | --- | --- |
| `test-writer` | sonnet | 定型的なテストコード生成 |
| `code-implementer` | sonnet | 実装は高速・低コストで十分 |
| `code-reviewer` | sonnet | コードレビューは sonnet で十分 |
| `design-reviewer` | opus | 設計判断は高精度モデルが適切 |
| `design-advisor` | opus | 設計相談は高精度モデルが適切 |
| `codebase-explorer` | haiku (design-plan) / sonnet (codebase-audit) | 調査・検索は軽量モデルで十分。ただし複雑な調査を伴う場合は sonnet |

## 判断に迷う場合のガイド

- **定型的・機械的なタスク** (コード生成、テスト記述、検索): `haiku` または `sonnet`
- **品質・精度が重要なタスク** (コードレビュー、実装): `sonnet`
- **複雑な判断・設計が必要なタスク** (設計レビュー、設計相談): `opus`

## 注意事項

- モデル選択基準はあくまで初期案。実際の品質・コストを見て調整する。
- `codebase-explorer` を haiku にする場合、複雑な調査タスクで精度が落ちる可能性があるため、必要に応じて sonnet を選択すること。
