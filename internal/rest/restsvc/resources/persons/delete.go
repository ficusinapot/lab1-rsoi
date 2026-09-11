package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ficusinapot/ds/internal/models/core/persons"
	"github.com/joomcode/errorx"
)

func (h *Handler) Delete(ctx context.Context, input *IDInput) (*struct{}, error) {
	if err := h.personUseCase.DeletePerson(ctx, input.ID); err != nil {
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			h.logger.DebugContext(ctx, "person not found", "error", err, "id", input.ID)
			return nil, huma.Error404NotFound("persons not found")
		}

		h.logger.ErrorContext(ctx, "delete person failed", "error", err, "id", input.ID)
		return nil, huma.Error500InternalServerError("delete persons failed")
	}

	return &struct{}{}, nil
}
