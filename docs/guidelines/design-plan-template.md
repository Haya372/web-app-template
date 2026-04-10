# 実装設計書テンプレート

`/design-plan` スキルで使用する実装設計書のテンプレート。
Issueの内容・対象の層に応じて不要なセクションは省略してよい。

---

## 実装設計書: <Issue タイトル>

### 概要
<このIssueで実現する機能・変更の1〜3文サマリー>

### 設計方針
<主要な設計上の判断と、その理由を箇条書きで>
- 例: Aggregate RootをXXXとする理由
- 例: Clean Architectureのどの層に何を配置するか

---

### アーキテクチャ

#### ドメイン層 (`go-backend/internal/domain/`)

| 対象 | 変更種別 | 内容 |
|------|----------|------|
| エンティティ名 | 新規 / 変更 / 変更なし | 変更の概要 |

**モデリングの意図:**
<このドメインモデルをこう設計した理由。集約境界・不変条件・ライフサイクルの考え方など>
- 例: `Order` を Aggregate Root とし、`OrderItem` を内包する理由（整合性を Order 単位で保つため）
- 例: `Status` を値オブジェクトにする理由（状態遷移ルールをドメイン層で表現するため）

**新規/変更するドメインモデル（シグネチャレベル）:**
```go
type Foo struct {
    ID     FooID
    Name   string
    Status FooStatus
}

// コンストラクタ: 生成時の不変条件・バリデーションをここで保証する
func NewFoo(name string) (*Foo, error) {
    // バリデーション例（ここに記載するルールが実装の仕様になる）
    if name == "" {
        return nil, ErrFooNameRequired
    }
    if len(name) > 100 {
        return nil, ErrFooNameTooLong
    }
    return &Foo{ID: NewFooID(), Name: name, Status: FooStatusActive}, nil
}
```

**ドメインメソッド設計:**

| メソッド | シグネチャ | 責務 | 実装するロジック / バリデーション |
|---------|-----------|------|--------------------------------|
| DoSomething | `(f *Foo) DoSomething(arg Type) error` | <何をするか> | <状態遷移ルール・事前条件チェック・副作用など> |

```go
// メソッドのシグネチャと主要ロジックを示す
func (f *Foo) DoSomething(arg Type) error {
    // 事前条件（不変条件の維持）
    if f.Status != FooStatusActive {
        return ErrFooNotActive
    }
    // ビジネスロジック
    // ...
    return nil
}
```

**バリデーション・不変条件まとめ:**
- <フィールド/操作>: <ルール> — 例: `Name` は空文字不可・100文字以内
- <フィールド/操作>: <ルール> — 例: `Publish()` は `Status == Draft` のときのみ呼び出し可

**DBスキーマ (`go-backend/db/`):**
```sql
-- ドメインモデルと対応するテーブル定義・マイグレーション
-- 変更がない場合はこのブロックを省略
ALTER TABLE xxx ADD COLUMN yyy TEXT NOT NULL;
```

#### ユースケース層 (`go-backend/internal/usecase/`)

| ユースケース | 変更種別 | 説明 |
|-------------|----------|------|
| XxxUseCase | 新規 / 変更 | 何をするか |

**処理フロー:**
```mermaid
sequenceDiagram
    actor Client
    participant Handler
    participant UseCase
    participant Repository
    participant DB

    Client->>Handler: POST /api/v1/xxx
    Handler->>UseCase: Execute(input)
    UseCase->>Repository: FindByID(id)
    Repository->>DB: SELECT ...
    DB-->>Repository: row
    Repository-->>UseCase: Entity
    UseCase->>Repository: Save(entity)
    Repository->>DB: INSERT / UPDATE ...
    UseCase-->>Handler: output
    Handler-->>Client: 200 OK
```

**インターフェース定義:**
```go
// 追加・変更するユースケースのメソッドシグネチャ
```

#### インターフェース層 (`go-backend/internal/interface/`)

| エンドポイント | メソッド | 説明 | 認証 |
|--------------|---------|------|------|
| /api/v1/xxx  | GET / POST / PUT / DELETE | 何をするか | 要 / 不要 |

**リクエスト/レスポンス定義:**
```json
// POST /api/v1/xxx
// Request
{
  "field": "value"
}

// Response 200
{
  "id": "uuid",
  "field": "value"
}
```

#### フロントエンド (`apps/react-frontend/`)

**ページ一覧:**

| パス | ページ名 | 説明 |
|------|---------|------|
| `/xxx` | Xxx一覧 | 何ができるページか |
| `/xxx/:id` | Xxx詳細 | 何ができるページか |

**ページ遷移:**
```mermaid
stateDiagram-v2
    [*] --> 一覧ページ: /xxx
    一覧ページ --> 詳細ページ: 行クリック
    一覧ページ --> 作成ページ: 新規作成ボタン
    作成ページ --> 一覧ページ: 保存成功 / キャンセル
    詳細ページ --> 一覧ページ: 戻るボタン
```

**各ページの機能:**

- **一覧ページ (`/xxx`)**
  - [ ] <できること1>
  - [ ] <できること2>

- **詳細ページ (`/xxx/:id`)**
  - [ ] <できること1>
  - [ ] <できること2>

**UIコンポーネント設計:**

各ページで使用するコンポーネントを `packages/ui` の既存コンポーネントから選定する。
既存コンポーネントで要件を満たせない場合は `packages/ui` に新規追加し、その設計もここに記載する。

| UI要素 | 使用コンポーネント | 出典 | 備考 |
|--------|------------------|------|------|
| 一覧テーブル | `<DataTable>` | `@repo/ui` 既存 | ソート・ページネーション付き |
| 作成ボタン | `<Button variant="default">` | `@repo/ui` 既存 | |
| <UI要素名> | `<NewComponent>` | `@repo/ui` **新規追加** | 既存にないため追加 |

**新規追加コンポーネント（`packages/ui` への追加が必要な場合）:**

> 既存コンポーネントで対応できる場合はこのブロックを省略する。

```tsx
// packages/ui/src/components/new-component.tsx

// Props 定義
type NewComponentProps = {
  // ...
}

// 使用例
<NewComponent prop="value" />
```

- **追加理由:** <既存コンポーネントでは対応できない理由>
- **Radix UI プリミティブ:** <使用する Radix プリミティブ（例: `@radix-ui/react-dialog`）、不要なら省略>

#### インフラ層 (`go-backend/internal/infrastructure/`)

| 対象 | 変更種別 | 内容 |
|------|----------|------|
| XxxRepository | 新規 / 変更 | 変更の概要 |

---

### エラーハンドリング

| エラーケース | HTTPステータス | エラーコード | 対応方法 |
|-------------|--------------|------------|---------|
| <ケース> | 400 / 404 / 500 等 | <コード> | <対応> |

---

### テスト方針

| テスト種別 | 対象 | テストシナリオ |
|-----------|------|--------------|
| ユニットテスト | ドメイン層 | 正常系・異常系・バリデーション |
| 統合テスト | リポジトリ層 | DB操作の正確性 |
| E2Eテスト | APIエンドポイント | ユーザーシナリオ全体 |

---

### 実装順序

依存関係に従い、下位層から上位層の順に実装する。

1. [ ] ドメイン層 + DBマイグレーション: <具体的なタスク>
2. [ ] インフラ層（Repository実装）: <具体的なタスク>
3. [ ] ユースケース層: <具体的なタスク>
4. [ ] インターフェース層: <具体的なタスク>

---

### 影響範囲

- 変更によって影響を受ける既存機能・API
- 破壊的変更の有無（ある場合は `docs/operations` への記載が必要）
- 依存するチーム・システムへの影響

---

### 受入テスト項目

実装完了後にユーザー視点で確認する項目。UT/IT/E2E とは別に、「この機能が正しく動いている」と判断できる手動確認のチェックリスト。

- [ ] <シナリオ1: 例）新規作成フォームで必須項目を入力して送信すると、一覧に追加されて表示される>
- [ ] <シナリオ2: 例）不正な値を入力した場合、エラーメッセージが表示されて送信できない>
- [ ] <シナリオ3: 例）権限のないユーザーがアクセスすると 403 が返る>

---

### 未解決事項

- [ ] <判断が必要な事項1>
- [ ] <判断が必要な事項2>
