git の変更を解析し、このリポジトリの規約に沿ったコミットメッセージを生成してコミットを実行してください。

**引数:** $ARGUMENTS（省略可 — 省略時は変更内容から自動推定）

---

## 実行手順

### 1. 変更内容を把握する

以下を並行して取得する:
- `git status` — ステージ済み・未ステージの変更一覧
- `git diff --staged` — ステージ済み差分（コミット対象）
- `git diff` — 未ステージの差分（参考）
- `git log --oneline -5` — 直近のコミット履歴（文体を合わせるため）

未ステージの変更があり、かつステージ済みの変更がない場合は、ユーザーに確認してからステージするか判断する。

### 2. コミットメッセージを生成する

**形式:**
```
<type>(<scope>): <summary> (#<issue>)
```

**type の選び方:**

| type | 用途 |
|------|------|
| feat | 新機能 |
| fix | バグ修正 |
| refactor | 動作を変えないリファクタリング |
| chore | ビルド・依存関係・設定変更 |
| docs | ドキュメントのみの変更 |
| test | テスト追加・修正 |
| perf | パフォーマンス改善 |
| ci | CI/CD 設定の変更 |

**scope:** 変更対象のコンポーネント（例: `go-backend`, `auth`, `docker`）。単一コンポーネントなら付ける。複数にまたがる場合は省略可。

**summary:** 命令形・英語で、50 文字以内。"add", "fix", "update", "remove" など動詞で始める。

**issue 番号:** `$ARGUMENTS` で渡された場合はそれを使う。渡されていない場合は、ブランチ名やコミット履歴から推定できれば付ける。不明な場合は省略する。

**例:**
- `feat(go-backend): add JWT authentication middleware (#42)`
- `fix(auth): correct token expiry calculation`
- `chore(docker): bump postgres from 18.2 to 18.3`
- `refactor: extract repository layer into separate package`

### 3. ユーザーに提案して確認する

生成したコミットメッセージをユーザーに示し、承認を求める。その際:
- 変更の要約を 1〜2 文で説明する
- 複数の候補がある場合は選択肢を提示する

### 4. コミットを実行する

承認を得たら、以下の形式で実行する:

```bash
git commit -m "$(cat <<'EOF'
<type>(<scope>): <summary> (#<issue>)

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
EOF
)"
```

コミット後、`git status` で成功を確認して結果を報告する。

---

## 注意事項

- `.env`、シークレット、認証情報が含まれるファイルはコミットしない
- コミットが空の場合（変更なし）は実行しない
- `--no-verify` は使わない。フックが失敗した場合は原因を調査する
- `git add -A` や `git add .` より、対象ファイルを明示的にステージする
