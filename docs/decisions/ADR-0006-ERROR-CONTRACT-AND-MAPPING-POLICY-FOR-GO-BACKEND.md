# ADR-0006: Error Contract and Mapping Policy for go-backend

Date: 2026-02-15  
Status: Accepted

---

## Context
`go-backend` の HTTP API で返すエラー契約を統一し、クライアント互換性と安全性を確保する必要がある。

### Background
- API 利用者は、エラー原因を機械判定できる安定した形式を必要とする。
- エラー形式がエンドポイントごとに異なると、クライアント側実装と保守コストが増える。
- 公開エラーに内部情報を含めない制約を設計段階で明確化する必要がある。

### Scope
- 対象: HTTP エラーレスポンス形式、エラーコード拡張、ステータスマッピング規約。
- 対象外: ログフォーマット、監視アラート設計、内部例外クラスの実装詳細。

### Constraints
- 公開仕様として標準的なエラー表現を採用すること。
- クライアントがエラー種別を安定して分岐できること。
- 秘匿情報（内部識別子、スタック、SQL詳細）をレスポンスに含めないこと。

## Decision
`go-backend` の HTTP エラー契約として **RFC 9457 (Problem Details for HTTP APIs) を部分的に採用**する。

- エラー時の `Content-Type` は `application/problem+json` を使用する。
- RFC 9457 の標準フィールドを使用する。
  - `type`: 問題種別を表す安定した `ErrorCode`（必須、内部向け API では URI を使わない）
  - `title`: 問題の短い要約（必須）
  - `status`: HTTP ステータス（必須）
  - `detail`: 利用者向け詳細説明（任意）
  - `instance`: 問題インスタンス識別子（任意）
- 内部向け API では機械判定キーを `type` に一本化し、`code` は原則使用しない。
- 外部公開 API では `type` に URI を採用するか、internal/external でレスポンス構造を分離する。
- バリデーションエラーの項目別情報は拡張メンバー `errors` に格納する。
- 未知のエラーは `500` と汎用 `type`/`title` に正規化し、内部原因は返さない。

運用ルール:
- `type`（ErrorCode）は安定したカタログで管理し、意味変更を禁止する。
- 外部公開 API で URI 型 `type` を採用する場合は、内部 ErrorCode とのマッピングを明示管理する。
- internal 向けと external 向けで形式を分ける場合は、それぞれの契約を明示し、同一公開面内では一貫した形式を維持する。

## Options

### Option A: RFC 9457 採用（完全準拠を目指す）
- 概要
  - `type` を URI とし、RFC 9457 の想定に沿って運用する。
- Pros
  - 仕様説明がシンプルで、標準との整合性が高い。
  - 外部公開 API としての相互運用性を確保しやすい。
- Cons
  - 内部システムでは URI 公開が不要なケースでも公開設計が必要。
  - 問題種別 URI の運用負荷が増える。
- 想定ユースケース / 制約
  - 外部公開を主目的とする API。

### Option B: RFC 9457 部分的採用（採用）
- 概要
  - 標準 Problem Details を基盤にしつつ、内部向けでは `type=ErrorCode` として運用する。
- Pros
  - 標準仕様とドメイン要件の両立ができる。
  - クライアント実装の共通化がしやすい。
  - 将来の外部連携で説明コストを下げられる。
- Cons
  - 既存の独自エラー形式からの移行コストが発生する。
  - `type` 設計と運用ルールの整備が必要。
- 想定ユースケース / 制約
  - 内部 API を含みつつ、将来の外部公開も見据える長期運用 API。

### Option C: 独自 JSON 形式（`code/message/details`）
- 概要
  - 独自スキーマを全 API で統一する。
- Pros
  - 導入が容易。
  - 既存ドメインコードに寄せやすい。
- Cons
  - 標準仕様との互換性が低い。
  - 外部連携やドキュメント標準化で不利。
- 想定ユースケース / 制約
  - 閉域利用のみで標準化要求が低いケース。

### Option D: ハンドラごとに自由形式
- 概要
  - エンドポイントごとに必要な形式を返す。
- Pros
  - 局所最適化が可能。
- Cons
  - クライアント・テスト・運用の複雑性が増大する。
  - 互換性保証が困難になる。
- 想定ユースケース / 制約
  - 短命な試験 API。

## Rationale
判断軸は、標準化、機械可読性、運用一貫性である。  
RFC 9457 は HTTP API エラー表現の標準であり、`type/title/status` を軸にしつつ拡張でドメイン要件を表現できる。  
ただし `type` の URI 公開は内部システムでは不要かつ過剰な情報公開につながるため、内部向けは `ErrorCode` を採用する。  
このため本 ADR は RFC 9457 の完全準拠ではなく、要件に合わせた部分的採用を選択する。  
外部向けは URI 型 `type` または internal/external 分離を選択可能とし、公開境界での情報露出を制御する。  
独自形式は短期導入は容易だが、標準化と長期運用の観点で不利である。

## Consequences
- Positive
  - API 利用者は共通フォーマットでエラー処理を実装できる。
  - RFC 9457 の主要要素に沿うため外部連携・ドキュメント整備がしやすくなる。
  - セキュアな情報公開境界を統一しやすい。

- Negative
  - 既存の独自エラー応答を RFC 9457 形式へ移行する作業が必要。
  - `type`（ErrorCode）カタログ、および必要に応じて外部 URI とのマッピング運用が必要。

- Migration / Follow-up
  - 全ハンドラのエラー応答を `application/problem+json` に統一する。
  - `type`（ErrorCode）の一覧を `docs/operations` または API 仕様に追加する。
  - 外部公開 API がある場合は URI 型 `type` とのマッピング表、または internal/external 応答差分を明記する。
  - 統合テストで `type/title/status` と `errors` の契約検証を追加する。

## References
- RFC 9457: Problem Details for HTTP APIs
- [ADR-0004: Clean Architecture for go-backend](./ADR-0004-CLEAN-ARCHITECTURE-FOR-GO-BACKEND.md)
- [ADR-0005: API Versioning and Route Design for go-backend](./ADR-0005-API-VERSIONING-AND-ROUTE-DESIGN-FOR-GO-BACKEND.md)
