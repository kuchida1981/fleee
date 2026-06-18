## Context

勘定科目の CRUD が完成しており、次のステップとして仕訳（journal entry）の入力・管理を実装する。フリーランスの経費・売上記帳が主な用途だが、複合仕訳にも対応する設計とする。

現在のアーキテクチャ:
- Backend: `model/` → `store/` → `handler/` の3層
- Frontend: `components/` に AccountForm, AccountTable 等のコンポーネント
- DB: SQLite、マイグレーションは `migrations/` に連番 SQL

## Goals / Non-Goals

**Goals:**
- 仕訳の CRUD（作成・一覧・取得・更新・削除）を API + UI で提供する
- 複合仕訳（N借方 : M貸方）をデータモデルでサポートする
- 単一仕訳のケースを UI のデフォルトモードとして簡単に入力できるようにする
- 借方合計 = 貸方合計のバランスチェックをアプリ層で保証する

**Non-Goals:**
- 仕訳のCSV/TSVインポート・エクスポート（将来の別 change で対応）
- 試算表・残高計算・財務諸表の生成（将来の別 change で対応）
- 仕訳テンプレート・定期仕訳の自動起票
- 仕訳番号の自動採番（ID で代用）

## Decisions

### 1. ヘッダー + 明細行のテーブル分離

**選択:** `journal_entries`（ヘッダー）+ `journal_lines`（明細行）の2テーブル構成

**代替案:**
- 1テーブル（debit_account_id, credit_account_id, amount）: 単純だが複合仕訳に対応できない
- JSON カラムに明細行を格納: クエリ・集計が困難

**理由:** 複合仕訳（N:M）に対応しつつ、SQL での集計（試算表等）を将来容易にする

### 2. 金額の表現: debit_amount / credit_amount の2カラム

**選択:** 各明細行に `debit_amount INTEGER` と `credit_amount INTEGER` を持つ。一方が正の値で他方は 0。

**代替案:**
- amount + side ("debit"/"credit"): GROUP BY 集計時に CASE 式が必要になる

**理由:** `SUM(debit_amount)`, `SUM(credit_amount)` で直感的に集計でき、CHECK 制約で不正データを防げる

### 3. 金額は INTEGER（円単位）

**選択:** 小数を扱わない。フリーランスの日本円取引では十分。

### 4. バランスチェックはアプリ層で実施

**選択:** Go のトランザクション内で明細行の借方合計と貸方合計を検証してから COMMIT する。

**代替案:**
- SQLite TRIGGER: 複雑で保守しにくい
- DB CHECK 制約: 複数行にまたがるチェックは不可

**理由:** シンプルで十分。テストも書きやすい。

### 5. 勘定科目の削除拒否（FK 制約）

**選択:** `journal_lines.account_id` に `REFERENCES accounts(id)` を設定し、使用中の科目の削除を DB レベルで拒否する。`PRAGMA foreign_keys = ON` は既に有効。

### 6. UI: 単一仕訳デフォルト + 複合仕訳モード切替

**選択:** デフォルトは借方1科目・貸方1科目・金額1つの簡易フォーム。「複合仕訳に切替」ボタンで明細行テーブル形式に展開する。

**理由:** フリーランスの日常仕訳の大半は単一仕訳。入力の手軽さを優先しつつ、複合仕訳もサポートする。

### 7. API 設計: ネストした JSON でヘッダー + 明細行を一括操作

**選択:**
```json
POST /api/journal-entries
{
  "date": "2026-06-18",
  "description": "携帯電話料金6月分",
  "receipt_required": true,
  "memo": "",
  "lines": [
    { "account_id": 5, "debit_amount": 10000, "credit_amount": 0 },
    { "account_id": 1, "debit_amount": 0, "credit_amount": 10000 }
  ]
}
```

ヘッダーと明細行を1リクエストで送受信する。部分更新はサポートしない（PUT で全体置換）。

### 8. 既存コードパターンへの準拠

account と同じレイヤー構成を踏襲する:
- `internal/model/journal_entry.go` — ドメインモデル
- `internal/store/journal_entry.go` — DB 操作
- `internal/handler/journal_entry.go` — HTTP ハンドラ
- `web/src/api/journalEntries.ts` — API クライアント
- `web/src/components/JournalEntry*.tsx` — UI コンポーネント

## Risks / Trade-offs

- **[Risk] 明細行のカスケード削除**: 仕訳ヘッダー削除時に明細行も連鎖削除される（`ON DELETE CASCADE`）。意図的な設計だが、誤削除のリスクはある → UI で削除確認ダイアログを表示する
- **[Risk] 大量明細行のパフォーマンス**: 1仕訳に数十行の明細がある場合の API レスポンスサイズ → フリーランスのユースケースでは現実的に問題にならない
- **[Trade-off] 単一仕訳モードと複合仕訳モードの状態管理**: UI に2つのモードがあることで状態管理が複雑になる → モード切替時にデータを変換する明確なロジックを持つ
