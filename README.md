# fleee

フリーランスエンジニア向けの複式簿記アプリケーション。セルフホスト可能な単一バイナリとして配布する。

## 前提条件

- Go 1.26+
- Node.js 22+
- [golangci-lint](https://golangci-lint.run/) v2
- [lefthook](https://github.com/evilmartians/lefthook) (pre-commit hooks)

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/evilmartians/lefthook@latest
```

## セットアップ

```bash
git clone https://github.com/kosuke/fleee.git
cd fleee
cd web && npm install && cd ..
lefthook install
```

## ビルド

```bash
make build
```

フロントエンドのビルドと Go のコンパイルが順に実行され、単一バイナリ `fleee` が生成される。

## 開発

バックエンドとフロントエンドを別ターミナルで起動する:

```bash
# ターミナル 1: API サーバー
make dev-api

# ターミナル 2: フロントエンド (Vite dev server)
make dev-web
```

Vite dev server は `/api/*` へのリクエストを Go サーバー (localhost:8080) にプロキシする。

## テスト

```bash
make test          # Go + Frontend のテスト実行
make cover         # カバレッジレポート付きで実行
make lint          # golangci-lint + ESLint + Prettier チェック
```

## 起動

```bash
./fleee serve              # デフォルト: port 8080, DB fleee.db
./fleee serve -port 3000   # ポート指定
./fleee serve -db my.db    # DB ファイル指定
```
