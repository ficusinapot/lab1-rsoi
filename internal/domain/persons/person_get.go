package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/domain/persons/metrics"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (u *PersonUseCase) GetPerson(ctx context.Context, id int) (*entities.Person, error) {
	u.logger.DebugContext(ctx, "get person", "id", id)
	ctx, span := u.tracer.Start(ctx, "Domain: GetPerson")
	defer span.End()

	metrics.GetRequestTotal.Inc()

	person, err := u.personRepository.GetByID(ctx, id)
	if err != nil {
		metrics.GetRequestFailedTotal.Inc()
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			return nil, persons.ErrPersonNotFound.Wrap(err, "get person by id")
		}

		return nil, errorx.InternalError.Wrap(err, "get person by id")
	}

	return person, nil
}
