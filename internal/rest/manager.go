package rest

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/deps"
	managercore "github.com/ficusinapot/ds/internal/manager"
	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

type Manager struct {
	managercore.GenericManager
	base *managerBase
}

type managerBase struct {
	appInfo         restsvc.AppInfo
	config          Config
	logger          *slog.Logger
	metricsRegistry *metrics.Registry
	service         *restsvc.Service
	imports         Imports
}

func NewManager(
	appInfo restsvc.AppInfo,
	restConfig Config,
	logger *slog.Logger,
	metricsRegistry *metrics.Registry,
) *Manager {
	base := &managerBase{
		appInfo:         appInfo,
		config:          restConfig,
		logger:          logger.With("subsystem", "rest", "manager", "api"),
		metricsRegistry: metricsRegistry,
		service:         nil,
		imports:         Imports{},
	}

	return &Manager{
		GenericManager: managercore.NewGenericManager(base),
		base:           base,
	}
}

func (m *managerBase) Name() string {
	return "REST API"
}

func (m *managerBase) PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error) {
	_ = ctx

	personUseCase, err := deps.Resolve(imports, deps.PersonUseCaseKey)
	if err != nil {
		return deps.Container{}, errorx.IllegalState.Wrap(err, "resolve person use case")
	}
	statusUseCase, err := deps.Resolve(imports, deps.StatusUseCaseKey)
	if err != nil {
		return deps.Container{}, errorx.IllegalState.Wrap(err, "resolve status use case")
	}

	m.imports = Imports{
		PersonUseCase: personUseCase,
		StatusUseCase: statusUseCase,
	}
	handlers := newHandlers(m.imports)
	m.service = restsvc.NewService(
		m.appInfo,
		m.config.Service,
		m.config.Metrics,
		handlers.person,
		handlers.status,
		m.logger,
		m.metricsRegistry,
	)

	return deps.New(), nil
}

func (m *managerBase) RunInner(ctx context.Context) error {
	if m.service == nil {
		return errorx.IllegalState.New("REST API service is not prepared")
	}

	if err := m.service.Run(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "run REST API service")
	}

	return nil
}

func (m *managerBase) StopCommunicationsInner(ctx context.Context) error {
	if m.service == nil {
		return nil
	}

	if err := m.service.Stop(ctx); err != nil {
		return errorx.ExternalError.Wrap(err, "stop REST API service")
	}

	return nil
}

func (m *managerBase) ShutdownInner(ctx context.Context) error {
	_ = ctx

	return nil
}

func (m *managerBase) DisposeInner() {
	if m.service != nil {
		m.service.Dispose()
	}
}

func (m *Manager) ServiceTerminationNotification() <-chan struct{} {
	if m.base.service == nil {
		ch := make(chan struct{})
		close(ch)
		return ch
	}

	return m.base.service.TerminationNotification()
}
