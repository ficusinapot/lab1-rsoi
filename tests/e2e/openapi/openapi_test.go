//go:build e2e

package openapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ficusinapot/ds/tests/e2e/infrastructure"

	"github.com/stretchr/testify/require"
)

func TestOpenAPIAndSwagger(t *testing.T) {
	service := infrastructure.NewService(t)
	client := http.Client{Timeout: 5 * time.Second}

	openAPI := infrastructure.DoJSON[openAPIResponse](t, client, http.MethodGet, service.BaseURL+"/api/v1/openapi.json", nil, http.StatusOK)
	require.Equal(t, "3.1.0", openAPI.OpenAPI)
	require.Contains(t, openAPI.Paths, "/persons")
	require.Contains(t, openAPI.Paths, "/status")

	docsResponse := infrastructure.Do(t, client, http.MethodGet, service.BaseURL+"/api/v1/docs", nil, http.StatusOK)
	defer docsResponse.Body.Close()
	docsBody, err := io.ReadAll(docsResponse.Body)
	require.NoError(t, err)
	docs := string(docsBody)
	require.Contains(t, strings.ToLower(docs), "swagger")
	require.Contains(t, docs, `data-url="/api/v1/openapi.json"`)

	swaggerClient := http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	swaggerResponse := infrastructure.Do(t, swaggerClient, http.MethodGet, service.BaseURL+"/api/v1/swagger", nil, http.StatusFound)
	defer swaggerResponse.Body.Close()
	require.Equal(t, "/api/v1/docs", swaggerResponse.Header.Get("Location"))
}

type openAPIResponse struct {
	OpenAPI string                     `json:"openapi"`
	Paths   map[string]json.RawMessage `json:"paths"`
}
