package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ficusinapot/ds/internal/config"

	"github.com/joomcode/errorx"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg config.LoggingConfig) (*slog.Logger, func(), error) {
	handlers := make([]slog.Handler, 0, 6)
	closers := make([]io.Closer, 0, 4)

	if cfg.Stdout.Level != "" {
		handlers = append(
			handlers,
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLevel(cfg.Stdout.Level)}), //nolint:exhaustruct_v5 // Handler defaults are enough.
		)
	}

	for _, file := range cfg.Files {
		if err := os.MkdirAll(filepath.Dir(file.Path), 0o750); err != nil {
			return nil, nil, errorx.InitializationFailed.Wrap(err, "create log directory")
		}

		writer := newRotator(file)

		handlers = append(handlers, levelHandler{
			level:   parseLevel(file.Level),
			handler: slog.NewJSONHandler(writer, nil),
		})
		closers = append(closers, writer)
	}

	if len(handlers) == 0 {
		handlers = append(handlers, slog.NewJSONHandler(io.Discard, nil))
	}

	logger := slog.New(multiHandler{
		level:    parseLevel(cfg.Level),
		handlers: handlers,
	})

	closeFunc := func() {
		for _, closer := range closers {
			_ = closer.Close()
		}
	}

	return logger, closeFunc, nil
}

func newRotator(cfg config.LoggingFileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{ //nolint:exhaustruct_v5 // Optional lumberjack fields keep library defaults.
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxFileSize.Megabytes(),
		MaxBackups: cfg.MaxFilesCount,
		MaxAge:     cfg.MaxFileAgeInDays,
	}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type multiHandler struct {
	level    slog.Level
	handlers []slog.Handler
}

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < h.level {
		return false
	}

	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (h multiHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level < h.level {
		return nil
	}

	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}

		if err := handler.Handle(ctx, record); err != nil {
			return errorx.InternalError.Wrap(err, "write log record")
		}
	}

	return nil
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithAttrs(attrs))
	}

	return multiHandler{
		level:    h.level,
		handlers: handlers,
	}
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithGroup(name))
	}

	return multiHandler{
		level:    h.level,
		handlers: handlers,
	}
}

type levelHandler struct {
	level   slog.Level
	handler slog.Handler
}

func (h levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level == h.level && h.handler.Enabled(ctx, level)
}

func (h levelHandler) Handle(ctx context.Context, record slog.Record) error {
	if err := h.handler.Handle(ctx, record); err != nil {
		return errorx.InternalError.Wrap(err, "write level log record")
	}

	return nil
}

func (h levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return levelHandler{
		level:   h.level,
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h levelHandler) WithGroup(name string) slog.Handler {
	return levelHandler{
		level:   h.level,
		handler: h.handler.WithGroup(name),
	}
}
