package persons

import (
	"log/slog"

	"github.com/ficusinapot/ds/internal/models/dbifc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type PersonUseCase struct {
	personRepository dbifc.PersonRepository
	logger           *slog.Logger
	tracer           trace.Tracer
}

func NewPersonUseCase(personRepository dbifc.PersonRepository) *PersonUseCase {
	return &PersonUseCase{
		personRepository: personRepository,
		logger:           slog.Default().With("subsystem", "coreifc", "use_case", "person"),
		tracer:           otel.Tracer("github.com/ficusinapot/ds/internal/coreifc/persons"),
	}
}
