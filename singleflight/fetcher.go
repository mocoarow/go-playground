// Package main demonstrates singleflight usage to suppress duplicate concurrent calls.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	logKeyKey           = "key"
	logKeyShared        = "shared"
	logKeyValue         = "value"
	logKeyError         = "error"
	logKeyGoroutines    = "goroutines"
	logKeyActualFetches = "actual_fetches"
)

// FetchResult represents data fetched from an external source.
type FetchResult struct {
	Value     string
	FetchedAt time.Time
}

// DataFetcher fetches data using singleflight to deduplicate concurrent requests.
type DataFetcher struct {
	group      singleflight.Group
	fetchCount atomic.Int64
	logger     *slog.Logger
}

// NewDataFetcher creates a new DataFetcher.
func NewDataFetcher(logger *slog.Logger) *DataFetcher {
	return &DataFetcher{
		group:      singleflight.Group{},
		fetchCount: atomic.Int64{},
		logger:     logger,
	}
}

// Fetch retrieves data for the given key. Concurrent calls with the same key
// are deduplicated so that only one actual fetch executes.
func (f *DataFetcher) Fetch(ctx context.Context, key string) (FetchResult, error) {
	v, err, shared := f.group.Do(key, func() (any, error) {
		return f.doFetch(ctx, key)
	})
	if err != nil {
		return FetchResult{}, fmt.Errorf("fetch %s: %w", key, err)
	}

	result, ok := v.(FetchResult)
	if !ok {
		return FetchResult{}, fmt.Errorf("fetch %s: unexpected result type", key)
	}

	f.logger.Info("fetch completed", logKeyKey, key, logKeyShared, shared)

	return result, nil
}

// FetchCount returns the number of actual fetch operations performed.
func (f *DataFetcher) FetchCount() int64 {
	return f.fetchCount.Load()
}

func (f *DataFetcher) doFetch(_ context.Context, key string) (FetchResult, error) {
	f.fetchCount.Add(1)
	f.logger.Info("performing actual fetch", logKeyKey, key)

	// Simulate slow external call.
	// time.Sleep(100 * time.Millisecond)
	time.Sleep(2 * time.Second)

	if key == "error" {
		return FetchResult{}, fmt.Errorf("simulated fetch error for key %s", key)
	}

	return FetchResult{
		Value:     fmt.Sprintf("data-for-%s", key),
		FetchedAt: time.Now(),
	}, nil
}
