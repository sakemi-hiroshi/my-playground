---
name: go-design-principles
description: Go で API を実装する際の設計原則。Always-Valid Domain、panic と error の役割分離、Tell Don't Ask、Sentinel Error と詳細エラー型の使い分け、層間でのエラー変換。Go コードのレビュー時、新規実装時、リファクタリング時に参照する。
---

# Go 設計原則

## 1. Always-Valid Domain

ドメイン型はコンストラクタ経由でのみ生成する。生 struct リテラルでの生成は禁止。コンストラクタ内で `Validate()` を呼び、不変条件を検証する。

**やる:**
```go
func NewBook(title, author, isbn string) (Book, error) {
    b := Book{Title: title, Author: author, ISBN: isbn}
    if err := b.Validate(); err != nil {
        return Book{}, err
    }
    return b, nil
}

func (b Book) Validate() error {
    if b.Title == "" {
        return errors.New("title is required")
    }
    if b.ISBN == "" {
        return errors.New("isbn is required")
    }
    return nil
}
```

**やらない:**
```go
// 生 struct リテラルは不変条件を素通りさせる
b := Book{Title: "", ISBN: ""}  // invalid state that slips through
```

---

## 2. panic と error の役割分離

- **error**: 外部入力のバリデーション失敗、DBアクセス失敗、業務上ありえる失敗
- **panic**: 契約違反（前提条件が満たされていない）、設計上ありえない状態、初期化バグ

HTTP ハンドラに `recover` ミドルウェアを置いてすべての panic を飲み込むことはしない。panic は意図的に止める。

**やる:**
```go
// 外部入力のエラーは error で返す
func NewBook(title, author, isbn string) (Book, error) {
    if title == "" {
        return Book{}, ErrInvalidBook
    }
    return Book{Title: title}, nil
}

// プログラミングエラーは panic
func (r *InMemoryRepository) FindByID(id int) Book {
    if id <= 0 {
        panic(fmt.Sprintf("id must be positive, got %d", id))
    }
    // ...
}
```

**やらない:**
```go
// 業務エラーを panic にしない
func NewBook(title string) Book {
    if title == "" {
        panic("title required") // 呼び出し元でハンドリングできない
    }
    return Book{Title: title}
}
```

---

## 3. Tell, Don't Ask

業務判定はドメインに語らせ、ハンドラ／ユースケース層でドメイン内部の状態を取り出して判断しない。ハンドラは薄いオーケストレーターに留める。

**やる:**
```go
// ドメインが業務判定を持つ
func (b Book) IsAvailableForLoan() bool {
    return b.ISBN != "" && b.Author != ""
}

// ハンドラは結果を使うだけ
func (h *Handler) handleLoan(w http.ResponseWriter, r *http.Request) {
    book, err := h.repo.FindByID(id)
    if err != nil { /* ... */ }
    if !book.IsAvailableForLoan() {
        http.Error(w, "book not available", http.StatusConflict)
        return
    }
    // ...
}
```

**やらない:**
```go
// ハンドラが業務判定を知りすぎている
func (h *Handler) handleLoan(w http.ResponseWriter, r *http.Request) {
    book, _ := h.repo.FindByID(id)
    if book.ISBN == "" || book.Author == "" { // ドメイン知識がハンドラに漏れている
        http.Error(w, "book not available", http.StatusConflict)
        return
    }
}
```

---

## 4. Sentinel Error と詳細エラー型を対で持つ

`Err*` 変数は `errors.Is` での比較用、struct 型は `errors.As` で詳細情報を取り出すために使う。`Unwrap()` でエラーチェーンを維持する。

**やる:**
```go
// sentinel（errors.Is 用）
var ErrBookNotFound = errors.New("book not found")

// 詳細型（errors.As 用）
type BookNotFoundError struct {
    ID int
}

func (e *BookNotFoundError) Error() string {
    return fmt.Sprintf("book %d not found", e.ID)
}

func (e *BookNotFoundError) Unwrap() error {
    return ErrBookNotFound
}

// 生成
return &BookNotFoundError{ID: id}

// 判定
if errors.Is(err, ErrBookNotFound) { /* ... */ }
var nfe *BookNotFoundError
if errors.As(err, &nfe) { /* nfe.ID が使える */ }
```

**やらない:**
```go
// 文字列比較でエラーを判定しない
if err.Error() == "book not found" { /* 壊れやすい */ }

// fmt.Errorf だけで文字列として表現すると、ID 付きのカスタム型を返せない
return fmt.Errorf("book not found: id=%d", id) // errors.As で *BookNotFoundError は取り出せない
```

---

## 5. 層間でエラーを変換する

下位層（repository）のエラーをそのまま上位層（handler）に返さない。各層は自層のエラー型に変換してから返す。handler は repository の実装詳細を知らなくてよい。

**やる:**
```go
// repository は自層のエラーを返す
func (r *InMemoryRepository) FindByID(id int) (Book, error) {
    book, ok := r.data[id]
    if !ok {
        return Book{}, &BookNotFoundError{ID: id}
    }
    return book, nil
}

// handler は repository のエラー型だけを判断材料にする
func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
    book, err := h.repo.FindByID(id)
    if errors.Is(err, ErrBookNotFound) {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(book)
}
```

**やらない:**
```go
// repository の生エラーを handler まで素通しにしない
func (r *InMemoryRepository) FindByID(id int) (Book, error) {
    book, ok := r.data[id]
    if !ok {
        return Book{}, fmt.Errorf("key %d missing in map", id) // 実装詳細が漏れる
    }
    return book, nil
}
```
