package config

import (
	"fmt"
	"strings"

	"github.com/docker/go-units"
	"github.com/ficusinapot/ds/internal/db"
	"github.com/ficusinapot/ds/internal/rest"
	"github.com/ficusinapot/ds/internal/rest/restsvc"
)

type Config struct {
	Rest     rest.Config   `mapstructure:"rest"`
	Database db.Config     `mapstructure:"database"`
	Logging  LoggingConfig `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level         string              `mapstructure:"level"`
	Stdout        LoggingStdoutConfig `mapstructure:"stdout"`
	Files         []LoggingFileConfig `mapstructure:"files"`
	OpenTelemetry OpenTelemetryConfig `mapstructure:"opentelemetry"`
}

type LoggingStdoutConfig struct {
	Level string `mapstructure:"level"`
}

type LoggingFileConfig struct {
	Path             string   `mapstructure:"path"`
	Level            string   `mapstructure:"level"`
	MaxFileSize      FileSize `mapstructure:"max_file_size"`
	MaxFilesCount    int      `mapstructure:"max_files_count"`
	MaxFileAgeInDays int      `mapstructure:"max_file_age_in_days"`
}

type OpenTelemetryConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Address     string `mapstructure:"address"`
	QueueSize   int    `mapstructure:"queue_size"`
	Level       string `mapstructure:"level"`
}

type FileSize int

const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
)

func (s *FileSize) UnmarshalText(text []byte) error {
	bytes, err := units.RAMInBytes(string(text))
	if err != nil {
		return fmt.Errorf("parse file size: %w", err)
	}

	*s = FileSize(bytes / units.MiB)

	return nil
}

func (s FileSize) Megabytes() int {
	return int(s)
}

func CreateDefaultConfig() Config {
	return Config{
		Rest:     rest.CreateDefaultConfig(),
		Database: db.CreateDefaultConfig(),
		Logging:  DefaultLoggingConfig(),
	}
}

func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level: logLevelDebug,
		Stdout: LoggingStdoutConfig{
			Level: logLevelDebug,
		},
		Files: []LoggingFileConfig{
			{
				Path:             "logs/info.log",
				Level:            logLevelInfo,
				MaxFileSize:      32,
				MaxFilesCount:    20,
				MaxFileAgeInDays: 15,
			},
			{
				Path:             "logs/error.log",
				Level:            "error",
				MaxFileSize:      16,
				MaxFilesCount:    100,
				MaxFileAgeInDays: 90,
			},
			{
				Path:             "logs/debug.log",
				Level:            logLevelDebug,
				MaxFileSize:      16,
				MaxFilesCount:    100,
				MaxFileAgeInDays: 90,
			},
		},
		OpenTelemetry: OpenTelemetryConfig{
			QueueSize: 2048,
			Level:     logLevelInfo,
		},
	}
}

func (c Config) Validate() error {
	if err := validateRestConfig(c.Rest); err != nil {
		return fmt.Errorf("rest: %w", err)
	}

	if err := validateDatabaseConfig(c.Database); err != nil {
		return fmt.Errorf("database: %w", err)
	}

	if err := c.Logging.Validate(); err != nil {
		return fmt.Errorf("logging: %w", err)
	}

	return nil
}

func validateRestConfig(c rest.Config) error {
	if err := validateHTTPConfig(c.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}

	if err := validateMetricsConfig(c.Metrics); err != nil {
		return fmt.Errorf("metrics: %w", err)
	}

	return nil
}

func validateHTTPConfig(c restsvc.ServiceConfig) error {
	if c.Addr == "" {
		return fmt.Errorf("addr is required")
	}

	if c.Timeout.Shutdown < 0 {
		return fmt.Errorf("timeout.shutdown must not be negative")
	}

	if c.Timeout.Read < 0 {
		return fmt.Errorf("timeout.read must not be negative")
	}

	if c.Timeout.Write < 0 {
		return fmt.Errorf("timeout.write must not be negative")
	}

	if c.Timeout.Idle < 0 {
		return fmt.Errorf("timeout.idle must not be negative")
	}

	return nil
}

func validateMetricsConfig(c restsvc.MetricsConfig) error {
	if c.Enabled && c.Addr == "" {
		return fmt.Errorf("addr is required")
	}

	return nil
}

func validateDatabaseConfig(c db.Config) error {
	if c.DSN == "" {
		return fmt.Errorf("dsn is required")
	}

	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("max_open_conns must be greater than zero")
	}

	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must not be negative")
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("max_idle_conns must not exceed max_open_conns")
	}

	if c.ConnMaxLifetime < 0 {
		return fmt.Errorf("conn_max_lifetime must not be negative")
	}

	if c.HealthCheck.CheckPeriod <= 0 {
		return fmt.Errorf("health_check.check_period must be greater than zero")
	}

	return nil
}

func (c LoggingConfig) Validate() error {
	if c.Level == "" {
		return fmt.Errorf("level is required")
	}

	if !isLogLevel(c.Level) {
		return fmt.Errorf("unsupported level: %s", c.Level)
	}

	if err := c.Stdout.Validate(); err != nil {
		return fmt.Errorf("stdout: %w", err)
	}

	for index, file := range c.Files {
		if err := file.Validate(); err != nil {
			return fmt.Errorf("files[%d]: %w", index, err)
		}
	}

	if err := c.OpenTelemetry.Validate(); err != nil {
		return fmt.Errorf("opentelemetry: %w", err)
	}

	return nil
}

func (c LoggingStdoutConfig) Validate() error {
	if c.Level != "" && !isLogLevel(c.Level) {
		return fmt.Errorf("unsupported level: %s", c.Level)
	}

	return nil
}

func (c LoggingFileConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("path is required")
	}

	if c.Level == "" {
		return fmt.Errorf("level is required")
	}

	if !isLogLevel(c.Level) {
		return fmt.Errorf("unsupported level: %s", c.Level)
	}

	if c.MaxFileSize <= 0 {
		return fmt.Errorf("max_file_size must be greater than zero")
	}

	if c.MaxFilesCount < 0 {
		return fmt.Errorf("max_files_count must not be negative")
	}

	if c.MaxFileAgeInDays < 0 {
		return fmt.Errorf("max_file_age_in_days must not be negative")
	}

	return nil
}

func (c OpenTelemetryConfig) Validate() error {
	if c.Address == "" {
		return nil
	}

	if c.QueueSize < 0 {
		return fmt.Errorf("queue_size must not be negative")
	}

	if c.Level != "" && !isLogLevel(c.Level) {
		return fmt.Errorf("unsupported level: %s", c.Level)
	}

	return nil
}

func isLogLevel(level string) bool {
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}
