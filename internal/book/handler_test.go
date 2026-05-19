package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_list_query(t *testing.T) {
	repo := NewRepository()
	repo.Save(&Book{Title: "Go言語入門", Author: "田中太郎", ISBN: "111"})
	repo.Save(&Book{Title: "分散システム", Author: "Go研究会", ISBN: "222"})
	repo.Save(&Book{Title: "Rust完全ガイド", Author: "山田花子", ISBN: "333"})

	handler := NewHandler(repo)

	tests := []struct {
		name     string
		url      string
		wantCode int
		wantLen  int
		wantID   int
	}{
		{
			name:     "正常系: qでタイトルまたは著者の部分一致が返る",
			url:      "/books?q=研究会",
			wantCode: http.StatusOK,
			wantLen:  1,
			wantID:   2,
		},
		{
			name:     "正常系: 空マッチは空配列を返す",
			url:      "/books?q=Python",
			wantCode: http.StatusOK,
			wantLen:  0,
		},
		{
			name:     "正常系: qなしは全件を返す",
			url:      "/books",
			wantCode: http.StatusOK,
			wantLen:  3,
		},
		{
			name:     "正常系: q空文字は全件を返す",
			url:      "/books?q=",
			wantCode: http.StatusOK,
			wantLen:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)

			handler.ServeHTTP(w, r)

			if w.Code != tt.wantCode {
				t.Fatalf("status code = %d, want %d", w.Code, tt.wantCode)
			}

			var got []*Book
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if len(got) != tt.wantLen {
				t.Fatalf("response length = %d, want %d", len(got), tt.wantLen)
			}

			if tt.wantID != 0 {
				if got[0].ID != tt.wantID {
					t.Fatalf("first id = %d, want %d", got[0].ID, tt.wantID)
				}
			}
		})
	}
}
