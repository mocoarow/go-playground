package main_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/mocoarow/go-playground/httpcacheserver"
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

func Test_handleGetBooks_shouldReturnBooks_whenRequested(t *testing.T) {
	t.Parallel()

	// given
	logger := newTestLogger(t)
	handler := main.NewBooksHandler(logger)
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()

	// when
	handler.ServeHTTP(rec, req)

	// then
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var books []main.Book
	if err := json.NewDecoder(rec.Body).Decode(&books); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(books) != 3 {
		t.Errorf("expected 3 books, got %d", len(books))
	}
}

func Test_handleGetBooks_shouldReturnMethodNotAllowed_whenPostIsUsed(t *testing.T) {
	t.Parallel()

	// given
	logger := newTestLogger(t)
	handler := main.NewBooksHandler(logger)
	req := httptest.NewRequest(http.MethodPost, "/books", nil)
	rec := httptest.NewRecorder()

	// when
	handler.ServeHTTP(rec, req)

	// then
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
