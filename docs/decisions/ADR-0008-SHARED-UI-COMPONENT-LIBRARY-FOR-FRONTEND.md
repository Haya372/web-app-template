# ADR-0008: Shared UI Component Library for Frontend

Date: 2026-03-08
Status: Accepted

---

## Context

フロントエンドアプリケーションは `apps/*` 配下に複数実装していく想定である。各アプリが独自に UI コンポーネントを持つと、デザインの一貫性が失われ、同一コンポーネントの重複実装が生じやすい。そのため、共通 UI コンポーネントを集約するパッケージが必要である。

### Background

- monorepo 構成（pnpm workspace + Turborepo）を採用しており、`packages/*` に共通パッケージを置く構成がすでに整っている。
- `apps/*` に複数のフロントエンドアプリを追加していく予定がある。
- 一貫したデザインとコンポーネント再利用性を担保したい。

### Scope

- 対象: `packages/ui` パッケージの方針（コンポーネント基盤・ビルド設定・公開範囲）。
- 対象外: 各 `apps/*` のフレームワーク選定、グローバル状態管理、ルーティング方針。

### Constraints

- `packages/ui` は外部 npm 公開を行わず、monorepo 内の `apps/*` からのみ利用する。
- コンポーネントはカスタマイズ可能であることが望ましい。
- 実装コストを抑えつつ、高品質なコンポーネントをすぐに利用できる状態にする。

## Decision

`packages/ui` を **shadcn/ui ベースの内部共通 UI コンポーネントライブラリ** として定義する。ビルドには **Vite（library mode）** を使用する。

運用ルール:

- `packages/ui` は pnpm workspace の内部パッケージとして管理し、外部 npm レジストリへは公開しない。
- shadcn/ui のコンポーネントをベースとして取り込み、プロジェクト固有のカスタマイズを加えた上で管理する。
- ビルドには Vite の library mode を使用し、ESM 形式で出力する。
- `apps/*` からは `@repo/ui`（または workspace package 名）として import する。
- shadcn/ui が提供しないカスタムコンポーネントも同パッケージに追加してよい。
- コンポーネントのスタイルは Tailwind CSS を使用する。

## Options

### Option A: shadcn/ui ベースの内部パッケージ（採用）

- 概要
  - shadcn/ui のコンポーネントコードを `packages/ui` に取り込み、Vite library mode でビルドする。
- Pros
  - アクセシビリティ対応済みの高品質なコンポーネントをすぐに利用できる。
  - コードが手元にあるため、プロジェクト固有のカスタマイズが容易。
  - Radix UI + Tailwind CSS ベースで、デザインシステムとの統合がしやすい。
- Cons
  - shadcn/ui 本体の更新を手動で追跡・マージする必要がある。
  - 初期セットアップに取り込み作業が必要。

### Option B: MUI・Chakra UI 等の外部コンポーネントライブラリをそのまま利用

- 概要
  - `packages/ui` を作らず、各 `apps/*` が直接 MUI 等の外部ライブラリを依存に持つ。
- Pros
  - 初期導入が最も速い。
  - ライブラリ側でメンテナンスされるため、管理コストが低い。
- Cons
  - カスタマイズに制約があり、デザイン変更時の対応コストが高くなりやすい。
  - 各アプリが独自にライブラリを選ぶと、アプリ間で一貫性が失われるリスクがある。

### Option C: コンポーネントをゼロから自作

- 概要
  - shadcn/ui 等に依存せず、`packages/ui` 内のコンポーネントをすべて自作する。
- Pros
  - 完全な制御が可能で、外部依存を最小化できる。
- Cons
  - 実装・メンテナンスコストが非常に高い。
  - アクセシビリティ対応を自前で行う必要がある。

## Rationale

最終判断で重視した軸は以下。

- カスタマイズ性（プロジェクト固有のデザインに追従できるか）
- 初期コスト（すぐに使えるコンポーネントが揃うか）
- 保守性（依存の更新・管理が現実的か）

外部ライブラリをそのまま利用する Option B は導入は速いが、長期的なカスタマイズ自由度が低く、アプリ間の一貫性も担保しにくい。ゼロ自作の Option C は制御性が高い反面、コストが過大である。

shadcn/ui はコンポーネントコードを直接手元に持つ設計であるため、プロジェクト固有のカスタマイズと統一管理を両立できる。Vite library mode は monorepo 内パッケージのビルドとして実績があり、ESM 出力により各アプリからの利用も容易である。これらの理由から Option A を採用する。

## Consequences

- Positive
  - 全 `apps/*` で一貫した UI コンポーネントを共有でき、デザインの分散を防げる。
  - shadcn/ui ベースにより、アクセシビリティ対応済みのコンポーネントをすぐに利用できる。
  - コードが手元にあるため、プロジェクト固有の要件に柔軟に対応できる。

- Negative
  - shadcn/ui の更新を手動で取り込む必要があり、メンテナンス運用が必要。
  - Vite library mode のビルド設定・型定義出力の初期整備が必要。

- Migration / Follow-up
  - `packages/ui` の初期セットアップ（Vite library mode 設定、shadcn/ui 取り込み、Tailwind CSS 設定）を行う。
  - `pnpm-workspace.yaml` に `apps/*` を追加し、`apps/*` からの参照を確認する。

## References

- [shadcn/ui](https://ui.shadcn.com/)
- [Vite Library Mode](https://vite.dev/guide/build#library-mode)
