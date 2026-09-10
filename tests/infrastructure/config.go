//go:build integration || e2e

package infrastructure

import (
	"time"

	"github.com/ficusinapot/ds/internal/db"
)

func testDatabaseConfig(dsn string) db.Config {
	return db.Config{
		DSN:             dsn,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		HealthCheck: db.HealthCheckConfig{
			CheckPeriod: time.Second,
		},
	}
}
