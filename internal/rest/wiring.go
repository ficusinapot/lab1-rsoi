package rest

import (
	personcontracts "github.com/ficusinapot/ds/internal/models/core/persons"
	statuscontracts "github.com/ficusinapot/ds/internal/models/core/status"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/persons"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/status"
)

type handlers struct {
	person *persons.Handler
	status *status.Handler
}

func newHandlers(personUseCase personcontracts.PersonUseCase, statusUseCase statuscontracts.StatusUseCase) handlers {
	personHandler := persons.NewHandler(personUseCase)
	statusHandler := status.NewHandler(statusUseCase)

	return handlers{
		person: personHandler,
		status: statusHandler,
	}
}
