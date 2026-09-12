package rest

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/models/core/persons"
	"github.com/ficusinapot/ds/internal/models/core/status"
	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

type Manager struct {
	appInfo         restsvc.AppInfo
	config          Config
	logger          *slog.Logger
	metricsRegistry *metrics.Registry
	service         *restsvc.Service
}

func NewManager(
	appInfo restsvc.AppInfo,
	restConfig Config,
	logger *slog.Logger,
	metricsRegistry *metrics.Registry,
) *Manager {
	return &Manager{
		appInfo:         appInfo,
		config:          restConfig,
		logger:          logger.With("subsystem", "rest", "manager", "api"),
		metricsRegistry: metricsRegistry,
		service:         nil,
	}
}

func (m *Manager) Name() string {
	return "REST API"
}

func (m *Manager) Prepare(ctx context.Context, personUseCase persons.PersonUseCase, statusUseCase status.StatusUseCase) error {
	_ = ctx

	if personUseCase == nil {
		return errorx.IllegalState.New("person use case is required")
	}
	if statusUseCase == nil {
		return errorx.IllegalState.New("status use case is required")
	}

	handlers := newHandlers(personUseCase, statusUseCase)
	m.service = restsvc.NewService(
		m.appInfo,
		m.config.Service,
		m.config.Metrics,
		handlers.person,
		handlers.status,
		m.logger,
		m.metricsRegistry,
	)

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	if m.service == nil {
		return errorx.IllegalState.New("REST API service is not prepared")
	}

	if err := m.service.Run(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "run REST API service")
	}

	return nil
}

func (m *Manager) StopCommunications(ctx context.Context) error {
	if m.service == nil {
		return nil
	}

	if err := m.service.Stop(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "stop REST API service")
	}

	return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
	_ = ctx

	return nil
}

func (m *Manager) Dispose() {
	if m.service != nil {
		m.service.Dispose()
	}
}

func (m *Manager) ServiceTerminationNotification() <-chan struct{} {
	if m.service == nil {
		ch := make(chan struct{})
		close(ch)
		return ch
	}

	return m.service.TerminationNotification()
}
