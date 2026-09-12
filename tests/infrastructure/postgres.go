//go:build integration || e2e

package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	appdb "github.com/ficusinapot/ds/internal/db"
	dbent "github.com/ficusinapot/ds/internal/db/ent"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresImage    = "postgres:18-alpine"
	postgresDatabase = "persons"
	postgresUser     = "postgres"
	postgresPassword = "postgres"
)

type Postgres struct {
	DSN string
	DB  *dbent.Client
}

func NewPostgres(t *testing.T) Postgres {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	t.Cleanup(cancel)

	container, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.WithDatabase(postgresDatabase),
		tcpostgres.WithUsername(postgresUser),
		tcpostgres.WithPassword(postgresPassword),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp").WithStartupTimeout(time.Minute)),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(context.Background()))
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbClient, err := appdb.Open(testDatabaseConfig(dsn))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, appdb.Close(dbClient))
	})

	require.NoError(t, applyMigrations(ctx, dsn, filepath.Join(projectRoot(t), "migrations")))

	return Postgres{
		DSN: dsn,
		DB:  dbClient.Ent(),
	}
}

func projectRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func applyMigrations(ctx context.Context, dsn, migrationsDir string) error {
	cmd := exec.CommandContext(
		ctx,
		"atlas",
		"migrate",
		"apply",
		"--dir",
		"file://"+migrationsDir,
		"--url",
		dsn,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apply atlas migrations: %w: %s", err, stderr.String())
	}

	return nil
}
