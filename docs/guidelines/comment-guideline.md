# Comment Guideline

コードはコメントなしで意図が伝わるよう、説明的な名前と明確な構造で書く。

- コメントは最小限にする。何をしているかを説明するコメントは書かない
- 残してよいコメントは「なぜこうなっているか」が非自明な場合のみ
- 残す必要があるコメントには `NOTE:` プレフィックスを付ける（非自明な分岐・特殊ケースのメソッド・やむを得ないワークアラウンドなど）

```go
// bad: what を説明するコメント（コードを読めばわかる）
// increment counter
count++

// good: NOTE: で why を説明するコメント
// NOTE: pgx returns a zero-value struct on conflict; treat as existing row to keep idempotency.
row, err := q.UpsertUser(ctx, params)
```

```tsx
// bad: what を説明するコメント（コードを読めばわかる）
// set loading to true
setIsLoading(true)

// good: NOTE: で why を説明するコメント
// NOTE: beforeLoad runs outside React; cannot use useAuth() here, so read token directly.
if (!getToken()) throw redirect({ to: '/login' })
```
