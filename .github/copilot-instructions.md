# Copilot Instructions

## リポジトリ概要

Go製のシンプルなBook CRUD APIを持つ実験用リポジトリ。
AI agent（GitHub Copilot / Claude Code）を使ったissue駆動開発の試し撃ち場として使用する。

## 技術スタック

- Language: Go 1.22+
- HTTP: `net/http`（標準ライブラリのみ）
- Storage: インメモリ（外部DBなし）

## ドメインモデル

`internal/book/` に Book CRUD を実装。

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    ISBN   string `json:"isbn"`
}
```

## API エンドポイント

| Method | Path        | 説明         |
|--------|-------------|--------------|
| GET    | /books      | 一覧取得     |
| POST   | /books      | 新規作成     |
| GET    | /books/{id} | 1件取得      |
| PUT    | /books/{id} | 更新         |
| DELETE | /books/{id} | 削除         |

## 言語

コミットメッセージ・PR タイトル・PR 本文・issue コメントはすべて**日本語**で書くこと。

## Copilot agentへの期待役割

- シンプルな機能追加・修正issueを担当する
- 実装は `internal/book/` 配下に集約する
- 外部依存ライブラリは追加しない（標準ライブラリのみ）
- テストを追加する場合は `_test.go` ファイルに記述する

**実装タスクを担当する場合は、必ず `.github/skills/implementing-feature/SKILL.md` を読み、観点に従うこと。**

## 設計が必要な複雑なissue

設計の壁打ちが必要な場合は Copilot ではなく Claude Code（issue-design スキル）で対応する。

## PRレビュー観点

**PRレビューを行う際は、必ず `.github/skills/code-review-lenses/SKILL.md` を読み込み、そこに記載された指示にすべて従うこと。**
**Tell, Don't Ask 違反の疑いがあれば `.github/skills/tell-dont-ask/SKILL.md` を併用すること。**
**簡素化・YAGNI・dead code の指摘を行う場合は `.github/skills/simplifying-code/SKILL.md` を参照すること。**

レビューは4観点（バグ／セキュリティ／パフォーマンス／保守性）と Critical / Important / Suggestion の重大度分類で行う。
