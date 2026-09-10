package db

import (
	"database/sql"

	"github.com/ficusinapot/ds/internal/models/dbifc"
)

type Exports struct {
	DB               *sql.DB
	PersonRepository dbifc.PersonRepository
	StatusProvider   dbifc.StatusProvider
}
