## Purpose

仕訳の管理機能。ヘッダー（日付・摘要・領収書要否・メモ）+ 明細行（科目・借方金額・貸方金額）の複合仕訳モデルによる CRUD 操作を提供する。

## Requirements

### Requirement: Journal entry data model
仕訳はヘッダーと明細行で構成されなければならない。ヘッダーは日付（必須）、摘要（必須）、領収書の要否（必須、デフォルト false）、メモ（任意）を持つ。明細行は勘定科目（必須）、借方金額（整数、0以上）、貸方金額（整数、0以上）を持ち、各行は借方か貸方のいずれか一方のみに正の値を持たなければならない。

#### Scenario: Single journal entry structure
- **WHEN** 通信費 10,000円を普通預金から支払う仕訳を作成する
- **THEN** ヘッダー1件と明細行2件（借方: 通信費 10,000、貸方: 普通預金 10,000）が保存される

#### Scenario: Compound journal entry structure
- **WHEN** 出張経費として旅費交通費 5,000円と会議費 10,000円を普通預金から支払う仕訳を作成する
- **THEN** ヘッダー1件と明細行3件（借方: 旅費交通費 5,000、借方: 会議費 10,000、貸方: 普通預金 15,000）が保存される

### Requirement: Balance validation
仕訳の明細行の借方合計と貸方合計は一致しなければならない。不一致の場合は保存を拒否しエラーを返さなければならない。

#### Scenario: Balanced entry accepted
- **WHEN** 借方合計 10,000 と貸方合計 10,000 の仕訳を保存する
- **THEN** 仕訳が正常に保存される

#### Scenario: Unbalanced entry rejected
- **WHEN** 借方合計 10,000 と貸方合計 8,000 の仕訳を保存する
- **THEN** エラーが返され、仕訳は保存されない

### Requirement: Minimum lines validation
仕訳は最低2件の明細行を持たなければならない。

#### Scenario: Entry with less than 2 lines rejected
- **WHEN** 明細行が1件のみの仕訳を保存しようとする
- **THEN** エラーが返され、仕訳は保存されない

### Requirement: Create journal entry
ユーザーは日付、摘要、明細行（科目・借方金額・貸方金額）、領収書の要否、メモを指定して仕訳を作成できなければならない。

#### Scenario: Successful creation
- **WHEN** ユーザーが日付「2026-06-18」、摘要「携帯電話料金6月分」、借方: 通信費 10,000、貸方: 普通預金 10,000 で仕訳を作成する
- **THEN** 仕訳がヘッダーと明細行を含めて保存され、作成された仕訳が返される

#### Scenario: Missing required fields rejected
- **WHEN** 日付または摘要が空の仕訳を作成しようとする
- **THEN** エラーが返され、仕訳は保存されない

#### Scenario: Invalid account ID rejected
- **WHEN** 存在しない勘定科目 ID を明細行に指定して仕訳を作成しようとする
- **THEN** FK 制約エラーが返され、仕訳は保存されない

### Requirement: List journal entries
ユーザーは登録済みの仕訳を一覧表示できなければならない。一覧は日付の降順で表示される。

#### Scenario: List all entries
- **WHEN** ユーザーが仕訳一覧画面を開く
- **THEN** 全ての仕訳がヘッダー情報（日付、摘要、合計金額）とともに日付降順で表示される

#### Scenario: Empty state
- **WHEN** 仕訳が1件も登録されていない状態で一覧画面を開く
- **THEN** 仕訳が未登録である旨のメッセージが表示される

### Requirement: Get journal entry detail
ユーザーは個別の仕訳の詳細（ヘッダー + 全明細行）を取得できなければならない。

#### Scenario: Successful retrieval
- **WHEN** ユーザーが仕訳 ID を指定して詳細を取得する
- **THEN** ヘッダー情報と全明細行（科目名を含む）が返される

#### Scenario: Not found
- **WHEN** 存在しない仕訳 ID を指定する
- **THEN** 404 エラーが返される

### Requirement: Update journal entry
ユーザーは既存の仕訳のヘッダーと明細行を全体置換で更新できなければならない。

#### Scenario: Successful update
- **WHEN** ユーザーが既存の仕訳の摘要を変更し、明細行を更新する
- **THEN** ヘッダーと明細行が全体置換で更新される

#### Scenario: Update with unbalanced lines rejected
- **WHEN** 借方合計と貸方合計が一致しない明細行で仕訳を更新しようとする
- **THEN** エラーが返され、仕訳は更新されない

### Requirement: Delete journal entry
ユーザーは仕訳を削除できなければならない。ヘッダーの削除時に明細行もカスケード削除される。

#### Scenario: Successful deletion
- **WHEN** ユーザーが仕訳を削除する
- **THEN** ヘッダーと全明細行が削除される

#### Scenario: Delete not found
- **WHEN** 存在しない仕訳 ID で削除を試みる
- **THEN** 404 エラーが返される

### Requirement: Journal entry input UI
仕訳入力画面は単一仕訳モード（デフォルト）と複合仕訳モードの2つのモードを持たなければならない。

#### Scenario: Single entry mode (default)
- **WHEN** ユーザーが仕訳入力画面を開く
- **THEN** 借方科目、貸方科目、金額の各1つを入力するシンプルなフォームが表示される

#### Scenario: Switch to compound mode
- **WHEN** ユーザーが「複合仕訳に切替」ボタンをクリックする
- **THEN** 明細行テーブル形式のフォームに切り替わり、行の追加・削除が可能になる

#### Scenario: Switch from compound to single mode
- **WHEN** 複合仕訳モードで明細行が借方1行・貸方1行の状態で「単一仕訳に切替」ボタンをクリックする
- **THEN** 単一仕訳モードに戻り、既存の入力値が引き継がれる

#### Scenario: Balance indicator in compound mode
- **WHEN** 複合仕訳モードで明細行を入力する
- **THEN** 借方合計・貸方合計がリアルタイムで表示され、一致/不一致が視覚的に示される
