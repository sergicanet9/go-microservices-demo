package ports

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
)

// TaskManagerV1HTTPClient interface for a Task Manager API v1 HTTP Client
type TaskManagerV1HTTPClient interface {
	Health(ctx context.Context) error
	CreateTask(ctx context.Context, task models.CreateTaskReq) (models.CreateTaskResp, error)
	GetTasks(ctx context.Context, userID string) ([]models.GetTaskResp, error)
	DeleteTask(ctx context.Context, userID, taskID string) error
}
