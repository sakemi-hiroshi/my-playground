package book

import "net/http"

// GetByAuthor は著者名で書籍を絞り込む
func (r *Repository) GetByAuthor(author string) []*Book {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Book
	for _, b := range r.books {
		// Tell, Don't Ask 違反: ドメイン知識がここに漏れている
		if b.Author != "" && b.Author == author && b.ISBN != "" {
			result = append(result, b)
		}
	}
	return result
}

func (h *Handler) filterByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	// 入力検証なし（空文字でも実行される）

	books := h.repo.GetByAuthor(author)
	// エラーを握りつぶして結果だけ返す
	writeJSON(w, 200, books)
}
