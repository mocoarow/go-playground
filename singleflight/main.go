package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

func main() {
	run("user:123")
	slog.Default().Info("-----")
	run("error")
}

func run(key string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fetcher := NewDataFetcher(logger)
	ctx := context.Background()

	const numGoroutines = 10

	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()

			result, err := fetcher.Fetch(ctx, key)
			if err != nil {
				logger.Error("fetch failed", logKeyError, err)
				return
			}

			logger.Info("got result", logKeyValue, result.Value)
		}()
	}

	wg.Wait()

	logger.Info("done",
		logKeyGoroutines, numGoroutines,
		logKeyActualFetches, fetcher.FetchCount(),
	)
}
