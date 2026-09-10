package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/domain/persons/metrics"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (u *PersonUseCase) UpdatePerson(ctx context.Context, person *entities.Person) error {
	u.logger.DebugContext(ctx, "update person", "id", person.ID)
	ctx, span := u.tracer.Start(ctx, "Domain: UpdatePerson")
	defer span.End()

	metrics.UpdateRequestTotal.Inc()

	if err := u.personRepository.Update(ctx, person); err != nil {
		metrics.UpdateRequestFailedTotal.Inc()
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			return persons.ErrPersonNotFound.Wrap(err, "update person")
		}

		return errorx.InternalError.Wrap(err, "update person")
	}

	return nil
}
