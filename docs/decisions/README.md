# Decision (ADR: Architecture Decision Records)

このディレクトリには、本リポジトリにおける **重要な設計判断（Architecture Decision Records）** を保存します。

ADR は以下を目的とします。

- なぜその設計・方針が選ばれたのかを将来に残す
- 設計判断を一箇所に集約し、ドキュメントを分散させない
- 過去の判断を前提に、次の判断を行いやすくする

本リポジトリでは、**設計の理由・比較・トレードオフはすべて ADR に記述**します。

---

## 基本方針

- ADR は **設計判断の一次情報源** とする
- 他のドキュメント（guidelines / operations 等）には理由を書かない
- 他ドキュメントからは、関連する ADR を参照するのみとする
- 設計判断を変更する場合は、既存 ADR を編集せず **新しい ADR を追加**する

---

## ADR の書き方

各 ADR は、以下のテンプレートを使用して作成します。

👉 **[`ADR-TEMPLATE.md`](./ADR-TEMPLATE.md)**

テンプレートの構成は次の通りです。

- Context  
  - Background / Scope / Constraints
- Decision  
  - 採用した方針・仕様・ルール
- Options  
  - 検討した選択肢と比較
- Rationale  
  - なぜその Decision を選んだか
- Consequences  
  - 採用による影響・トレードオフ
- References  
  - 関連 ADR / Docs / PR

---

## ステータスについて

ADR は以下のステータスを持ちます。

- **Proposed**  
  まだ議論中、またはレビュー待ちの状態
- **Accepted**  
  採用が決定し、実装してよい状態
- **Superseded**  
  後続の ADR によって置き換えられた状態
- **Deprecated**  
  現在は使用していない
---

## ファイル命名規則

ADR ファイルは、以下の形式で命名します。

```
ADR-{NUMBER}-{TITLE}.md
```


### ルール

- `ADR-` は固定
- 番号は **4桁・時系列**
- `TITLE` は **簡潔で検索しやすい英語の名詞句**
- 単語区切りは **ハイフン（`-`）** を使用
- すべて **大文字** で統一する

### 例

- `ADR-0001-MONOREPO.md`
- `ADR-0002-DOCS-STRUCTURE.md`
- `ADR-0003-SNAPSHOT.md`
- `ADR-0010-MIGRATION-STRATEGY.md`

---

## ADR を追加すべきタイミング

以下のような検討・変更を行う場合は、ADR の追加を検討してください。

- アーキテクチャや構成の変更
- ドメインモデルの設計方針決定
- パッケージ境界・依存関係の定義
- 永続化・データ構造の方針決定
- CI / CD・デプロイ方式の決定
- 運用・マイグレーション方針の決定

---

## 運用ルール（要約）

- Decision には仕様を書いてよい
- 理由・比較は Rationale に集約する
- ADR は履歴として積み上げ、上書きしない
