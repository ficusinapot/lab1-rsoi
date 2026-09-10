package domain

import (
	"github.com/ficusinapot/ds/internal/models/dbifc"
)

type Imports struct {
	PersonRepository dbifc.PersonRepository
	DBStatusProvider dbifc.StatusProvider
}
