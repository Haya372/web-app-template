# ADR-0003: DB Framework and Tools for go-backend

Date: 2026-02-15  
Status: Accepted

---

## Context
go-backend の永続化層において、DBアクセス方式と関連ツールの選定が必要となった。

### Background
- `go-backend` は PostgreSQL を前提に API を提供する。
- 現在の実装は、`internal/infrastructure/db` と `internal/infrastructure/repository` に DB 関連処理を集約している。
- チーム開発を前提に、クエリ変更時の安全性とレビュー容易性を重視する。

### Scope
- 対象: go-backend の DB アクセス層（接続、クエリ実行、トランザクション）とスキーマ適用手段。
- 対象外: DB 製品の再選定（PostgreSQL 以外への移行判断）、運用基盤の構成そのもの。

### Constraints
- ドメイン層に ORM 固有 API を漏らさないこと。
- SQL の挙動を明示的に管理し、型不整合を早期検知できること。
- ローカル開発で再現可能なスキーマ適用手順を維持すること。

## Decision
go-backend の DB 周辺は以下を採用する。

- DB ドライバ / 接続管理: **pgx/v5**（`pgxpool` を利用）
- クエリ実装: **sqlc** による SQL からの Go コード生成
- スキーマ適用（ローカル）: **psqldef** による `db/schema/schema.sql` の適用

運用ルール:
- SQL は `go-backend/db/query/query.sql` に定義し、`sqlc` 生成コードを `internal/infrastructure/sqlc` で利用する。
- リポジトリ層は `DbManager` 経由でクエリ実行し、トランザクション文脈は `TransactionManager` で管理する。
- スキーマの正本は `go-backend/db/schema/schema.sql` とし、`make migrate-local` でローカル DB に適用する。
- DB 実装方式を変更する場合（例: ORM 導入、マイグレーション方式変更）は後続 ADR を追加する。

## Options

### Option A: pgx + sqlc + psqldef（採用）
- 概要
  - ドライバは `pgx`、SQL は手書きし `sqlc` で型付きコード生成、スキーマ適用は `psqldef` で管理する。
- Pros
  - SQL を明示的に管理でき、実行クエリの可読性が高い。
  - `sqlc` によりクエリ結果・引数の型不整合をビルド時に検知しやすい。
  - リポジトリ境界で DB 依存を閉じ込めやすく、ドメイン層への漏れを抑制できる。
  - `schema.sql` を単一の正本として扱えるため、ローカル再現が容易。
- Cons
  - SQL 変更時にローカル/CIで生成実行が必要になり、運用手順が増える（生成物自体は Git 管理対象外）。
  - ORM と比較すると開発初期の CRUD 記述量が増える場合がある。
- 想定ユースケース / 制約
  - SQL の明示性と型安全性、責務分離を重視するバックエンド。

### Option B: ORM 中心（例: GORM / ent）
- 概要
  - ORM を主軸にクエリ構築とモデル管理を行う。
- Pros
  - 単純 CRUD の実装速度が高い。
  - スキーマ・モデルの一体運用をしやすい。
- Cons
  - ORM 独自 API への依存が強まり、レイヤー境界が曖昧になりやすい。
  - 複雑クエリで最終 SQL の把握と最適化が難しくなる場合がある。
  - 導入時の抽象化方針を誤ると、将来の移行コストが増える。
- 想定ユースケース / 制約
  - CRUD 中心で、ORM 依存を許容できる体制。

### Option C: 生 SQL 手書き（生成なし）+ 手動管理
- 概要
  - `database/sql` または `pgx` で SQL を都度実行し、マッピングを手動実装する。
- Pros
  - 追加ツール依存を最小化できる。
  - SQL 制御の自由度が高い。
- Cons
  - 行マッピングや引数取り扱いのヒューマンエラーが増えやすい。
  - 型不整合検知が実行時寄りになり、レビュー負荷が上がる。
  - コード重複が増え、保守性が低下しやすい。
- 想定ユースケース / 制約
  - クエリ数が非常に少なく、単純構成を優先するケース。

## Rationale
最終判断で重視した軸は以下。

- 型安全性（実行前に不整合を検知できるか）
- SQL の明示性（クエリを読みやすく、最適化しやすいか）
- レイヤー分離（ドメイン層を DB 実装詳細から守れるか）
- 運用コスト（生成・適用手順をチームで回せるか）

Option B は実装速度に優位がある一方、go-backend の方針である境界明確化と SQL の可視性を維持しづらい。Option C は依存最小だが、型安全性と保守性のコストが高くなりやすい。Option A は生成ステップの運用を必要とするが、`pgx` + `sqlc` で型安全性と明示性を両立でき、`psqldef` によりローカルスキーマ適用も標準化できるため、現状のチーム開発とスコープに最も適合する。

## Consequences
- Positive
  - SQL とスキーマの正本が明確になり、レビュー時の判断がしやすい。
  - DB 変更による型不整合をビルド時に検知しやすい。
  - トランザクション境界をアプリケーション側で統一的に扱える。

- Negative
  - `sqlc` 生成漏れや `schema.sql` 更新漏れがあると、ローカル実行や CI で不整合が発生する。
  - ORM 比較では、開発者に SQL 記述力が要求される。

- Migration / Follow-up
  - DB スキーマ・クエリ変更時は `make generate-db-client` を実行し、生成物整合を確認する。
  - ローカル DB 反映は `make migrate-local` を標準手順とする。
  - 将来 ORM 導入やマイグレーション戦略変更を行う場合は、既存 ADR を編集せず新規 ADR で置換関係を明記する。

## References
- pgx: https://github.com/jackc/pgx
- sqlc: https://sqlc.dev/
- psqldef: https://github.com/sqldef/sqldef
- GORM: https://gorm.io/
- ent: https://entgo.io/
