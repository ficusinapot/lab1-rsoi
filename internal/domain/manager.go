package domain

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/deps"
	managercore "github.com/ficusinapot/ds/internal/manager"

	"github.com/joomcode/errorx"
)

type Manager struct {
	managercore.GenericManager
	base *managerBase
}

type managerBase struct {
	imports    Imports
	components components
	exports    Exports
	logger     *slog.Logger
}

func NewManager() *Manager {
	base := &managerBase{
		imports:    Imports{},
		components: components{},
		exports:    Exports{},
		logger:     slog.Default().With("subsystem", "coreifc", "manager", "coreifc"),
	}

	return &Manager{
		GenericManager: managercore.NewGenericManager(base),
		base:           base,
	}
}

func (m *managerBase) Name() string {
	return "Domain"
}

func (m *managerBase) PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error) {
	_ = ctx

	personRepository, err := deps.Resolve(imports, deps.PersonRepositoryKey)
	if err != nil {
		return deps.Container{}, errorx.IllegalState.Wrap(err, "resolve person repository")
	}
	statusProvider, err := deps.Resolve(imports, deps.StatusProviderKey)
	if err != nil {
		return deps.Container{}, errorx.IllegalState.Wrap(err, "resolve status provider")
	}

	m.imports = Imports{
		PersonRepository: personRepository,
		DBStatusProvider: statusProvider,
	}
	m.components = newComponents(m.imports)
	m.exports = m.components.exports()

	return m.exportsContainer(), nil
}

func (m *managerBase) RunInner(ctx context.Context) error {
	_ = ctx

	if m.exports.PersonUseCase == nil {
		return errorx.IllegalState.New("coreifc manager is not prepared")
	}
	if m.exports.StatusUseCase == nil {
		return errorx.IllegalState.New("coreifc manager is not prepared")
	}

	return nil
}

func (m *managerBase) StopCommunicationsInner(ctx context.Context) error {
	_ = ctx

	return nil
}

func (m *managerBase) ShutdownInner(ctx context.Context) error {
	_ = ctx

	if err := m.components.Shutdown(); err != nil {
		return errorx.InternalError.Wrap(err, "shutdown coreifc components")
	}

	m.imports = Imports{}
	m.components = components{}
	m.exports = Exports{}

	return nil
}

func (m *managerBase) DisposeInner() {
	m.components.Dispose()
	m.imports = Imports{}
	m.components = components{}
	m.exports = Exports{}
}

func (m *managerBase) exportsContainer() deps.Container {
	container := deps.New()
	deps.Add(&container, deps.PersonUseCaseKey, m.exports.PersonUseCase)
	deps.Add(&container, deps.StatusUseCaseKey, m.exports.StatusUseCase)

	return container
}
