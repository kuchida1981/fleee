## 1. Go Lint 環境

- [ ] 1.1 `.golangci.yml` を作成し、gofmt, govet, staticcheck, errcheck, unused を有効にする
- [ ] 1.2 `golangci-lint run ./...` が既存コードでパスすることを確認し、既存の lint エラーがあれば修正する

## 2. Go テスト基盤

- [ ] 2.1 `internal/testutil/testutil.go` を作成: `NewTestDB(t)` (`:memory:` DB 作成、マイグレーション実行、`t.Cleanup()` 登録)
- [ ] 2.2 `internal/store/account_test.go` を作成: Create, GetByID, ListAll, Update, Delete の正常系・異常系テスト
- [ ] 2.3 `internal/store/db_test.go` を作成: Migrate の動作テスト
- [ ] 2.4 `internal/handler/account_test.go` を作成: 全エンドポイント (GET list/get, POST create/import, PUT update, DELETE) の正常系・エラー系テスト
- [ ] 2.5 `go test -coverprofile=coverage.out ./...` でカバレッジ 80% 以上を確認する

## 3. Frontend テスト・Format 環境

- [ ] 3.1 Vitest + React Testing Library + jsdom + @vitest/coverage-v8 を devDependencies に追加する
- [ ] 3.2 `web/vitest.config.ts` を作成: jsdom 環境、coverage 設定（閾値 80%）、setup ファイルの指定
- [ ] 3.3 `web/src/test/setup.ts` を作成: `@testing-library/jest-dom` のインポート
- [ ] 3.4 Prettier + prettier-plugin-tailwindcss を devDependencies に追加し、`web/.prettierrc` を作成する
- [ ] 3.5 既存コードに `npx prettier --write .` を適用してフォーマットを統一する
- [ ] 3.6 各コンポーネントのテストファイル (*.test.tsx) を作成する
- [ ] 3.7 `npx vitest run --coverage` でカバレッジ 80% 以上を確認する

## 4. Makefile 拡張

- [ ] 4.1 `lint` ターゲットを追加: `golangci-lint run ./...` + `cd web && npx eslint .` + `cd web && npx prettier --check .`
- [ ] 4.2 `test` ターゲットを追加: `go test ./...` + `cd web && npx vitest run`
- [ ] 4.3 `cover` ターゲットを追加: `go test -coverprofile` + `cd web && npx vitest run --coverage`

## 5. Pre-commit (lefthook)

- [ ] 5.1 `lefthook.yml` を作成: Go lint, Go test, Frontend typecheck, Frontend lint, Frontend format, Frontend test を並列実行する設定
- [ ] 5.2 `lefthook install` でフックをインストールし、テストコミットで動作確認する

## 6. GitHub Actions CI

- [ ] 6.1 `.github/workflows/ci.yml` を作成: backend job (setup-go → golangci-lint → go test -cover → coverage 閾値チェック → go build)
- [ ] 6.2 同ファイルに frontend job を追加: setup-node → npm ci → tsc --noEmit → eslint → prettier --check → vitest --coverage → vite build
- [ ] 6.3 backend / frontend を並列 job として構成し、キャッシュ設定を含める

## 7. ドキュメント・設定ファイル更新

- [ ] 7.1 `README.md` を作成: プロジェクト概要、前提条件 (Go, Node.js, golangci-lint, lefthook)、ビルド・開発・テスト手順
- [ ] 7.2 `.gitignore` にカバレッジ出力 (`coverage.out`, `web/coverage/`) を追加する
- [ ] 7.3 `CLAUDE.md` にテスト方針セクションを追加: テスト戦略、カバレッジ目標、テストの書き方ガイドライン
