## 1. Go プロジェクト基盤

- [ ] 1.1 Go モジュール初期化（`go mod init`）、ディレクトリ構成（cmd/, internal/, migrations/）の作成
- [ ] 1.2 chi ルーター、modernc.org/sqlite の依存追加とエントリポイント（cmd/fleee/main.go）の実装
- [ ] 1.3 SQLite 接続管理と マイグレーション実行の仕組み（internal/store/db.go）の実装
- [ ] 1.4 accounts テーブルのマイグレーション SQL 作成（migrations/001_create_accounts.sql）

## 2. 勘定科目 API

- [ ] 2.1 勘定科目のドメインモデル定義（internal/model/account.go）— account_type からの normal_balance, statement_type 導出ロジックを含む
- [ ] 2.2 勘定科目リポジトリの実装（internal/store/account.go）— CRUD 操作、一覧取得（display_order 順）
- [ ] 2.3 勘定科目 API ハンドラの実装（internal/handler/account.go）— GET/POST/PUT/DELETE /api/accounts
- [ ] 2.4 CSV/TSV インポートロジックの実装（internal/importer/account_csv.go）— 科目貸借タイプ+精算種別 → account_type 変換、重複スキップ
- [ ] 2.5 インポート API ハンドラの実装（POST /api/accounts/import）— multipart/form-data でファイル受信

## 3. フロントエンド基盤

- [ ] 3.1 React + Vite + TypeScript プロジェクトの初期化（web/）
- [ ] 3.2 Tailwind CSS と shadcn/ui のセットアップ
- [ ] 3.3 Vite の API プロキシ設定（開発時に /api/* を Go サーバーへ転送）
- [ ] 3.4 API クライアントの共通処理（web/src/api/）の実装

## 4. 勘定科目 UI

- [ ] 4.1 勘定科目一覧ページの実装 — テーブル表示、空状態メッセージ
- [ ] 4.2 勘定科目作成フォームの実装 — 科目名、勘定区分（セレクト）、表示順序の入力
- [ ] 4.3 勘定科目編集フォームの実装 — 既存値のプリフィル、バリデーション
- [ ] 4.4 勘定科目削除機能の実装 — 確認ダイアログ付き
- [ ] 4.5 CSV/TSV インポート機能の実装 — ファイル選択、インポート結果の表示（成功件数、スキップ件数）

## 5. ビルドと配布

- [ ] 5.1 go:embed によるフロントエンド静的ファイルの埋め込み設定
- [ ] 5.2 Makefile の作成（dev, build, clean ターゲット）
- [ ] 5.3 単一バイナリでのフロントエンド配信と API の統合動作確認
