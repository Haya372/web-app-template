# ADR-0007: Transaction Boundary and Propagation for go-backend

Date: 2026-02-15  
Status: Accepted

---

## Context
`go-backend` でユースケース単位の原子性を維持するため、トランザクション境界と伝搬方式を統一する必要がある。

### Background
- 複数リポジトリ更新を伴うユースケースでは、部分コミット防止のため一貫した境界管理が必要である。
- 境界開始責務が曖昧な場合、レイヤー間でトランザクション制御が分散しやすい。
- ネストトランザクションの扱いを明示しないと、意図しないコミット/ロールバック挙動を招く。

### Scope
- 対象: トランザクション開始責務、下位レイヤーへの伝搬規約、ネストの扱い。
- 対象外: 分散トランザクション、Outbox、2PC、DB 製品差し替え時の詳細。

### Constraints
- `usecase` 層で業務単位の原子性を担保できること。
- `domain/usecase` 層へインフラ型を漏らさないこと。
- **ネストトランザクションを禁止し、挙動を単純かつ予測可能に保つこと。**

## Decision
`go-backend` のトランザクション方針として、以下を採用する。

- トランザクション境界は `usecase` 層で開始し、`TransactionManager` を唯一の入口とする。
- `repository` 層は `DbManager` 経由でのみ DB 操作を行う。
- トランザクション伝搬は `context` により行い、同一ユースケース内は同一トランザクションを利用する。
- トランザクション内でエラーが発生した場合はロールバックし、エラーを上位へ返す。
- 読み取り専用ユースケースは、原則として明示トランザクションを開始しない。
- **`TransactionManager.Do` のネスト呼び出しを禁止する。**
- ネスト要求が発生した場合は、savepoint 方針を含む後続 ADR が承認されるまで実装しない。

運用ルール:
- 複数更新を伴うユースケースでは、境界開始箇所を PR 説明に明記する。
- ネストに該当する設計はレビューで差し戻す。
- `DbManager` 迂回の直接接続取得を禁止する。

## Options

### Option A: usecase 境界 + context 伝搬 + ネスト禁止（採用）
- 概要
  - 境界開始責務を usecase に固定し、単一トランザクションを伝搬する。
- Pros
  - 原子性の責務が明確でレビューしやすい。
  - レイヤー分離を維持しやすい。
  - ネスト禁止により挙動が予測可能になる。
- Cons
  - `context` 伝搬規約の順守が必要。
  - savepoint など高度制御には別設計が必要。
- 想定ユースケース / 制約
  - 単一 DB の標準的な業務 API。

### Option B: repository ごとにトランザクション制御
- 概要
  - 各 repository が begin/commit/rollback を管理する。
- Pros
  - 局所的には実装が単純。
- Cons
  - 複数 repository を跨ぐ原子性が崩れやすい。
  - 境界責務が分散し、整合性レビューが難しい。
- 想定ユースケース / 制約
  - 単一更新中心で整合性要件が低い処理。

### Option C: explicit Unit of Work 導入
- 概要
  - UoW オブジェクトを明示受け渡しして管理する。
- Pros
  - 伝搬が明示的で追跡しやすい。
  - 将来の savepoint 拡張に対応しやすい。
- Cons
  - インターフェースと実装コストが増える。
  - 現時点では過剰設計になりやすい。
- 想定ユースケース / 制約
  - 複雑なトランザクション制御が常時必要なシステム。

## Rationale
判断軸は、原子性保証、責務明確化、運用単純性である。  
Option A は境界責務を固定しつつ挙動を単純化でき、チームでのレビューと運用に適している。  
ネスト禁止を明示することで、不確定な savepoint 挙動やコミット順序の曖昧さを排除できる。

## Consequences
- Positive
  - ユースケース単位で一貫した原子性を担保できる。
  - レイヤー境界とレビュー観点が明確になる。
  - ネスト禁止により障害時の原因切り分けが容易になる。

- Negative
  - ネストが必要な複雑処理は別途設計が必要。
  - `context` 伝搬の規約違反には注意が必要。

- Migration / Follow-up
  - PR チェック項目に「トランザクションネストなし」を追加する。
  - ネスト要件が出た場合は savepoint ADR を先行作成する。
  - 主要ユースケースで commit/rollback と境界の統合テストを維持する。

## References
- [ADR-0003: DB Framework and Tools for go-backend](./ADR-0003-DB-FRAMEWORK-AND-TOOLS-FOR-GO-BACKEND.md)
- [ADR-0004: Clean Architecture for go-backend](./ADR-0004-CLEAN-ARCHITECTURE-FOR-GO-BACKEND.md)
