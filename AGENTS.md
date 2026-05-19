# AGENTS.md

Go 製の Book CRUD API を持つ実験用リポジトリ。AI agent 駆動開発の試し撃ち場として使用する。

リポジトリ詳細・制約・HARD-GATE は [`.github/copilot-instructions.md`](.github/copilot-instructions.md) を先に読むこと。

## ブランチ運用

- `main` ブランチへの直接 push は禁止
- ブランチ名: `feat/<概要>` / `fix/<概要>` / `refactor/<概要>` / `experiment/<概要>`
- Copilot Coding Agent が作成する場合は `copilot/<issue番号>-<概要>`

## PR ルール

- タイトル: Conventional Commits 形式（`feat:` / `fix:` / `refactor:` / `test:`）
- 本文: 変更内容のサマリー + `Closes #N`
- コミットメッセージ・PR タイトル・PR 本文・issue コメントはすべて**日本語**
- 基本は **draft PR** として作成（Ready for Review は人間が判断）

## ビルド・テスト要件

PR 作成前に以下を必ず通す:

```bash
go build ./...   # コンパイルエラーゼロ
go test ./...    # テスト失敗ゼロ
```

## 実装制約

- 外部依存ライブラリの追加禁止（標準ライブラリのみ）
- 実装は `internal/book/` 配下に集約する
- 既存テストの削除・無効化禁止
- 設計が必要な複雑な issue は実装せず、issue にコメントで設計レビューを要請する

## 詳細規約

path-specific な詳細規約は `.github/instructions/` を参照:

- Go 実装: [`.github/instructions/go-implementation.instructions.md`](.github/instructions/go-implementation.instructions.md)
- テスト: [`.github/instructions/go-tests.instructions.md`](.github/instructions/go-tests.instructions.md)
