package db

import (
	"context"
	"sync"
	"time"

	dbmetrics "github.com/ficusinapot/ds/internal/db/metrics"

	"github.com/joomcode/errorx"
)

type monitor struct {
	checker interface {
		Check(context.Context) bool
	}
	checkPeriodStop chan struct{}
	done            chan struct{}
	startOnce       sync.Once
	stopOnce        sync.Once
	started         bool
}

func newMonitor(checker interface {
	Check(context.Context) bool
},
) *monitor {
	return &monitor{
		checker:         checker,
		checkPeriodStop: nil,
		done:            nil,
		startOnce:       sync.Once{},
		stopOnce:        sync.Once{},
		started:         false,
	}
}

func (m *monitor) Start(ctx context.Context, cfg HealthCheckConfig) error {
	var err error
	m.startOnce.Do(func() {
		if cfg.CheckPeriod <= 0 {
			err = errorx.IllegalState.New("database health check period must be greater than zero")
			return
		}

		stop := make(chan struct{})
		done := make(chan struct{})
		m.checkPeriodStop = stop
		m.done = done
		m.started = true

		m.update(ctx)
		go m.run(ctx, cfg, stop, done)
	})

	if err != nil {
		return errorx.IllegalState.Wrap(err, "start database monitor")
	}

	return nil
}

func (m *monitor) Stop() error {
	if !m.started {
		return nil
	}

	m.stopOnce.Do(func() {
		close(m.checkPeriodStop)
		<-m.done
		m.started = false
	})

	return nil
}

func (m *monitor) Dispose() {
	_ = m.Stop()
}

func (m *monitor) run(ctx context.Context, cfg HealthCheckConfig, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)

	ticker := time.NewTicker(cfg.CheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.update(ctx)
		case <-stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (m *monitor) update(ctx context.Context) {
	if m.checker.Check(ctx) {
		dbmetrics.AliveConnection.Set(1)
		return
	}

	dbmetrics.AliveConnection.Set(0)
}
