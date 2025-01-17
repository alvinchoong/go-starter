package slogr

import (
	"context"
	"log/slog"
)

// ctxKey is an unexported type to prevent context key collisions
type ctxKey struct{}

// loggerCtxKey is the context key used to store and retrieve the *slog.Logger
var loggerCtxKey = ctxKey{}

// FromContext retrieves the *slog.Logger from the provided context
// If no logger is found, it returns the default logger provided by slog
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	if logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger); ok && logger != nil {
		return logger
	}

	return slog.Default()
}

// ToContext returns a new context with the provided *slog.Logger embedded
// If the provided logger is nil, the default logger is used instead
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		logger = slog.Default()
	}
	return context.WithValue(ctx, loggerCtxKey, logger)
}
