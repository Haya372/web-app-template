# ADR-0010: Connect-RPC を BFF〜Backend 間通信に採用する

Date: 2026-03-28
Status: Accepted

---

## Context

### Background

- 本リポジトリは Next.js BFF（`apps/next-frontend`）と Go バックエンド（`go-backend`）を持つモノレポ構成を想定している。
- BFF はサーバサイドで Go バックエンドと通信する必要がある。
- バックエンドにはすでに Protocol Buffers のスキーマが用意される予定であり（`proto/` ディレクトリ）、型安全な RPC が望ましい。
- ADR-0009 で REST + OpenAPI による BFF〜クライアント間の契約は確立したが、BFF〜Backend 間は別途設計が必要だった。

### Scope

- BFF（Next.js サーバコンポーネント / Route Handlers）から Go バックエンドへの呼び出し
- Go バックエンドのハンドラ実装と RPC サーバ設定
- TypeScript クライアントコード生成（`connect-es` / `@connectrpc/connect`）

### Constraints

- BFF は Node.js 上で動作する（HTTP/2 フルデュプレックスが必要な生の gRPC は Node.js での利用が複雑）
- ブラウザから直接 Go バックエンドを呼び出さない（BFF 経由のみ）
- 既存の Echo v5 の REST エンドポイントと共存させる必要がある

## Decision

BFF〜Backend 間の通信プロトコルとして **Connect-RPC** を採用する。

### バックエンド（Go）

- `connectrpc.com/connect` を使用して gRPC ハンドラを実装する
- `proto/` 配下の `.proto` ファイルから `buf generate` でサーバスタブを生成する
- Connect ハンドラは Echo v5 とは独立したポート（またはパスプレフィックス `/connect/`）でマウントする
- Connect ハンドラは `net/http` の `http.Handler` として動作するため、Echo とは `http.ServeMux` 経由で共存させる

### フロントエンド（Next.js BFF）

- `@connectrpc/connect` + `@connectrpc/connect-web` を使用する
- `buf generate` で TypeScript クライアントスタブを生成し、`apps/next-frontend/src/generated/` に出力する
- クライアントは `createConnectTransport` を用いて HTTP/1.1 でも動作させる（gRPC-Web トランスポートは使用しない）

### proto / コード生成

- `.proto` ファイルは `proto/` ディレクトリに配置し、リポジトリルートで管理する
- `buf.gen.yaml` で Go サーバスタブと TypeScript クライアントスタブを同時生成する
- 生成ファイルは `.gitignore` に追加し、CI で `buf generate` を実行して整合性を検証する

## Options

### Option A: Connect-RPC — 採用

- **概要**: Buf 社が開発した RPC プロトコル。HTTP/1.1 と HTTP/2 の両方をサポートし、gRPC・gRPC-Web・Connect の 3 プロトコルと相互通信できる
- **Pros**:
  - Node.js クライアント（`@connectrpc/connect`）が標準的な `fetch` API ベースで動作するため、BFF への組み込みが容易
  - gRPC-Web プロキシ（Envoy 等）が不要
  - Protocol Buffers による型安全なインターフェース定義
  - Go サーバも TypeScript クライアントも `buf generate` で自動生成可能
  - HTTP/1.1 動作可能なため、ローカル開発環境での設定が単純
- **Cons**:
  - Connect プロトコル自体は比較的新しく、エコシステムが gRPC より小さい
  - チームに Proto / Buf の学習コストが発生する

### Option B: 純 gRPC（`@grpc/grpc-js`）— 却下

- **概要**: Node.js 公式の gRPC クライアントを BFF に組み込む
- **Pros**: gRPC エコシステムが成熟している
- **Cons**:
  - `@grpc/grpc-js` は HTTP/2 ネイティブ実装であり、Node.js の HTTP/2 サポートが限定的な環境（Edge Runtime 等）では動作しない
  - バンドルサイズが大きく、Next.js との相性が悪い
  - TLS 設定やチャンネル管理が煩雑

### Option C: gRPC-Web — 却下

- **概要**: gRPC をブラウザや HTTP/1.1 環境向けに変換する規格。Envoy / grpc-gateway などのプロキシが必要
- **Pros**: ブラウザから直接 gRPC バックエンドを呼べる
- **Cons**:
  - BFF〜Backend 間のユースケースではプロキシが不要な中間層になる
  - Connect-RPC は gRPC-Web プロトコルとも互換性があるため、gRPC-Web 単体を選ぶメリットがない
  - インフラ構成が複雑になる

### Option D: REST（ADR-0009 の継続）— 却下

- **概要**: BFF も REST / OpenAPI でバックエンドを呼ぶ
- **Pros**: 既存の OpenAPI 基盤（ADR-0009）と統一できる
- **Cons**:
  - BFF〜Backend 間は型安全な RPC が望ましく、REST は手書きクライアントが必要になる可能性がある
  - ストリーミング RPC が必要になった際に対応できない

## Rationale

BFF は Node.js 環境で動作するため、純 gRPC（`@grpc/grpc-js`）の HTTP/2 依存は障壁になる。gRPC-Web はプロキシが必要で構成が複雑になる。REST（Option D）はストリーミングへの拡張性に欠ける。

Connect-RPC は HTTP/1.1 フォールバックを持つため BFF での利用が容易であり、gRPC・gRPC-Web との互換性も保持している。`buf` ツールチェーンによる型安全なコード生成は、REST + OpenAPI（ADR-0009）と同様の開発体験を BFF〜Backend 間にも提供できる。

## Consequences

- **Positive**:
  - BFF〜Backend 間が型安全な RPC インターフェースで定義される
  - gRPC-Web プロキシが不要でインフラ構成がシンプルになる
  - 将来的にサーバストリーミングが必要になった場合も Connect で対応可能
  - `buf generate` で Go・TypeScript 両側のコードを自動生成できる
- **Negative**:
  - チームに Protocol Buffers と `buf` ツールチェーンの習得コストが発生する
  - Proto スキーマの追加・変更時に生成コマンドの実行が必要になる
- **Migration**:
  - 既存の REST エンドポイント（ADR-0009 の OpenAPI）はフロントエンド〜BFF 間の通信に引き続き利用する
  - BFF〜Backend 間のみ段階的に Connect-RPC へ移行する

## References

- ADR-0001: Web Framework (Echo v5 採用)
- ADR-0009: OpenAPI Single Source of Truth and Code Generation Strategy
- [connectrpc.com — Connect-RPC ドキュメント](https://connectrpc.com)
- [buf.build — Buf ツールチェーン](https://buf.build)
- Issue #57（親 Issue）
- Issue #102
