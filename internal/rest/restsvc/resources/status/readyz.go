package status

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Readyz(ctx context.Context, input *struct{}) (*StatusOutput, error) {
	_ = input

	if err := h.statusUseCase.Readyz(ctx); err != nil {
		h.logger.ErrorContext(ctx, "readyz failed", "error", err)
		return nil, huma.Error503ServiceUnavailable("service is not ready")
	}

	return &StatusOutput{Body: StatusResponse{Status: "ok"}}, nil
}
