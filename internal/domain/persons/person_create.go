package persons

import (
	"context"

	"github.com/ficusinapot/ds/internal/domain/persons/metrics"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
)

func (u *PersonUseCase) CreatePerson(ctx context.Context, person *entities.Person) (int, error) {
	u.logger.DebugContext(ctx, "create person")
	ctx, span := u.tracer.Start(ctx, "Domain: CreatePerson")
	defer span.End()

	metrics.CreateRequestTotal.Inc()

	id, err := u.personRepository.Create(ctx, person)
	if err != nil {
		metrics.CreateRequestFailedTotal.Inc()
		return 0, errorx.InternalError.Wrap(err, "create person")
	}

	return id, nil
}
