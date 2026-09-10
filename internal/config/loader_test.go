package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadParsesConfig(t *testing.T) {
	t.Parallel()

	path := writeConfig(t, `
rest:
  service:
    addr: ":8080"
    timeout:
      shutdown: "30s"
    openapi:
      ui: true
  metrics:
    enabled: true
    addr: "0.0.0.0:11190"
database:
  dsn: "host=localhost user=program password=test dbname=persons port=5432 sslmode=disable"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "1h"
  health_check:
    check_period: "10s"
logging:
  level: Debug
  stdout:
    level: Debug
  files:
    - path: "info.log"
      level: Info
      max_file_size: 32MB
      max_files_count: 20
      max_file_age_in_days: 15
  opentelemetry:
    service_name: "aboba"
    address: "192.168.31.244:4317"
    queue_size: 2048
    level: Info
`)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if !cfg.Rest.Metrics.Enabled {
		t.Fatalf("metrics must be enabled")
	}
	if cfg.Rest.Metrics.Addr != "0.0.0.0:11190" {
		t.Fatalf("unexpected metrics address: %q", cfg.Rest.Metrics.Addr)
	}
	if !cfg.Rest.Service.OpenAPI.EnableSwaggerUI {
		t.Fatalf("OpenAPI UI must be enabled")
	}
	if cfg.Logging.Files[0].MaxFileSize.Megabytes() != 32 {
		t.Fatalf("unexpected max file size: %d", cfg.Logging.Files[0].MaxFileSize.Megabytes())
	}
	if cfg.Logging.OpenTelemetry.ServiceName != "aboba" {
		t.Fatalf("unexpected service name: %q", cfg.Logging.OpenTelemetry.ServiceName)
	}
}

func TestLoadRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	path := writeConfig(t, `
rest:
  service:
    addr: ":8080"
    timeout:
      shutdown: "30s"
    openapi:
      ui: false
  metrics:
    enabled: false
database:
  dsn: "postgres://program:test@localhost:5432/persons?sslmode=disable"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "1h"
  health_check:
    check_period: "10s"
logging:
  level: Info
  stdout:
    level: Info
  unknown: true
`)

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "decode config") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return path
}
