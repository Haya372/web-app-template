# Domain Layer

アプリケーションのビジネスロジックを実装する

```
domain
├── aggregate
│   └── repository
├── entity
│   └── repository
└── vo
```

---

## `vo/`

**Value Object**

* 識別子を持たない
* 不変（immutable）
* 値の等価性で比較される
* ドメインの型安全性・表現力を高めるために使う

### 例

* `PermissionCode`
* `Email`

> Value Object は Entity や Aggregate の構成要素として使用される。

---

## `entity/`

**永続化されるドメイン概念（Entity / Aggregate Root）**

* 一意な識別子を持つ
* ライフサイクルを持つ（生成・更新・削除）
* 永続化の対象となる
* 状態変更と不変条件（invariant）を責務とする
* 直接ミューテーションせず、`UpdateXxx` などのメソッドは常に **新しい Entity インスタンス** を返す（Immutable）
* リポジトリ経由で取得・更新する

### 例

* `User`
* `Role`
* `Organization`

> Entity は「更新の単位」であり、読み取り専用の用途や一時的な状態の束は含めない。

---

## `aggregate/`

**複数の Entity をまとめたドメインの集約（Aggregate）**

* Aggregate Root（通常は主要な Entity）を中心に構成される
* 複数の Entity・VO を束ねた整合性境界を持つ
* 不変条件（invariant）を集約単位で維持する
* 不変（immutable）— I/O を行わない
* 「状態から自明に導ける」軽量な判定ロジックを持ってよい
* Repository 経由で取得する（**更新は各 Entity のリポジトリを使う**）

Aggregate は、**複数の UseCase で共通して利用される整合した状態表現**として定義される。

### 例

* `UserPermissionAggregate`（User + Permissions の整合した集約）

```go
type UserPermissionAggregate struct {
    UserId      uuid.UUID
    User        entity.User
    Permissions []vo.Permission
}

func (a *UserPermissionAggregate) HasPermission(p vo.Permission) bool {
    return slices.Contains(a.Permissions, p)
}
```

> Aggregate の更新が必要な場合は、UseCase 内で各 Entity のリポジトリを使う：

```go
userId := agg.UserId
user := userRepository.FindById(ctx, userId)

if _, err := userRepository.Update(ctx, user); err != nil {
    return err
}
```

---

## 設計上の方針

* ドメイン層は **I/O（DB / 外部API / Cache など）に依存しない**
* 永続化や読み取り最適化に関する「接続口」はドメイン層に **インターフェース（Port）として定義**する

  * `entity/repository/`: Entity（集約ルート）を取得・保存するための I/F
  * `aggregate/repository/`: Aggregate を取得するための I/F
* 実装（DBクエリ、ORM、HTTPクライアント等）はインフラ層に置き、依存関係は **domain → interface（port） ← infra** の形にする
* Entity を肥大化させないため、複数の Entity にまたがる状態表現は Aggregate として切り出す
* Aggregate は DTO ではなく、**ドメインの言葉として命名**する

---

## 置かないもの

* Repository / Query の実装
* 外部 API / DB / Cache へのアクセス
* UI・フレームワーク依存の概念
* 可変な読み取りモデル

---

この構成により、
**「更新の概念（Entity）」「値の表現（VO）」「整合した状態の集約（Aggregate）」**
を明確に分離し、ドメイン層をシンプルかつ拡張しやすく保つ。
