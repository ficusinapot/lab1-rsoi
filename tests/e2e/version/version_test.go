//go:build e2e

package version_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ficusinapot/ds/tests/e2e/infrastructure"

	"github.com/stretchr/testify/require"
)

func TestVersionAPI(t *testing.T) {
	service := infrastructure.NewService(t)
	client := http.Client{Timeout: 5 * time.Second}

	version := infrastructure.DoJSON[versionResponse](t, client, http.MethodGet, service.BaseURL+"/api/v1/version", nil, http.StatusOK)

	require.Equal(t, "0.1.0", version.Version.Version)
}

type versionResponse struct {
	Version versionData `json:"data"`
}

type versionData struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
}
