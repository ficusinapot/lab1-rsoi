package app

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/config"
	"github.com/ficusinapot/ds/internal/core/persons"
	"github.com/ficusinapot/ds/internal/core/status"
	"github.com/ficusinapot/ds/internal/db"
	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/rest"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

type Manager struct {
	dbManager       *db.Manager
	restManager     *rest.Manager
	metricsRegistry *metrics.Registry
	dbPrepared      bool
	restPrepared    bool
	running         bool
}

func NewManager(cfg config.Config, logger *slog.Logger, metricsRegistry *metrics.Registry, appInfo restsvc.AppInfo) *Manager {
	dbManager := db.NewManager(cfg.Database)
	restManager := rest.NewManager(
		appInfo,
		cfg.Rest,
		logger,
		metricsRegistry,
	)

	return &Manager{
		dbManager:       dbManager,
		restManager:     restManager,
		metricsRegistry: metricsRegistry,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	if err := m.Prepare(ctx); err != nil {
		return errorx.InitializationFailed.Wrap(err, "prepare application manager")
	}
	if err := m.Run(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "run application manager")
	}

	return nil
}

func (m *Manager) Name() string {
	return "Application"
}

func (m *Manager) Prepare(ctx context.Context) error {
	if err := m.dbManager.Prepare(ctx); err != nil {
		return errorx.InitializationFailed.Wrap(err, "prepare %s manager", m.dbManager.Name())
	}
	m.dbPrepared = true
	dbExports := m.dbManager.Exports()

	if err := m.metricsRegistry.RegisterDatabase(dbExports.DB); err != nil {
		_ = m.Shutdown(ctx)
		return errorx.InitializationFailed.Wrap(err, "register database metrics")
	}

	personUseCase := persons.NewPersonUseCase(dbExports.PersonRepository)
	statusUseCase := status.NewStatusUseCase(dbExports.StatusProvider)
	if err := m.restManager.Prepare(ctx, personUseCase, statusUseCase); err != nil {
		_ = m.Shutdown(ctx)
		return errorx.InitializationFailed.Wrap(err, "prepare %s manager", m.restManager.Name())
	}
	m.restPrepared = true

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	if m.dbPrepared {
		if err := m.dbManager.Run(ctx); err != nil {
			return errorx.ExternalError.Wrap(err, "run %s manager", m.dbManager.Name())
		}
	}
	if m.restPrepared {
		if err := m.restManager.Run(ctx); err != nil {
			return errorx.ExternalError.Wrap(err, "run %s manager", m.restManager.Name())
		}
	}

	m.running = true
	return nil
}

func (m *Manager) StopCommunications(ctx context.Context) error {
	if !m.running {
		return nil
	}

	var result error
	if m.restPrepared {
		if err := m.restManager.StopCommunications(ctx); err != nil {
			result = errorx.DecorateMany("stop communications", result, err)
		}
	}
	if m.dbPrepared {
		if err := m.dbManager.StopCommunications(ctx); err != nil {
			result = errorx.DecorateMany("stop communications", result, err)
		}
	}
	m.running = false

	return result //nolint:wrapcheck
}

func (m *Manager) Shutdown(ctx context.Context) error {
	var result error
	if err := m.StopCommunications(ctx); err != nil {
		result = errorx.DecorateMany("shutdown managers", result, err)
	}

	if m.restPrepared {
		if err := m.restManager.Shutdown(ctx); err != nil {
			result = errorx.DecorateMany("shutdown managers", result, err)
		}
		m.restPrepared = false
	}
	if m.dbPrepared {
		if err := m.dbManager.Shutdown(ctx); err != nil {
			result = errorx.DecorateMany("shutdown managers", result, err)
		}
		m.dbPrepared = false
	}
	m.running = false

	return result //nolint:wrapcheck // Result is already accumulated with contextual shutdown errors.
}

func (m *Manager) Dispose() {
	m.restManager.Dispose()
	m.dbManager.Dispose()
	m.restPrepared = false
	m.dbPrepared = false
	m.running = false
}

func (m *Manager) RestServiceTerminationNotification() <-chan struct{} {
	return m.restManager.ServiceTerminationNotification()
}
