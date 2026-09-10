package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/domain/persons/metrics"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"

	"github.com/joomcode/errorx"
)

func (u *PersonUseCase) DeletePerson(ctx context.Context, id int) error {
	u.logger.DebugContext(ctx, "delete person", "id", id)
	ctx, span := u.tracer.Start(ctx, "Domain: DeletePerson")
	defer span.End()

	metrics.DeleteRequestTotal.Inc()

	if err := u.personRepository.Delete(ctx, id); err != nil {
		metrics.DeleteRequestFailedTotal.Inc()
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			return persons.ErrPersonNotFound.Wrap(err, "delete person")
		}

		return errorx.InternalError.Wrap(err, "delete person")
	}

	return nil
}
