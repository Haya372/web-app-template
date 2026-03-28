# Frontend Coding Guideline

## 目的
React を用いてフロントエンドを実装する際に共通して守るべき設計原則をまとめた。特定のフレームワーク（Next.js, TanStack Start 等）やツールチェーンに依存しない普遍的なルールを定義し、今後追加するフロントエンド実装でも再利用できるようにする。

## 基本原則
- 小さな単位の Pull Request を心掛け、変更目的・影響範囲・動作確認を明記する
- 静的解析・自動テストを常にパスした状態でレビューに出す
- DRY/YAGNI や単一責任の原則を守り、重複・過剰設計・多責務化を避ける。抽象化は複数の利用箇所が確認できてから導入する
- コンポーネントは「何を表示するか」と「どうデータを取得・変換するか」を分離し、UI ロジックとビジネスロジックを混在させない
- ガイドラインを逸脱する場合は理由をドキュメント化し、チーム全体で再合意する

## プロジェクト構造

Feature ベースのディレクトリ構成を採用する。機能ごとにコンポーネント・フック・API クライアント・型を同じディレクトリに集約し、複数の機能から参照される共通部品のみ `src/components/` や `src/hooks/` に置く。

```
src/
  features/            # 機能単位のディレクトリ（Feature ベース構成）
    <feature-name>/
      components/      # その機能専用のコンポーネント
      hooks/           # その機能専用のカスタムフック
      api/             # その機能の API クライアント（fetch 関数等）
      types/           # その機能のドメイン型・レスポンス型定義
  components/          # 複数機能で共有する UI コンポーネント
  hooks/               # 複数機能で共有するカスタムフック
  utils/               # 純粋関数ユーティリティ（副作用を含まない）
```

- feature 内部および feature をまたぐ参照ともに、対象ファイルを直接 import する
- `index.ts` による re-export は禁止する（ファイルが増えるだけで追跡を困難にするため）
- 特定フレームワークが要求するディレクトリ（`pages/`, `routes/`, `app/` 等）はこの構成に追加する形で共存させる

## packages/ui の利用方針

アプリ側（`apps/`）のコンポーネントは `packages/ui` のコンポーネントを組み合わせて構築する。

- `packages/ui` コンポーネントに `className` prop は渡さない。デザインの揺れを防ぐため、見た目の変更は `variant` / `size` などのセマンティックな props で表現する
- レイアウト調整（余白・幅・配置など）が必要な場合は、ラップする要素に Tailwind クラスを当てる

```tsx
// good — ラッパーでレイアウトのみ調整
<div className="mt-4 w-full">
  <Button variant="primary">送信</Button>
</div>

// bad — UI コンポーネントに直接 className を渡してデザインを上書き
<Button className="mt-4 bg-red-500">送信</Button>
```

- `packages/ui` に存在しないパターンが必要になった場合は、`className` で個別対応するのではなく `packages/ui` に新しい variant を追加する

## コンポーネント設計

- **1 ファイル 1 コンポーネント** を原則とし、150 行を超える場合は分割を検討する
- ルートコンポーネントとレイアウトコンポーネントのみ default export を使い、それ以外は named export にする
- Props のドリルダウンが 2 段階を超える場合は Context または状態管理ライブラリへの移行を検討する
- UI の状態（開閉・ホバーなど）はコンポーネント内部に閉じ込め、アプリケーション状態（認証・グローバル設定など）は上位に引き上げる

**Good:**
```tsx
// features/user/components/UserCard.tsx — packages/ui を組み合わせて構築
import { Card, Avatar, Text } from "@repo/ui"

export function UserCard({ name, avatarUrl }: UserCardProps) {
  return (
    <Card>
      <Avatar src={avatarUrl} alt={name} />
      <Text>{name}</Text>
    </Card>
  )
}
```

**Bad:**
```tsx
// bad: packages/ui を使わず Tailwind を直書き、データ取得も混在、型が any
export default function UserCard({ id }: { id: any }) {
  const [user, setUser] = React.useState<any>(null)
  useEffect(() => {
    fetch(`/api/users/${id}`).then(r => r.json()).then(setUser)
  }, [id])
  return <div style={{ borderRadius: 12, padding: 16 }}>{user?.name}</div>
}
```

## 状態管理

- **ローカル状態:** `useState` / `useReducer` をコンポーネント内に閉じ込める。ただし副作用（`useEffect` など）を伴う場合はカスタムフックに切り出し、コンポーネントは返り値を使うだけにする
- **サーバー状態（キャッシュ・再フェッチ）:** TanStack Query を使用する。`useEffect` + `useState` での手動フェッチは禁止
- **グローバル UI 状態（テーマ・モーダル等）:** React Context か軽量ストア（Zustand 等）で管理する。Redux などの重量級ライブラリは導入前にチームで合意する
- 状態を重複して持たない（derived state は `useMemo` で算出し、コピーして別 state に持たない）

## 認証状態管理

> 設計判断の詳細は [ADR-0011](../decisions/ADR-0011-FRONTEND-AUTH-STATE-MANAGEMENT.md) を参照。

### 設計方針

認証状態（ログイン済みフラグ・JWT トークン）はグローバルなクライアント状態として React Context で管理する。
サーバー状態（ユーザープロフィール等）は TanStack Query で別途取得し、認証 Context に混在させない。

### 状態管理の階層

| 状態 | 管理場所 | 理由 |
|------|---------|------|
| JWT トークン | `AuthContext`（メモリ）+ `tokenStorage`（localStorage） | アプリ全体で参照が必要 |
| ログイン済みフラグ | `AuthContext` の派生値（`isAuthenticated`） | トークン有無から決定する |
| ユーザープロフィール | TanStack Query（`useCurrentUser`） | サーバー状態として別管理 |
| フォームの入力値 | `react-hook-form` | ローカル UI 状態 |

### AuthContext の責務

`features/auth/contexts/AuthContext.tsx` に以下を定義する:

- `token: string | null` — JWT トークン（メモリ上のキャッシュ）
- `isAuthenticated: boolean` — `token !== null` から派生
- `login(token: string): void` — トークンを localStorage に保存しメモリを更新
- `logout(): void` — トークンを localStorage から削除しメモリをクリア

### tokenStorage.ts の責務

`features/auth/utils/tokenStorage.ts` は localStorage への純粋な読み書き API のみを担う。
コンポーネント・フックからは必ず `useAuth()` 経由で操作し、`tokenStorage` を直接呼ばない。

```typescript
// good: Context 経由でトークン操作
const { login, logout } = useAuth()
login(token) // tokenStorage の保存 + React state の更新を一括で行う

// bad: tokenStorage を直接呼び出す（React の再レンダリングがトリガーされない）
saveToken(token)
```

### ルートガード（保護ルートへのアクセス制御）

保護ルートは TanStack Router の `beforeLoad` でアクセス制御を行う。
`beforeLoad` は React 外部で実行されるため `useAuth()` は使用不可で、`getToken()` を直接参照する。

```typescript
// src/routes/_authenticated.tsx
export const Route = createFileRoute('/_authenticated')({
  beforeLoad: () => {
    if (!getToken()) {
      throw redirect({ to: '/login' })
    }
  },
})
```

`Route.loader` ではなく `beforeLoad` を使う理由は、データ取得前にアクセス制御を行いたいためである。

### Provider の配置

`AuthProvider` は `src/routes/__root.tsx` の `RootLayout` 内に配置し、全ルートで `useAuth()` を参照できるようにする。

### フックの利用

コンポーネント・フックから認証状態を参照する場合は必ず `useAuth()` を経由する。

```typescript
// good
const { isAuthenticated, login, logout } = useAuth()

// bad: Context を直接 useContext で消費する（Provider なし時に null になりエラーが不明瞭）
const auth = useContext(AuthContext)
```

## API 呼び出し方針

- データ取得には必ず **TanStack Query** を使用する。`useEffect` + `useState` による手動フェッチは禁止
- `api/` ディレクトリには TanStack Query の `queryFn` として呼び出す純粋な fetch 関数を置き、コンポーネントから直接 `fetch` を呼ばない
- レスポンス型は TypeScript の型で定義し、`unknown` として受け取りバリデーション後に使用する
- エラーハンドリングは `useQuery` / `useMutation` のエラー状態を通じて UI に伝播する
- 環境変数でベース URL やトークンを管理し、コードにハードコードしない

```tsx
// good: features/user/api/users.ts — 純粋な fetch 関数
export async function fetchUser(id: string): Promise<User> {
  const res = await fetch(`${import.meta.env.VITE_API_BASE}/users/${id}`)
  if (!res.ok) throw new Error(`fetchUser failed: ${res.status}`)
  const data: unknown = await res.json()
  return parseUser(data) // zod や valibot で検証
}

// good: features/user/hooks/useUser.ts — TanStack Query でラップ
export function useUser(id: string) {
  return useQuery({ queryKey: ["user", id], queryFn: () => fetchUser(id) })
}

// bad: useEffect で手動フェッチ
useEffect(() => {
  fetch(`/api/users/${id}`).then(r => r.json()).then(setUser)
}, [id])
```

## TypeScript

- strict モードを有効にし、`any` の使用を禁止する。型が不明な場合は `unknown` を使って明示的に絞り込む
- `as` キャストは原則使用しない。型ガード関数（`is` / `asserts`）やバリデーションライブラリ（zod 等）で安全に型を確定させる。やむを得ず使う場合は理由をコメントで明示し、`as any` は一切禁止とする
- コンポーネントの Props は `interface` で定義する（合併型など `interface` で表現できない場合のみ `type`）
- エクスポートする関数には明示的な戻り値の型を付ける

## 命名規則

### スコープに応じた命名の詳細度

識別子の名前は**スコープの広さに比例して詳細に**する。スコープが広いほど名前だけで意図が伝わる必要があり、省略は禁止する。

| スコープ | 方針 | 例 |
|---|---|---|
| エクスポートされるコンポーネント・関数・型 | 省略禁止。役割が名前だけで伝わること | `UserProfileCard`, `fetchUserById`, `AuthTokenPayload` |
| モジュールスコープの変数・定数 | 省略禁止。文脈から切り離しても意味が伝わること | `maxRetryCount`, `defaultPageSize` |
| 関数内のローカル変数 | 役割が明確なら短縮可。ただし 1 文字は原則禁止 | `user`, `token`, `error` は可。`u`, `t`, `e` は不可 |
| ループカウンタ・ごく短いコールバック | 慣用的な 1 文字（`i`, `v`）のみ許容 | `items.map((item, i) => ...)` |

```tsx
// good
export function UserProfileCard({ userId }: UserProfileCardProps) {
  const { data: userProfile, isLoading } = useUserProfile(userId)
  ...
}

// bad: 省略により何を表すか不明
export function UPC({ uid }: UPCProps) {
  const { data: d, isLoading: ld } = useUP(uid)
  ...
}
```

### 名前が表すべき内容

- **bool 型の変数・Props:** `is` / `has` / `can` / `should` プレフィックスを付ける（例: `isLoading`, `hasError`, `isDisabled`）
- **イベントハンドラ:** `handle` + イベント対象 + 動作（例: `handleSubmitForm`, `handleDeleteItem`）
- **カスタムフック:** `use` プレフィックス必須。フックが返す状態・操作を端的に表す名前（例: `useUserProfile`, `useAuthToken`）
- **型・インターフェイス:** 役割を表す名前。`I` プレフィックスや `Type` サフィックスは使わない（例: `UserProfileCardProps`, `AuthTokenPayload`）

## 単一責任の原則（SRP）

「コンポーネント・関数・フックが変更される理由は 1 つであるべき」という原則を守る。

### 判断基準

- 関数・コンポーネントの説明に **「〜して、さらに〜する」** という `and` が必要になったら分割を検討する
- コンポーネントが「表示」と「データ取得」の両方に責任を持つ場合は、フックに切り出してコンポーネントは表示のみに集中させる
- Props が 5 個を超えてきたら、コンポーネントが複数の関心事を扱っていないか見直す（目安）

### 責任の分離パターン

| 単位 | 責任 | 含めてはいけないもの |
|---|---|---|
| コンポーネント | 受け取ったデータの表示・ユーザー操作のハンドリング | `fetch`・`useQuery` 等のデータ取得、ビジネスロジック |
| カスタムフック | データ取得・状態管理・副作用のカプセル化 | JSX のレンダリング |
| `api/` の fetch 関数 | 1 エンドポイントとの通信と型変換のみ | キャッシュ管理、UI の状態 |
| ユーティリティ関数 | 純粋な変換・計算 | 副作用・状態・ネットワーク通信 |

```tsx
// good: コンポーネントは表示のみ、データ取得はフックに分離
export function UserProfile({ userId }: UserProfileProps) {
  const { userProfile, isLoading } = useUserProfile(userId)
  if (isLoading) return <Spinner />
  return <ProfileCard name={userProfile.name} email={userProfile.email} />
}

// bad: データ取得と表示が混在している
export function UserProfile({ userId }: UserProfileProps) {
  const [user, setUser] = useState(null)
  useEffect(() => {          // データ取得はフックの責務
    fetch(`/api/users/${userId}`).then(r => r.json()).then(setUser)
  }, [userId])
  return <div>{user?.name}</div>
}
```

## コーディングスタイル

- コンポーネント名・ファイル名は PascalCase、フック・ユーティリティ・非コンポーネントファイルは camelCase
- フックは `use` プレフィックス必須
- 静的解析ツール（ESLint / Biome 等）のルールは常にパスした状態を維持する

## 開発コマンドと品質ゲート

各プロジェクトは以下のスクリプトを定義し、`CLAUDE.md` にコマンドを明記する。

| スクリプト | 役割 | react-frontend での例 |
|---|---|---|
| `lint` | 静的解析（read-only） | `pnpm lint` (Biome + knip) |
| `fmt` | 自動修正 | `pnpm fmt` |
| `typecheck` | 型チェック | `pnpm typecheck` |
| `test` | ユニットテスト一括実行 | `pnpm test` |

**実行タイミング:**

- コミット前: `lint` と `typecheck` を必ず通す
- PR 作成前: `test` も実行し、すべてグリーンであることを確認する
- CI でも同じコマンドを実行し、ローカルと CI の結果が一致する状態を保つ

## テストと品質

- コンポーネントと hooks のユニットテストを実装する
- テストファイルはソースファイルと同じディレクトリに置く（例: `UserCard.test.tsx`）
- ユーティリティ関数はテーブルドリブンテストで境界値・エラー系を網羅する
- UI スナップショットテストは原則採用しない（変更追跡が困難になるため）
- バグ修正では必ず再発防止テストを用意し、失敗シナリオがテストに表現されるまでマージしない

## セキュリティ

- API エンドポイント・トークン・シークレットは環境変数で管理し、リポジトリにハードコードしない
- ユーザー入力を DOM に挿入する前にサニタイズし、XSS を防止する
- 外部リンクには `rel="noopener noreferrer"` を付与する
- 依存ライブラリの CVE を定期的に確認し、Renovate/Dependabot で更新を自動化する

## 運用フロー

- main ブランチへの変更は Pull Request + レビュー + 自動テスト成功を必須とする
- 新しいフロントエンドを追加する場合は本ガイドラインを初期チェックリストとして利用し、差異があればドキュメントで明示する
