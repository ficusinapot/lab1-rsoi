package persons

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const (
	path   = "/persons"
	idPath = "/persons/{id}"
)

var PersonsTag = huma.Tag{
	Name:        "Persons",
	Description: "Person CRUD operations.",
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "list-persons",
		Method:      http.MethodGet,
		Path:        path,
		Tags:        []string{PersonsTag.Name},
	}, h.List)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID: "get-persons",
		Method:      http.MethodGet,
		Path:        idPath,
		Tags:        []string{PersonsTag.Name},
	}, h.Get)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID:   "create-persons",
		Method:        http.MethodPost,
		Path:          path,
		DefaultStatus: http.StatusCreated,
		Tags:          []string{PersonsTag.Name},
	}, h.Create)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID:   "update-persons",
		Method:        http.MethodPatch,
		Path:          idPath,
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{PersonsTag.Name},
	}, h.Update)
	huma.Register(api, huma.Operation{ //nolint:exhaustruct_v5
		OperationID:   "delete-persons",
		Method:        http.MethodDelete,
		Path:          idPath,
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{PersonsTag.Name},
	}, h.Delete)
}
