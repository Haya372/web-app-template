# ADR-0012: Worktree Port Isolation for Parallel Local Development

<!-- cSpell:ignore worktree worktrees WORKTREE portless sindresorhus direnv envrc -->

Date: 2026-03-29
Status: Accepted

---

## Context

git worktree を使って複数ブランチを同時進行する際、Go バックエンド・React フロントエンド・PostgreSQL のローカル開発サーバーがデフォルトポートを固定値で占有するため、2 つ目以降の worktree でサービスを起動しようとするとポート競合が発生する。

### Background

- 本リポジトリは git worktree を活用した並列ブランチ開発を推奨している。
- Go バックエンド（`:8080`）・React フロントエンド（`:3000`）・PostgreSQL（`:55432`）はすべてポートをハードコードしていた。
- 2 つの worktree を同時起動すると即座に競合が発生し、2 つ目の worktree のサービスが起動できない。

### Scope

- 対象: ローカル開発サーバーのポート管理（Go バックエンド・React フロントエンド・PostgreSQL）
- 対象外: CI / E2E 環境のポート設定、本番環境、Storybook（別途検討）

### Constraints

- 単一 worktree での既定の動作を変更しない（後方互換）
- 新しいアプリをモノレポに追加した場合にも同じ仕組みで対応できること
- mise をツールチェーンとして利用しており、`.env.local` は mise によって自動ロードされる

## Decision

**WORKTREE_OFFSET 環境変数 + git 共通ディレクトリのオフセット台帳** を採用する。

### 仕組み

1. **WORKTREE_OFFSET**: 各 worktree に固有の整数オフセット（0, 100, 200, …）を割り当てる。各サービスのポートは `BASE_PORT + WORKTREE_OFFSET` で決定する。

   | サービス | ベースポート | offset=0 | offset=100 | offset=200 |
   |---------|------------|---------|-----------|-----------|
   | Go バックエンド | 8080 | 8080 | 8180 | 8280 |
   | React フロントエンド | 3000 | 3000 | 3100 | 3200 |
   | PostgreSQL | 55432 | 55432 | 55532 | 55632 |

2. **git 共通ディレクトリのレジストリ**: `git rev-parse --git-common-dir` で取得できる `.git/` ディレクトリ（全 worktree で共有）に `worktree-ports` ファイルを置き、worktree パス → オフセットのマッピングを管理する。

3. **`mise run ports:init`**: worktree を台帳に自動登録し、未使用の最小オフセットを割り当てて `.env.local` に書き込む（冪等）。

4. **`mise run ports:clean`**: worktree 削除時に台帳からエントリを削除してオフセットを解放する。

5. **`mise run ports:list`**: 現在の割り当て状況を表示する。

### 運用ルール

- 新しい worktree を作成したら `mise run ports:init` を実行する。
- worktree を削除する前に `mise run ports:clean` を実行する。
- `.env.local` 内のポート関連キー（`WORKTREE_OFFSET`, `APP_PORT`, `VITE_PORT`, `DB_PORT`, `CORS_ALLOW_ORIGINS`）は `ports:init` が管理するため手動編集しない。
- 新しいサービスを追加する場合は `BASE_PORT + WORKTREE_OFFSET` の規約に従う。
- 最大 10 worktree（オフセット 0〜900）を同時サポートする。

## Options

### Option A: WORKTREE_OFFSET + git 共通ディレクトリレジストリ（採用）

- 概要
  - 全サービスに共通のオフセット変数を使い、git 共通ディレクトリで台帳管理する。
- Pros
  - 外部ライブラリ不要（mise + shell script のみ）
  - 1 変数で全サービスのポートを制御できるためスケーラブル
  - git 共通ディレクトリは全 worktree で共有されるため、専用の外部ファイルや設定が不要
  - 冪等性があり、何度実行しても安全
  - 後方互換（`.env.local` がなければデフォルトポートで動作）
- Cons
  - git 共通ディレクトリへの依存（git リポジトリ外では動作しない）
  - 最大 10 worktree という上限がある
  - worktree 削除時に `ports:clean` を忘れるとオフセットが枯渇する可能性
- 想定ユースケース / 制約
  - git worktree を活用した並列ブランチ開発

### Option B: 動的ポート割り当てライブラリ（portless / get-port / detect-port）

- 概要
  - `portless`（Node.js）、`get-port`（sindresorhus/npm）、`detect-port`（alibaba/npm）などのライブラリを使い、起動時に空きポートを自動取得する。
- Pros
  - ポートの事前登録が不要で、実行時に動的に空きポートを取得できる
  - 上限なく worktree を並列起動できる
- Cons
  - **Node.js エコシステム限定**: Go バックエンドや Docker Compose に直接適用できず、Go / Node 両対応の統一ラッパーが別途必要になる
  - **起動時のポートが確定しない**: サービスが毎回異なるポートで起動するため、CORS 設定・フロントエンドの API 向き先・Playwright の baseURL を毎回動的に解決する仕組みが必要になる
  - **npm 依存の追加**: `get-port` / `detect-port` を devDependencies に追加することになり、Go バックエンド側では別途ラッパーが必要
  - メンテナンス状況（`portless` は更新が止まっている）
- 想定ユースケース / 制約
  - Node.js のみで構成されたプロジェクト、または起動時のポート通知機構（サービスディスカバリ）が整備されているプロジェクト

### Option C: OS の port 0 バインディング（ランダム空きポート）

- 概要
  - サービスを `:0` でバインドさせ、OS がランダムな空きポートを割り当てる。
- Pros
  - 最もシンプル。OS が空きを保証するため競合が起きない
  - 外部依存なし
- Cons
  - 割り当てられたポートを他のサービスや開発者が知る手段がない（ログを見るか、プロセス管理ツールが必要）
  - Go の Echo や Vite は起動後に確定ポートを通知する標準的な仕組みを持たない
  - CORS 設定・`.env.local` への反映・Playwright の baseURL の動的更新が毎回必要になる
- 想定ユースケース / 制約
  - テスト用のエフェメラルサーバーや、ポートを自動検出するオーケストレーション環境（Kubernetes / systemd socket activation など）

### Option D: direnv + `.envrc` による worktree ごとのポート固定

- 概要
  - `direnv` を導入し、各 worktree のルートに `.envrc` を置いてポートを固定する。
- Pros
  - ディレクトリ移動時に自動で環境変数が切り替わる
  - 明示的でわかりやすい
- Cons
  - **direnv の追加インストールが必要**: 全開発者に `brew install direnv` と shell フック設定が必要
  - mise が既に `.env.local` を自動ロードするため、direnv は機能的に重複する
  - ポートの重複管理は手動のまま（台帳機能はない）
- 想定ユースケース / 制約
  - mise を使用していない、または複数プロジェクトで direnv を統一管理しているチーム

### Option E: Docker Compose project name + port range override

- 概要
  - `docker-compose.override.yml` や `COMPOSE_PROJECT_NAME` + `.env` の組み合わせでコンテナのポートを worktree ごとに分離する。
- Pros
  - Docker Compose の標準機能で対応できる
  - コンテナ名・ボリューム名の競合も同時に解決できる
- Cons
  - **Go バックエンドは Docker なしで起動するケースが多く**、アプリ側のポートは別途管理が必要
  - React フロントエンドのポート管理には別の手段が必要
  - 設定ファイルが worktree ごとに増えてリポジトリが煩雑になる
- 想定ユースケース / 制約
  - すべてのサービスを Docker 経由で起動し、ホストプロセスを直接使わないプロジェクト

### Option F: 手動で `.env.local` にポートを記述

- 概要
  - 開発者が各 worktree の `.env.local` に手動でポート番号を設定する。
- Pros
  - シンプルで追加ツール不要
- Cons
  - 開発者がポート番号の重複を自分で管理する必要がある
  - worktree が増えるほど手間とエラーが増える
  - 新メンバーへのオンボーディングが煩雑になる
- 想定ユースケース / 制約
  - worktree を使う開発者が 1〜2 人のみで頻度が低い場合

## Rationale

重視した判断軸:

- **Go / Node 両エコシステムへの統一対応**: Option B（動的ライブラリ）と Option C（port 0）は Node.js 側での動的ポート取得は容易だが、Go バックエンドへの適用に追加の仕組みが必要で、さらに毎回変わるポートを CORS 設定・フロントエンドの API 向き先に伝達する仕組みが重くなる。Option A はどのサービスも同じ環境変数を参照するだけで対応できる。

- **起動前にポートが確定していること**: 開発中は `.env.local` を確認するだけでどのポートで何が動いているかを把握したい。Option B / C の動的割り当てではログを調べるまでポートがわからず、CORS 設定も都度更新が必要になる。Option A は `ports:init` を実行した時点でポートが確定し、`.env.local` に記録される。

- **外部依存の最小化と mis ツールチェーンの活用**: Option B は npm パッケージの追加、Option D は direnv のインストールが必要。mise が既に `.env.local` を自動ロードしているため、Option A は既存ツールチェーンのみで完結する。

- **スケーラビリティ**: `WORKTREE_OFFSET` 一変数方式なら、新サービス追加時も同じ規約に従うだけでよい。Option E の docker-compose override や Option F の手動設定はサービスが増えるほど管理コストが増大する。

git 共通ディレクトリを台帳として使う発想は、git 自体のメタデータ管理機能を活用する設計で、外部依存ゼロで全 worktree に状態を共有できる。

## Consequences

- Positive
  - `mise run ports:init` 1 コマンドで worktree のポート設定が完了する。
  - ポート番号の手動管理から開発者を解放し、競合によるセットアップ失敗を防ぐ。
  - 新サービスも `BASE_PORT + WORKTREE_OFFSET` 規約に従うだけでスケールできる。
  - 後方互換が保たれるため、既存の開発環境を変更しなくてよい。

- Negative
  - git 共通ディレクトリに依存するため、git リポジトリ外での利用はできない。
  - `ports:clean` を実行せずに worktree を削除するとオフセットが解放されない（最大 10 の制約に影響）。
  - `CORS_ALLOW_ORIGINS` など関連する設定も `ports:init` で管理されるため、手動設定との混在に注意が必要。

- Migration / Follow-up
  - Storybook（port 6006）の対象化: 本 ADR のスコープ外だが、同じ規約（`6006 + WORKTREE_OFFSET`）で自然に拡張できる。
  - worktree を 10 個以上同時運用する必要が生じた場合は、オフセット刻みを小さくするか上限を引き上げる改訂 ADR を策定する。
  - E2E テストの並列 worktree 実行: `E2E_BASE_URL` / `E2E_API_BASE_URL` と `VITE_PORT` / `APP_PORT` の連携は今後のタスクで対応する。

## References

- `docs/guidelines/worktree-parallel-dev.md` — 並列開発セットアップガイド
- git worktree: https://git-scm.com/docs/git-worktree
- Issue #122: worktree 並列作業時のポート競合解消
