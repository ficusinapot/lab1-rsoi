package repos

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (r *PersonRepository) Create(ctx context.Context, person *entities.Person) (int, error) {
	r.logger.DebugContext(ctx, "create person")
	ctx, span := r.tracer.Start(ctx, "DB: CreatePerson")
	defer span.End()

	created, err := r.client.Person.
		Create().
		SetName(person.Name).
		SetAge(person.Age).
		SetAddress(person.Address).
		SetWork(person.Work).
		Save(ctx)
	if err != nil {
		return 0, errorx.ExternalError.Wrap(err, "insert person")
	}

	return created.ID, nil
}
