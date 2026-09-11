package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ficusinapot/ds/internal/models/core/persons"
	"github.com/joomcode/errorx"
)

func (h *Handler) Update(ctx context.Context, input *WithIDInput) (*Output, error) {
	person, err := h.personUseCase.GetPerson(ctx, input.ID)
	if errorx.IsOfType(err, persons.ErrPersonNotFound) {
		h.logger.DebugContext(ctx, "person not found", "error", err, "id", input.ID)
		return nil, huma.Error404NotFound("persons not found")
	}
	if err != nil {
		h.logger.ErrorContext(ctx, "get person before update failed", "error", err, "id", input.ID)
		return nil, huma.Error500InternalServerError("get persons failed")
	}

	person = mergeRequest(person, input.Body)
	if err := h.personUseCase.UpdatePerson(ctx, person); err != nil {
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			h.logger.DebugContext(ctx, "person not found", "error", err, "id", input.ID)
			return nil, huma.Error404NotFound("persons not found")
		}

		h.logger.ErrorContext(ctx, "update person failed", "error", err, "id", input.ID)
		return nil, huma.Error500InternalServerError("update persons failed")
	}

	return &Output{Body: modelToResponse(person)}, nil
}
