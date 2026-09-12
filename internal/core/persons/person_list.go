package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/core/persons/metrics"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (u *PersonUseCase) ListPersons(ctx context.Context) ([]entities.Person, error) {
	u.logger.DebugContext(ctx, "list persons")
	ctx, span := u.tracer.Start(ctx, "Domain: ListPersons")
	defer span.End()

	metrics.ListRequestTotal.Inc()

	persons, err := u.personRepository.List(ctx)
	if err != nil {
		metrics.ListRequestFailedTotal.Inc()
		return nil, errorx.InternalError.Wrap(err, "list persons")
	}

	return persons, nil
}
