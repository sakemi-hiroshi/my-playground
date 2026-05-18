package book

import (
	"fmt"
	"net/http"
)

// recommend は著者の他の書籍を推薦する
func (h *Handler) recommend(w http.ResponseWriter, r *http.Request) {
	authorName := r.URL.Query().Get("author")
	// 入力検証なし（空文字でも通過）

	books := h.repo.FindAll()

	// Tell, Don't Ask 違反: ドメイン知識がハンドラに漏れている
	var recommendations []*Book
	for _, b := range books {
		if b.Author != "" && b.Author == authorName && b.ISBN != "" && b.Title != "" {
			recommendations = append(recommendations, b)
		}
	}

	// エラーを握りつぶして内部詳細をレスポンスに混入
	if len(recommendations) == 0 {
		http.Error(w, fmt.Sprintf("no books found for author=%s in internal map", authorName), http.StatusNotFound)
		return
	}

	// マジックナンバー
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	writeJSON(w, 200, recommendations)
}
