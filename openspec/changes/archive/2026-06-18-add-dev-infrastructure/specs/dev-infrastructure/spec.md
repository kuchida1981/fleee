## ADDED Requirements

### Requirement: Go test infrastructure
Go のテストは SQLite `:memory:` を使ったインテグレーションテストで実行しなければならない。テストヘルパー `internal/testutil` が in-memory DB の作成・マイグレーション・クリーンアップを提供する。

#### Scenario: Test helper creates in-memory database
- **WHEN** テストが `testutil.NewTestDB(t)` を呼び出す
- **THEN** SQLite `:memory:` データベースが作成され、全マイグレーションが適用された状態で返される

#### Scenario: Test database is cleaned up after test
- **WHEN** テスト関数が終了する
- **THEN** `t.Cleanup()` により DB 接続が自動的にクローズされる

### Requirement: Go store tests
store パッケージの全 CRUD 操作がインテグレーションテストでカバーされなければならない。

#### Scenario: Store test uses real database
- **WHEN** `store/account_test.go` のテストが実行される
- **THEN** 実際の SQLite クエリが実行され、モックは使用されない

#### Scenario: Store tests cover all operations
- **WHEN** `go test ./internal/store/...` を実行する
- **THEN** Create, GetByID, ListAll, Update, Delete の各操作が正常系・異常系ともにテストされる

### Requirement: Go handler tests
handler パッケージの全 API エンドポイントが httptest を使ってテストされなければならない。

#### Scenario: Handler test uses httptest with real database
- **WHEN** `handler/account_test.go` のテストが実行される
- **THEN** `httptest.NewRecorder` と実 DB を使い、HTTP リクエスト/レスポンスの完全なサイクルがテストされる

#### Scenario: Handler tests cover all endpoints
- **WHEN** `go test ./internal/handler/...` を実行する
- **THEN** GET (list/get), POST (create/import), PUT (update), DELETE の各エンドポイントが正常系・エラー系ともにテストされる

### Requirement: Go lint with golangci-lint
Go コードは golangci-lint による静的解析を通過しなければならない。有効 linter は gofmt, govet, staticcheck, errcheck, unused とする。

#### Scenario: Lint passes on clean code
- **WHEN** `golangci-lint run ./...` を実行する
- **THEN** lint エラーが 0 件で正常終了する

#### Scenario: Lint catches unhandled error
- **WHEN** Go コードにエラーを無視する箇所がある（`_ = someFunc()` ではなく戻り値を捨てている）
- **THEN** errcheck が検出し、lint が失敗する

### Requirement: Frontend test infrastructure
フロントエンドのテストは Vitest + React Testing Library + jsdom で実行しなければならない。

#### Scenario: Vitest runs component tests
- **WHEN** `npx vitest run` を web ディレクトリで実行する
- **THEN** `*.test.tsx` ファイルのテストが jsdom 環境で実行される

#### Scenario: Test setup provides DOM matchers
- **WHEN** テストファイルで `@testing-library/jest-dom` のマッチャーを使用する
- **THEN** `toBeInTheDocument()` 等の DOM アサーションが利用可能である

### Requirement: Frontend lint and format
フロントエンドコードは ESLint によるリント + Prettier によるフォーマットチェックを通過しなければならない。Prettier は Tailwind CSS クラスの自動整列を含む。

#### Scenario: ESLint passes on clean code
- **WHEN** `npx eslint .` を web ディレクトリで実行する
- **THEN** lint エラーが 0 件で正常終了する

#### Scenario: Prettier checks formatting
- **WHEN** `npx prettier --check .` を web ディレクトリで実行する
- **THEN** フォーマットが統一されていることが確認される

#### Scenario: Prettier sorts Tailwind classes
- **WHEN** JSX 内の className に Tailwind クラスが記述されている
- **THEN** prettier-plugin-tailwindcss により推奨順序に整列される

### Requirement: Coverage enforcement
Go と Frontend の両方でコードカバレッジを計測し、閾値 80% を CI で enforce しなければならない。

#### Scenario: Go coverage meets threshold
- **WHEN** `go test -coverprofile=coverage.out ./...` を実行し、カバレッジを集計する
- **THEN** 全体のカバレッジが 80% 以上であること

#### Scenario: Frontend coverage meets threshold
- **WHEN** `npx vitest run --coverage` を実行する
- **THEN** lines, functions, branches, statements のすべてが 80% 以上であること

#### Scenario: CI fails on low coverage
- **WHEN** テストカバレッジが 80% を下回るコードが PR に含まれる
- **THEN** CI ジョブが失敗し、マージがブロックされる

### Requirement: Pre-commit hooks with lefthook
lefthook による pre-commit hook が、コミット前に lint とテストをローカルで実行しなければならない。Go と Frontend のチェックは並列に実行される。

#### Scenario: Pre-commit runs all checks
- **WHEN** 開発者が `git commit` を実行する
- **THEN** 以下のチェックが並列に実行される: Go lint, Go test, Frontend typecheck, Frontend lint, Frontend format check, Frontend test

#### Scenario: Pre-commit blocks on failure
- **WHEN** いずれかのチェックが失敗する
- **THEN** コミットがブロックされ、失敗したチェックのエラーが表示される

#### Scenario: golangci-lint does not auto-fix in pre-commit
- **WHEN** pre-commit で golangci-lint が実行される
- **THEN** `--fix` フラグは使用されず、エラー報告のみが行われる

### Requirement: GitHub Actions CI
GitHub Actions で PR ごとに backend と frontend のジョブが並列実行されなければならない。

#### Scenario: CI triggers on push and PR to main
- **WHEN** main ブランチへの push または pull request が作成される
- **THEN** CI ワークフローが起動する

#### Scenario: Backend job runs lint, test, build
- **WHEN** backend ジョブが実行される
- **THEN** golangci-lint → go test (カバレッジ付き、閾値チェック) → go build の順で実行される

#### Scenario: Frontend job runs typecheck, lint, format, test, build
- **WHEN** frontend ジョブが実行される
- **THEN** tsc --noEmit → eslint → prettier --check → vitest --coverage (閾値チェック) → vite build の順で実行される

#### Scenario: Backend and frontend jobs run in parallel
- **WHEN** CI ワークフローが起動する
- **THEN** backend ジョブと frontend ジョブが並列に実行される

### Requirement: Developer documentation
README.md に一般開発者向けのプロジェクト概要、前提条件、ビルド・開発・テスト手順が記載されなければならない。

#### Scenario: README covers setup and build
- **WHEN** 新しい開発者がリポジトリをクローンする
- **THEN** README.md の手順に従って、前提条件のインストール → ビルド → 開発サーバー起動 → テスト実行ができる
