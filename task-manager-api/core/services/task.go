package services

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
)

// taskService adapter of an task service
type taskService struct {
	config config.Config
}

// NewTaskService creates a new task service
func NewTaskService(cfg config.Config) ports.TaskService {
	return &taskService{
		config: cfg,
	}
}

// Create task
func (t *taskService) Create(ctx context.Context, userID string, task models.CreateTaskReq) (models.CreateTaskResp, error) {
	panic("unimplemented")
}

// GetByUserID tasks
func (t *taskService) GetByUserID(ctx context.Context, userID string) ([]models.GetTaskResp, error) {
	panic("unimplemented")
}

// Delete task
func (t *taskService) Delete(ctx context.Context, userID string, taskID string) error {
	panic("unimplemented")
}
