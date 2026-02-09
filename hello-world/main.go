// Package main provides a simple hello world application.
package main

import "log/slog"

func main() {
	slog.Default().Info("Hello World")
}
