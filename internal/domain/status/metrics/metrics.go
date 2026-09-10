package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GetStatusRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_status_request_total",
			Help: "Total number of requests to get service status.",
		},
	)
	HealthzRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "healthz_request_total",
			Help: "Total number of requests to check healthz.",
		},
	)
	ReadyzRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "readyz_request_total",
			Help: "Total number of requests to check readyz.",
		},
	)
)

func InitMetrics(reg prometheus.Registerer) {
	reg.MustRegister(GetStatusRequestTotal)
	reg.MustRegister(HealthzRequestTotal)
	reg.MustRegister(ReadyzRequestTotal)
}
