package status

import (
	"log/slog"

	statuscontracts "github.com/ficusinapot/ds/internal/models/coreifc/status"
)

type Handler struct {
	statusUseCase statuscontracts.StatusUseCase
	logger        *slog.Logger
}

func NewHandler(statusUseCase statuscontracts.StatusUseCase) *Handler {
	return &Handler{
		statusUseCase: statusUseCase,
		logger:        slog.Default().With("subsystem", "rest", "handler", "status"),
	}
}
