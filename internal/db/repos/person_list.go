package repos

import (
	"context"

	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (r *PersonRepository) List(ctx context.Context) ([]entities.Person, error) {
	r.logger.DebugContext(ctx, "list persons")
	ctx, span := r.tracer.Start(ctx, "DB: ListPersons")
	defer span.End()

	persons, err := r.client.Person.Query().All(ctx)
	if err != nil {
		return nil, errorx.ExternalError.Wrap(err, "select persons")
	}

	return personsFromEnt(persons), nil
}
