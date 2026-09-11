package db

import (
	"context"
	"log/slog"

	dbmetrics "github.com/ficusinapot/ds/internal/db/metrics"
	"github.com/ficusinapot/ds/internal/db/repos"
	"github.com/ficusinapot/ds/internal/db/resources"

	"github.com/joomcode/errorx"
)

type Manager struct {
	config  Config
	db      *Client
	exports Exports
	monitor *monitor
	logger  *slog.Logger
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		config:  cfg,
		db:      nil,
		exports: Exports{},
		monitor: nil,
		logger:  slog.Default().With("subsystem", "db", "manager", "database"),
	}
}

func (m *Manager) Name() string {
	return "DB"
}

func (m *Manager) Prepare(ctx context.Context) error {
	_ = ctx

	database, err := Open(m.config)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "open database")
	}

	m.db = database
	statusProvider := resources.NewStatusProvider(database.SQLDB())
	m.exports = Exports{
		DB:               database.SQLDB(),
		PersonRepository: repos.NewPersonRepository(database.Ent()),
		StatusProvider:   statusProvider,
	}
	m.monitor = newMonitor(statusProvider)
	dbmetrics.TotalConnectionsCreated.Inc()

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	_ = ctx

	if m.db == nil {
		return errorx.IllegalState.New("database manager is not prepared")
	}
	m.updateConnectionMetrics()
	if err := m.monitor.Start(ctx, m.config.HealthCheck); err != nil {
		return errorx.ExternalError.Wrap(err, "start database monitor")
	}

	return nil
}

func (m *Manager) StopCommunications(ctx context.Context) error {
	_ = ctx

	if m.monitor != nil {
		if err := m.monitor.Stop(); err != nil {
			return errorx.InternalError.Wrap(err, "stop database monitor")
		}
	}

	return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
	_ = ctx

	if m.db == nil {
		return nil
	}

	if m.monitor != nil {
		if err := m.monitor.Stop(); err != nil {
			return errorx.InternalError.Wrap(err, "stop database monitor")
		}
	}
	if err := Close(m.db); err != nil {
		return errorx.InternalError.Wrap(err, "close database")
	}
	m.reset()

	return nil
}

func (m *Manager) Dispose() {
	if m.db == nil {
		return
	}

	if err := Close(m.db); err != nil {
		m.logger.Warn("dispose database", "error", err)
	}
	if m.monitor != nil {
		m.monitor.Dispose()
	}
	m.reset()
}

func (m *Manager) reset() {
	dbmetrics.TotalConnectionsDestroyed.Inc()
	dbmetrics.ActiveConnections.Set(0)
	dbmetrics.AliveConnection.Set(0)
	m.db = nil
	m.exports = Exports{}
	m.monitor = nil
}

func (m *Manager) Exports() Exports {
	return m.exports
}

func (m *Manager) updateConnectionMetrics() {
	stats := m.db.SQLDB().Stats()
	dbmetrics.ActiveConnections.Set(float64(stats.InUse))
}
