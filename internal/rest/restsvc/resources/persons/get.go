package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/joomcode/errorx"
)

func (h *Handler) Get(ctx context.Context, input *IDInput) (*Output, error) {
	person, err := h.personUseCase.GetPerson(ctx, input.ID)
	if errorx.IsOfType(err, persons.ErrPersonNotFound) {
		h.logger.DebugContext(ctx, "person not found", "error", err, "id", input.ID)
		return nil, huma.Error404NotFound("persons not found")
	}
	if err != nil {
		h.logger.ErrorContext(ctx, "get person failed", "error", err, "id", input.ID)
		return nil, huma.Error500InternalServerError("get persons failed")
	}

	return &Output{Body: modelToResponse(person)}, nil
}
