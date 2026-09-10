package version

import "context"

func (h *Handler) Get(ctx context.Context, input *struct{}) (*AppVersionOutput, error) {
	_ = ctx
	_ = input

	return &AppVersionOutput{
		Body: AppVersionBody{
			Version: h.version,
		},
	}, nil
}
