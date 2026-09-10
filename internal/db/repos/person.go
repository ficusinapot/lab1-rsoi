package repos

import (
	"log/slog"

	"github.com/ficusinapot/ds/internal/db/ent"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type PersonRepository struct {
	client *ent.Client
	logger *slog.Logger
	tracer trace.Tracer
}

func NewPersonRepository(client *ent.Client) *PersonRepository {
	return &PersonRepository{
		client: client,
		logger: slog.Default().With("subsystem", "db", "repo", "persons"),
		tracer: otel.Tracer("github.com/ficusinapot/ds/internal/db/repos"),
	}
}
