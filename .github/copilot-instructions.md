# Copilot Instructions

## リポジトリ概要

Go 製の Book CRUD API を持つ実験用リポジトリ。AI agent（GitHub Copilot Coding Agent / Claude Code）を使った issue 駆動開発の試し撃ち場として使用する。

## 技術スタック

- Language: Go 1.26.1
- HTTP: `net/http`（標準ライブラリのみ）
- Storage: インメモリ（外部 DB なし）
- kinesis sample のみ `aws-sdk-go-v2` を使用（Copilot agent の作業範囲外）

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

## パッケージ構成

| ファイル | 役割 |
|---|---|
| `internal/book/model.go` | Book 構造体の定義 |
| `internal/book/repository.go` | インメモリ CRUD 操作 |
| `internal/book/handler.go` | HTTP ルーティングと処理 |
| `internal/book/recommendation.go` | おすすめ機能 |
| `cmd/api/main.go` | サーバー起動エントリーポイント |

`internal/envelope/` と `cmd/kinesis_sample/` は kinesis 実験用の別系統。**Copilot agent は触らない**。

## 言語

コミットメッセージ・PR タイトル・PR 本文・issue コメントはすべて**日本語**で書くこと。

## Copilot agent への期待役割

- シンプルな機能追加・修正 issue を担当する
- 実装は `internal/book/` 配下に集約する
- 外部依存ライブラリは追加しない（標準ライブラリのみ）
- テストを追加する場合は `*_test.go` ファイルに記述する
- 設計が必要な複雑な issue は Copilot で対応せず、issue にコメントで設計レビューを要請する

## HARD-GATE（対応してはいけない条件）

- issue のコメントに `<!-- m11n:state:workpad -->` マーカーがある場合: Claude Code が設計済みの issue。実装の詳細は Workpad コメントに従う
- issue に `needs-design` ラベルがある場合: 設計の壁打ちが必要。Copilot は対応しない
