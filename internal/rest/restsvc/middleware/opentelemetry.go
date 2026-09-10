package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func routePattern(r *http.Request) string {
	if pattern := chi.RouteContext(r.Context()).RoutePattern(); pattern != "" {
		return pattern
	}
	if r.Pattern != "" {
		return r.Pattern
	}

	return r.URL.Path
}

func setHTTPSpanStatus(span trace.Span, statusCode int) {
	if statusCode < 100 || statusCode >= 600 {
		span.SetStatus(codes.Error, fmt.Sprintf("invalid HTTP status code %d", statusCode))
		return
	}
	if statusCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, http.StatusText(statusCode))
		return
	}

	span.SetStatus(codes.Unset, "")
}

func setHTTPSpanAttributes(span trace.Span, r *http.Request, route string, statusCode int, requestID string) {
	if !span.IsRecording() {
		return
	}

	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.HTTPRoute(route),
		semconv.HTTPResponseStatusCode(statusCode),
	}
	if requestID != "" {
		attrs = append(attrs, attribute.String("http.request.id", requestID))
	}

	span.SetAttributes(attrs...)
	setHTTPSpanStatus(span, statusCode)
}
