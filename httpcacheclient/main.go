package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

const (
	logKeyAttempt = "attempt"
	logKeyID      = "id"
	logKeyTitle   = "title"
	logKeyAuthor  = "author"

	fetchCount = 2
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := &http.Client{} //nolint:exhaustruct
	ctx := context.Background()

	baseURL := "http://localhost:8080"

	for i := range fetchCount {
		logger.Info("fetching books", logKeyAttempt, i+1)

		books, err := FetchBooks(ctx, client, baseURL, logger)
		if err != nil {
			logger.Error("fetch books", logKeyAttempt, i+1, logKeyError, err)
			os.Exit(1)
		}

		for _, book := range books {
			logger.Info("book", logKeyID, book.ID, logKeyTitle, book.Title, logKeyAuthor, book.Author)
		}
	}

	logger.Info("done")
}
