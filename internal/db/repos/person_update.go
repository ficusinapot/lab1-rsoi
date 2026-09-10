package repos

import (
	"context"

	dbent "github.com/ficusinapot/ds/internal/db/ent"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (r *PersonRepository) Update(ctx context.Context, person *entities.Person) error {
	r.logger.DebugContext(ctx, "update person", "id", person.ID)
	ctx, span := r.tracer.Start(ctx, "DB: UpdatePerson")
	defer span.End()

	if _, err := r.client.Person.
		UpdateOneID(person.ID).
		SetName(person.Name).
		SetAge(person.Age).
		SetAddress(person.Address).
		SetWork(person.Work).
		Save(ctx); err != nil {
		if dbent.IsNotFound(err) {
			return persons.ErrPersonNotFound.Wrap(err, "update person")
		}

		return errorx.ExternalError.Wrap(err, "update person")
	}

	return nil
}
