//go:build e2e

package status_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ficusinapot/ds/tests/e2e/infrastructure"

	"github.com/stretchr/testify/require"
)

func TestStatusAPI(t *testing.T) {
	service := infrastructure.NewService(t)
	client := http.Client{Timeout: 5 * time.Second}

	status := infrastructure.DoJSON[serviceStatusResponse](t, client, http.MethodGet, service.BaseURL+"/api/v1/status", nil, http.StatusOK)
	require.True(t, status.DBAvailable)
	require.Equal(t, "FullyOperational", status.OperatingStatus)

	healthz := infrastructure.DoJSON[statusResponse](t, client, http.MethodGet, service.BaseURL+"/api/v1/healthz", nil, http.StatusOK)
	require.Equal(t, "ok", healthz.Status)

	readyz := infrastructure.DoJSON[statusResponse](t, client, http.MethodGet, service.BaseURL+"/api/v1/readyz", nil, http.StatusOK)
	require.Equal(t, "ok", readyz.Status)
}

type serviceStatusResponse struct {
	DBAvailable     bool   `json:"db_available"`
	OperatingStatus string `json:"operating_status"`
}

type statusResponse struct {
	Status string `json:"status"`
}
