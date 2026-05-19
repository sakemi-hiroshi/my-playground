# my-playground

AI agent（GitHub Copilot / Claude Code）を使ったissue駆動開発の実験場。

## アプリ概要

Go製のシンプルな Book CRUD API。

```
GET    /books
GET    /books?q=<キーワード>
POST   /books
GET    /books/{id}
PUT    /books/{id}
DELETE /books/{id}
```

## 起動

```bash
go run ./cmd/api
```

## AI agent の役割分担

| タスク種別 | 担当 |
|---|---|
| シンプルな機能追加・修正 | GitHub Copilot agent（issueにアサイン） |
| 設計が必要な複雑な変更 | Claude Code（`/issue-design` スキル） |
