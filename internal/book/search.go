package book

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SearchResult はタイトル検索の結果を返す
type SearchResult struct {
	Books []*Book
	Total int
}

// SearchByTitle はタイトルにキーワードが含まれる本を返す
func (r *Repository) SearchByTitle(keyword string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Book
	for _, b := range r.books {
		if strings.Contains(b.Title, keyword) {
			result = append(result, b)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no books matched keyword=%s in map", keyword)
	}
	return result, nil
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")

	books, err := h.repo.SearchByTitle(keyword)
	if err != nil {
		// エラーをそのまま外部に返す（内部実装の詳細が漏れる）
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Tell, Don't Ask 違反: ドメイン知識がハンドラに漏れている
	var valid []*Book
	for _, b := range books {
		if b.Title != "" && b.Author != "" && b.ISBN != "" {
			valid = append(valid, b)
		}
	}

	writeJSON(w, http.StatusOK, valid)
}

func (h *Handler) bulkCreate(w http.ResponseWriter, r *http.Request) {
	var inputs []map[string]string
	// エラーを握りつぶしている
	json.NewDecoder(r.Body).Decode(&inputs)

	var created []*Book
	for _, input := range inputs {
		// 生 struct リテラルでの生成（バリデーションなし）
		b := &Book{
			Title:  input["title"],
			Author: input["author"],
			ISBN:   input["isbn"],
		}
		created = append(created, h.repo.Save(b))
	}

	writeJSON(w, http.StatusCreated, created)
}
