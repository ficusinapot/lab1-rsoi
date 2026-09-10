package restsvc

import (
	"log/slog"
	"net/http"

	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/persons"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/status"
)

func newHTTPServer(
	appConfig AppInfo,
	serviceConfig ServiceConfig,
	personHandler *persons.Handler,
	statusHandler *status.Handler,
	logger *slog.Logger,
	metricsRegistry *metrics.Registry,
) *http.Server {
	return &http.Server{ //nolint:exhaustruct_v5
		Addr:              serviceConfig.Addr,
		Handler:           NewRouter(appConfig, serviceConfig.OpenAPI, personHandler, statusHandler, logger, metricsRegistry),
		ReadTimeout:       serviceConfig.Timeout.Read,
		ReadHeaderTimeout: serviceConfig.Timeout.Read,
		WriteTimeout:      serviceConfig.Timeout.Write,
		IdleTimeout:       serviceConfig.Timeout.Idle,
	}
}

func newMetricsServer(metricsConfig MetricsConfig, metricsRegistry *metrics.Registry) *http.Server {
	if !metricsConfig.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", metricsRegistry.Handler())

	return &http.Server{ //nolint:exhaustruct_v5
		Addr:              metricsConfig.Addr,
		Handler:           mux,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
}
