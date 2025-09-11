package ports

import (
	"context"
)

// TaskManagerV1HTTPClient interface for a Task Manager API v1 HTTP Client
type TaskManagerV1HTTPClient interface {
	Health(ctx context.Context) error
}
