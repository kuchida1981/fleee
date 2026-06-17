## Context

fleee は Go + React の単一バイナリアプリケーション。現在 Go ファイル 6 個、主要パッケージ 4 つ (store, handler, model, importer) の小規模なコードベース。テスト・lint・CI は未整備で、ESLint の設定のみが存在する。

コードベースが小さい今、開発基盤を整えておくことで、今後の機能追加（仕訳入力、レポート出力など）で品質を維持する仕組みを作る。

## Goals / Non-Goals

**Goals:**
- Go / Frontend の両方でテスト・lint・format を自動化する
- カバレッジ 80% を CI で enforce し、テストなしの機能追加を防ぐ
- pre-commit hook でコミット前にローカルで品質チェックを実行する
- GitHub Actions で PR ごとに自動検証する
- 一般開発者向けの README.md を整備する

**Non-Goals:**
- E2E テスト (Playwright 等) の導入 — UI が安定してから検討
- API ドキュメント自動生成 (swaggo 等) — 現段階では不要
- Docker / コンテナベースの CI — シンプルな GitHub Actions で十分
- テストカバレッジの PR コメント自動投稿 — ログ確認で十分

## Decisions

### D1: Go テストは全レイヤー実 DB (SQLite `:memory:`)

store も handler も SQLite `:memory:` で実 DB を使う。interface によるモックは導入しない。

**Alternatives considered:**
- store を interface 化して handler テストでモック → コードベースが小さく store 実装が 1 つしかないため premature abstraction。モックと実 DB の乖離リスクもある。

**Rationale:** `:memory:` SQLite はミリ秒単位で起動するため速度の問題がない。全レイヤーで本物の SQL が通ることを保証でき、モック/実装の乖離によるバグを排除する。

### D2: テストヘルパーは `internal/testutil` パッケージ

`internal/testutil/testutil.go` に `NewTestDB(t *testing.T) *store.DB` を配置。`:memory:` DB の作成、マイグレーション実行、`t.Cleanup()` による自動クローズを担う。

**Rationale:** 全テストファイルで DB セットアップのボイラープレートを排除する。`internal/` 配下なので外部に公開されない。

### D3: golangci-lint の linter セット

初期有効 linter: `gofmt`, `govet`, `staticcheck`, `errcheck`, `unused`。

**Alternatives considered:**
- gofumpt (gofmt の厳密版) → 既存コードとの差分が大きくなりうるため初期導入は見送り。
- exhaustive, gocritic 等 → 有用だが初期には過剰。必要に応じて後から追加。

**Rationale:** 最小限だが実効性の高いセット。フォーマット、型安全性、未使用コード、エラーハンドリングをカバーする。

### D4: lefthook を pre-commit ツールとして採用

**Alternatives considered:**
- pre-commit (Python ベース) → Python ランタイムが必要。Go + Node プロジェクトに追加の言語依存は避けたい。
- husky + lint-staged (Node ベース) → Go 側のカバーが弱い。lefthook は Go / Node 両方を統一的に扱える。

**Rationale:** Go 製の単一バイナリで、YAML 1 ファイルで Go + Node 両方の hook を定義できる。parallel 実行もサポート。

### D5: pre-commit で golangci-lint の `--fix` は使わない

lint エラーの報告のみ行い、自動修正はしない。

**Rationale:** `--fix` が staged されていないファイルを変更すると、git の staging 状態と実ファイルの不一致が起きる。修正は開発者が明示的に行い、再度 stage してからコミットする。

### D6: Frontend テストは Vitest + React Testing Library

**Alternatives considered:**
- Jest → Vite プロジェクトでは Vitest の方が設定がシンプルで ESM との親和性が高い。

**Rationale:** Vite のビルド設定を再利用でき、設定の二重管理を避けられる。

### D7: Frontend API モックは `vi.mock` から開始

**Alternatives considered:**
- msw (Mock Service Worker) → ネットワークレイヤーのモックとしては優秀だが、現段階では API 呼び出しが少なく over-engineering。

**Rationale:** API 呼び出しが複雑化してきたら msw に移行する。現段階では `vi.mock` で十分。

### D8: カバレッジ閾値 80%

Go は `go test -coverprofile` + シェルスクリプトで閾値チェック。Frontend は Vitest の `coverage.thresholds` 設定で enforce。

**Rationale:** コードベースが小さい今は余裕で達成できる値。高すぎ (95%+) るとエラーハンドリングの全分岐など投資効率の低いテストを強制する。成長に応じて調整する。

### D9: GitHub Actions は backend / frontend を別 job で並列実行

**Rationale:** 互いに依存がないため並列実行で CI 時間を短縮できる。`actions/setup-go` と `actions/setup-node` のキャッシュをそれぞれ活用。

## Risks / Trade-offs

- **[pre-commit にテストを含めることでコミットが遅くなる]** → 現状の SQLite `:memory:` テスト + Vitest は数秒で完了する見込み。遅くなったらテストを CI のみに移す。
- **[lefthook のインストールが開発者ごとに必要]** → README に手順を記載。`go install` で入るため Go 開発者には障壁が低い。
- **[カバレッジ 80% が将来的に足かせになる]** → 閾値は `.golangci.yml` や `vitest.config.ts` で管理されているため調整は容易。
- **[golangci-lint と ESLint のバージョン固定]** → CI の `golangci-lint-action` でバージョンを指定。Frontend は `package.json` で lock。
