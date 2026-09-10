package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const requestIDHeader = "X-Request-ID"

func Observability(logger *slog.Logger, metrics HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := requestID(r)
			w.Header().Set(requestIDHeader, requestID)

			startedAt := time.Now()
			defer func() {
				if !recoverPanic(w, r, logger) {
					return
				}

				writeObservability(logger, metrics, r, routePattern(r), http.StatusInternalServerError, time.Since(startedAt), requestID)
			}()

			httpMetrics := httpsnoop.CaptureMetrics(next, w, r)
			writeObservability(logger, metrics, r, routePattern(r), httpMetrics.Code, httpMetrics.Duration, requestID)
		})
	}
}

func requestID(r *http.Request) string {
	if requestID := r.Header.Get(requestIDHeader); requestID != "" {
		return requestID
	}

	return uuid.NewString()
}

func writeObservability(
	logger *slog.Logger,
	metrics HTTPMetrics,
	r *http.Request,
	route string,
	statusCode int,
	duration time.Duration,
	requestID string,
) {
	if metrics != nil {
		metrics.ObserveHTTPRequest(r.Method, route, statusCode, duration)
	}

	span := trace.SpanFromContext(r.Context())
	setHTTPSpanAttributes(span, r, route, statusCode, requestID)

	level := slog.LevelDebug
	message := "REQUEST " + r.Method + " " + r.URL.Path
	attrs := []slog.Attr{
		slog.String("uri", r.URL.RequestURI()),
		slog.String("method", r.Method),
		slog.String("route", route),
		slog.Int("status", statusCode),
		slog.Duration("duration", duration),
		slog.String("x-request-id", requestID),
	}
	if statusCode >= http.StatusInternalServerError {
		level = slog.LevelError
		message = "REQUEST_ERROR " + r.Method + " " + r.URL.Path
	}

	logger.LogAttrs(r.Context(), level, message, attrs...)
}
