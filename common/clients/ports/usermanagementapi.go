package ports

import (
	"context"
)

// UserManagementV1GRPCClient interface for a User Management API v1 gRPC Client
type UserManagementV1GRPCClient interface {
	Close() error
	Health(ctx context.Context) error
	Exists(ctx context.Context, token, userID string) (bool, error)
}
