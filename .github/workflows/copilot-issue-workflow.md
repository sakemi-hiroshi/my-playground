# Copilot Issue Workflow

## issueをアサインされたときの手順

1. **issue内容の確認**
   - issue本文・コメントを読む
   - `.github/copilot-instructions.md` の制約（標準ライブラリのみ等）を確認

2. **コードベースの調査**
   - `internal/book/` 配下の既存実装を確認
   - 変更が必要なファイルを特定

3. **実装**
   - `internal/book/` に実装を集約する
   - 外部ライブラリを追加しない（標準ライブラリのみ）
   - テストは `_test.go` ファイルに記述する

4. **PR作成**
   - タイトル形式: `feat: ...` / `fix: ...` / `refactor: ...`
   - 本文に `Closes #N` でissue番号を記載する

---

## HARD-GATE（対応してはいけない条件）

- issueのコメントに `<!-- m11n:state:workpad -->` マーカーがある場合：
  Claude Codeが設計済みのissue。実装の詳細はWorkpadコメントに従う。
- issueに `needs-design` ラベルがある場合：
  設計の壁打ちが必要。Copilotは対応せず、Claude Code（`issue-design` スキル）に委ねる。

---

## このリポジトリの実装マップ

| ファイル | 役割 |
|---|---|
| `internal/book/model.go` | Book 構造体の定義 |
| `internal/book/repository.go` | インメモリCRUD操作 |
| `internal/book/handler.go` | HTTP リクエストのルーティングと処理 |
| `cmd/api/main.go` | サーバー起動エントリーポイント |
