package persons

import (
	"log/slog"

	"github.com/ficusinapot/ds/internal/models/core/persons"
)

type Handler struct {
	personUseCase persons.PersonUseCase
	logger        *slog.Logger
}

func NewHandler(personUseCase persons.PersonUseCase) *Handler {
	return &Handler{
		personUseCase: personUseCase,
		logger:        slog.Default().With("subsystem", "rest", "handler", "persons"),
	}
}
