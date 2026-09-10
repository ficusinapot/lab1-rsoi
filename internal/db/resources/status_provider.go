package resources

import (
	"context"
	"database/sql"
	"log/slog"
	"sync/atomic"

	"github.com/ficusinapot/ds/internal/models/dbifc"
)

type StatusProvider struct {
	db        *sql.DB
	logger    *slog.Logger
	available atomic.Bool
}

func NewStatusProvider(db *sql.DB) *StatusProvider {
	provider := &StatusProvider{
		db:     db,
		logger: slog.Default().With("subsystem", "db", "resource", "status_provider"),
	}
	provider.available.Store(false)

	return provider
}

func (p *StatusProvider) GetDBStatus(ctx context.Context) dbifc.ConnectionStatus {
	_ = ctx

	if p.IsConnectionAvailable(ctx) {
		return dbifc.Connected
	}

	return dbifc.Disconnected
}

func (p *StatusProvider) IsConnectionAvailable(ctx context.Context) bool {
	_ = ctx

	return p.available.Load()
}

func (p *StatusProvider) Check(ctx context.Context) bool {
	if err := p.db.PingContext(ctx); err != nil {
		p.logger.WarnContext(ctx, "ping database", "error", err)
		p.available.Store(false)
		return false
	}

	p.available.Store(true)
	return true
}
