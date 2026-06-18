# CLAUDE.md — fleee プロジェクト

## プロジェクト概要

フリーランスエンジニア向けの複式簿記アプリケーション。
仕入・在庫・消費税は扱わないシンプルなスコープ。セルフホスト可能な単一バイナリとして配布する。

## 技術スタック

| レイヤー | 技術 |
|----------|------|
| バックエンド | Go, chi (ルーター) |
| データベース | SQLite (modernc.org/sqlite, pure Go) |
| フロントエンド | React, Vite, TypeScript, Tailwind CSS, shadcn/ui |
| 配布形態 | 単一バイナリ (go:embed でフロント埋め込み) |

## テスト方針

### 戦略

- **Go**: SQLite `:memory:` を使ったインテグレーションテスト。モックは使わず全レイヤーで実 DB を通す
- **Frontend**: Vitest + React Testing Library + jsdom。API モックは `vi.mock` を使用
- **カバレッジ目標**: Go / Frontend ともに 80% 以上を CI で enforce

### Go テストの書き方

- `testutil.NewTestDB(t)` で in-memory DB を取得（マイグレーション済み、`t.Cleanup()` で自動クローズ）
- store テスト: 正常系 + エラー系（ErrNotFound, ErrDuplicateName 等）
- handler テスト: `httptest.NewRecorder` + `chi.NewRouter()` で HTTP サイクル全体をテスト
- テストファイルは対象ファイルと同じディレクトリに `*_test.go` で配置

### Frontend テストの書き方

- テストファイルは対象コンポーネントと同じディレクトリに `*.test.tsx` で配置
- `src/components/ui/` (shadcn/ui) はテスト・カバレッジ対象外
- API 呼び出しは `vi.mock('@/api/...')` でモック

### Lint / Format

- Go: golangci-lint (govet, staticcheck, errcheck, unused, gofmt)
- Frontend: ESLint + Prettier (prettier-plugin-tailwindcss)

## AI ツール役割分担

このプロジェクトでは Claude Code と Antigravity CLI (agy) を役割に応じて使い分ける。

| フェーズ | ツール | 理由 |
|---------|--------|------|
| 設計・調査・探索 | **Claude Code** | 設計能力・コードベース理解が高い |
| 実装 | **Antigravity CLI** (`agy`) | Gemini クレジットを活用、Claude クレジットを節約 |
| コードレビュー・PR作成 | **Claude Code** | `/code-review` スキルを使う |
| レビュー指摘の修正 | **Antigravity CLI** (`agy`) | 実装担当が修正。Claude Code は修正内容を指示するのみ |

### 難易度の目安

- **低**: 既存パターンのコピー・追記、2ファイル以内、マイグレーション不要
- **中**: 複数ファイルにまたがる変更、新しいロジックの追加、スキーマ変更を伴わない
- **高**: アーキテクチャ変更、DBマイグレーション、複雑な状態管理、セキュリティ要件

難易度が高い実装や、agy が行き詰まった場合は Claude Code に切り替える。

### Claude Code から agy を呼び出す方法

実装は Claude Code から次のコマンドで agy に委譲する:

```bash
GIT_TERMINAL_PROMPT=0 CI=true \
  agy --dangerously-skip-permissions --print-timeout 3m --print "<実装プロンプト>" 2>&1
```

**タイムアウト値の目安:**

| タスク規模 | ファイル数 | タイムアウト |
|-----------|-----------|------------|
| 小 | 1-2 | `--print-timeout 3m` |
| 中 | 3-5 | `--print-timeout 5m` |
| 大 | 6+ | `--print-timeout 8m`（さらなる分割を検討） |

`CI=true` の副作用でツールの動作が変わった場合は、当該環境変数を除外してプロンプト内ルールのみでハング予防する。

**タスク単位の分割呼び出し:**
一括委譲ではなく、OpenSpec tasks.md の各タスクを個別の agy ワンショットで実行する。
十分に小さいタスク（1ファイルの小変更）は Claude Code の判断でまとめてよい。
各タスク完了後に `git status` / `git diff` で結果を確認してから次のタスクに進む。

プロンプトには以下を含める:
- 実装対象ファイルと変更内容（具体的に）
- 既存コードのパターン（コピーすべき書き方）
- スコープ外の制約（「この2ファイルだけ触れ」など）
- OpenSpec の change ディレクトリへの参照（`openspec/changes/<name>/` 以下）
- **実装完了後に `git add <実装ファイル> && git commit -m "feat: <変更内容>"` でコミットすること**
- **不明点・判断に迷う点がある場合は実装せず、`[QUESTION] ...` の形式で質問を出力すること**
- **ハング予防ルール（以下の「agy プロンプトの必須ルール」セクション参照）**

実装完了後は Claude Code で `git diff HEAD` または `git log` を確認してからレビューに進む。

### agy プロンプトの必須ルール

agy へのすべてのプロンプトに、以下のルールを必ず含める:

```
制約:
- 対話的入力（y/n, パスワード等）を求めるコマンドは絶対に実行しないこと
- 必ず非対話フラグ（--yes, -y, --no-input 等）を付けること
- git push, npm publish など外部サービスへの送信は行わないこと
- 対話的入力が必要な状況に遭遇したら、実行せず [QUESTION] で報告すること
```

### agy との対話ループ

agy の出力に `[QUESTION]` が含まれる場合、以下のループで対応する:

```
agy 実行 (1回目)
  ↓ 出力を解析
  ├─ [QUESTION] なし → 実装完了、レビューへ
  └─ [QUESTION] あり → 質問ごとに判定:
       ├─ Claude Code が回答できる → 回答をまとめる
       └─ 判断つかない / センシティブ → ユーザーに確認
       ↓
     agy 実行 (2回目: 回答を追加コンテキストとして渡す)
       ↓ 再度出力を解析（同じループ）
```

**ループの上限は 3 回**。3 回で解決しない場合は Claude Code が直接実装に切り替える。

### agy タイムアウト時のリカバリ

agy がタイムアウトまたはハングした場合、以下のフローで対応する:

```
タイムアウト発生
  ↓
git status / git diff で途中成果を確認
  ├─ 成果あり → --continue で残り作業を指示
  │   agy --dangerously-skip-permissions --print-timeout 3m \
  │     --continue --print "前回タイムアウトした。対話的入力を避けて続行しろ"
  └─ 成果なし → 新規セッションで別アプローチを試行
      ↓
再試行も失敗（2回連続）
  ↓
Claude Code が直接実装に切り替える
```

**注意:** `--print-timeout` が効かないケースがある。Bash ツールの `timeout` パラメータも併用し、
agy プロセスが応答しない場合は `kill` して対処する。

**Claude Code がユーザーに確認すべきケース:**
- ビジネスロジックの仕様判断（「この場合どう振る舞うべきか」）
- 破壊的変更や後方互換性に関わる判断
- セキュリティ・認証に関わる設計判断
- 複数の妥当な選択肢がありトレードオフが明確でない場合

**Claude Code が自分で回答してよいケース:**
- 既存コードのパターンやコンベンションに関する質問
- API の型やインターフェースの確認
- ファイル構成やインポートパスの確認
- OpenSpec の design.md / specs から読み取れる仕様

---

## 開発ワークフロー

### 機能追加・バグ修正の標準フロー

```
1. 設計 (Claude Code)
   /opsx:explore  → 問題を探索し設計を固める
   /opsx:propose  → change proposal を生成（proposal.md / design.md / specs / tasks.md）

2. トピックブランチ作成 & proposal コミット (Claude Code)
   git checkout -b feature/<change-name>
   git add openspec/changes/<change-name>/
   git commit -m "docs(openspec): propose <change-name>"

3. 実装 (agy タスク単位 × N)
   tasks.md の各タスクを個別の agy ワンショットで実行
   GIT_TERMINAL_PROMPT=0 CI=true \
     agy --dangerously-skip-permissions --print-timeout 3m --print "..." 2>&1
   → 各タスク完了後に git status / git diff で確認
   → [QUESTION] があれば Claude Code が回答 or ユーザーに確認し再実行
   → タイムアウト時は --continue で再開 or Claude Code が引き継ぐ
   → agy が実装コードをコミットする（"feat: <変更内容>"）

3.5. agy 別セッションレビュー (agy 新規セッション)
   実装とは別の新規 agy セッションで diff をレビューさせる
   → 変更要約（変更ファイル一覧、問題点、設計判断、テスト充足度）を出力させる
   → Claude Code は要約をベースに最終レビューの判断材料にする
   → 変更が 1 ファイル・20 行以下の場合はこのステップを省略してよい

4. コードレビュー (Claude Code)
   /code-review（agy レビュー結果も参考にしつつ最終判断）
   → 指摘があれば、修正内容を具体的にまとめる

5. レビュー指摘の修正 (agy)
   Claude Code がレビュー結果から修正プロンプトを作成し agy に委譲する
   → 修正対象ファイル・行・具体的な変更内容を指示する
   → agy が修正コードをコミットする（"fix: <修正内容>"）
   → 指摘がなければこのステップはスキップ

6. PR 作成 & OpenSpec アーカイブ (Claude Code)
   /opsx:archive  → change をアーカイブ（delta spec sync を含む）
   gh pr create
   → アーカイブと spec sync のコミットを PR に含める

7. CI 確認 (Claude Code)
   gh pr checks --watch
   → 失敗したら是正してプッシュし、再度 watch する
```

### ドキュメント更新の原則

**設計・実装・レビューのすべての場面で常に検討すること。**

- 仕様変更・機能追加があれば `openspec/specs/` の対応する spec.md を更新する
- API の変更・新規エンドポイントがあれば関連ドキュメントに反映する
- CLAUDE.md 自体のワークフローや規約が変わったときは即座に更新する
- agy への実装プロンプトにも「関連ドキュメントの更新が必要か検討すること」を明示する

### CI / pre-commit ルール

**pre-commit のスキップ禁止**
`--no-verify` や `SKIP=...` による pre-commit フックのスキップは一切行わない。
フックが失敗した場合は、スキップせず根本原因を修正してから再コミットする。

**GitHub Actions の失敗は必ず是正する**
pre-commit が通っても GitHub Actions が失敗した場合は問題とみなす。
PR 作成後は必ず次のコマンドで CI 結果を待ち、失敗があれば修正してプッシュすること:

```bash
gh pr checks --watch   # 成功: exit 0 / 失敗: exit 非0
```

失敗した場合は `gh pr checks` でどのジョブが落ちたかを確認し、修正 → コミット → プッシュ → 再度 watch のサイクルで是正する。
CI がすべてグリーンになったことを確認してから作業完了を報告する。

### OpenSpec チートシート

| スキル | 用途 |
|--------|------|
| `/opsx:explore` | アイデア・問題を探索する（実装しない） |
| `/opsx:propose` | change proposal を一括生成 |
| `/opsx:apply` | tasks.md のタスクを順に実装（Claude Code が担当する場合） |
| `/opsx:sync` | delta spec を main spec にマージ |
| `/opsx:archive` | 完了した change をアーカイブ（PR 作成前に実施） |

### ブランチ命名規則

```
feature/<change-name>   # 機能追加（OpenSpec change 名と一致させる）
fix/<change-name>       # バグ修正
```
