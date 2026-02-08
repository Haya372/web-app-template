# Domain Layer

アプリケーションのビジネスロジックを実装する

```
domain
├── entity
│   └── repository
├── snapshot
│   └── query
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

> Value Object は Entity や Snapshot の構成要素として使用される。

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

## `snapshot/`

**不変な「ある時点の状態」を表すドメイン概念**

* 複数の Entity を束ねた読み取り専用の構造
* 不変（immutable）
* I/O を行わない
* 「状態から自明に導ける」軽量な判定ロジックを持ってよい
* QueryService経由で取得する（**更新はしない**）

Snapshot は、**複数の UseCase で共通して利用される状態表現**として定義される。

### 例

* `UserWithPermission`（User + Permissions の状態）
* `AccountStatusSnapshot`
* `OrganizationSettingsSnapshot`

```go
type UserWithPermission interface {
  UserId()        str
  HasPermission(pc vo.PermissionCode) bool
}

type userWithPermissionImpl struct {
  user        User
  permissions []Permission
}

func (u *userWithPermissionImpl) HasPermission(pc vo.PermissionCode) bool {
  for _, p := range u.permissions {
    if p.code == pc {
      return true
    }
  }
  return false
}

// 更新する場合はUseCase内で以下のように更新する
...
userId := userWithPermission.UserId()
user := userRepository.FindById(ctx, userId)

if _, err := userRepository.Update(ctx, user); err != nil {
  return err
}
```

> Snapshot は Entity の代替ではなく、
> **「ある時点の状態を切り取ったドメイン概念」**として扱う。

---

## 設計上の方針

* ドメイン層は **I/O（DB / 外部API / Cache など）に依存しない**
* 永続化や読み取り最適化に関する「接続口」はドメイン層に **インターフェース（Port）として定義**する

  * `repository/`: Entity（集約）を取得・保存するための I/F
  * `snapshot/query/`: Snapshot を構築するための読み取り専用 I/F
* 実装（DBクエリ、ORM、HTTPクライアント等）はインフラ層に置き、依存関係は **domain → interface（port） ← infra** の形にする
* Entity を肥大化させないため、読み取り専用の状態表現は Snapshot として切り出す
* Snapshot は DTO ではなく、**ドメインの言葉として命名**する

---

## 置かないもの

* Repository / Query の実装
* 外部 API / DB / Cache へのアクセス
* UI・フレームワーク依存の概念
* 可変な読み取りモデル

---

この構成により、
**「更新の概念（Entity）」「値の表現（VO）」「状態の切り取り（Snapshot）」**
を明確に分離し、ドメイン層をシンプルかつ拡張しやすく保つ。
