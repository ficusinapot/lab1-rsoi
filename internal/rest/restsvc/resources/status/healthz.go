package status

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Healthz(ctx context.Context, input *struct{}) (*StatusOutput, error) {
	_ = input

	if err := h.statusUseCase.Healthz(ctx); err != nil {
		h.logger.ErrorContext(ctx, "healthz failed", "error", err)
		return nil, huma.Error503ServiceUnavailable("service is unhealthy")
	}

	return &StatusOutput{Body: StatusResponse{Status: "ok"}}, nil
}
