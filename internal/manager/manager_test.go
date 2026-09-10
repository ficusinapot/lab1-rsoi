package manager

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ficusinapot/ds/internal/deps"

	"github.com/joomcode/errorx"
)

func TestGenericManagerLifecycle(t *testing.T) {
	t.Parallel()

	inner := &fakeInner{}
	manager := NewGenericManager(inner)

	if _, err := manager.Prepare(context.Background(), deps.New()); err != nil {
		t.Fatalf("prepare: %v", err)
	}
	if err := manager.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	if err := manager.StopCommunications(context.Background()); err != nil {
		t.Fatalf("stop communications: %v", err)
	}
	if err := manager.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	manager.Dispose()

	if manager.State() != Disposed {
		t.Fatalf("unexpected state: %s", manager.State())
	}

	wantCalls := []string{"prepare", "run", "stop", "shutdown", "dispose"}
	if len(inner.calls) != len(wantCalls) {
		t.Fatalf("unexpected calls: %#v", inner.calls)
	}
	for i := range wantCalls {
		if inner.calls[i] != wantCalls[i] {
			t.Fatalf("unexpected calls: %#v", inner.calls)
		}
	}
}

func TestGenericManagerRejectsInvalidTransition(t *testing.T) {
	t.Parallel()

	manager := NewGenericManager(&fakeInner{})

	if err := manager.Run(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestGenericManagerWrapsInnerError(t *testing.T) {
	t.Parallel()

	innerErr := errors.New("boom")
	inner := &fakeInner{prepareErr: innerErr}
	manager := NewGenericManager(inner)

	_, err := manager.Prepare(context.Background(), deps.New())
	if !errorx.IsOfType(err, PrepareError) || !strings.Contains(err.Error(), innerErr.Error()) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenericManagerCleansUpAfterFailedPrepare(t *testing.T) {
	t.Parallel()

	innerErr := errors.New("boom")
	inner := &fakeInner{prepareErr: innerErr}
	manager := NewGenericManager(inner)

	_, err := manager.Prepare(context.Background(), deps.New())
	if !errorx.IsOfType(err, PrepareError) || !strings.Contains(err.Error(), innerErr.Error()) {
		t.Fatalf("unexpected error: %v", err)
	}
	if manager.State() != Created {
		t.Fatalf("unexpected state: %s", manager.State())
	}

	wantCalls := []string{"prepare", "shutdown"}
	if len(inner.calls) != len(wantCalls) {
		t.Fatalf("unexpected calls: %#v", inner.calls)
	}
	for i := range wantCalls {
		if inner.calls[i] != wantCalls[i] {
			t.Fatalf("unexpected calls: %#v", inner.calls)
		}
	}
}

type fakeInner struct {
	calls      []string
	prepareErr error
}

func (f *fakeInner) Name() string {
	return "fake"
}

func (f *fakeInner) PrepareInner(ctx context.Context, imports deps.Container) (deps.Container, error) {
	_ = ctx
	_ = imports
	f.calls = append(f.calls, "prepare")
	return deps.New(), f.prepareErr
}

func (f *fakeInner) RunInner(ctx context.Context) error {
	_ = ctx
	f.calls = append(f.calls, "run")
	return nil
}

func (f *fakeInner) StopCommunicationsInner(ctx context.Context) error {
	_ = ctx
	f.calls = append(f.calls, "stop")
	return nil
}

func (f *fakeInner) ShutdownInner(ctx context.Context) error {
	_ = ctx
	f.calls = append(f.calls, "shutdown")
	return nil
}

func (f *fakeInner) DisposeInner() {
	f.calls = append(f.calls, "dispose")
}
