package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/books")
	switch {
	case path == "" || path == "/":
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		idStr := strings.TrimPrefix(path, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.get(w, r, id)
		case http.MethodPut:
			h.update(w, r, id)
		case http.MethodDelete:
			h.delete(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author != "" {
		writeJSON(w, http.StatusOK, h.repo.FindByAuthor(author))
		return
	}
	writeJSON(w, http.StatusOK, h.repo.FindAll())
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, h.repo.Save(&b))
}

func (h *Handler) get(w http.ResponseWriter, _ *http.Request, id int) {
	b, err := h.repo.FindByID(id)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, id int) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b.ID = id
	if err := h.repo.Update(&b); errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, &b)
}

func (h *Handler) delete(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.repo.Delete(id); errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
