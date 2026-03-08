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
