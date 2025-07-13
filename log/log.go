// Package log is a helper for the slog package for configuring and connecting
// to command line application.
package log

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"endobit.io/clog"
	"endobit.io/clog/ansi"
)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Options holds the results of the command line flags.
type Options struct {
	Writer   io.Writer
	Filename string
	level    string
	file     string
	doJSON   bool
}

type flagSet interface {
	StringVar(*string, string, string, string)
	BoolVar(*bool, string, bool, string)
}

// NewOptions add flags to the command line. Works for cobra and the standard
// library.
func NewOptions(f flagSet) *Options {
	var opt Options

	f.StringVar(&opt.level, "log-level", LevelInfo.String(), "set the log level")
	f.StringVar(&opt.file, "log-file", "", "log to a file (implies json)")
	f.BoolVar(&opt.doJSON, "log-json", false, "log in json format")

	return &opt
}

func New(opts *Options) (*slog.Logger, error) {
	if opts == nil {
		return nil, errors.New("options for logger cannot be nil")
	}

	if opts.Writer == nil && opts.file != "" {
		return nil, errors.New("writer cannot be nil if file is set (caller manages the file)")
	}

	if opts.Writer == nil {
		opts.Writer = os.Stderr
	}

	level, err := parseLevel(opts.level)
	if err != nil {
		return nil, err
	}

	hopts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	switch {
	case opts.file == "" && opts.doJSON:
		w := niceJSON{Writer: opts.Writer}
		handler = slog.NewJSONHandler(&w, hopts)
	case opts.file != "":
		handler = slog.NewJSONHandler(opts.Writer, hopts)
	default:
		hopts := clog.HandlerOptions{Level: level}
		handler = hopts.NewHandler(opts.Writer, clog.WithLevel(LevelTrace, "TRC", ansi.Cyan))
	}

	return slog.New(handler), nil
}

func Format(format, key string, value any) slog.Attr {
	return slog.String(key, fmt.Sprintf(format, value))
}

func parseLevel(level string) (slog.Level, error) {
	var l slog.Level

	level = strings.ToLower(level)

	if err := l.UnmarshalText([]byte(level)); err != nil {
		switch level {
		case "trace":
			l = LevelTrace
		default:
			return l, err
		}
	}

	return l, nil
}
