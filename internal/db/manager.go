package db

import (
	"context"
	"log/slog"

	dbmetrics "github.com/ficusinapot/ds/internal/db/metrics"
	"github.com/ficusinapot/ds/internal/deps"
	managercore "github.com/ficusinapot/ds/internal/manager"

	"github.com/joomcode/errorx"
)

type Manager struct {
	managercore.GenericManager
	base *managerBase
}

type managerBase struct {
	config     Config
	db         *Client
	components components
	exports    Exports
	monitor    *monitor
	logger     *slog.Logger
}

func NewManager(cfg Config) *Manager {
	base := &managerBase{
		config:     cfg,
		db:         nil,
		components: components{},
		exports:    Exports{},
		monitor:    nil,
		logger:     slog.Default().With("subsystem", "db", "manager", "database"),
	}

	return &Manager{
		GenericManager: managercore.NewGenericManager(base),
		base:           base,
	}
}

func (m *managerBase) Name() string {
	return "DB"
}

func (m *managerBase) PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error) {
	_ = ctx
	_ = imports

	database, err := Open(m.config)
	if err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "open database")
	}

	m.db = database
	m.components = newComponents(database.Ent(), database.SQLDB())
	m.exports = m.components.exports(database.SQLDB())
	m.monitor = newMonitor(m.components.statusProvider)
	dbmetrics.TotalConnectionsCreated.Inc()

	return m.exportsContainer(), nil
}

func (m *managerBase) RunInner(ctx context.Context) error {
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

func (m *managerBase) StopCommunicationsInner(ctx context.Context) error {
	_ = ctx

	if m.monitor != nil {
		if err := m.monitor.Stop(); err != nil {
			return errorx.InternalError.Wrap(err, "stop database monitor")
		}
	}

	return nil
}

func (m *managerBase) ShutdownInner(ctx context.Context) error {
	_ = ctx

	if m.db == nil {
		return nil
	}

	if m.monitor != nil {
		if err := m.monitor.Stop(); err != nil {
			return errorx.InternalError.Wrap(err, "stop database monitor")
		}
	}
	if err := m.components.Shutdown(); err != nil {
		return errorx.InternalError.Wrap(err, "shutdown database components")
	}
	if err := Close(m.db); err != nil {
		return errorx.InternalError.Wrap(err, "close database")
	}
	dbmetrics.TotalConnectionsDestroyed.Inc()
	dbmetrics.ActiveConnections.Set(0)
	dbmetrics.AliveConnection.Set(0)
	m.db = nil
	m.components = components{}
	m.exports = Exports{}
	m.monitor = nil

	return nil
}

func (m *managerBase) DisposeInner() {
	if m.db == nil {
		return
	}

	if err := Close(m.db); err != nil {
		m.logger.Warn("dispose database", "error", err)
	}
	if m.monitor != nil {
		m.monitor.Dispose()
	}
	m.components.Dispose()
	dbmetrics.TotalConnectionsDestroyed.Inc()
	dbmetrics.ActiveConnections.Set(0)
	dbmetrics.AliveConnection.Set(0)
	m.db = nil
	m.components = components{}
	m.exports = Exports{}
	m.monitor = nil
}

func (m *managerBase) exportsContainer() deps.Container {
	container := deps.New()
	deps.Add(&container, deps.DBKey, m.exports.DB)
	deps.Add(&container, deps.PersonRepositoryKey, m.exports.PersonRepository)
	deps.Add(&container, deps.StatusProviderKey, m.exports.StatusProvider)

	return container
}

func (m *managerBase) updateConnectionMetrics() {
	stats := m.db.SQLDB().Stats()
	dbmetrics.ActiveConnections.Set(float64(stats.InUse))
}
