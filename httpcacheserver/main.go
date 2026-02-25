package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second

	logKeyAddr = "addr"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := NewBooksHandler(logger)

	server := &http.Server{ //nolint:exhaustruct
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		logger.Info("server starting", logKeyAddr, server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down server")
	case err := <-errCh:
		logger.Error("listen and serve", logKeyError, err)
		os.Exit(1)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown", logKeyError, err)
	}

	logger.Info("server stopped")
}
