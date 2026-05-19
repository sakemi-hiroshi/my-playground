---
name: implementing-feature
description: Go プロジェクトでの機能実装ガイド。issue の要件把握、既存規約の踏襲、テスト駆動の実装、ビルド・テスト確認までの観点を整理する。Copilot agent が issue 実装を担当する時に参照する。
---

# 機能実装ガイド

## 着手前に必ず読む

実装を始める前に以下を順番に Read して、リポジトリの規約・コードスタイルを把握する。

1. `.github/copilot-instructions.md` — 技術スタック・配置ルール・期待役割
2. `internal/book/handler.go` — ハンドラの構造・receiver 名・エラーハンドリング・JSON 返却パターン
3. `internal/book/repository.go` — リポジトリのメソッド名・エラー定義（`ErrNotFound` 等）
4. `internal/book/model.go` — 型定義
5. 既存の `*_test.go` — テスト関数命名・テーブル駆動・httptest パターン

把握した上で実装を始め、**既存スタイルと一貫させる**。

---

## 実装ガイドライン

### パッケージ・配置

- 実装は `internal/book/` 配下に集約する
- ファイル名はドメイン語（例: `search.go`, `recommendation.go`）
- 外部依存ライブラリは追加しない（標準ライブラリのみ）

### ハンドラの書き方

既存パターンに合わせる:

```go
// receiver は常に h *Handler
func (h *Handler) methodName(w http.ResponseWriter, r *http.Request) {
    // 入力取得
    // 処理
    // エラー時: http.Error(w, "メッセージ", http.StatusXxx)
    // 正常時: writeJSON(w, http.StatusOK, result)
}
```

- HTTP ステータスは必ず定数（`http.StatusOK` 等）を使う — 生の `200` は書かない
- エラー判定は `errors.Is(err, ErrNotFound)` パターンを踏襲する
- 未使用パラメータは `_` に揃える（`_ *http.Request` 等、既存に合わせる）

### 命名規則

| 対象 | 規則 | 例 |
|------|------|----|
| ハンドラメソッド | 小文字・動詞 | `search`, `recommend` |
| リポジトリメソッド | `FindXxx` / `Save` / `Update` / `Delete` | `FindByAuthor` |
| エラー変数 | `ErrXxx` | `ErrNotFound` |
| テスト関数 | `TestHandlerName_条件` | `TestHandler_search_found` |

---

## テスト方針

新規機能には必ず `*_test.go` を添える。

- パッケージ: `package book` （ホワイトボックス）または `package book_test` — 既存に合わせる
- テーブル駆動テストで複数ケース（正常・境界値・エラー）を網羅する
- HTTP ハンドラのテストは `net/http/httptest` の `httptest.NewRecorder()` と `httptest.NewRequest()` を使う
- テスト実行:

```bash
go test ./...
```

すべて `PASS` になるまで修正を続ける。

---

## ビルド・テスト確認（実装完了前に必ず実行）

```bash
go build ./...   # コンパイルエラーゼロ
go test ./...    # テスト失敗ゼロ
```

どちらかが失敗した場合は修正してリトライする。確認が取れてから PR を作成する。

---

## PR 作成時のセルフチェック

実装完了後、PR を出す前に以下を確認する:

- [ ] 関係ないファイルを変更していない
- [ ] `copilot-instructions.md` の役割範囲を逸脱していない（設計が必要な大規模変更はしない）
- [ ] 標準ライブラリのみを使っている（`go.mod` の依存が増えていない）
- [ ] マジックナンバーがない（`10` や `200` 等をハードコードしていない）
- [ ] dead code を残していない（呼ばれない関数・未使用変数を削除した）
- [ ] ドメイン判断がハンドラに漏れていない（Tell, Don't Ask 観点）

---

## 避けるべきこと

- 外部依存ライブラリの追加
- 既存テストの削除・無効化
- 設計判断が必要な変更（アーキテクチャ変更・型設計変更等）— その場合は実装せず issue にコメントで設計レビュー要請する
- `main` ブランチへの直接 push
