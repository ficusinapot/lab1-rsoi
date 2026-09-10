package domain

import (
	personmetrics "github.com/ficusinapot/ds/internal/domain/persons/metrics"
	"github.com/ficusinapot/ds/internal/domain/status"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics(reg prometheus.Registerer) {
	personsRegisterer := prometheus.WrapRegistererWithPrefix("domain_", reg)
	personmetrics.InitMetrics(personsRegisterer)
	status.InitMetrics(reg)
}
