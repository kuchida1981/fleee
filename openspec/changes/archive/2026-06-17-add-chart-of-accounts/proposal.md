## Why

複式簿記で帳簿を記録するアプリケーション「fleee」の最初の機能として、勘定科目の管理基盤を構築する。勘定科目は仕訳・元帳・財務諸表すべての前提となるマスタデータであり、これがなければ他の機能を作れない。

## What Changes

- Go + SQLite バックエンド、React + Vite フロントエンドによるプロジェクト初期構造を構築
- `go:embed` でフロントエンドを埋め込み、単一バイナリとして配布可能にする
- 勘定科目の CRUD API およびUI を実装
- 勘定科目の5区分モデル（資産・負債・純資産・収益・費用）を採用し、貸借タイプと精算種別は account_type から導出する
- TSV/CSV インポート機能（既存マスタの「科目名, 科目貸借タイプ, 出力順番, 精算種別」形式に対応）

## Capabilities

### New Capabilities

- `project-foundation`: Go + SQLite + React + Vite によるプロジェクト基盤。単一バイナリ配布、DB マイグレーション、API ルーティング、フロントエンドビルド＆埋め込み
- `account-management`: 勘定科目の CRUD と CSV/TSV インポート。5区分モデルによるデータ管理

### Modified Capabilities

(なし — 初回の change のため既存 capability はない)

## Impact

- 新規プロジェクト構造の作成（Go モジュール、npm プロジェクト、ビルドスクリプト）
- SQLite スキーマの新規作成（accounts テーブル）
- 外部依存: Go SQLite ドライバー、React、Vite、UI コンポーネントライブラリ
