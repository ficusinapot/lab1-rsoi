package metrics

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/joomcode/errorx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Registry struct {
	namespace       string
	registry        *prometheus.Registry
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

const (
	restSubsystem = "rest_api"
	dbSubsystem   = "db"
)

func NewRegistry(appName, appVersion string) *Registry {
	registry := prometheus.NewRegistry()
	namespace := metricNameComponent(appName)

	appInfo := prometheus.NewGauge(prometheus.GaugeOpts{ //nolint:exhaustruct_v5
		Namespace: namespace,
		Name:      "app_info",
		Help:      "Application information.",
		ConstLabels: prometheus.Labels{
			"app":     appName,
			"version": appVersion,
		},
	})
	appInfo.Set(1)

	requestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{ //nolint:exhaustruct_v5
		Namespace: namespace,
		Subsystem: restSubsystem,
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "route", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{ //nolint:exhaustruct_v5
		Namespace: namespace,
		Subsystem: restSubsystem,
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route", "status"})

	registry.MustRegister(appInfo, requestsTotal, requestDuration)

	return &Registry{
		namespace:       namespace,
		registry:        registry,
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
	}
}

func (r *Registry) RegisterDatabase(db *sql.DB) error {
	collectors := []prometheus.Collector{
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{ //nolint:exhaustruct_v5
			Namespace: r.namespace,
			Subsystem: dbSubsystem,
			Name:      "connections_open",
			Help:      "Number of established database connections.",
		}, func() float64 {
			return float64(db.Stats().OpenConnections)
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{ //nolint:exhaustruct_v5
			Namespace: r.namespace,
			Subsystem: dbSubsystem,
			Name:      "connections_in_use",
			Help:      "Number of database connections currently in use.",
		}, func() float64 {
			return float64(db.Stats().InUse)
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{ //nolint:exhaustruct_v5
			Namespace: r.namespace,
			Subsystem: dbSubsystem,
			Name:      "connections_idle",
			Help:      "Number of idle database connections.",
		}, func() float64 {
			return float64(db.Stats().Idle)
		}),
	}

	for _, collector := range collectors {
		if err := r.registry.Register(collector); err != nil {
			return errorx.InternalError.Wrap(err, "register database metric collector")
		}
	}

	return nil
}

func (r *Registry) Handler() http.Handler {
	return promhttp.HandlerFor(r.registry, promhttp.HandlerOpts{}) //nolint:exhaustruct_v5
}

func (r *Registry) Registerer() prometheus.Registerer {
	return r.registry
}

func (r *Registry) ObserveHTTPRequest(method, route string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)

	r.requestsTotal.WithLabelValues(method, route, status).Inc()
	r.requestDuration.WithLabelValues(method, route, status).Observe(duration.Seconds())
}

var invalidMetricNameChar = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func metricNameComponent(value string) string {
	value = strings.Trim(invalidMetricNameChar.ReplaceAllString(value, "_"), "_")
	if value == "" {
		return "app"
	}
	first := rune(value[0])
	if first != '_' && !unicode.IsLetter(first) {
		return "app_" + value
	}

	return value
}
