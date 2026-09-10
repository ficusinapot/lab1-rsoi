package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Create(ctx context.Context, input *Input) (*Output, error) {
	if input.Body.Name == nil || input.Body.Age == nil || input.Body.Address == nil || input.Body.Work == nil {
		return nil, huma.Error400BadRequest("name, age, address and work are required")
	}

	person := requestToModel(input.Body)
	id, err := h.personUseCase.CreatePerson(ctx, person)
	if err != nil {
		h.logger.ErrorContext(ctx, "create person failed", "error", err)
		return nil, huma.Error500InternalServerError("create persons failed")
	}

	person.ID = id
	return &Output{Body: modelToResponse(person)}, nil
}
