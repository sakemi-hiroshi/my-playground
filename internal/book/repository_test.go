package book

import "testing"

func TestRepository_FindByAuthor(t *testing.T) {
	tests := []struct {
		name       string
		author     string
		wantCount  int
		wantAuthor string
		wantTitles map[string]bool
	}{
		{
			name:       "正常系: 完全一致した著者の書籍だけ返す",
			author:     "夏目漱石",
			wantCount:  2,
			wantAuthor: "夏目漱石",
			wantTitles: map[string]bool{
				"こころ":   true,
				"坊っちゃん": true,
			},
		},
		{
			name:      "境界値: 空文字列なら空文字列の著者だけ返す",
			author:    "",
			wantCount: 1,
			wantTitles: map[string]bool{
				"無名の本": true,
			},
		},
		{
			name:      "境界値: 一致なしなら空配列を返す",
			author:    "太宰治",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository()
			repo.Save(&Book{Title: "こころ", Author: "夏目漱石", ISBN: "9784003101018"})
			repo.Save(&Book{Title: "坊っちゃん", Author: "夏目漱石", ISBN: "9784101010137"})
			repo.Save(&Book{Title: "舞姫", Author: "森鴎外", ISBN: "9784101020013"})
			repo.Save(&Book{Title: "無名の本", Author: "", ISBN: "9784000000000"})

			got := repo.FindByAuthor(tt.author)

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
