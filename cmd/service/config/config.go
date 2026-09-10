package config

import (
	"fmt"
	"os"

	internalconfig "github.com/ficusinapot/ds/internal/config"
)

const (
	DefaultConfigPath = "./configs/config.yaml"
	EnvConfigPathKey  = "RSOI_CONFIG_PATH"
	EnvPortKey        = "PORT"
	EnvDatabaseURLKey = "DATABASE_URL"
)

func LoadConfig(path string) (*internalconfig.Config, error) {
	cfg, err := internalconfig.LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load internal config: %w", err)
	}

	applyEnvOverrides(&cfg)

	return &cfg, nil
}

func applyEnvOverrides(cfg *internalconfig.Config) {
	if port := os.Getenv(EnvPortKey); port != "" {
		cfg.Rest.Service.Addr = "0.0.0.0:" + port
	}

	if databaseURL := os.Getenv(EnvDatabaseURLKey); databaseURL != "" {
		cfg.Database.DSN = databaseURL
	}
}
