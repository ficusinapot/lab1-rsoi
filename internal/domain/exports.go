package domain

import (
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/coreifc/status"
)

type Exports struct {
	PersonUseCase persons.PersonUseCase
	StatusUseCase status.StatusUseCase
}
