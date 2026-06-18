## MODIFIED Requirements

### Requirement: Development environment
開発時はフロントエンドとバックエンドを独立して起動し、ホットリロード可能な環境を提供しなければならない。Makefile に lint, test, cover ターゲットを追加し、開発ワークフローを支援する。

#### Scenario: Frontend hot reload during development
- **WHEN** 開発者が `web/` 配下のソースを変更する
- **THEN** Vite dev server がホットリロードでブラウザに反映する

#### Scenario: API proxy in development
- **WHEN** Vite dev server が起動している状態でフロントエンドから `/api/*` にリクエストする
- **THEN** リクエストは Go サーバーにプロキシされる

#### Scenario: Makefile lint target
- **WHEN** 開発者が `make lint` を実行する
- **THEN** Go の golangci-lint と Frontend の ESLint + Prettier チェックが実行される

#### Scenario: Makefile test target
- **WHEN** 開発者が `make test` を実行する
- **THEN** Go テストと Frontend テストの両方が実行される

#### Scenario: Makefile cover target
- **WHEN** 開発者が `make cover` を実行する
- **THEN** Go と Frontend のカバレッジレポートが生成される
