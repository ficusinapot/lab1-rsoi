package repos

import (
	"context"

	dbent "github.com/ficusinapot/ds/internal/db/ent"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (r *PersonRepository) GetByID(ctx context.Context, id int) (*entities.Person, error) {
	r.logger.DebugContext(ctx, "get person by id", "id", id)
	ctx, span := r.tracer.Start(ctx, "DB: GetPersonByID")
	defer span.End()

	person, err := r.client.Person.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, persons.ErrPersonNotFound.Wrap(err, "select person")
		}

		return nil, errorx.ExternalError.Wrap(err, "select person")
	}

	return personFromEnt(person), nil
}
