## 1. DB マイグレーション

- [x] 1.1 `migrations/002_create_journal_entries.sql` を作成: `journal_entries` テーブル（id, date, description, receipt_required, memo, created_at, updated_at）と `journal_lines` テーブル（id, journal_entry_id FK, account_id FK, debit_amount, credit_amount, created_at, updated_at）を定義。journal_lines に CHECK 制約（金額の非負、借方/貸方の排他）と `ON DELETE CASCADE` を設定

## 2. Backend モデル

- [x] 2.1 `internal/model/journal_entry.go` を作成: JournalEntry 構造体（ID, Date, Description, ReceiptRequired, Memo, Lines, CreatedAt, UpdatedAt）と JournalLine 構造体（ID, JournalEntryID, AccountID, AccountName, DebitAmount, CreditAmount）を定義。JSON シリアライズ対応

## 3. Backend Store

- [x] 3.1 `internal/store/journal_entry.go` を作成: JournalEntryStore に Create, GetByID, ListAll, Update, Delete を実装。Create/Update ではトランザクション内でバランスチェック（借方合計 = 貸方合計）と最低2行チェックを行う。GetByID では明細行に科目名を JOIN で取得
- [x] 3.2 `internal/store/journal_entry_test.go` を作成: 正常系（CRUD）+ エラー系（バランス不一致、行数不足、存在しない科目 ID、存在しない仕訳 ID）のテスト

## 4. Backend Handler

- [x] 4.1 `internal/handler/journal_entry.go` を作成: JournalEntryHandler に Routes（GET /, POST /, GET /{id}, PUT /{id}, DELETE /{id}）を実装。リクエストバリデーション（日付・摘要の必須チェック）を含む
- [x] 4.2 `internal/handler/journal_entry_test.go` を作成: 全エンドポイントの HTTP テスト（正常系 + バリデーションエラー）
- [x] 4.3 `internal/server/server.go` にルーティング追加: `/api/journal-entries` に JournalEntryHandler をマウント。`cmd/fleee/main.go` で JournalEntryStore と JournalEntryHandler を初期化

## 5. Account 削除の FK 制約対応

- [x] 5.1 `internal/handler/account.go` の Delete ハンドラで FK 制約エラー（FOREIGN KEY constraint failed）を検出し、409 Conflict を返すように修正。既存の account_test.go に仕訳使用中の科目削除テストを追加

## 6. Frontend API クライアント

- [x] 6.1 `web/src/api/journalEntries.ts` を作成: 型定義（JournalEntry, JournalLine, CreateJournalEntryRequest, UpdateJournalEntryRequest）と API 関数（fetchJournalEntries, fetchJournalEntry, createJournalEntry, updateJournalEntry, deleteJournalEntry）を実装

## 7. Frontend コンポーネント

- [x] 7.1 `web/src/components/JournalEntryForm.tsx` を作成: 単一仕訳モード（デフォルト）のフォーム。日付、借方科目セレクト、貸方科目セレクト、金額、摘要、領収書チェックボックス、メモを入力
- [x] 7.2 JournalEntryForm に複合仕訳モードを追加: モード切替ボタン、明細行テーブル（科目・借方金額・貸方金額の行追加/削除）、借方/貸方合計のリアルタイム表示とバランスインジケーター
- [x] 7.3 `web/src/components/JournalEntryTable.tsx` を作成: 仕訳一覧テーブル（日付、摘要、合計金額を表示）。空状態メッセージ対応
- [x] 7.4 `web/src/App.tsx` にルーティング追加: 仕訳一覧・入力画面への導線を追加

## 8. Frontend テスト

- [x] 8.1 `web/src/components/JournalEntryForm.test.tsx` を作成: 単一仕訳モードのフォーム表示・入力・送信テスト、複合仕訳モードの切替・行追加削除・バランス表示テスト
- [x] 8.2 `web/src/components/JournalEntryTable.test.tsx` を作成: 一覧表示・空状態のテスト
