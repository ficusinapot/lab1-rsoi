package db

import (
	"database/sql"

	"github.com/ficusinapot/ds/internal/db/ent"
	"github.com/ficusinapot/ds/internal/db/repos"
	"github.com/ficusinapot/ds/internal/db/resources"
	"github.com/ficusinapot/ds/internal/models/dbifc"
)

type components struct {
	personRepository dbifc.PersonRepository
	statusProvider   *resources.StatusProvider
}

func newComponents(client *ent.Client, sqlDB *sql.DB) components {
	return components{
		personRepository: repos.NewPersonRepository(client),
		statusProvider:   resources.NewStatusProvider(sqlDB),
	}
}

func (c components) exports(sqlDB *sql.DB) Exports {
	return Exports{
		DB:               sqlDB,
		PersonRepository: c.personRepository,
		StatusProvider:   c.statusProvider,
	}
}

func (c components) Shutdown() error {
	return nil
}

func (c components) Dispose() {
}
