## Why

複式簿記アプリとして最も基本的な機能である仕訳の入力・管理がまだ存在しない。勘定科目の CRUD は完成しているため、次のステップとして仕訳入力を実装し、実際の記帳が行える状態にする。

## What Changes

- 仕訳（journal entry）のヘッダー + 明細行テーブルを追加し、複合仕訳に対応する
- 仕訳の CRUD API エンドポイントを追加する（`/api/journal-entries`）
- 仕訳入力画面を追加する（単一仕訳モードをデフォルト、複合仕訳モードに切替可能）
- 仕訳一覧画面を追加する
- 勘定科目の削除時に、仕訳で使用中の科目は FK 制約で削除を拒否する

## Capabilities

### New Capabilities
- `journal-entry-management`: 仕訳の作成・一覧表示・編集・削除。ヘッダー（日付・摘要・領収書要否・メモ）+ 明細行（科目・借方金額・貸方金額）の複合仕訳モデル。

### Modified Capabilities
- `account-management`: 仕訳で使用中の勘定科目の削除を FK 制約で拒否するようになる

## Impact

- **DB**: マイグレーション追加（`journal_entries`, `journal_lines` テーブル）
- **Backend**: `internal/model/`, `internal/store/`, `internal/handler/` に仕訳関連のファイルを追加
- **Frontend**: 仕訳入力フォーム・一覧コンポーネントを追加、ルーティング追加
- **API**: `GET/POST /api/journal-entries`, `GET/PUT/DELETE /api/journal-entries/{id}`
- **既存コード**: `accounts` テーブルの DELETE が FK 制約の影響を受ける
