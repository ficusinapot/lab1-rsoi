package manager

import (
	"context"

	"github.com/ficusinapot/ds/internal/deps"

	"github.com/joomcode/errorx"
)

type Manager interface {
	Name() string
	Prepare(ctx context.Context, imports deps.Container) (deps.Container, error)
	Run(ctx context.Context) error
	StopCommunications(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Dispose()
}

type Inner interface {
	Name() string
	PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error)
	RunInner(ctx context.Context) error
	StopCommunicationsInner(ctx context.Context) error
	ShutdownInner(ctx context.Context) error
	DisposeInner()
}

type GenericManager struct {
	inner Inner
	state State
}

func NewGenericManager(inner Inner) GenericManager {
	return GenericManager{
		inner: inner,
		state: Created,
	}
}

func (m *GenericManager) Name() string {
	return m.inner.Name()
}

func (m *GenericManager) Prepare(ctx context.Context, imports deps.Container) (deps.Container, error) {
	if err := m.state.ValidateTransition(Prepared); err != nil {
		return deps.Container{}, err
	}

	exports, err := m.inner.PrepareInner(ctx, imports)
	if err != nil {
		if shutdownErr := m.inner.ShutdownInner(ctx); shutdownErr != nil {
			err = errorx.DecorateMany("cleanup after failed prepare", err, shutdownErr)
		}

		return deps.Container{}, PrepareError.Wrap(err, "prepare %s manager", m.Name())
	}

	m.state = Prepared
	return exports, nil
}

func (m *GenericManager) Run(ctx context.Context) error {
	if err := m.state.ValidateTransition(Running); err != nil {
		return err
	}

	if err := m.inner.RunInner(ctx); err != nil {
		return RunError.Wrap(err, "run %s manager", m.Name())
	}

	m.state = Running
	return nil
}

func (m *GenericManager) StopCommunications(ctx context.Context) error {
	if m.state != Running {
		return nil
	}

	if err := m.state.ValidateTransition(CommunicationStopped); err != nil {
		return err
	}

	if err := m.inner.StopCommunicationsInner(ctx); err != nil {
		return StopCommunicationsError.Wrap(err, "stop %s communications", m.Name())
	}

	m.state = CommunicationStopped
	return nil
}

func (m *GenericManager) Shutdown(ctx context.Context) error {
	if m.state == Shutdown || m.state == Disposed {
		return nil
	}

	if m.state == Running {
		if err := m.StopCommunications(ctx); err != nil {
			return err
		}
	}

	if err := m.state.ValidateTransition(Shutdown); err != nil {
		return err
	}

	if err := m.inner.ShutdownInner(ctx); err != nil {
		return ShutdownError.Wrap(err, "shutdown %s manager", m.Name())
	}

	m.state = Shutdown
	return nil
}

func (m *GenericManager) Dispose() {
	if m.state == Disposed {
		return
	}

	m.inner.DisposeInner()
	m.state = Disposed
}

func (m *GenericManager) State() State {
	return m.state
}
