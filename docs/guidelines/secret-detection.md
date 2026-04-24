# Secret Detection ガイドライン

## 概要

gitleaks を使い、シークレット（API キー・トークン等）のコミットをプッシュ前と CI の 2 段階で検出する。
設計判断の詳細は [ADR-0013](../decisions/ADR-0013-SECRET-DETECTION-STRATEGY.md) を参照。

## セットアップ（初回のみ）

```sh
mise install    # gitleaks バイナリを取得
pnpm install    # Husky フックを有効化（prepare スクリプトが husky install を実行）
```

確認:

```sh
ls .husky/pre-push    # ファイルが存在すること
gitleaks version      # バージョンが表示されること
```

## 動作

### ローカル（pre-push フック）

`git push` 時に `gitleaks protect --staged --redact --verbose` が自動実行される。

- シークレット検出時: プッシュがブロックされ、ファイル・行番号・ルール名が表示される
- 検出なし: プッシュが通常通り完了する

### CI（GitHub Actions）

PR の作成・更新時に `gitleaks detect` が自動実行される。

- シークレット検出時: `gitleaks-scan` ジョブが FAIL し、PR の check run に結果が表示される
- 検出なし: ジョブが PASS する

## false positive（誤検出）の対応手順

誤検出が発生した場合は `.gitleaksignore` に fingerprint を追記して PR を作成する。
**追記には必ずレビューを経ること（security review として扱う）。**

### fingerprint の取得方法

```sh
gitleaks detect --report-format json --report-path /tmp/gitleaks-report.json
cat /tmp/gitleaks-report.json | jq '.[].Fingerprint'
```

### `.gitleaksignore` への追記フォーマット

```
<fingerprint>  # refs PR #<番号>
```

例:

```
abc123def456789  # refs PR #157
```

### PR レビュー時の確認チェックリスト

- [ ] この finding は本物のシークレットではないか（実際に機能するキーでないか）
- [ ] `.gitleaksignore` コメントに PR 番号が記載されているか
- [ ] シークレット値や個人情報がコメントに含まれていないか

> **注意**: regex/paths による一括除外は検知漏れリスクがあるため、`.gitleaks.toml` への allowlists 追加は原則行わない。個別 finding のみ `.gitleaksignore` で対処する。

## トラブルシューティング

### フックが実行されない

```sh
mise install          # gitleaks が PATH に存在するか確認
gitleaks version      # 表示されなければ mise install が必要
pnpm install          # Husky が未初期化の場合
```

### 緊急時の回避方法

```sh
git push --no-verify  # ローカルフックをスキップ（CI は回避不可）
```

> **警告**: `--no-verify` で回避しても CI の `gitleaks-scan` ジョブでブロックされる。
> 使用前に `.gitleaksignore` に適切な fingerprint を追加するか、該当コミットを修正すること。

## `.gitleaksignore` の棚卸し（四半期ごと）

1. `.gitleaksignore` の各 fingerprint が引き続き有効か確認する
2. 不要になった fingerprint（対応済み・削除済みのコードに関するもの）を削除して PR で報告する
