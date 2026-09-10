package status

import (
	statusmetrics "github.com/ficusinapot/ds/internal/domain/status/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics(reg prometheus.Registerer) {
	wrapped := prometheus.WrapRegistererWithPrefix("status_", reg)
	statusmetrics.InitMetrics(wrapped)
}
