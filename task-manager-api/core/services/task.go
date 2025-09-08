package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/entities"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
)

// taskService adapter of an task service
type taskService struct {
	config     config.Config
	repository ports.TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(cfg config.Config, repo ports.TaskRepository) ports.TaskService {
	return &taskService{
		config:     cfg,
		repository: repo,
	}
}

// Create task
func (t *taskService) Create(ctx context.Context, userID string, task models.CreateTaskReq) (resp models.CreateTaskResp, err error) {
	// TODO check userID exists in user management api?
	entity := entities.Task(task)
	entity.UserID = userID

	id, err := t.repository.Create(ctx, entity)
	if err != nil {
		return
	}

	resp = models.CreateTaskResp{
		ID: id,
	}

	return
}

// GetByUserID tasks
func (t *taskService) GetByUserID(ctx context.Context, userID string) (resp []models.GetTaskResp, err error) {
	filter := map[string]interface{}{"user_id": userID}

	result, err := t.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = nil
		}
		return
	}

	resp = make([]models.GetTaskResp, len(result))
	for i, v := range result {
		resp[i] = models.GetTaskResp(*(v.(*entities.Task)))
	}

	return
}

// Delete task
func (t *taskService) Delete(ctx context.Context, userID string, taskID string) (err error) {
	result, err := t.repository.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("TaskID %s not found", taskID))
		}
		return
	}
	task := models.GetTaskResp(*result.(*entities.Task))
	if task.UserID != userID {
		err = wrappers.NewUnauthenticatedErr(fmt.Errorf("UserID %s not allowed to delete TaskID %s", userID, taskID))
		return
	}

	err = t.repository.Delete(ctx, taskID)
	return
}
