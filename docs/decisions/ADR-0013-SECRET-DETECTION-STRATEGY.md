# ADR-0013: Secret Detection Strategy with gitleaks

Date: 2026-04-24
Status: Accepted

---

## Context

### Background

- シークレット（API キー、トークン等）が誤ってコミットされても検知する仕組みがなかった
- git 履歴にシークレットが残ると、公開後の取り消しが困難でセキュリティインシデントに直結する
- backend/frontend コーディングガイドラインには hardcode 禁止が明記されているが、技術的強制はなかった
- 既存のセキュリティツール（zizmor, shellcheck）は CI レベルで mise 管理されており、同じパターンを踏襲できる

### Scope

- 対象: 新規コミットのシークレット検出（pre-push フック + CI）
- 対象外: 既存 git 履歴全体のスキャン・クリーンアップ（別チケットで対応）

### Constraints

- ツールは `mise` で管理（既存パターンとの統一）
- Git フックは既存 Husky（v9）で管理
- CI は GitHub Actions（PR ベース）
- 全 GitHub Actions アクションはコミットハッシュで固定するプロジェクト方針

## Decision

gitleaks を採用し、**pre-push フック（ローカル）** と **GitHub Actions（CI）** の 2 段階でスキャンを実施する。

### フック責務分離

| フック | 目的 | ツール | 失敗時の動作 |
|--------|------|--------|------------|
| pre-commit | コード品質（lint, format） | lint-staged | コミットをブロック |
| pre-push | セキュリティ（secret detection） | gitleaks protect | プッシュをブロック |

### false positive 管理方針

- `.gitleaks.toml` は `useDefault = true` のみを設定する
- regex/paths による一括除外は検知漏れリスクがあるため原則使用しない
- false positive は `.gitleaksignore` で fingerprint 単位の個別無視で管理する
- `.gitleaksignore` への追記は必ず PR レビューを経ること（security review として扱う）

## Options

### Option A: gitleaks（pre-push + CI）— 採用

| 観点 | 評価 |
|------|------|
| Pros | 多層防御、開発者即時フィードバック、mise 管理、`.gitleaks.toml` で柔軟設定 |
| Cons | 初回セットアップ（`mise install`）が必要、false positive 対応コスト |

### Option B: TruffleHog

| 観点 | 評価 |
|------|------|
| Pros | より高精度なエントロピー検出 |
| Cons | mise 非対応、GitHub Actions 統合が複雑、チーム習熟コスト高 |

### Option C: GitHub Advanced Security（Secret Scanning）

| 観点 | 評価 |
|------|------|
| Pros | 設定不要、プラットフォーム統合 |
| Cons | 有償（GitHub Advanced Security が必要）、ローカル pre-push に使えない |

### Option D: 何もしない

採用 ×。シークレット漏洩リスクが継続する。

## Rationale

Option A（gitleaks）を採用。

- **多層防御**: push 前の早期発見で修正コストを最小化し、CI は `--no-verify` 回避を防ぐ最終防御として機能
- **mise との親和性**: 既存ツールと同じ管理方式でチーム全員に同一バージョンを配布できる
- **設定の柔軟性**: `.gitleaks.toml` + `.gitleaksignore` で false positive を fine-tune できる
- **GitHub 統合**: `gitleaks-action` で PR check run・コメント連携が標準提供される

## Consequences

- **Positive**: 技術的にシークレット漏洩をブロックし、ガイドライン違反の抑制効果が生まれる
- **Negative**: 初回セットアップ手順（`mise install`）が必要、false positive 対応コスト、`.gitleaksignore` の継続的なメンテナンスが必要
- **Migration**: 導入直後に既存リポジトリの全履歴スキャン（`gitleaks detect`）を手動実行し、既存 finding を `.gitleaksignore` に登録するか対応する（別チケット）
- **Follow-up**: `.gitleaksignore` の定期棚卸し（四半期ごと）を運用ルールとして定める

## References

- Issue: #156
- Guideline: docs/guidelines/secret-detection.md
- gitleaks: https://github.com/gitleaks/gitleaks
- gitleaks-action: https://github.com/gitleaks/gitleaks-action
