package repos

import (
	"context"

	dbent "github.com/ficusinapot/ds/internal/db/ent"
	"github.com/ficusinapot/ds/internal/models/core/persons"

	"github.com/joomcode/errorx"
)

func (r *PersonRepository) Delete(ctx context.Context, id int) error {
	r.logger.DebugContext(ctx, "delete person", "id", id)
	ctx, span := r.tracer.Start(ctx, "DB: DeletePerson")
	defer span.End()

	if err := r.client.Person.DeleteOneID(id).Exec(ctx); err != nil {
		if dbent.IsNotFound(err) {
			return persons.ErrPersonNotFound.Wrap(err, "delete person")
		}

		return errorx.ExternalError.Wrap(err, "delete person")
	}

	return nil
}
