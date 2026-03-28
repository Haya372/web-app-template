# テストガイドライン

## 目的

バックエンド・フロントエンド・E2E を横断した統一的なテスト方針を定義する。全レイヤーで **命令網羅（Statement Coverage / C0）** を基本指標とし、すべての実行可能なステートメントが少なくとも 1 つのテストケースで到達されることを最低保証として設ける。

---

## カバレッジ方針

| 指標 | 目標値 | 適用範囲 |
|---|---|---|
| 命令網羅（C0） | ≥ 80% | 全レイヤー共通の最低ライン |
| 命令網羅（C0） | ≥ 90% | クリティカルなユースケース（認証・決済・データ変換等） |

**命令網羅（C0）の定義:** コード中のすべての実行可能なステートメント（代入・関数呼び出し・return 等）が、少なくとも 1 つのテストケースで実行されること。

テストを実装する前に対象ファイルの**全実行可能ステートメントをリストアップし**、各ステートメントに対応するテストケースが存在するかを確認する。カバレッジ目標を下回る場合は、PR の説明に理由を記載してレビュアーの合意を得ること。

---

## バックエンド（Go / testify）

### 層別テスト戦略

| 層 | テスト種別 | Make ターゲット | ビルドタグ |
|---|---|---|---|
| `internal/domain` | ユニットテスト | `make test-unit` | なし |
| `internal/usecase` | ユニットテスト | `make test-unit` | なし |
| `internal/infrastructure` | 統合テスト | `make test-integration` | `//go:build integration` |

### 共通規約

- **パッケージ名:** 外部テストパッケージ（`package foo_test`）を使い、公開 API のみをテストする
- **ファイル配置:** ソースファイルと同じディレクトリに隣接して置く（`foo.go` → `foo_test.go`）
- **テスト関数名:** `TestFunctionName_HappyCase` / `TestFunctionName_FailureCase` / `TestFunctionName_ErrorCase`
- **テストケース名・コメント:** 英語で記述する
- **複数ケースにはテーブルドリブンを必須とする**
- `require` はテストを即座に中断すべきエラーチェックに使い、`assert` は非致命的なアサーションに使う

### ユニットテスト（Domain 層）

モック不要。純粋なロジックのみ。境界値・無効入力をテーブルドリブンで網羅する。

```go
package vo_test

import (
    "testing"
    "github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPassword_FailureCase(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {name: "password length under 8 characters", input: "passwor"},
        {name: "empty string", input: ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            password, err := vo.NewPassword(tt.input)
            require.Error(t, err)
            assert.Nil(t, password)
        })
    }
}
```

### ユニットテスト（UseCase 層）

`go.uber.org/mock/gomock` で外部依存をモックする。モックは `test/mock/` に格納されている。

- `gomock.NewController(t)` を使う（クリーンアップ自動）
- `EXPECT().Method(...).Return(...).Times(n)` で呼び出し回数を明示する

```go
package user_test

import (
    "context"
    "testing"
    mock_repository "github.com/Haya372/web-app-template/go-backend/test/mock/domain/entity/repository"
    mock_shared "github.com/Haya372/web-app-template/go-backend/test/mock/usecase/shared"
    "go.uber.org/mock/gomock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSignupUseCase_HappyCase(t *testing.T) {
    tests := []struct {
        name  string
        input SignupInput
    }{
        {name: "valid email and password", input: SignupInput{Email: "a@example.com", Password: "password1"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            ctx := context.Background()

            repo := mock_repository.NewMockUserRepository(ctrl)
            repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(seedUser, nil).Times(1)

            txManager := mock_shared.NewMockTransactionManager(ctrl)
            uc := user.NewSignupUseCase(repo, txManager)

            output, err := uc.Execute(ctx, tt.input)

            require.NoError(t, err)
            assert.Equal(t, tt.input.Email, output.Email)
        })
    }
}
```

### 統合テスト（Infrastructure 層）

- 先頭行に `//go:build integration` を必ず記述する
- `testDb` フィクスチャ（同パッケージの `utils_test.go` が提供）を使用する
- 各テスト関数の末尾で `testDb.Cleanup()` を呼ぶ（`t.Cleanup` ではない）
- テストデータはテスト内でインラインでシードし、テスト間で状態を共有しない

```go
//go:build integration

package repository_test

import (
    "context"
    "testing"
    "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/repository"
    "github.com/stretchr/testify/assert"
)

func TestUserRepository_Create_HappyCase(t *testing.T) {
    target := repository.NewUserRepository(testDb.DbManager())
    tests := []struct {
        name string
        user entity.User
    }{
        {name: "create success", user: entity.ReconstructUser(...)},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            result, err := target.Create(ctx, tt.user)
            assert.Nil(t, err)
            assert.NotNil(t, result)
        })
    }
    testDb.Cleanup()
}
```

### テスト実行コマンド

```bash
# go-backend/ ディレクトリから

make test-unit          # Domain + UseCase（DB 不要）
make test-integration   # Infrastructure（Postgres on :55432 が必要）
make test-coverage      # 統合テスト + カバレッジプロファイル出力

# Postgres が必要な場合
docker compose -f docker-compose.yml up -d db
```

---

## フロントエンド（React / Vitest）

### テスト戦略

| 対象 | パターン | ランナー |
|---|---|---|
| コンポーネント・ページ | `createRoot` + ネイティブ DOM API | `pnpm test:agent` |
| カスタムフック | 最小限のホストコンポーネント内で呼び出す | `pnpm test:agent` |
| ユーティリティ関数 | 純粋なユニットテスト、テーブルドリブン | `pnpm test:agent` |

`@testing-library/react` はインストールされていない。`react-dom/client` の `createRoot` とネイティブ DOM API を使うこと。

### 共通規約

- **ファイル配置:** ソースファイルに隣接して置く（`LoginPage.tsx` → `LoginPage.test.tsx`）
- **テストケース名:** 英語で記述する
- UI スナップショットテストは採用しない（変更追跡が困難になるため）

### モック宣言の順序（Vitest の必須ルール）

1. `vi.hoisted()` — `vi.mock` ファクトリをまたいで共有する値を定義
2. `vi.mock(...)` — モジュールスコープで宣言（モック対象の import より前）
3. モック対象モジュールの import

```tsx
import React, { act } from "react"
import { createRoot } from "react-dom/client"
import { afterEach, describe, expect, it, vi } from "vitest"

const { mockNavigate } = vi.hoisted(() => ({ mockNavigate: vi.fn() }))

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) =>
    React.createElement("a", { href: to }, children),
}))

import { MyComponent } from "./MyComponent"

function mount(): HTMLDivElement {
  const container = document.createElement("div")
  document.body.appendChild(container)
  act(() => { createRoot(container).render(<MyComponent />) })
  return container
}

afterEach(() => {
  while (document.body.firstChild) document.body.removeChild(document.body.firstChild)
  vi.unstubAllEnvs()
})

describe("MyComponent", () => {
  it("renders a submit button", () => {
    const container = mount()
    expect(container.querySelector("button")).not.toBeNull()
  })
})
```

### ユーティリティ関数のテーブルドリブンパターン

```ts
describe("formatDate", () => {
  const cases = [
    { name: "ISO string",   input: "2026-01-01T00:00:00Z", expected: "Jan 1, 2026" },
    { name: "empty string", input: "",                      expected: "" },
  ]
  for (const { name, input, expected } of cases) {
    it(name, () => { expect(formatDate(input)).toBe(expected) })
  }
})
```

### テスト実行コマンド

```bash
# apps/react-frontend/ ディレクトリから

pnpm test:agent                                              # 全テスト
pnpm test:agent src/features/auth/pages/LoginPage.test.tsx  # ファイル指定
pnpm test:agent -t "LoginPage — rendering"                  # テスト名フィルタ
```

---

## E2E テスト（Playwright）

E2E テストはブラウザ・フロントエンド・バックエンド・PostgreSQL を横断したユーザーシナリオを検証する。詳細は [`e2e-testing.md`](./e2e-testing.md) を参照。

主要なルールのみ再掲する:

- スペックファイルは機能単位で作成する（例: `signup.spec.ts`）
- テストデータは API 経由でセットアップし、テスト間で状態を共有しない
- セレクタは `getByRole` / `getByLabel` を優先する
- 未実装 UI に依存するテストは `test.fixme()` でマークし、前提条件をコメントに残す
- 新機能実装時は同一 PR 内で E2E テストを追加・更新する

---

## テスト実装チェックリスト

テスト提出前に以下を確認する:

### 共通

- [ ] 対象コードのすべての実行可能ステートメントをリストアップした
- [ ] すべての実行可能ステートメントが少なくとも 1 つのテストで到達される（命令網羅 / C0）
- [ ] エクスポートされた関数・メソッドすべてにテストが存在する
- [ ] ハッピーパスと失敗・エラーケースの両方をカバーしている
- [ ] テストケース名は英語で、シナリオを明確に説明している
- [ ] バグ修正には再発防止テストが含まれている（修正前に失敗すること）

### バックエンド（Go）固有

- [ ] 複数ケースにはテーブルドリブン形式を使っている
- [ ] 統合テストには `//go:build integration` タグが付いている
- [ ] 統合テストの各関数末尾で `testDb.Cleanup()` を呼んでいる
- [ ] モックはコード外部の依存にのみ使い、テスト対象コードそのものをモックしていない
- [ ] `require` は致命的なエラーチェックに、`assert` は非致命的なアサーションに使っている

### フロントエンド（React）固有

- [ ] `vi.mock` はモック対象モジュールの import より前に宣言している
- [ ] `afterEach` で `document.body` をクリーンアップしている
- [ ] 環境変数をスタブした場合は `vi.unstubAllEnvs()` を呼んでいる
- [ ] `@testing-library/react` を import していない（未インストール）
- [ ] `as any` を使っていない（`as ReturnType<typeof vi.fn>` のように型付きで使う）
