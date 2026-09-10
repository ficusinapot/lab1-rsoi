package dbifc

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"
)

//go:generate go tool mockgen -source=persons_repo.go -destination=mocks/person_repository.go -package=mocks

type PersonRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Person, error)
	List(ctx context.Context) ([]entities.Person, error)
	Create(ctx context.Context, person *entities.Person) (int, error)
	Update(ctx context.Context, person *entities.Person) error
	Delete(ctx context.Context, id int) error
}
