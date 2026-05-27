package book

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateAcceptsPublishedAt(t *testing.T) {
	repo := NewRepository()
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title":"t","author":"a","isbn":"i","published_at":"2026-01-01"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var got Book
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.PublishedAt != "2026-01-01" {
		t.Fatalf("published_at = %q, want %q", got.PublishedAt, "2026-01-01")
	}
}

func TestHandler_UpdateAcceptsPublishedAt(t *testing.T) {
	repo := NewRepository()
	h := NewHandler(repo)

	saved := repo.Save(&Book{Title: "before", Author: "a", ISBN: "i", PublishedAt: "2025-01-01"})
	req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewBufferString(`{"title":"after","author":"a2","isbn":"i2","published_at":"2026-02-02"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	got, err := repo.FindByID(saved.ID)
	if err != nil {
		t.Fatalf("find book: %v", err)
	}
	if got.PublishedAt != "2026-02-02" {
		t.Fatalf("published_at = %q, want %q", got.PublishedAt, "2026-02-02")
	}
}
