package version

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var VersionTag = huma.Tag{
	Name:        "Version",
	Description: "Application version and build metadata.",
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "version",
		Method:      http.MethodGet,
		Path:        "/version",
		Summary:     "Get application version",
		Tags:        []string{VersionTag.Name},
	}, h.Get)
}
