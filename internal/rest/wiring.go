package rest

import (
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/persons"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/status"
)

type handlers struct {
	person *persons.Handler
	status *status.Handler
}

func newHandlers(imports Imports) handlers {
	personUseCase := imports.PersonUseCase
	personHandler := persons.NewHandler(personUseCase)
	statusHandler := status.NewHandler(imports.StatusUseCase)

	return handlers{
		person: personHandler,
		status: statusHandler,
	}
}
