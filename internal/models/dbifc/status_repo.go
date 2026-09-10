package dbifc

import "context"

type ConnectionStatus string

const (
	Connected    ConnectionStatus = "connected"
	Disconnected ConnectionStatus = "disconnected"
)

type StatusProvider interface {
	GetDBStatus(ctx context.Context) ConnectionStatus
	IsConnectionAvailable(ctx context.Context) bool
}
