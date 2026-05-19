package book

import "net/http"

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	var books []*Book
	if keyword == "" {
		books = h.repo.FindAll()
	} else {
		books = h.repo.FindByKeyword(keyword)
	}
	writeJSON(w, http.StatusOK, books)
}
