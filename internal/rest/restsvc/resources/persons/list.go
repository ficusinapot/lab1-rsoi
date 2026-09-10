package persons

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) List(ctx context.Context, input *struct{}) (*ListOutput, error) {
	persons, err := h.personUseCase.ListPersons(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "list persons failed", "error", err)
		return nil, huma.Error500InternalServerError("list persons failed")
	}

	return &ListOutput{Body: modelsToResponses(persons)}, nil
}
