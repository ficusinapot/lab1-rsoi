package config

import (
	"fmt"

	internalconfig "github.com/ficusinapot/ds/internal/config"
)

const (
	DefaultConfigPath = "./configs/config.yaml"
	EnvConfigPathKey  = "RSOI_CONFIG_PATH"
)

func LoadConfig(path string) (*internalconfig.Config, error) {
	cfg, err := internalconfig.LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load internal config: %w", err)
	}

	return &cfg, nil
}
