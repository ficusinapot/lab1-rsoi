package status

import (
	"context"
	"log/slog"

	"github.com/ficusinapot/ds/internal/core/status/metrics"
	statuscontracts "github.com/ficusinapot/ds/internal/models/core/status"
	"github.com/ficusinapot/ds/internal/models/dbifc"
	"github.com/ficusinapot/ds/internal/models/entities"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type StatusUseCase struct {
	dbStatusProvider dbifc.StatusProvider
	logger           *slog.Logger
	tracer           trace.Tracer
}

func NewStatusUseCase(dbStatusProvider dbifc.StatusProvider) *StatusUseCase {
	return &StatusUseCase{
		dbStatusProvider: dbStatusProvider,
		logger:           slog.Default().With("subsystem", "status", "use_case", "status"),
		tracer:           otel.Tracer("github.com/ficusinapot/ds/internal/core/status/usecases"),
	}
}

func (u *StatusUseCase) GetStatus(ctx context.Context) (*entities.StatusData, error) {
	u.logger.DebugContext(ctx, "get status")
	ctx, span := u.tracer.Start(ctx, "Status: GetStatus")
	defer span.End()

	metrics.GetStatusRequestTotal.Inc()

	dbAvailable := u.dbStatusProvider.IsConnectionAvailable(ctx)
	operatingStatus := entities.NonOperational
	if dbAvailable {
		operatingStatus = entities.FullyOperational
	}

	return &entities.StatusData{
		DBAvailable:     dbAvailable,
		OperatingStatus: operatingStatus,
	}, nil
}

func (u *StatusUseCase) Healthz(ctx context.Context) error {
	u.logger.DebugContext(ctx, "healthz")
	ctx, span := u.tracer.Start(ctx, "Status: Healthz")
	defer span.End()

	metrics.HealthzRequestTotal.Inc()

	if !u.dbStatusProvider.IsConnectionAvailable(ctx) {
		return statuscontracts.ErrConnectionNotAvailable.New("database connection is not available")
	}

	return nil
}

func (u *StatusUseCase) Readyz(ctx context.Context) error {
	u.logger.DebugContext(ctx, "readyz")
	ctx, span := u.tracer.Start(ctx, "Status: Readyz")
	defer span.End()

	metrics.ReadyzRequestTotal.Inc()

	if !u.dbStatusProvider.IsConnectionAvailable(ctx) {
		return statuscontracts.ErrConnectionNotAvailable.New("database connection is not available")
	}

	return nil
}
