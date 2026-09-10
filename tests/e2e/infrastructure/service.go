//go:build e2e

package infrastructure

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ficusinapot/ds/cmd/service/version"
	"github.com/ficusinapot/ds/internal/app"
	"github.com/ficusinapot/ds/internal/config"
	"github.com/ficusinapot/ds/internal/rest/restsvc"
	testinfra "github.com/ficusinapot/ds/tests/infrastructure"

	"github.com/stretchr/testify/require"
)

type Service struct {
	BaseURL    string
	MetricsURL string
}

func NewService(t *testing.T) Service {
	t.Helper()

	pg := testinfra.NewPostgres(t)
	httpAddr := freeTCPAddress(t)
	metricsAddr := freeTCPAddress(t)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- app.Run(ctx, testConfig(pg.DSN, httpAddr, metricsAddr), testAppInfo())
	}()

	service := Service{
		BaseURL:    "http://" + httpAddr,
		MetricsURL: "http://" + metricsAddr,
	}
	waitForService(t, service.BaseURL)

	t.Cleanup(func() {
		cancel()

		select {
		case err := <-errCh:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("service shutdown timeout")
		}
	})

	return service
}

func testAppInfo() restsvc.AppInfo {
	info := version.GetInfo()

	return restsvc.AppInfo{
		Name:      "service",
		Version:   info.Version,
		BuildTime: info.BuildTime,
		Branch:    info.Branch,
		Commit:    info.Commit,
	}
}

func testConfig(dsn, httpAddr, metricsAddr string) config.Config {
	cfg := config.CreateDefaultConfig()
	cfg.Rest.Service.Addr = httpAddr
	cfg.Rest.Service.OpenAPI.EnableSwaggerUI = true
	cfg.Database.DSN = dsn
	cfg.Rest.Metrics.Enabled = true
	cfg.Rest.Metrics.Addr = metricsAddr
	cfg.Logging.Stdout.Level = ""
	cfg.Logging.Files = nil
	cfg.Logging.OpenTelemetry.Address = ""

	return cfg
}

func freeTCPAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, listener.Close())
	}()

	return listener.Addr().String()
}

func waitForService(t *testing.T, baseURL string) {
	t.Helper()

	client := http.Client{Timeout: time.Second}
	deadline := time.Now().Add(10 * time.Second)
	url := fmt.Sprintf("%s/api/v1/persons", baseURL)

	for time.Now().Before(deadline) {
		response, err := client.Get(url)
		if err == nil {
			_ = response.Body.Close()
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("service did not start at %s", baseURL)
}
