## Why

テスト・lint・CI が未整備の状態で機能開発が進んでいる。コードベースが小さい今のうちに開発基盤を導入し、高いカバレッジを維持する仕組みを作ることで、今後の機能追加時のリグレッションリスクを最小化する。

## What Changes

- Go テストの導入: store (インテグレーション) + handler (API) テスト、テストヘルパー (`internal/testutil`)
- Go lint の導入: golangci-lint (gofmt, govet, staticcheck, errcheck, unused)
- Frontend テストの導入: Vitest + React Testing Library + jsdom
- Frontend formatter の導入: Prettier (prettier-plugin-tailwindcss 含む)
- カバレッジ計測: Go / Frontend ともに閾値 80% で enforce
- pre-commit hook: lefthook で lint + test をローカル実行
- GitHub Actions CI: backend / frontend 並列 job で lint, test, build を自動実行
- README.md: 一般開発者向けのセットアップ・ビルド・テスト手順
- Makefile 拡張: lint, test, cover ターゲット追加
- CLAUDE.md にテスト方針セクション追加

## Capabilities

### New Capabilities

- `dev-infrastructure`: 開発基盤（テスト、lint、format、pre-commit、CI、カバレッジ）の設定と方針

### Modified Capabilities

- `project-foundation`: Makefile 拡張、.gitignore 更新、README.md 追加、CLAUDE.md にテスト方針追加

## Impact

- 新規ファイル: `.golangci.yml`, `lefthook.yml`, `.github/workflows/ci.yml`, `README.md`, `internal/testutil/testutil.go`, テストファイル群, `web/vitest.config.ts`, `web/src/test/setup.ts`, `web/.prettierrc`
- 変更ファイル: `Makefile`, `web/package.json`, `.gitignore`, `CLAUDE.md`
- 新規 devDependencies (frontend): vitest, @testing-library/react, @testing-library/jest-dom, @vitest/coverage-v8, jsdom, prettier, prettier-plugin-tailwindcss
- 新規ツール (Go): golangci-lint, lefthook (開発者がローカルインストール)
