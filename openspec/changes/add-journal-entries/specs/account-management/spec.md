## MODIFIED Requirements

### Requirement: Delete account
ユーザーは勘定科目を削除できなければならない。ただし、仕訳の明細行で使用されている勘定科目は削除できない。

#### Scenario: Successful account deletion
- **WHEN** ユーザーが仕訳で使用されていない勘定科目「雑費」を削除する
- **THEN** 勘定科目が削除され、一覧から消える

#### Scenario: Deletion rejected when used in journal entries
- **WHEN** ユーザーが仕訳の明細行で使用されている勘定科目を削除しようとする
- **THEN** FK 制約エラーが返され、勘定科目は削除されない
