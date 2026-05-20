package book

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("book not found")

type Repository struct {
	mu     sync.RWMutex
	books  map[int]*Book
	nextID int
}

func NewRepository() *Repository {
	return &Repository{
		books:  make(map[int]*Book),
		nextID: 1,
	}
}

func (r *Repository) FindAll() []*Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result
}

func (r *Repository) FindByAuthor(author string) []*Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Book, 0, len(r.books))
	for _, b := range r.books {
		if b.Author == author {
			result = append(result, b)
		}
	}
	return result
}

func (r *Repository) FindByID(id int) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (r *Repository) Save(b *Book) *Book {
	r.mu.Lock()
	defer r.mu.Unlock()
	b.ID = r.nextID
	r.nextID++
	r.books[b.ID] = b
	return b
}

func (r *Repository) Update(b *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[b.ID]; !ok {
		return ErrNotFound
	}
	r.books[b.ID] = b
	return nil
}

func (r *Repository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return ErrNotFound
	}
	delete(r.books, id)
	return nil
}
