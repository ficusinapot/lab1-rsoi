package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func recoverPanic(w http.ResponseWriter, r *http.Request, logger *slog.Logger) bool {
	recovered := recover()
	if recovered == nil {
		return false
	}

	timestamp := time.Now()
	stack := string(debug.Stack())

	logger.ErrorContext(
		r.Context(),
		"recovering from panic",
		"error", recovered,
		"stacktrace", stack,
	)

	span := trace.SpanFromContext(r.Context())
	if span.IsRecording() {
		message := "recovering from panic"
		span.AddEvent(
			"log",
			trace.WithTimestamp(timestamp),
			trace.WithAttributes(
				attribute.String("log.time", timestamp.Format(time.RFC3339Nano)),
				attribute.String("log.level", slog.LevelError.String()),
				attribute.String("log.message", message),
				attribute.String("exception.stacktrace", stack),
				attribute.String("exception.message", fmt.Sprint(recovered)),
				attribute.Bool("panic", true),
			),
		)
		span.SetStatus(codes.Error, message)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	return true
}
