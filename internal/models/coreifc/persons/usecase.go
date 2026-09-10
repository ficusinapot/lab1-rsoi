package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"
)

type PersonUseCase interface {
	GetPerson(ctx context.Context, id int) (*entities.Person, error)
	ListPersons(ctx context.Context) ([]entities.Person, error)
	CreatePerson(ctx context.Context, person *entities.Person) (int, error)
	UpdatePerson(ctx context.Context, person *entities.Person) error
	DeletePerson(ctx context.Context, id int) error
}
