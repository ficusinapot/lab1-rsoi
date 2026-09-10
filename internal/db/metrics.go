package db

import (
	dbmetrics "github.com/ficusinapot/ds/internal/db/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics(reg prometheus.Registerer) {
	wrapped := prometheus.WrapRegistererWithPrefix("db_", reg)
	dbmetrics.InitMetrics(wrapped)
}
