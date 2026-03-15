# ADR-0009: OpenAPI Single Source of Truth and Code Generation Strategy

Date: 2026-03-15
Status: Accepted

---

## Context

### Background

- バックエンド (Go) とフロントエンド (React) の間の API 契約が明文化されておらず、バックエンドとフロントエンドの型定義が独立して手書きされていた。
- フロントエンドは `fetch` を直接呼び出し、バックエンドのレスポンス型を zod スキーマとして手書きで定義していた。
- バックエンドのハンドラ型もすべて手書きであり、スペックとコードの乖離が容易に起きる状態だった。

### Scope

- バックエンドサーバ（Echo v5）のハンドラインターフェース生成
- フロントエンド API クライアント（型付き fetch ラッパー）生成
- CI での生成整合性チェック

### Constraints

- バックエンドは Echo v5 を使用している（`func(*echo.Context) error` シグネチャ）
- `oapi-codegen` の `echo-server` ターゲットは Echo v4 向けであり、Echo v5 とは互換性がない
- OpenAPI 3.1 は `oapi-codegen` v2.4.1 でサポートされていない

## Decision

`openapi/openapi.yaml`（OpenAPI 3.0.3）をリポジトリルートに配置し、バックエンドとフロントエンドの Single Source of Truth とする。

### バックエンド

- `oapi-codegen` v2 の `chi-server` + `strict-server` ターゲットで生成する
  - `strict-server: true` — 型付きリクエスト/レスポンスオブジェクトを持つ `StrictServerInterface` を生成する
  - `chi-server: true` — `net/http` 互換の `ServerInterfaceWrapper` を生成する（`StrictServerInterface` のアダプタ）
- `echo-server` を使わない理由: oapi-codegen の echo テンプレートは Echo v4 を対象としており、本プロジェクトの Echo v5 と互換性がない
- Echo v5 のルーティングは `ServerInterfaceWrapper` の各メソッド（`func(http.ResponseWriter, *http.Request)` シグネチャ）を `echo.HandlerFunc` にラップする手書きアダプタ（`router.go`）で実現する
- 生成ファイルの出力先: `go-backend/internal/infrastructure/http/generated/server.gen.go`
- 生成ファイルは `.gitignore` に追加し、`make generate` で常に再生成する

### フロントエンド

- `@hey-api/openapi-ts` + `@hey-api/client-fetch` で型付き fetch クライアントを生成する
- 生成ファイルの出力先: `apps/react-frontend/src/generated/`
- 生成ファイルは `.gitignore` に追加し、`pnpm generate:api` で再生成する
- 各 API 関数は `baseUrl: import.meta.env.VITE_API_BASE_URL` をオプションで受け取り、`@hey-api/client-fetch` のシングルトンクライアントの設定より優先させる

### CI

- バックエンド CI (`ci-backend.yml`): `openapi/**` のパス変更もトリガーに追加。`setup-backend` アクションで `oapi-codegen` をインストールし、`make generate` で生成コードを再生成してから `go build` で検証する
- フロントエンド CI (`ci-apps-react-frontend.yml`): `openapi/**` のパス変更もトリガーに追加。既存の `pnpm prebuild` ステップ（`generate` + `typecheck`）により生成整合性を検証する

## Options

### Option A: OpenAPI + oapi-codegen (chi-server + strict-server) — 採用

- **概要**: `strict-server: true` で型安全な `StrictServerInterface` を生成、`chi-server: true` で net/http アダプタを生成。Echo v5 向けの手書きラッパーと組み合わせる
- **Pros**: 型安全、Echo 非依存の生成コード、StrictServerInterface で全エンドポイントの実装漏れをコンパイル時に検出できる
- **Cons**: Echo v5 アダプタを手書きする必要がある

### Option B: OpenAPI + oapi-codegen (echo-server)

- **概要**: `echo-server: true` ターゲットで Echo 用コードを生成する
- **Pros**: Echo 用コードが自動生成される
- **Cons**: Echo v4 専用であり、本プロジェクトの Echo v5 と互換性がない。利用不可

### Option C: OpenAPI スキーマなし（手書き型のまま）

- **概要**: 現状維持
- **Pros**: 追加ツール不要
- **Cons**: バックエンドとフロントエンドで型が乖離するリスクが残る。ドキュメントも別途管理が必要

## Rationale

Option B は Echo v5 非互換のため除外。Option C は型乖離リスクが高い。Option A は net/http 互換の生成コードと Echo v5 向けの最小限のアダプタで、型安全性と保守性を両立できる。

## Consequences

- **Positive**:
  - API 契約が `openapi/openapi.yaml` 一点に集約される
  - バックエンドのハンドラ実装漏れをコンパイル時に検出できる
  - フロントエンドの型付き API クライアントを手書きする必要がなくなる
  - CI で spec と生成コードの整合性を自動検証できる
- **Negative**:
  - Echo v5 アダプタ（`router.go` の `wrap` ヘルパー）を手書きで維持する必要がある
  - oapi-codegen が OpenAPI 3.1 未対応のため、スキーマは 3.0.3 で記述する必要がある
- **Migration**:
  - 既存の手書き型（`auth.ts`、`post.ts` の zod スキーマ）は段階的に削除し、生成型で置き換えることを推奨する

## References

- ADR-0001: Web Framework (Echo v5 採用)
- ADR-0005: API Versioning and Route Design
- ADR-0006: Error Contract and Mapping Policy
- [oapi-codegen v2 ドキュメント](https://github.com/oapi-codegen/oapi-codegen)
- [@hey-api/openapi-ts ドキュメント](https://heyapi.dev)
- Issue #62
