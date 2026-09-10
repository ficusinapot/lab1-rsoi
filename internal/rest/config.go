package rest

import "github.com/ficusinapot/ds/internal/rest/restsvc"

type Config struct {
	Service restsvc.ServiceConfig `mapstructure:"service"`
	Metrics restsvc.MetricsConfig `mapstructure:"metrics"`
}

func CreateDefaultConfig() Config {
	return Config{
		Service: restsvc.CreateDefaultServiceConfig(),
		Metrics: restsvc.CreateDefaultMetricsConfig(),
	}
}
