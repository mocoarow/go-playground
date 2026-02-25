// Package main provides an HTTP server with a books API.
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	logKeyCount = "count"
	logKeyError = "error"

	bookID1 = 1
	bookID2 = 2
	bookID3 = 3
)

// Book represents a book resource.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func newBooks() []Book {
	return []Book{
		{ID: bookID1, Title: "The Go Programming Language", Author: "Alan A.A. Donovan"},
		{ID: bookID2, Title: "Learning Go", Author: "Jon Bodner"},
		{ID: bookID3, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
	}
}

// NewBooksHandler creates an HTTP handler that serves book resources.
func NewBooksHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", handleGetBooks(logger))

	return mux
}

func handleGetBooks(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		books := newBooks()

		data, err := json.Marshal(books)
		if err != nil {
			logger.Error("encode books", logKeyError, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(data); err != nil {
			logger.Error("write response", logKeyError, err)

			return
		}

		logger.Info("served books", logKeyCount, len(books))
	}
}
