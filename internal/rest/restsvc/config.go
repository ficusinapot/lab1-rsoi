package restsvc

import "time"

const defaultReadHeaderTimeout = 5 * time.Second

type ServiceConfig struct {
	Addr    string        `mapstructure:"addr"`
	Timeout TimeoutConfig `mapstructure:"timeout"`
	OpenAPI OpenAPIConfig `mapstructure:"openapi"`
}

type TimeoutConfig struct {
	Shutdown time.Duration `mapstructure:"shutdown"`
	Read     time.Duration `mapstructure:"read"`
	Write    time.Duration `mapstructure:"write"`
	Idle     time.Duration `mapstructure:"idle"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Addr    string `mapstructure:"addr"`
}

type OpenAPIConfig struct {
	EnableSwaggerUI bool `mapstructure:"ui"`
}

func CreateDefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Addr: ":8080",
		Timeout: TimeoutConfig{
			Shutdown: 30 * time.Second,
			Read:     0,
			Write:    0,
			Idle:     0,
		},
		OpenAPI: OpenAPIConfig{
			EnableSwaggerUI: false,
		},
	}
}

func CreateDefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled: true,
		Addr:    ":9090",
	}
}
