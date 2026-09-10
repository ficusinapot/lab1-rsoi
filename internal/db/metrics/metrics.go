package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "connections_active",
			Help: "Number of active database connections.",
		},
	)
	TotalConnectionsCreated = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "connections_created_total",
			Help: "Total number of database connections created.",
		},
	)
	TotalConnectionsDestroyed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "connections_destroyed_total",
			Help: "Total number of database connections destroyed.",
		},
	)
	AliveConnection = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "connections_alive",
			Help: "Database connection alive flag.",
		},
	)
)

func InitMetrics(reg prometheus.Registerer) {
	reg.MustRegister(ActiveConnections)
	reg.MustRegister(TotalConnectionsCreated)
	reg.MustRegister(TotalConnectionsDestroyed)
	reg.MustRegister(AliveConnection)
}
