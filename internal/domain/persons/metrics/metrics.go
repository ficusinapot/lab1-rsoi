package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	CreateRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_request_total",
			Help: "Total number of requests to create Person.",
		},
	)
	CreateRequestFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_request_failed_total",
			Help: "Total number of requests to create Person that failed.",
		},
	)

	GetRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_request_total",
			Help: "Total number of requests to get Person.",
		},
	)
	GetRequestFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_request_failed_total",
			Help: "Total number of requests to get Person that failed.",
		},
	)

	ListRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "list_request_total",
			Help: "Total number of requests to list Persons.",
		},
	)
	ListRequestFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "list_request_failed_total",
			Help: "Total number of requests to list Persons that failed.",
		},
	)

	UpdateRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "update_request_total",
			Help: "Total number of requests to update Person.",
		},
	)
	UpdateRequestFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "update_request_failed_total",
			Help: "Total number of requests to update Person that failed.",
		},
	)

	DeleteRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "delete_request_total",
			Help: "Total number of requests to delete Person.",
		},
	)
	DeleteRequestFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "delete_request_failed_total",
			Help: "Total number of requests to delete Person that failed.",
		},
	)
)

func InitMetrics(reg prometheus.Registerer) {
	reg.MustRegister(CreateRequestTotal)
	reg.MustRegister(CreateRequestFailedTotal)
	reg.MustRegister(GetRequestTotal)
	reg.MustRegister(GetRequestFailedTotal)
	reg.MustRegister(ListRequestTotal)
	reg.MustRegister(ListRequestFailedTotal)
	reg.MustRegister(UpdateRequestTotal)
	reg.MustRegister(UpdateRequestFailedTotal)
	reg.MustRegister(DeleteRequestTotal)
	reg.MustRegister(DeleteRequestFailedTotal)
}
