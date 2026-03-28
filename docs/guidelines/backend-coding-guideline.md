# Backend Coding Guideline

## 目的
web-app-template を含む各バックエンド実装で共通して守るべき設計原則をまとめた。特定の実装（例: `go-backend`）の構造を出発点にしつつ、今後追加するバックエンドでも再利用できる普遍的なルールを定義する。

## 基本原則
- 小さな単位の Pull Request を心掛け、変更目的・影響範囲・動作確認を明記する
- 静的解析・自動テストを常にパスした状態でレビューに出す（例: `make lint`, `make test`）
- DRY/YAGNI や SOLID（特に単一責任の原則）など一般的な設計原則を守り、重複・過剰設計・多責務化を避ける。抽象化は複数の利用箇所が確認できてから導入する
- ドメイン駆動設計 (DDD) の考え方を採用し、ユビキタス言語と境界づけられたコンテキストを意識して命名・構造化する
- テスト駆動開発 (TDD) を推奨し、失敗するテストを書いてから最小限の実装でグリーンにし、リファクタリングで品質を高める
- 責務の境界を明確にし、曖昧なコードはどの層に置くかを合意してから実装する
- ガイドラインを逸脱する場合は理由をドキュメント化し、チーム全体で再合意する

## プロジェクト構造
- `cmd/<service>`: アプリケーションのエントリーポイント。DI コンテナや設定の初期化、インフラ立ち上げ、Graceful Shutdown を扱う
- `internal/domain`: ビジネスルール（Entity, Value Object, Domain Service, Snapshot 等）を定義する層。外部 I/O に依存せず純粋なロジックのみを置く
- `internal/usecase`: ドメインを組み合わせて操作するアプリケーション層。Command/Query のユースケース単位で構造化し、副作用は Transaction Manager などを介して制御する
- `internal/infrastructure`: DB・メッセージング・HTTP 等のアダプタ層。Port/Interface を実装し、ログやトレーシングをここで行う
- `internal/common` などの横断的モジュールは、ログ/設定/エラーハンドリングの共通実装のみを置き、ドメイン固有ロジックを含めない
- `db`, `test` 等は任意だが、スキーマ・コード生成・統合テストを分離しておくとサービス間で再利用しやすい

## アーキテクチャ指針
- クリーンアーキテクチャを基本とし、依存方向は `domain → usecase → infrastructure` のみ許可する
- 依存注入は interface/port を介して行い、DI コンテナ（例: Wire, Fx, Uber Dig 等）で束ねる
- ドメイン層は不変条件の保持とビジネスルールに集中し、永続化や API 形式に関する知識を持たせない
- UseCase 層は「何をいつ実行するか」を記述し、副作用を TransactionManager や Repository を通じて制御する
- Command と Query を分離する CQRS パターンを採用し、読み取り系と書き込み系のユースケース・モデルを明確に分離することで責務とスケーラビリティを最適化する
- Infrastructure 層は `*_impl` 等で実装を明示し、テストではモックや Testcontainers などで差し替えやすくする

## 実装例 (Good vs Bad)

### Domain 層
- **Immutable を徹底する:** Entity/VO を直接 mutate せず、状態変更メソッドは常に新しいインスタンスを返す（例: `func (u User) UpdateStatus(...) (User, error)`). これによりテスタビリティとスレッドセーフティを保ち、副作用をユースケース層に閉じ込める。
- **Good:** Value Object で検証を担わせ、永続化やフレームワークに依存しない純粋な構造にする

```go
// good: domain/vo/password.go
func NewPassword(raw string) (*Password, error) {
    if len(raw) < 8 {
        return nil, ErrTooShort
    }
    pwd := Password(raw)
    return &pwd, nil
}
```

- **Bad:** ドメインオブジェクトから直接 SQL を発行したり、HTTP レスポンスを知っている実装

```go
// bad: domain で DB/HTTP に依存してしまう例
func (u *User) Save(ctx context.Context, db *sql.DB) error {
    _, err := db.ExecContext(ctx, "INSERT ...")
    return err
}
```

### UseCase 層
- **Good:** 依存を interface で受け取り、トランザクションやログ/トレースをここで制御する

```go
func (uc *signupUseCaseImpl) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    return uc.txManager.Do(ctx, func(ctx context.Context) error {
        user, err := entity.NewUser(...)
        if err != nil { return err }
        _, err = uc.userRepository.Create(ctx, user)
        return err
    })
}
```

- **Bad:** HTTP レスポンスを返したり、SQL を直接書いたり、複数責務を抱えるユースケース

```go
func (uc *signupUseCaseImpl) Execute(ctx context.Context, req *echo.Context) error {
    // bad: HTTP 型やレスポンスコードをここで扱っている
    if err := db.ExecContext(ctx, "INSERT ..."); err != nil {
        return c.JSON(http.StatusInternalServerError, ...)
    }
    return c.JSON(http.StatusCreated, ...)
}
```

### Infrastructure 層
- **Good:** Port の実装として副作用を閉じ込め、トレーシング・ロギングを行う

```go
func (r *userRepositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
    return runInTx(ctx, func(q sqlc.Queries) error {
        return q.CreateUser(ctx, mapToParams(user))
    })
}
```

- **NOTE:** スキーマで外部キーや CHECK 制約を張っていても、永続化層では Value Object への変換時に防御的なバリデーションを残すこと。DB から予期せぬ値が返ってきた場合に早期検知ができ、ログやトレースに兆候が残せる。

- **Bad:** UseCase や Domain から直接呼び出されるグローバル変数的な実装、または interface を介さず依存が固定化された実装

```go
var globalDB *sql.DB

func SaveUser(ctx context.Context, u *entity.User) error {
    // bad: Port を実装せず、global 変数と強く結合している
    _, err := globalDB.ExecContext(ctx, "INSERT ...")
    return err
}
```

## 命名規則

### スコープに応じた命名の詳細度

識別子の名前は**スコープの広さに比例して詳細に**する。スコープが広いほど名前だけで意図が伝わる必要があり、省略は禁止する。

| スコープ | 方針 | 例 |
|---|---|---|
| エクスポート済み関数・メソッド・型 | 省略禁止。ドメイン用語をそのまま使う | `CalculateDiscountedPrice`, `UserRepository` |
| パッケージ変数・構造体フィールド | 省略禁止。文脈から切り離しても意味が伝わること | `maxRetryCount`, `CreatedAt` |
| 関数・メソッドのローカル変数 | 役割が明確なら短縮可。ただし 1 文字は原則禁止 | `user`, `repo`, `ctx` は可。`u`, `r` は不可 |
| ループカウンタ・ごく短いクロージャ | 慣用的な 1 文字（`i`, `v`, `k`）のみ許容 | `for i, v := range items` |

```go
// good
func (uc *signupUseCaseImpl) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
    hashedPassword, err := vo.NewHashedPassword(input.Password)
    ...
}

// bad: 省略により何を表すか不明
func (u *sUCImpl) Exec(c context.Context, in SIn) (*SOut, error) {
    hp, e := vo.NewHP(in.P)
    ...
}
```

### 名前が表すべき内容

- **bool 型の変数・フィールド:** `is` / `has` / `can` / `should` プレフィックスを付ける（例: `isActive`, `hasPermission`）
- **エラー変数:** `err` を使い、複数必要な場合のみ `createErr`, `validateErr` のように接頭語で区別する
- **インターフェイス:** 実装ではなく役割を表す名前にする（例: `UserRepository`, `TransactionManager`）。`IFoo` や `FooInterface` は使わない
- **ユビキタス言語:** ドメイン用語は設計書・会話・コードで統一する。略語への読み替えは禁止（例: `user` を `u` に、`repository` を `repo` に略さない）

## 単一責任の原則（SRP）

「クラス・関数・モジュールが変更される理由は 1 つであるべき」という原則を全レイヤーで守る。

### 判断基準

- 関数・メソッドの説明に **「〜して、さらに〜する」** という `and` が必要になったら分割を検討する
- 1 つの関数が 20 行を超えたら、単一責任を守っているか見直す（目安）
- 変更理由が複数想定される場合（例: 「バリデーションロジックが変わるとき」と「DB スキーマが変わるとき」の両方で触れる関数）は分割する

### レイヤー別の責任範囲

| レイヤー | 責任 | 含めてはいけないもの |
|---|---|---|
| Domain（Entity/VO） | 不変条件の保持とビジネスルール | SQL・HTTP・ログ・外部サービス呼び出し |
| UseCase | 「何を・いつ・どの順序で」実行するか | HTTPリクエスト/レスポンス型、SQL |
| Infrastructure | 外部 I/O の実装（DB・HTTP クライアント等） | ビジネスルール、ユースケース固有のロジック |
| Handler（HTTP 層） | リクエスト解析・UseCase 呼び出し・レスポンス整形 | ビジネスロジック、トランザクション管理 |

```go
// good: Handler は解析・委譲・整形のみ
func (h *userHandler) Signup(c echo.Context) error {
    var req SignupRequest
    if err := c.Bind(&req); err != nil {
        return err
    }
    output, err := h.signupUseCase.Execute(c.Request().Context(), toInput(req))
    if err != nil {
        return err
    }
    return c.JSON(http.StatusCreated, toResponse(output))
}

// bad: Handler にバリデーションや DB 操作が混在している
func (h *userHandler) Signup(c echo.Context) error {
    email := c.FormValue("email")
    if !strings.Contains(email, "@") { // バリデーションはドメイン層の責務
        return c.JSON(http.StatusBadRequest, ...)
    }
    h.db.ExecContext(c.Request().Context(), "INSERT ...") // DB アクセスはインフラ層の責務
    return c.JSON(http.StatusCreated, ...)
}
```

## コーディングスタイル
- `go fmt`, `goimports`（必要に応じて `gofumpt`）を保存時に適用し、`golangci-lint` 等の静的解析を常に走らせる
- エラーハンドリングは `errors.Is/As` や `%w` によるラップを徹底し、HTTP レイヤーではエラーコードとレスポンスボディのマッピングを一箇所に集約する
- Value Object を積極的に使い、入力検証や型安全性をドメイン層に閉じ込める
- コンフィグ値・秘密情報は環境変数や Secrets Manager から受け取り、リポジトリにハードコードしない
- ソースコード内のコメント・テストケース名・テスト文言は英語で記述し、国際化や他チーム連携時も一貫して理解できるようにする

## テストと品質
- Domain/UseCase 層はユニットテストを必須とし、テーブルドリブンで境界条件・例外系を明示する
- Infrastructure 層は外部依存（DB, HTTP, Queue 等）との統合を確認するインテグレーションテストを実装する。Testcontainers やローカルモックサーバーで実際の接続を再現する
- Branch coverage 80%以上を共通目標とし、クリティカルなユースケースは 90%以上を推奨する。カバレッジ低下が許容される場合は理由を PR で共有する
- バグ修正では必ず再発防止テストを用意し、失敗シナリオがテストに表現されるまでマージしない

## 観測可能性
- すべてのサーバーは構造化ログとトレース ID を出力し、リクエスト単位で追跡できるようにする
- OpenTelemetry など標準的な SDK を使い、Tracer/Logger/Metrics のセットアップをエントリーポイントに集約する
- 主要メトリクス（レイテンシ、エラー率、スループット、DB クエリ数など）を可視化し、SLO/SLA を定義して監視する

## セキュリティ
- 入力値はドメイン層に渡す前に検証し、SQL/コマンドインジェクションや CSRF など基本的な攻撃ベクトルを防ぐ
- DB や外部 API へのアクセスはプリペアドステートメントや公式 SDK を用い、資格情報は Secret Store で管理する
- 依存ライブラリの CVE を定期的に確認し、Renovate/Dependabot などを使って更新を自動化する

## 運用フロー
- main ブランチに変更を取り込む際は Pull Request + レビュー + 自動テスト成功を必須とする
- リリース前にロールバック手順とリリースノートを準備し、障害時はインシデントテンプレートで記録する
- 新しいバックエンドを追加する場合は本ガイドラインを初期チェックリストとして利用し、差異があればドキュメントで明示する
