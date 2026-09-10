package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/joomcode/errorx"
)

func (h *Handler) Update(ctx context.Context, input *WithIDInput) (*struct{}, error) {
	person := requestToModel(input.Body)
	person.ID = input.ID

	if err := h.personUseCase.UpdatePerson(ctx, person); err != nil {
		if errorx.IsOfType(err, persons.ErrPersonNotFound) {
			h.logger.DebugContext(ctx, "person not found", "error", err, "id", input.ID)
			return nil, huma.Error404NotFound("persons not found")
		}

		h.logger.ErrorContext(ctx, "update person failed", "error", err, "id", input.ID)
		return nil, huma.Error500InternalServerError("update persons failed")
	}

	return &struct{}{}, nil
}
