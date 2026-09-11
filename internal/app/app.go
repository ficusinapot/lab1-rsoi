package app

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/config"
	"github.com/ficusinapot/ds/internal/core"
	"github.com/ficusinapot/ds/internal/db"
	"github.com/ficusinapot/ds/internal/observability/logging"
	"github.com/ficusinapot/ds/internal/observability/metrics"
	"github.com/ficusinapot/ds/internal/observability/tracing"
	"github.com/ficusinapot/ds/internal/rest"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

type App struct {
	manager        *Manager
	logger         *slog.Logger
	closeLogger    func()
	tracerProvider interface {
		Shutdown(context.Context) error
	}
}

func Run(ctx context.Context, cfg config.Config, appInfo restsvc.AppInfo) error {
	application, err := New(ctx, cfg, appInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err := application.Close(ctx); err != nil {
			application.logger.Error("close application", "error", err)
		}
	}()

	return application.Run(ctx)
}

func New(ctx context.Context, cfg config.Config, appInfo restsvc.AppInfo) (*App, error) {
	logger, closeLogger, err := logging.NewLogger(cfg.Logging)
	if err != nil {
		return nil, errorx.InitializationFailed.Wrap(err, "create logger")
	}
	slog.SetDefault(logger)

	tracerProvider, err := tracing.NewTracerProvider(ctx, appInfo.Name, appInfo.Version, cfg.Logging.OpenTelemetry)
	if err != nil {
		closeLogger()
		return nil, errorx.InitializationFailed.Wrap(err, "create tracer provider")
	}
	logger.Info(
		"opentelemetry tracing configured",
		"enabled",
		cfg.Logging.OpenTelemetry.Address != "",
		"address",
		cfg.Logging.OpenTelemetry.Address,
	)

	metricsRegistry := metrics.NewRegistry(appInfo.Name, appInfo.Version)
	db.InitMetrics(metricsRegistry.Registerer())
	rest.InitMetrics(metricsRegistry.Registerer())
	core.InitMetrics(metricsRegistry.Registerer())
	manager := NewManager(cfg, logger, metricsRegistry, appInfo)

	return &App{
		manager:        manager,
		logger:         logger,
		closeLogger:    closeLogger,
		tracerProvider: tracerProvider,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.manager.Start(ctx); err != nil {
		return errorx.InternalError.Wrap(err, "start application manager")
	}

	select {
	case <-a.manager.RestServiceTerminationNotification():
		a.logger.Info("REST API service terminated")
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (a *App) Close(ctx context.Context) error {
	defer a.closeLogger()
	defer a.manager.Dispose()

	var result error
	if err := a.manager.Shutdown(ctx); err != nil {
		result = errorx.DecorateMany("shutdown application", result, err)
	}

	if err := a.tracerProvider.Shutdown(ctx); err != nil {
		result = errorx.DecorateMany("shutdown application", result, err)
	}

	return result //nolint:wrapcheck
}
