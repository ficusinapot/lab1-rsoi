package rest

import (
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/coreifc/status"
)

type Imports struct {
	PersonUseCase persons.PersonUseCase
	StatusUseCase status.StatusUseCase
}
