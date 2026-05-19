---
applyTo: "**/*.go"
---

## パッケージ・配置

- Book ドメインの実装は `internal/book/` 配下に集約する（`cmd/kinesis_sample/` や `internal/envelope/` には手を出さない）
- ファイル名はドメイン語（例: `search.go`, `recommendation.go`）
- 外部依存ライブラリは追加しない（標準ライブラリのみ）

## ハンドラの書き方

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

- HTTP ステータスは必ず定数（`http.StatusOK` 等）を使う（生の `200` はNG）
- エラー判定は `errors.Is(err, ErrNotFound)` パターンを踏襲する
- 未使用パラメータは `_` に揃える（`_ *http.Request` 等）
- パスパラメータ抽出は `strings.TrimPrefix` を使う（chi 等のルータ追加禁止）

## 命名規則

| 対象 | 規則 | 例 |
|---|---|---|
| ハンドラメソッド | 小文字・動詞 | `search`, `recommend` |
| リポジトリメソッド | `FindXxx` / `Save` / `Update` / `Delete` | `FindByAuthor` |
| エラー変数 | `ErrXxx` | `ErrNotFound` |
| テスト関数 | `TestHandlerName_条件` | `TestHandler_search_found` |

## 設計原則

- 既存の抽象化レベルを維持する（新しいインターフェースや層を追加しない・YAGNI）
- `Repository` を struct 直接使用しているのが現行スタイル。interface 化しない
- dead code を残さない（呼ばれない関数・未使用変数は削除）
