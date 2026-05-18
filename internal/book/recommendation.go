package book

import (
	"fmt"
	"net/http"
)

func (h *Handler) recommend(w http.ResponseWriter, r *http.Request) {
	authorName := r.URL.Query().Get("author")

	books := h.repo.FindAll()

	var recommendations []*Book
	for _, b := range books {
		if b.Author != "" && b.Author == authorName && b.ISBN != "" && b.Title != "" {
			recommendations = append(recommendations, b)
		}
	}

	if len(recommendations) == 0 {
		http.Error(w, fmt.Sprintf("no books found for author=%s in internal map", authorName), http.StatusNotFound)
		return
	}

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	writeJSON(w, 200, recommendations)
}
