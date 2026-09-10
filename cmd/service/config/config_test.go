package config

import (
	"testing"

	internalconfig "github.com/ficusinapot/ds/internal/config"
)

func TestApplyEnvOverridesUsesRenderPort(t *testing.T) {
	t.Setenv(EnvPortKey, "10000")
	t.Setenv(EnvDatabaseURLKey, "")

	cfg := internalconfig.CreateDefaultConfig()
	cfg.Rest.Service.Addr = ":8080"

	applyEnvOverrides(&cfg)

	if cfg.Rest.Service.Addr != "0.0.0.0:10000" {
		t.Fatalf("unexpected REST address: %q", cfg.Rest.Service.Addr)
	}
}

func TestApplyEnvOverridesUsesDatabaseURL(t *testing.T) {
	const dsn = "postgresql://host:5432/persons"

	t.Setenv(EnvPortKey, "")
	t.Setenv(EnvDatabaseURLKey, dsn)

	cfg := internalconfig.CreateDefaultConfig()
	cfg.Database.DSN = "postgres://program:test@localhost:5432/persons?sslmode=disable"

	applyEnvOverrides(&cfg)

	if cfg.Database.DSN != dsn {
		t.Fatalf("unexpected database DSN: %q", cfg.Database.DSN)
	}
}
