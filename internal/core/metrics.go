package core

import (
	personmetrics "github.com/ficusinapot/ds/internal/core/persons/metrics"
	"github.com/ficusinapot/ds/internal/core/status"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics(reg prometheus.Registerer) {
	personsRegisterer := prometheus.WrapRegistererWithPrefix("core_", reg)
	personmetrics.InitMetrics(personsRegisterer)
	status.InitMetrics(reg)
}
