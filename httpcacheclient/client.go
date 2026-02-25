// Package main provides an HTTP client that fetches books from the httpcacheserver.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	logKeyError = "error"
)

// Book represents a book resource returned by the server.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// FetchBooks fetches the list of books from the given base URL.
func FetchBooks(ctx context.Context, client *http.Client, baseURL string, logger *slog.Logger) ([]Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/books", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("close response body", logKeyError, err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var books []Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return books, nil
}
