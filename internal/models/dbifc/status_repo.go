package dbifc

import "context"

//go:generate go tool mockgen -source=status_repo.go -destination=mocks/status_provider.go -package=mocks

type ConnectionStatus string

const (
	Connected    ConnectionStatus = "connected"
	Disconnected ConnectionStatus = "disconnected"
)

type StatusProvider interface {
	GetDBStatus(ctx context.Context) ConnectionStatus
	IsConnectionAvailable(ctx context.Context) bool
}
