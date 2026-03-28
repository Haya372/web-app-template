# Worktree 並列開発ガイド

<!-- cspell:ignore worktree worktrees envrc direnv WORKTREE portless -->

git worktree で複数ブランチを並列開発する際のポート自動割り当てについて説明します。
設計の詳細は [ADR-0012](../decisions/ADR-0012-WORKTREE-PORT-ISOLATION.md) を参照してください。

## 仕組み

各 worktree に **WORKTREE_OFFSET**（0, 100, 200, …）が割り当てられ、サービスのポートは `BASE_PORT + WORKTREE_OFFSET` で決まります。

| サービス | ベースポート | offset=0 | offset=100 |
|---------|------------|---------|-----------|
| Go バックエンド | 8080 | 8080 | 8180 |
| React フロントエンド | 3000 | 3000 | 3100 |
| PostgreSQL | 55432 | 55432 | 55532 |

オフセット台帳は `git rev-parse --git-common-dir` で取得した git 共通ディレクトリ（全 worktree で共有）に保存されます。

## セットアップ

新しい worktree を作成したら、必ず `ports:init` を実行してください。

```bash
# 1. 新しい worktree を追加
git worktree add ../feature-branch feature/my-branch

# 2. worktree に移動
cd ../feature-branch

# 3. ポートを自動割り当て
mise run ports:init
# → New offset assigned: WORKTREE_OFFSET=100
# → Done: APP_PORT=8180 VITE_PORT=3100 DB_PORT=55532
```

`ports:init` は `.env.local` に以下を書き込みます（mise が自動ロード）:

```
WORKTREE_OFFSET=100
APP_PORT=8180
VITE_PORT=3100
DB_PORT=55532
CORS_ALLOW_ORIGINS=http://localhost:3100
```

## 割り当て状況の確認

```bash
# 現在の割り当て一覧
mise run ports:list
# Registered worktrees:
#   [offset=0]   APP_PORT=8080 VITE_PORT=3000 DB_PORT=55432  /path/to/main
#   [offset=100] APP_PORT=8180 VITE_PORT=3100 DB_PORT=55532  /path/to/feature-branch

# この worktree のポート確認
cat .env.local | grep -E '^(WORKTREE_OFFSET|APP_PORT|VITE_PORT|DB_PORT)='
```

## 冪等性

`ports:init` は何度実行しても安全です。既に登録済みの worktree では同じオフセットを再利用します。

```bash
mise run ports:init
# Already registered: WORKTREE_OFFSET=100
```

## クリーンアップ

worktree を削除する前に必ず `ports:clean` を実行してオフセットを解放してください。

```bash
# 1. この worktree のポート登録を削除
mise run ports:clean
# → Cleaned: /path/to/feature-branch

# 2. worktree を削除
cd ../main
git worktree remove ../feature-branch

# 3. 不要な参照を削除（任意）
git worktree prune
```

## 単一 worktree での動作

`.env.local` がない、または `WORKTREE_OFFSET` が未設定の場合、各サービスはデフォルトポートで起動します。
**既存の開発環境への変更は不要です。**

```bash
# ports:init を実行しない場合、デフォルトポートで動作
go run ./cmd/server  # → :8080
pnpm dev             # → :3000
docker compose up    # → :55432
```

## E2E テストとの連携

E2E テストは `E2E_BASE_URL`・`E2E_API_BASE_URL` 環境変数を参照します。
worktree でポートを変更した場合は、これらも合わせて設定してください。

```bash
# .env.local に追記（ports:init 後）
E2E_BASE_URL=http://localhost:3100
E2E_API_BASE_URL=http://localhost:8180
```

## トラブルシューティング

### ポート競合が解消されない

台帳を直接確認します:

```bash
cat "$(git rev-parse --git-common-dir)/worktree-ports"
```

重複エントリがある場合は手動で削除してから `ports:init` を再実行してください。

### `ports:clean` を忘れて worktree を削除してしまった

台帳から手動でエントリを削除します:

```bash
REGISTRY="$(git rev-parse --git-common-dir)/worktree-ports"
# 削除したい worktree のパスを確認
cat "$REGISTRY"
# 該当行を削除（例: /path/to/old-worktree のエントリを削除）
grep -v '^/path/to/old-worktree=' "$REGISTRY" > "$REGISTRY.tmp" && mv "$REGISTRY.tmp" "$REGISTRY"
```

### 最大 10 worktree を超えようとした

```
Error: Too many worktrees registered (max 10)
```

不要な worktree を削除して `ports:clean` を実行してからオフセットを解放してください。

## 制限事項

- 同時サポート最大 10 worktree（オフセット 0〜900）
- Storybook（port 6006）は現在対象外（今後 `6006 + WORKTREE_OFFSET` で拡張予定）
