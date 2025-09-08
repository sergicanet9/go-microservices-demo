package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/entities"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreate_Ok checks that Create returns the expected response when a valid request is received
func TestCreate_Ok(t *testing.T) {
	// Arrange
	req := models.CreateTaskReq{
		Title:       "test-title",
		Description: "test-description",
	}

	expectedResponse := models.CreateTaskResp{
		ID: "new-id",
	}

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Create), context.Background(), mock.AnythingOfType("entities.Task")).Return(expectedResponse.ID, nil).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	resp, err := service.Create(context.Background(), "user-123", req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestCreate_RepositoryError checks that Create returns an error when the repository fails
func TestCreate_RepositoryError(t *testing.T) {
	// Arrange
	req := models.CreateTaskReq{
		Title:       "test-title",
		Description: "test-description",
	}

	expectedError := "repository-error"

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Create), context.Background(), mock.AnythingOfType("entities.Task")).Return("", errors.New(expectedError)).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	_, err := service.Create(context.Background(), "user-123", req)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestGetByUserID_Ok checks that GetByUserID returns the expected response when a valid request is received
func TestGetByUserID_Ok(t *testing.T) {
	// Arrange
	now := time.Now().UTC()
	userID := "test-user-id"
	var repoResponse []interface{}
	tasks := []entities.Task{
		{
			ID:          "task-id-1",
			UserID:      userID,
			Title:       "title-1",
			Description: "description-1",
			CreatedAt:   now,
		},
		{
			ID:          "task-id-2",
			UserID:      userID,
			Title:       "title-2",
			Description: "description-2",
			CreatedAt:   now,
		},
	}
	repoResponse = append(repoResponse, &tasks[0], &tasks[1])
	expectedResponse := []models.GetTaskResp{
		models.GetTaskResp(tasks[0]),
		models.GetTaskResp(tasks[1]),
	}

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Get), context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(repoResponse, nil).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	resp, err := service.GetByUserID(context.Background(), userID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestGetByUserID_NoResourcesFound checks that GetByUserID does not return an error when the repository does not return any task for the received userID
func TestGetByUserID_NoResourcesFound(t *testing.T) {
	// Arrange
	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Get), context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil, wrappers.NonExistentErr).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	resp, err := service.GetByUserID(context.Background(), "user-123")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp))
}

// TestGetByUserID_RepositoryError checks that GetByID returns an error when the repository fails
func TestGetByUserID_RepositoryError(t *testing.T) {
	// Arrange
	expectedError := "repository-error"

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Get), context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(expectedError)).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	_, err := service.GetByUserID(context.Background(), "user-123")

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestDelete_Ok checks that Delete does not return an error when a valid request is received
func TestDelete_Ok(t *testing.T) {
	// Arrange
	task := entities.Task{
		ID:          "task-id",
		UserID:      "user-123",
		Title:       "title",
		Description: "description",
		CreatedAt:   time.Now(),
	}
	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.GetByID), context.Background(), mock.Anything).Return(&task, nil).Once()
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Delete), context.Background(), mock.Anything).Return(nil).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), task.UserID, task.ID)

	// Assert
	assert.Nil(t, err)
}

// TestDelete_NotFound checks that Delete returns an error when the provided taskID does not exist
func TestDelete_NotFound(t *testing.T) {
	// Arrange
	task := entities.Task{
		ID:     "task-id",
		UserID: "user-123",
	}
	expectedError := fmt.Sprintf("TaskID %s not found", task.ID)

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.GetByID), context.Background(), mock.Anything).Return(nil, wrappers.NonExistentErr).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), task.UserID, task.ID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestDelete_Invalid checks that Delete returns an error when the provided taskID does correspond to the received userID
func TestDelete_InvalidUser(t *testing.T) {
	// Arrange
	wrongUserID := "other-user"
	task := entities.Task{
		ID:     "task-id",
		UserID: "user-123",
	}
	expectedError := fmt.Sprintf("UserID %s not allowed to delete TaskID %s", wrongUserID, task.ID)

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.GetByID), context.Background(), mock.Anything).Return(&task, nil).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), wrongUserID, task.ID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.UnauthenticatedErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestDelete_RepositoryError checks that Delete returns an error when the repository fails
func TestDelete_RepositoryError(t *testing.T) {
	// Arrange
	task := entities.Task{
		ID:          "task-id",
		UserID:      "user-123",
		Title:       "title",
		Description: "description",
		CreatedAt:   time.Now(),
	}
	expectedError := "repository-error"

	taskRepositoryMock := mocks.NewTaskRepository(t)
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.GetByID), context.Background(), mock.Anything).Return(&task, nil).Once()
	taskRepositoryMock.On(testutils.FunctionName(t, ports.TaskRepository.Delete), context.Background(), mock.Anything).Return(errors.New(expectedError)).Once()

	service := &taskService{
		config:     config.Config{},
		repository: taskRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), task.UserID, task.ID)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}
