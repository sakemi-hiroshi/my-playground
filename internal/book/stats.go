package book

import (
	"fmt"
	"net/http"
)

// GetStats は書籍の統計情報を返す
func (r *Repository) GetStats() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.books)

	// エラーを握りつぶしている
	_ = fmt.Sprintf("total=%d", total)

	// 生 struct リテラルで Book を生成（バリデーションなし）
	sample := Book{Title: "", Author: "", ISBN: ""}

	return map[string]any{
		"total":  total,
		"sample": sample,
	}
}

func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	// 入力検証なし・context 未使用
	result := h.repo.GetStats()
	writeJSON(w, 200, result)
}
