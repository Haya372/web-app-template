# E2E テストガイド

## 概要

E2E テストは [Playwright](https://playwright.dev/) を使用し、`e2e/` ディレクトリに配置されています。ブラウザ・React フロントエンド・Go バックエンド・PostgreSQL を横断したユーザーシナリオを自動検証します。

## ディレクトリ構成

バックエンドイメージのビルドに使用する `go-backend/Dockerfile` は `go-backend/` ディレクトリに配置されています。

```
e2e/
  docker-compose.yml       # E2E 環境（DB + バックエンド + フロントエンド）
  Dockerfile.frontend      # フロントエンドのマルチステージビルド（ビルドコンテキスト: モノレポルート）
  nginx.e2e.conf           # フロントエンドコンテナ用 nginx SPA 設定
  playwright.config.ts     # Playwright 設定
  package.json
  tsconfig.json
  tests/
    signup.spec.ts         # サインアップ正常系・メール重複エラー
    login.spec.ts          # ログイン正常系・パスワード誤りエラー
    post.spec.ts           # 投稿作成（fixme: 投稿作成 UI の実装待ち）
    auth-guard.spec.ts     # 未認証リダイレクト（fixme: 認証ガードの実装待ち）
```

## ローカルでの実行方法

### 前提条件

- Docker Desktop が起動していること
- `pnpm` がインストール済みであること（Node 24.x、corepack 有効）

### クイックスタート

モノレポルートから:

```bash
pnpm test:e2e
# または直接:
pnpm --filter e2e test:e2e
```

npm-run-all2 の pre/post フックにより、以下の順序で自動実行されます:

| フック / スクリプト | 内容                                                      |
|---------------------|-----------------------------------------------------------|
| `pretest:e2e`       | `docker compose up --build -d` → `wait-on localhost:3000` |
| `test:e2e`          | `playwright test`                                         |
| `posttest:e2e`      | `docker compose down -v`                                  |

スタックを起動済みの状態でテストだけ実行する場合（開発時の高速ループ）:

```bash
pnpm --filter e2e docker:up    # スタックを起動
pnpm --filter e2e test          # テストのみ実行（pre/post フックなし）
pnpm --filter e2e docker:down  # スタックを停止
```

### レポートの表示

```bash
pnpm --filter e2e test:report
```

## 環境変数

| 変数名             | デフォルト値             | 説明                                       |
|--------------------|--------------------------|--------------------------------------------|
| `E2E_BASE_URL`     | `http://localhost:3000`  | Playwright がアクセスするベース URL        |
| `E2E_API_BASE_URL` | `http://localhost:8080`  | テストセットアップ時に使用するバックエンド API の URL |

## 新機能追加時の E2E テスト

新機能を実装する際は、同じ PR 内で `e2e/tests/` 配下のテストも追加・更新してください。

1. **スペックファイルは機能単位で作成する**（例: 投稿機能なら `post.spec.ts`、認証ガードなら `auth-guard.spec.ts`）
2. **未実装の UI が必要なテストは `test.fixme()` でマークする** — 不足している前提条件をコメントに記載する
3. **テストデータは API 経由でセットアップする** — `beforeEach` またはテスト内で `fetch()` を使ってバックエンド REST API を直接呼び出す。テスト間で状態を共有しない
4. **セレクタは `getByRole` や `getByLabel` を優先する** — CSS クラスや test ID より変更に強い
5. **テストは独立して実行できるようにする** — 各テストが単独で動作する状態を維持する

## CI

以下のパスが変更された PR でテストが自動実行されます:

- `go-backend/**`
- `apps/react-frontend/**`
- `packages/**`
- `e2e/**`

ワークフローは `.github/workflows/ci-e2e.yml` に定義されています。Playwright のレポートは毎回 CI アーティファクトとしてアップロードされます。
