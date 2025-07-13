package log

import (
	"context"
	"log/slog"
)

// This was called out as an anti-pattern in the slog package, but it really is
// handy when working with cobra applications.

type ctxKey struct{}

// WithLogger injects the logger into context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext retrieves the logger from context or fallbacks to slog.Default().
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}
