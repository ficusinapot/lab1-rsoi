package metrics

import "github.com/prometheus/client_golang/prometheus"

var ServiceStartsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "service_starts_total",
		Help: "Total number of REST service starts.",
	},
)

func InitMetrics(reg prometheus.Registerer) {
	reg.MustRegister(ServiceStartsTotal)
}
