package status

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var StatusTag = huma.Tag{
	Name:        "Status",
	Description: "Service status, liveness, and readiness operations.",
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        "/status",
		Tags:        []string{StatusTag.Name},
	}, h.Get)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "healthz",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Tags:        []string{StatusTag.Name},
	}, h.Healthz)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "readyz",
		Method:      http.MethodGet,
		Path:        "/readyz",
		Tags:        []string{StatusTag.Name},
	}, h.Readyz)
}
