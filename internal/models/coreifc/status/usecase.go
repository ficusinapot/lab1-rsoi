package status

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"
)

type StatusUseCase interface {
	GetStatus(ctx context.Context) (*entities.StatusData, error)
	Healthz(ctx context.Context) error
	Readyz(ctx context.Context) error
}
