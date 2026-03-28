# ADR-0011: Frontend Authentication State Management

Date: 2026-03-28
Status: Accepted

---

## Context

フロントエンドの認証状態（ログイン済みフラグ・JWT トークン）の管理方針が規約として存在しなかった。`features/auth/utils/tokenStorage.ts` にトークンの読み書きロジックはあるが、「どこで認証状態を保持し、どこで参照するか」が未定義であった。

### Background

- サーバー状態は TanStack Query で管理する規約があるが、クライアント状態（認証トークン・ログイン済みフラグ）の管理方針が未定義
- `tokenStorage.ts` で localStorage への低レベル操作は実装済みだが、React の状態と紐付いていない
- `_authenticated.tsx` にはルートガードの TODO コメントが残っており、保護ルートへの未認証アクセスが制御されていない

### Scope

- 対象: `apps/react-frontend/` のグローバル認証状態管理・ルートガード・tokenStorage の責務定義
- 対象外: バックエンドの認証実装、リフレッシュトークン機構、ソーシャルログイン、MFA

### Constraints

- Zustand 等のグローバルストアライブラリは未導入であり、導入にはチームの合意が必要
- TanStack Router の `beforeLoad` は React 外部で実行されるため、React hooks を使用できない
- 既存の `tokenStorage.ts` は純粋なユーティリティとして設計されており、この方針を維持する

## Decision

グローバル認証状態は **React Context（`AuthContext`）** で管理する。ルートガードは TanStack Router の **`beforeLoad`** で `getToken()` を直接参照して実装する。

### 管理場所

| 状態 | 管理場所 |
|------|---------|
| JWT トークン（メモリ） | `AuthContext` の `token` state |
| JWT トークン（永続化） | `tokenStorage.ts` / localStorage |
| ログイン済みフラグ | `AuthContext` の派生値（`isAuthenticated: token !== null`） |
| ユーザープロフィール | TanStack Query（別途管理） |

### AuthContext の責務

`features/auth/contexts/AuthContext.tsx` に `AuthContextValue` 型・`AuthContext`・`AuthProvider` を定義する:

- `token: string | null` — JWT トークンのメモリキャッシュ
- `isAuthenticated: boolean` — `token !== null` から派生
- `login(token: string): void` — `saveToken()` でlocalStorageに保存し、state を更新
- `logout(): void` — `removeToken()` でlocalStorageから削除し、state をクリア

`AuthProvider` は `src/routes/__root.tsx` の `RootLayout` 内に配置し、全ルートから `useAuth()` を参照可能にする。

### tokenStorage.ts の責務

`features/auth/utils/tokenStorage.ts` は localStorage への純粋な読み書き API のみを担う。コンポーネント・フックからは必ず `useAuth()` 経由で操作し、`tokenStorage` を直接呼ばない。ただし `beforeLoad`（React 外部）では `getToken()` を直接参照してよい。

### ルートガード

`_authenticated.tsx` の `beforeLoad` に以下を実装する:

```typescript
beforeLoad: () => {
  if (!getToken()) {
    throw redirect({ to: '/login' })
  }
}
```

`Route.loader` ではなく `beforeLoad` を使う理由: データ取得前にアクセス制御を行うため。`useAuth()` ではなく `getToken()` を直接参照する理由: `beforeLoad` は React レンダリング外で実行されるため hooks が使用不可であるため。

## Options

### Option A: React Context（採用）

- 概要: `features/auth/` に `AuthContext` と `AuthProvider` を定義し、`__root.tsx` で全体をラップする
- Pros
  - 標準の React API のみで実現でき、追加ライブラリ不要
  - 既存の `frontend-coding-guideline.md` の「グローバル UI 状態は Context で管理する」方針と一致
  - Provider 外での誤用を明示的なエラーで検出できる
- Cons
  - Context の更新が全子孫コンポーネントの再レンダリングを引き起こす可能性がある（認証状態は頻繁に変わらないため許容範囲）

### Option B: Zustand 等の外部ストア

- 概要: Zustand を導入してグローバルストアで認証状態を管理する
- Pros
  - セレクターによる部分購読で不要な再レンダリングを防げる
  - DevTools でのデバッグが容易
- Cons
  - 新規ライブラリの導入にはチームの合意が必要（現時点では未導入）
  - 認証状態程度の規模では過剰な抽象化になる

### Option C: TanStack Router Context のみで管理

- 概要: Router の `context` オプションに認証状態を渡し、React Context を使わない
- Pros
  - `beforeLoad` で直接参照でき、React と Router の状態を統一できる
- Cons
  - Router Context は React の外側で初期化されるため、状態変化が React の再レンダリングに反映されない
  - ログイン・ログアウト後の UI 更新に追加の仕組みが必要になる

## Rationale

認証状態の変化（ログイン・ログアウト）は UI（ヘッダーのログインボタン、保護ルートの表示など）への即座の反映が必要であり、React の再レンダリング機構と統合することが必須要件である。

Option C は Router Context のみでは React の再レンダリングに反映できないため不採用。Option B は将来的な選択肢だが、現状の規模では過剰であり、ライブラリ追加の合意コストが高い。Option A は標準 API で十分な機能を提供でき、既存のガイドラインとも一致するため採用する。

## Consequences

- Positive
  - 認証状態の保持場所が一箇所に集約され、実装者が迷わずコードを書ける
  - `tokenStorage.ts` の責務が明文化され、React state との二重管理が防止される
  - `beforeLoad` によるルートガードで、保護ルートへの未認証アクセスがフレームワークレベルで制御される

- Negative
  - 認証状態変化時に全 `AuthProvider` 配下のコンポーネントが再レンダリングされる（認証状態の変化は稀なため実害はほぼない）
  - `beforeLoad` と React Context で参照方法が異なる（`getToken()` vs `useAuth()`）ため、使い分けルールをドキュメント化する必要がある

- Migration / Follow-up
  - 既存の `useLoginForm.ts` の `saveToken()` 直接呼び出しを `login()` 経由に変更する
  - `_auth.tsx` の TODO コメントを削除する
  - リフレッシュトークン・トークン有効期限チェックは別 Issue で対応する

## References

- [docs/guidelines/frontend-coding-guideline.md](../guidelines/frontend-coding-guideline.md)
- [apps/react-frontend/CLAUDE.md](../../apps/react-frontend/CLAUDE.md)
- [TanStack Router - Authentication](https://tanstack.com/router/latest/docs/framework/react/guide/authenticated-routes)
- Issue #79: feat(docs): フロントエンドのクライアント状態管理（認証状態）方針をガイドラインに明記する
