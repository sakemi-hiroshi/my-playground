package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_list(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantCount  int
		wantAuthor string
		wantTitles map[string]bool
	}{
		{
			name:      "正常系: author指定なしなら全件返す",
			target:    "/books",
			wantCount: 3,
			wantTitles: map[string]bool{
				"こころ":   true,
				"坊っちゃん": true,
				"舞姫":    true,
			},
		},
		{
			name:       "正常系: author完全一致なら該当著者のみ返す",
			target:     "/books?author=夏目漱石",
			wantCount:  2,
			wantAuthor: "夏目漱石",
			wantTitles: map[string]bool{
				"こころ":   true,
				"坊っちゃん": true,
			},
		},
		{
			name:      "境界値: 一致なしなら空配列を返す",
			target:    "/books?author=太宰治",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository()
			repo.Save(&Book{Title: "こころ", Author: "夏目漱石", ISBN: "9784003101018"})
			repo.Save(&Book{Title: "坊っちゃん", Author: "夏目漱石", ISBN: "9784101010137"})
			repo.Save(&Book{Title: "舞姫", Author: "森鴎外", ISBN: "9784101020013"})

			handler := NewHandler(repo)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.target, nil)

			handler.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("status code = %d, want %d", w.Code, http.StatusOK)
			}

			var got []Book
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if len(got) != tt.wantCount {
				t.Fatalf("len(got) = %d, want %d", len(got), tt.wantCount)
			}

			for _, b := range got {
				if tt.wantAuthor != "" && b.Author != tt.wantAuthor {
					t.Errorf("book author = %q, want %q", b.Author, tt.wantAuthor)
				}
				if tt.wantTitles != nil && !tt.wantTitles[b.Title] {
					t.Errorf("unexpected title %q", b.Title)
				}
			}
		})
	}
}
