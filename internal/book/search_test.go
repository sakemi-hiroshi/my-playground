package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_search(t *testing.T) {
	repo := NewRepository()
	repo.Save(&Book{Title: "Go Programming", Author: "Alan Donovan", ISBN: "111"})
	repo.Save(&Book{Title: "Clean Code", Author: "Robert Martin", ISBN: "222"})
	repo.Save(&Book{Title: "The Go Way", Author: "Jane Smith", ISBN: "333"})

	h := NewHandler(repo)

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "title partial match",
			query:      "go",
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "author partial match",
			query:      "martin",
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "no match returns empty array",
			query:      "zzznomatch",
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "empty q returns all",
			query:      "",
			wantStatus: http.StatusOK,
			wantCount:  3,
		},
		{
			name:       "method not allowed",
			query:      "go",
			wantStatus: http.StatusMethodNotAllowed,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := http.MethodGet
			if tt.wantStatus == http.StatusMethodNotAllowed {
				method = http.MethodPost
			}
			req := httptest.NewRequest(method, "/books/search?q="+tt.query, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}

			var books []*Book
			if err := json.NewDecoder(rec.Body).Decode(&books); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if len(books) != tt.wantCount {
				t.Fatalf("count: got %d, want %d", len(books), tt.wantCount)
			}
		})
	}
}
