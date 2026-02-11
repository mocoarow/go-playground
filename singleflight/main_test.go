package main_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	main "github.com/mocoarow/go-playground/singleflight"
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

func Test_Fetch_shouldDeduplicateConcurrentCalls_whenSameKeyIsUsed(t *testing.T) {
	t.Parallel()

	// given
	logger := newTestLogger(t)
	fetcher := main.NewDataFetcher(logger)
	ctx := context.Background()

	const numGoroutines = 10

	// when
	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()

			result, err := fetcher.Fetch(ctx, "user:123")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Value != "data-for-user:123" {
				t.Errorf("unexpected value: %s", result.Value)
			}
		}()
	}

	wg.Wait()

	// then
	if fetcher.FetchCount() >= int64(numGoroutines) {
		t.Errorf("expected fewer fetches than goroutines, got %d fetches for %d goroutines",
			fetcher.FetchCount(), numGoroutines)
	}
}

func Test_Fetch_shouldFetchIndependently_whenDifferentKeysAreUsed(t *testing.T) {
	t.Parallel()

	// given
	logger := newTestLogger(t)
	fetcher := main.NewDataFetcher(logger)
	ctx := context.Background()

	// when
	result1, err := fetcher.Fetch(ctx, "user:1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result2, err := fetcher.Fetch(ctx, "user:2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// then
	if result1.Value != "data-for-user:1" {
		t.Errorf("unexpected value for key user:1: %s", result1.Value)
	}

	if result2.Value != "data-for-user:2" {
		t.Errorf("unexpected value for key user:2: %s", result2.Value)
	}

	if fetcher.FetchCount() != 2 {
		t.Errorf("expected 2 fetches, got %d", fetcher.FetchCount())
	}
}
