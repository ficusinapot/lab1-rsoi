package middleware

import "time"

type HTTPMetrics interface {
	ObserveHTTPRequest(method, route string, statusCode int, duration time.Duration)
}
