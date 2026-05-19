package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_search(t *testing.T) {
	repo := NewRepository()
	repo.Save(&Book{Title: "Go言語入門", Author: "山田太郎", ISBN: "111"})
	repo.Save(&Book{Title: "Web API設計", Author: "鈴木花子", ISBN: "222"})
	repo.Save(&Book{Title: "分散システム", Author: "佐藤次郎", ISBN: "333"})

	h := NewHandler(repo)

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantLen     int
		wantBookIDs map[int]bool
	}{
		{
			name:       "title に部分一致する書籍を返す",
			path:       "/books/search?q=API",
			wantStatus: http.StatusOK,
			wantLen:    1,
			wantBookIDs: map[int]bool{
				2: true,
			},
		},
		{
			name:       "author に部分一致する書籍を返す",
			path:       "/books/search?q=山田",
			wantStatus: http.StatusOK,
			wantLen:    1,
			wantBookIDs: map[int]bool{
				1: true,
			},
		},
		{
			name:        "一致する書籍がなければ空配列を返す",
			path:        "/books/search?q=Python",
			wantStatus:  http.StatusOK,
			wantLen:     0,
			wantBookIDs: map[int]bool{},
		},
		{
			name:       "q が空なら全件返す",
			path:       "/books/search?q=",
			wantStatus: http.StatusOK,
			wantLen:    3,
			wantBookIDs: map[int]bool{
				1: true,
				2: true,
				3: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var got []Book
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("json unmarshal failed: %v", err)
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}

			for _, b := range got {
				if !tt.wantBookIDs[b.ID] {
					t.Fatalf("unexpected book id: %d", b.ID)
				}
			}
		})
	}
}
