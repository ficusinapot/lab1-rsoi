package restsvc

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ficusinapot/ds/internal/rest/restsvc/middleware"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/persons"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/status"
	"github.com/ficusinapot/ds/internal/rest/restsvc/resources/version"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type AppInfo struct {
	Name      string
	Version   string
	BuildTime string
	Branch    string
	Commit    string
}

const (
	apiPrefix = "/api/v1"
)

func init() {
	defaultNewError := huma.NewError
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		statusError := defaultNewError(status, msg, errs...)
		var errorModel *huma.ErrorModel
		ok := errors.As(statusError, &errorModel)
		if !ok {
			return statusError
		}

		title, detail, ok := splitErrorMessage(msg)
		if !ok {
			return errorModel
		}

		errorModel.Title = title
		errorModel.Detail = detail

		return errorModel
	}
}

func NewRouter(
	app AppInfo,
	openAPIConfig OpenAPIConfig,
	personHandler *persons.Handler,
	statusHandler *status.Handler,
	logger *slog.Logger,
	metrics middleware.HTTPMetrics,
) http.Handler {
	router := chi.NewRouter()

	router.Route(apiPrefix, func(router chi.Router) {
		api := humachi.New(router, newHumaConfig(app, openAPIConfig))
		statusHandler.RegisterRoutes(api)
		personHandler.RegisterRoutes(api)
		version.NewHandler(app.Version, app.BuildTime, app.Branch, app.Commit).RegisterRoutes(api)

		if openAPIConfig.EnableSwaggerUI {
			router.Get("/swagger", func(response http.ResponseWriter, request *http.Request) {
				http.Redirect(response, request, apiPrefix+"/docs", http.StatusFound)
			})
		}
	})

	return otelhttp.NewHandler(
		middleware.Observability(logger, metrics)(router),
		"http.server",
	)
}

func splitErrorMessage(message string) (string, string, bool) {
	title, detail, found := strings.Cut(message, ":")
	if !found {
		return "", "", false
	}

	title = strings.TrimSpace(title)
	detail = strings.TrimSpace(detail)
	if title == "" || detail == "" {
		return "", "", false
	}

	return title, detail, true
}

func newHumaConfig(app AppInfo, openAPIConfig OpenAPIConfig) huma.Config {
	config := huma.DefaultConfig(app.Name, app.Version)
	config.OpenAPIPath = "/openapi"
	config.Servers = []*huma.Server{
		{URL: apiPrefix},
	}
	config.Tags = []*huma.Tag{
		&status.StatusTag,
		&persons.PersonsTag,
		&version.VersionTag,
	}
	if openAPIConfig.EnableSwaggerUI {
		config.DocsPath = "/docs"
		config.DocsRenderer = huma.DocsRendererSwaggerUI
	} else {
		config.DocsPath = ""
	}

	return config
}
