//go:build e2e

package metrics_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ficusinapot/ds/tests/e2e/infrastructure"

	"github.com/stretchr/testify/require"
)

func TestMetricsEndpoint(t *testing.T) {
	service := infrastructure.NewService(t)
	client := http.Client{Timeout: 5 * time.Second}

	infrastructure.DoNoBody(t, client, http.MethodGet, service.BaseURL+"/api/v1/persons", nil, http.StatusOK)
	infrastructure.DoNoBody(t, client, http.MethodGet, service.BaseURL+"/api/v1/status", nil, http.StatusOK)
	infrastructure.DoNoBody(t, client, http.MethodGet, service.BaseURL+"/api/v1/healthz", nil, http.StatusOK)
	infrastructure.DoNoBody(t, client, http.MethodGet, service.BaseURL+"/api/v1/readyz", nil, http.StatusOK)

	response, err := client.Get(service.MetricsURL + "/metrics")
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			require.NoError(t, err)
		}
	}(response.Body)
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	metrics := string(body)
	require.Contains(t, metrics, "service_rest_api_http_requests_total")
	require.Contains(t, metrics, "service_db_connections_open")
	require.Contains(t, metrics, "db_connections_alive")
	require.Contains(t, metrics, "rest_service_starts_total")
	require.Contains(t, metrics, "domain_persons_list_request_total")
	require.Contains(t, metrics, "status_get_status_request_total")
	require.Contains(t, metrics, "status_healthz_request_total")
	require.Contains(t, metrics, "status_readyz_request_total")
}
