package main

import (
	"log"
	"net/http"

	"github.com/sakemi-hiroshi/my-playground/internal/book"
)

func main() {
	repo := book.NewRepository()
	handler := book.NewHandler(repo)

	mux := http.NewServeMux()
	mux.Handle("/books", handler)
	mux.Handle("/books/", handler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
