## Context

fleee は複式簿記の帳簿記録アプリケーション。初回の change であり、プロジェクト構造からすべて新規に構築する。主な対象ユーザーはフリーランスエンジニアで、仕入・在庫・消費税は扱わないシンプルなスコープ。将来的にセルフホスト可能な OSS として公開することを見据え、単一バイナリでの配布を前提とする。

## Goals / Non-Goals

**Goals:**

- 勘定科目の CRUD と CSV/TSV インポートが動作する単一バイナリを配布可能にする
- 将来の機能追加（仕訳、元帳、財務諸表）に耐えるプロジェクト構造を確立する
- 開発者が `go run` / `npm run dev` で即座に開発を始められる環境を整える

**Non-Goals:**

- 認証・認可（後続の change で対応）
- マルチテナント対応
- PostgreSQL 対応（SQLite のみ）
- 消費税・仕入・在庫に関する勘定科目の特別な扱い
- モバイル対応

## Decisions

### 1. ディレクトリ構成

```
fleee/
├── cmd/
│   └── fleee/
│       └── main.go          # エントリポイント
├── internal/
│   ├── server/
│   │   └── server.go        # HTTP サーバー設定、ルーティング
│   ├── handler/
│   │   └── account.go       # 勘定科目 API ハンドラ
│   ├── model/
│   │   └── account.go       # ドメインモデル
│   ├── store/
│   │   ├── db.go            # DB 接続、マイグレーション
│   │   └── account.go       # 勘定科目リポジトリ
│   └── importer/
│       └── account_csv.go   # CSV/TSV インポートロジック
├── migrations/
│   └── 001_create_accounts.sql
├── web/                      # フロントエンド (React + Vite)
│   ├── src/
│   │   ├── App.tsx
│   │   ├── pages/
│   │   │   └── accounts/
│   │   ├── components/
│   │   └── api/
│   ├── package.json
│   └── vite.config.ts
├── go.mod
├── go.sum
└── Makefile
```

**理由:** Go の標準的なレイアウト。`internal/` で公開範囲を制限し、`cmd/` でバイナリのエントリポイントを分離する。フロントエンドは `web/` に配置し、ビルド成果物を `go:embed` で埋め込む。

### 2. HTTP ルーター: `chi`

標準ライブラリの `net/http` でも十分だが、chi は `net/http` 互換でありながらミドルウェアチェーンとURLパラメータのルーティングを簡潔に書ける。依存が軽量で、将来認証ミドルウェアを追加しやすい。

**代替案:** 標準 `net/http` — ルーティングの記述が冗長になる。echo/gin — 独自の Context 型を持ち `net/http` からの距離が生まれる。

### 3. SQLite ドライバー: `modernc.org/sqlite`

CGo 不要の pure Go 実装。クロスコンパイルが容易で、セルフホスト向け単一バイナリ配布と相性が良い。

**代替案:** `mattn/go-sqlite3` — CGo 必要。パフォーマンスは優れるが、ビルド環境に C コンパイラが必要になりクロスコンパイルが困難。

### 4. データモデル: 5区分の account_type

```sql
CREATE TABLE accounts (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT    NOT NULL UNIQUE,
    account_type  TEXT    NOT NULL CHECK (account_type IN ('asset', 'liability', 'equity', 'revenue', 'expense')),
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);
```

`account_type` から以下を導出:
- **normal_balance**: asset, expense → debit / liability, equity, revenue → credit
- **statement_type**: asset, liability, equity → balance_sheet / revenue, expense → income_statement

**理由:** 探索フェーズで議論した通り、貸借タイプ + 精算種別 の2カラムでは負債と純資産を区別できない。5区分を明示的に持つことで正確な財務諸表生成の基盤を作る。

### 5. CSV/TSV インポート時の account_type 判定

ユーザーの既存マスタ形式（科目名, 科目貸借タイプ, 出力順番, 精算種別）からの変換ルール:

| 科目貸借タイプ | 精算種別 | → account_type |
|---|---|---|
| 借方 | 貸借対照表 | asset |
| 借方 | 損益計算書 | expense |
| 貸方 | 損益計算書 | revenue |
| 貸方 | 貸借対照表 | **判定不可** → デフォルト liability、既知の科目名（元入金）は equity にマッピング |

「貸方 + 貸借対照表」の判定: 純資産に該当する科目は限定的（元入金 等）なので、既知の科目名リストでマッピングし、該当しなければ liability とする。インポート後にユーザーが UI で修正可能。

### 6. フロントエンド: React + Vite + Tailwind CSS + shadcn/ui

- **Tailwind CSS**: ユーティリティファーストで、コンポーネントライブラリとの相性が良い
- **shadcn/ui**: コピー&ペースト方式のUIコンポーネント集。テーブル、フォーム、ダイアログなど帳簿アプリに必要な部品が揃っている。依存としてではなくソースとして取り込むため、カスタマイズしやすい

**代替案:** MUI — 重厚で、単純なアプリにはオーバーヘッドが大きい。Ant Design — 充実しているがバンドルサイズが大きい。

### 7. API 設計

```
GET    /api/accounts          勘定科目一覧
POST   /api/accounts          勘定科目作成
GET    /api/accounts/:id      勘定科目取得
PUT    /api/accounts/:id      勘定科目更新
DELETE /api/accounts/:id      勘定科目削除
POST   /api/accounts/import   CSV/TSV インポート
```

JSON リクエスト/レスポンス。インポートは `multipart/form-data` でファイルを受け付ける。

### 8. フロントエンド埋め込みとdev モード

- **本番ビルド:** `web/dist/` を `go:embed` で Go バイナリに埋め込む。Go サーバーが API と静的ファイルの両方を配信
- **開発時:** Vite dev server (port 5173) + Go サーバー (port 8080) を並行起動。Vite の proxy 設定で `/api/*` を Go に転送

## Risks / Trade-offs

- **modernc.org/sqlite のパフォーマンス**: CGo 版より遅いが、個人の帳簿レベルのデータ量では問題にならない → 許容
- **フロントエンド埋め込みのビルド複雑性**: Make タスクで `npm run build` → `go build` の順序を管理する必要がある → Makefile で自動化
- **CSV インポートの account_type 自動判定**: 貸方 + 貸借対照表 の判定が完全ではない → 既知科目名マッピング + UI での修正で対応
