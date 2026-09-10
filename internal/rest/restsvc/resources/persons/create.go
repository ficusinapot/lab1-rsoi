package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Create(ctx context.Context, input *Input) (*Output, error) {
	person := requestToModel(input.Body)
	id, err := h.personUseCase.CreatePerson(ctx, person)
	if err != nil {
		h.logger.ErrorContext(ctx, "create person failed", "error", err)
		return nil, huma.Error500InternalServerError("create persons failed")
	}

	person.ID = id
	return &Output{Body: modelToResponse(person)}, nil
}
