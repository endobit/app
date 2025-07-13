package log

import (
	"context"
	"fmt"
	"log/slog"
)

// Legacy is a structured logger that implements common logger interfaces.
type Legacy struct {
	logger *slog.Logger
	level  slog.Level
	filter func(string) (string, []slog.Attr)
}

func WithLevel(level slog.Level) func(*Legacy) {
	return func(l *Legacy) {
		l.level = level
	}
}

func WithFilter(filter func(string) (string, []slog.Attr)) func(*Legacy) {
	return func(l *Legacy) {
		l.filter = filter
	}
}

// NewLegacy creates a new Legacy logger that write to the provided slog.Logger.
func NewLegacy(logger *slog.Logger, opts ...func(*Legacy)) Legacy {
	l := Legacy{
		logger: logger,
		filter: func(s string) (string, []slog.Attr) { return s, nil },
		level:  slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(&l)
	}

	return l
}

// Printf formats according to a format specifier and writes to the logger.
func (l Legacy) Printf(format string, v ...any) {
	if l.logger == nil {
		return
	}

	fmt.Printf("Printf %v\n", l.level)

	msg, attrs := l.filter(fmt.Sprintf(format, v...))

	l.logger.LogAttrs(context.Background(), l.level, msg, attrs...)
}

// Println formats using the default formats for its operands and writes to the
// logger.
func (l Legacy) Println(v ...any) {
	if l.logger == nil {
		return
	}

	msg, attrs := l.filter(fmt.Sprint(v...))

	l.logger.LogAttrs(context.Background(), l.level, msg, attrs...)
}
