package domain

import (
	"github.com/ficusinapot/ds/internal/domain/persons"
	"github.com/ficusinapot/ds/internal/domain/status"
)

type components struct {
	personUseCase *persons.PersonUseCase
	statusUseCase *status.StatusUseCase
}

func newComponents(imports Imports) components {
	return components{
		personUseCase: persons.NewPersonUseCase(imports.PersonRepository),
		statusUseCase: status.NewStatusUseCase(imports.DBStatusProvider),
	}
}

func (c components) exports() Exports {
	return Exports{
		PersonUseCase: c.personUseCase,
		StatusUseCase: c.statusUseCase,
	}
}

func (c components) Shutdown() error {
	return nil
}

func (c components) Dispose() {
}
