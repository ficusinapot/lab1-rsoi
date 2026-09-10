package status

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Get(ctx context.Context, input *struct{}) (*GetStatusOutput, error) {
	_ = input

	statusData, err := h.statusUseCase.GetStatus(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "get status failed", "error", err)
		return nil, huma.Error500InternalServerError("get status failed")
	}

	return &GetStatusOutput{Body: statusDataToResponse(statusData)}, nil
}
