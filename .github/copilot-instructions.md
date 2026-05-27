# Copilot Instructions

## リポジトリ概要

Go 製の実験用リポジトリ。AI agent（GitHub Copilot Coding Agent / Claude Code）を使った issue 駆動開発の試し撃ち場として使用する。

2 つのドメインを持つ:
1. **Book CRUD API** - シンプルな書籍管理 (`internal/book/`)
2. **注文処理 Actor システム** - Actor モデル + Saga パターンによる注文処理 (`internal/order/`, `internal/service/`)

## 技術スタック

- Language: Go 1.26.1
- HTTP: `net/http`（Book API）/ `github.com/labstack/echo/v4`（注文処理 API）
- Actor: `github.com/asynkron/protoactor-go`
- Test: `github.com/stretchr/testify`
- Storage: インメモリ（外部 DB なし）
- kinesis sample のみ `aws-sdk-go-v2` を使用（Copilot agent の作業範囲外）

## ドメインモデル

### Book CRUD (`internal/book/`)

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    ISBN   string `json:"isbn"`
}
```

### 注文処理 Actor システム (`internal/order/`, `internal/service/`)

注文処理を Actor モデルで実装。`OrderActor` が Saga のオーケストレーターとなり、
`PaymentActor`・`PointActor`・`CouponActor` の各サービス Actor と非同期メッセージで協調する。

- `internal/order/actor.go` - OrderActor（Saga オーケストレーター）
- `internal/order/saga.go` - Saga ステップ定義
- `internal/order/messages.go` - OrderActor が扱うメッセージ型
- `internal/service/payment/` - 決済サービス Actor
- `internal/service/point/` - ポイントサービス Actor
- `internal/service/coupon/` - クーポンサービス Actor
- `internal/http/handler.go` - Echo を使った HTTP ハンドラ（注文処理 API）
- `internal/tracing/` - OpenTelemetry トレーシング
- `internal/failmode/` - 障害モード制御（テスト用）
- `cmd/order-server/main.go` - 注文処理サーバー起動エントリーポイント

## API エンドポイント

### Book API (`cmd/api/`)

| Method | Path        | 説明         |
|--------|-------------|--------------|
| GET    | /books      | 一覧取得     |
| POST   | /books      | 新規作成     |
| GET    | /books/{id} | 1件取得      |
| PUT    | /books/{id} | 更新         |
| DELETE | /books/{id} | 削除         |

### 注文処理 API (`cmd/order-server/`)

| Method | Path     | 説明         |
|--------|----------|--------------|
| POST   | /orders  | 注文処理開始 |

`internal/envelope/` と `cmd/kinesis_sample/` は kinesis 実験用の別系統。**Copilot agent は触らない**。

## 言語

コミットメッセージ・PR タイトル・PR 本文・issue コメントはすべて**日本語**で書くこと。

## Copilot agent への期待役割

- シンプルな機能追加・修正・リファクタ issue を担当する
- 実装スコープ: `internal/book/`・`internal/order/`・`internal/service/`・`internal/http/`
- 外部依存ライブラリは既存のもの（echo, protoactor-go, testify）のみ使用可。新規追加は禁止
- テストを追加する場合は `*_test.go` ファイルに記述する
- ビルドとテストは必ず確認: `go build ./...` `go test ./...`
- 設計が必要な複雑な issue は Copilot で対応せず、issue にコメントで設計レビューを要請する

## HARD-GATE（対応してはいけない条件）

- issue のコメントに `<!-- m11n:state:workpad -->` マーカーがある場合: Claude Code が設計済みの issue。実装の詳細は Workpad コメントに従う
- issue に `needs-design` ラベルがある場合: 設計の壁打ちが必要。Copilot は対応しない
