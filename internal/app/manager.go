package app

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/config"
	"github.com/ficusinapot/ds/internal/db"
	"github.com/ficusinapot/ds/internal/deps"
	"github.com/ficusinapot/ds/internal/domain"
	managercore "github.com/ficusinapot/ds/internal/manager"
	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/rest"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

type Manager struct {
	managercore.GenericManager
	base *managerBase
}

type managerBase struct {
	dbManager        *db.Manager
	domainManager    *domain.Manager
	restManager      *rest.Manager
	managers         []managercore.Manager
	preparedManagers []managercore.Manager
	metricsRegistry  *metrics.Registry
}

func NewManager(cfg config.Config, logger *slog.Logger, metricsRegistry *metrics.Registry, appInfo restsvc.AppInfo) *Manager {
	dbManager := db.NewManager(cfg.Database)
	domainManager := domain.NewManager()
	restManager := rest.NewManager(
		appInfo,
		cfg.Rest,
		logger,
		metricsRegistry,
	)

	base := &managerBase{
		dbManager:       dbManager,
		domainManager:   domainManager,
		restManager:     restManager,
		managers:        []managercore.Manager{dbManager, domainManager, restManager},
		metricsRegistry: metricsRegistry,
	}

	return &Manager{
		GenericManager: managercore.NewGenericManager(base),
		base:           base,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	if _, err := m.Prepare(ctx, deps.New()); err != nil {
		return errorx.InitializationFailed.Wrap(err, "prepare application manager")
	}
	if err := m.Run(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "run application manager")
	}

	return nil
}

func (m *managerBase) Name() string {
	return "Application"
}

func (m *managerBase) PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error) {
	_ = imports

	dbExports, err := m.dbManager.Prepare(ctx, deps.New())
	if err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "prepare %s manager", m.dbManager.Name())
	}
	m.preparedManagers = append(m.preparedManagers, m.dbManager)
	container := deps.New().Merge(dbExports)

	sqlDB, err := deps.Resolve(container, deps.DBKey)
	if err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "resolve database dependency")
	}
	if err := m.metricsRegistry.RegisterDatabase(sqlDB); err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "register database metrics")
	}

	domainExports, err := m.domainManager.Prepare(ctx, container)
	if err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "prepare %s manager", m.domainManager.Name())
	}
	m.preparedManagers = append(m.preparedManagers, m.domainManager)
	container = container.Merge(domainExports)

	if _, err := m.restManager.Prepare(ctx, container); err != nil {
		return deps.Container{}, errorx.InitializationFailed.Wrap(err, "prepare %s manager", m.restManager.Name())
	}
	m.preparedManagers = append(m.preparedManagers, m.restManager)

	return deps.New(), nil
}

func (m *managerBase) RunInner(ctx context.Context) error {
	for _, manager := range m.preparedManagers {
		if err := manager.Run(ctx); err != nil {
			return errorx.ExternalError.Wrap(err, "run %s manager", manager.Name())
		}
	}

	return nil
}

func (m *managerBase) StopCommunicationsInner(ctx context.Context) error {
	var result error
	for i := len(m.preparedManagers) - 1; i >= 0; i-- {
		if err := m.preparedManagers[i].StopCommunications(ctx); err != nil {
			result = errorx.DecorateMany("stop communications", result, err)
		}
	}

	return result //nolint:wrapcheck // Result is already accumulated with contextual shutdown errors.
}

func (m *managerBase) ShutdownInner(ctx context.Context) error {
	var result error
	for i := len(m.preparedManagers) - 1; i >= 0; i-- {
		if err := m.preparedManagers[i].Shutdown(ctx); err != nil {
			result = errorx.DecorateMany("shutdown managers", result, err)
		}
	}
	m.preparedManagers = nil

	return result //nolint:wrapcheck // Result is already accumulated with contextual shutdown errors.
}

func (m *managerBase) DisposeInner() {
	for i := len(m.managers) - 1; i >= 0; i-- {
		m.managers[i].Dispose()
	}
	m.preparedManagers = nil
}

func (m *Manager) RestServiceTerminationNotification() <-chan struct{} {
	return m.base.restManager.ServiceTerminationNotification()
}
