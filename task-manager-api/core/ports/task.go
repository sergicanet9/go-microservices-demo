package ports

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/scv-go-tools/v4/repository"
)

// TaskRepository interface
type TaskRepository interface {
	repository.Repository
}

// TaskService interface
type TaskService interface {
	Create(ctx context.Context, userID string, task models.CreateTaskReq) (models.CreateTaskResp, error)
	GetByUserID(ctx context.Context, userID string) ([]models.GetTaskResp, error)
	Delete(ctx context.Context, userID, taskID string) error
}
