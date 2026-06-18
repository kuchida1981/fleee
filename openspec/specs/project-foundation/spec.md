## Purpose

fleee アプリケーションのプロジェクト基盤。単一バイナリ配布、SQLite データベース初期化、開発環境を定義する。

## Requirements

### Requirement: Single binary distribution
アプリケーションは単一の Go バイナリとして配布可能でなければならない。フロントエンドの静的ファイルは `go:embed` によりバイナリに埋め込まれる。

#### Scenario: Binary serves both API and frontend
- **WHEN** ユーザーがビルド済みバイナリを `./fleee serve` で起動する
- **THEN** 単一のポートで API エンドポイントとフロントエンド画面の両方が配信される

#### Scenario: Build produces single artifact
- **WHEN** `make build` を実行する
- **THEN** フロントエンドのビルドと Go のコンパイルが順に実行され、単一のバイナリファイルが生成される

### Requirement: SQLite database initialization
アプリケーションは起動時に SQLite データベースファイルを自動的に初期化しなければならない。

#### Scenario: First launch creates database
- **WHEN** データベースファイルが存在しない状態でアプリケーションを起動する
- **THEN** SQLite データベースファイルが作成され、必要なテーブルが自動的にマイグレーションされる

#### Scenario: Existing database is preserved
- **WHEN** 既存のデータベースファイルがある状態でアプリケーションを起動する
- **THEN** 既存のデータは保持され、未適用のマイグレーションのみが実行される

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
