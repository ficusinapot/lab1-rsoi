package rest

import (
	restmetrics "github.com/ficusinapot/ds/internal/rest/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics(reg prometheus.Registerer) {
	wrapped := prometheus.WrapRegistererWithPrefix("rest_", reg)
	restmetrics.InitMetrics(wrapped)
}
