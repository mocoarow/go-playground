package main_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	main "github.com/mocoarow/go-playground/httpcacheclient"
)

func newTestLogger(t *testing.T) *slog.Logger {
	t.Helper()

	return slog.New(slog.NewTextHandler(&testWriter{t: t}, nil))
}

type testWriter struct {
	t *testing.T
}

func (w *testWriter) Write(p []byte) (int, error) {
	w.t.Helper()
	w.t.Log(string(p))

	return len(p), nil
}

func Test_FetchBooks_shouldReturnBooks_whenServerReturnsValidResponse(t *testing.T) {
	t.Parallel()

	// given
	books := []main.Book{
		{ID: 1, Title: "Test Book", Author: "Test Author"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(books); err != nil {
			t.Errorf("encode books: %v", err)
		}
	}))
	defer server.Close()

	logger := newTestLogger(t)
	client := server.Client()
	ctx := context.Background()

	// when
	result, err := main.FetchBooks(ctx, client, server.URL, logger)

	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 book, got %d", len(result))
	}

	if result[0].Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got %s", result[0].Title)
	}
}

func Test_FetchBooks_shouldReturnError_whenServerReturns500(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	logger := newTestLogger(t)
	client := server.Client()
	ctx := context.Background()

	// when
	_, err := main.FetchBooks(ctx, client, server.URL, logger)

	// then
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func Test_FetchBooks_shouldCallServerTwice_whenCalledTwice(t *testing.T) {
	t.Parallel()

	// given
	var callCount atomic.Int64

	books := []main.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(books); err != nil {
			t.Errorf("encode books: %v", err)
		}
	}))
	defer server.Close()

	logger := newTestLogger(t)
	client := server.Client()
	ctx := context.Background()

	// when
	_, err1 := main.FetchBooks(ctx, client, server.URL, logger)
	_, err2 := main.FetchBooks(ctx, client, server.URL, logger)

	// then
	if err1 != nil {
		t.Fatalf("first call: unexpected error: %v", err1)
	}

	if err2 != nil {
		t.Fatalf("second call: unexpected error: %v", err2)
	}

	if callCount.Load() != 2 {
		t.Errorf("expected 2 calls to server, got %d", callCount.Load())
	}
}
