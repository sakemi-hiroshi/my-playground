---
applyTo: "**/*_test.go"
---

## テスト実行

```bash
go test ./...   # 全テスト実行
```

## テストの書き方

- パッケージは既存に合わせる（`package book` または `package book_test`）
- テーブル駆動テストで正常・境界値・エラーの各ケースを網羅する:

```go
tests := []struct {
    name string
    // ...
}{
    {"正常系: ...", ...},
    {"エラー系: ...", ...},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

- HTTP ハンドラのテストは `net/http/httptest` を使う:

```go
w := httptest.NewRecorder()
r := httptest.NewRequest(http.MethodGet, "/books/1", nil)
handler.ServeHTTP(w, r)
```

- アサーションは標準の `if got != want { t.Errorf(...) }` で十分（`testify` 等の外部ライブラリは追加しない）

## 制約

- 既存テストを削除・無効化しない
- `t.Skip()` で無条件スキップする実装を追加しない
- 外部モックライブラリ追加禁止
