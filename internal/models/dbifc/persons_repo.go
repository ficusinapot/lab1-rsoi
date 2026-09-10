package dbifc

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"
)

type PersonRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Person, error)
	List(ctx context.Context) ([]entities.Person, error)
	Create(ctx context.Context, person *entities.Person) (int, error)
	Update(ctx context.Context, person *entities.Person) error
	Delete(ctx context.Context, id int) error
}
