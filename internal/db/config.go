package db

import "time"

type Config struct {
	DSN             string            `mapstructure:"dsn"`
	MaxOpenConns    int               `mapstructure:"max_open_conns"`
	MaxIdleConns    int               `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration     `mapstructure:"conn_max_lifetime"`
	HealthCheck     HealthCheckConfig `mapstructure:"health_check"`
}

type HealthCheckConfig struct {
	CheckPeriod time.Duration `mapstructure:"check_period"`
}

func CreateDefaultConfig() Config {
	return Config{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		HealthCheck: HealthCheckConfig{
			CheckPeriod: 10 * time.Second,
		},
	}
}
