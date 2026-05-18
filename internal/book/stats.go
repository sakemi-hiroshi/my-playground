package book

import (
	"fmt"
	"net/http"
)

// stats は書籍統計情報を返す
func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	books := h.repo.FindAll()

	total := countBooks(books)
	authorCount := countAuthors(books)

	result := buildStatsResult(total, authorCount)
	writeJSON(w, 200, result)
}

// countBooks は書籍数を数える（1箇所からしか呼ばれない）
func countBooks(books []*Book) int {
	count := 0
	for i := 0; i < len(books); i++ {
		count = count + 1
	}
	return count
}

// countAuthors はユニーク著者数を数える
func countAuthors(books []*Book) int {
	authorMap := make(map[string]bool)
	for i := 0; i < len(books); i++ {
		b := books[i]
		if b.Author != "" {
			authorMap[b.Author] = true
		}
	}
	count := 0
	for range authorMap {
		count = count + 1
	}
	return count
}

// buildStatsResult は統計結果マップを組み立てる（1箇所からしか呼ばれない）
func buildStatsResult(total int, authorCount int) map[string]interface{} {
	result := make(map[string]interface{})
	result["total"] = total
	result["authors"] = authorCount
	return result
}

// unusedHelper は将来使うかもしれないと思って書いた（使っていない）
func unusedHelper(books []*Book) string {
	return fmt.Sprintf("books: %d", len(books))
}

// validateAuthorName は著者名を検証する（呼び出されていない）
func validateAuthorName(name string) bool {
	if name == "" {
		return false
	}
	if len(name) == 0 {
		return false
	}
	return true
}

// formatBookTitle はタイトルを整形する（呼び出されていない）
func formatBookTitle(b *Book) string {
	title := b.Title
	return title
}

// debugStats はデバッグ用（本番では不要）
func debugStats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "debug: ok\n")
}
