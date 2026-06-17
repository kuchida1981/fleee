## ADDED Requirements

### Requirement: Account data model
勘定科目は以下の属性を持たなければならない: 名前（一意）、勘定区分（asset, liability, equity, revenue, expense の5区分）、表示順序。

#### Scenario: Account type determines normal balance
- **WHEN** 勘定区分が asset または expense である
- **THEN** その科目の通常残高は借方（debit）として導出される

#### Scenario: Account type determines statement type
- **WHEN** 勘定区分が asset, liability, または equity である
- **THEN** その科目は貸借対照表（balance sheet）に分類される

#### Scenario: Account type determines statement type for PL
- **WHEN** 勘定区分が revenue または expense である
- **THEN** その科目は損益計算書（income statement）に分類される

### Requirement: Create account
ユーザーは新しい勘定科目を作成できなければならない。科目名、勘定区分、表示順序を指定する。

#### Scenario: Successful account creation
- **WHEN** ユーザーが科目名「通信費」、勘定区分「expense」、表示順序「3」で勘定科目を作成する
- **THEN** 勘定科目が保存され、一覧に表示される

#### Scenario: Duplicate name rejected
- **WHEN** 既に「通信費」という科目が存在する状態で同名の科目を作成しようとする
- **THEN** エラーが返され、重複した科目は作成されない

### Requirement: List accounts
ユーザーは登録済みの全勘定科目を一覧表示できなければならない。

#### Scenario: List all accounts
- **WHEN** ユーザーが勘定科目一覧画面を開く
- **THEN** 全ての勘定科目が表示順序に従って表示される

#### Scenario: Empty state
- **WHEN** 勘定科目が1件も登録されていない状態で一覧画面を開く
- **THEN** 科目が未登録である旨のメッセージが表示される

### Requirement: Update account
ユーザーは既存の勘定科目の名前、勘定区分、表示順序を編集できなければならない。

#### Scenario: Successful account update
- **WHEN** ユーザーが「通信費」の表示順序を「3」から「5」に変更する
- **THEN** 変更が保存され、一覧の表示順序が更新される

#### Scenario: Update rejects duplicate name
- **WHEN** 「通信費」を「消耗品費」（既存の科目名）にリネームしようとする
- **THEN** エラーが返され、変更は適用されない

### Requirement: Delete account
ユーザーは勘定科目を削除できなければならない。

#### Scenario: Successful account deletion
- **WHEN** ユーザーが勘定科目「雑費」を削除する
- **THEN** 勘定科目が削除され、一覧から消える

### Requirement: Import accounts from CSV/TSV
ユーザーは CSV または TSV ファイルから勘定科目を一括インポートできなければならない。ファイル形式は「科目名, 科目貸借タイプ, 出力順番, 精算種別」のヘッダ付き4カラム。

#### Scenario: Successful TSV import
- **WHEN** ユーザーが以下の形式の TSV ファイルをインポートする:
  ```
  科目名	科目貸借タイプ	出力順番	精算種別
  普通預金	借方	0	貸借対照表
  売上	貸方	1	損益計算書
  ```
- **THEN** 各行が account_type に変換されて保存される（借方+貸借対照表→asset、貸方+損益計算書→revenue）

#### Scenario: Import resolves ambiguous type
- **WHEN** インポートデータに「貸方 + 貸借対照表」の科目がある
- **THEN** 既知の科目名（元入金 等）は equity に、それ以外は liability にマッピングされる

#### Scenario: Import skips duplicate names
- **WHEN** インポートファイルに既に登録済みの科目名が含まれている
- **THEN** 重複した行はスキップされ、新規の行のみがインポートされる。スキップされた件数がユーザーに通知される

#### Scenario: Import rejects invalid format
- **WHEN** 必須カラムが不足したファイルをインポートする
- **THEN** エラーが返され、インポートは実行されない
