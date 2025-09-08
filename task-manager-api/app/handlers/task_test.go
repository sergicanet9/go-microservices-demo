package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v4/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
)

// TestCreateTask_Ok checks that CreateTask handler returns the expected response when a valid request is received
func TestCreateTask_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	expectedResp := models.CreateTaskResp{ID: "new-task-id"}
	taskService.On(testutils.FunctionName(t, ports.TaskService.Create), mock.Anything, "user-123", mock.Anything).Return(expectedResp, nil).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	body := models.CreateTaskReq{Title: "test task"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.CreateTaskResp
	_ = json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, expectedResp, response)
}

// TestCreateTask_InvalidRequest checks that CreateTask handler returns an error when the received request is not valid
func TestCreateTask_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, nil)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	invalidBody := []byte(`{"Email":invalid-type}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(invalidBody))
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestCreateTask_ServiceError checks that CreateTask handler returns an error when the service fails
func TestCreateTask_ServiceError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	taskService.On(testutils.FunctionName(t, ports.TaskService.Create), mock.Anything, "user-123", mock.Anything).Return(models.CreateTaskResp{}, errors.New("service failure")).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	body := models.CreateTaskReq{Title: "bad task"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestGetTasks_Ok checks that GetTasks handler returns the expected response when a valid request is received
func TestGetTasks_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	expectedResp := []models.GetTaskResp{{ID: "task-123", Title: "test task"}}
	taskService.On(testutils.FunctionName(t, ports.TaskService.GetByUserID), mock.Anything, "user-123").Return(expectedResp, nil).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.GetTaskResp
	_ = json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, expectedResp, response)
}

// TestGetTasks_ServiceError checks that GetTasks handler returns an error when the service fails
func TestGetTasks_ServiceError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	taskService.On(testutils.FunctionName(t, ports.TaskService.GetByUserID), mock.Anything, "user-123").Return([]models.GetTaskResp{}, errors.New("service failure")).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestDeleteTask_Ok checks that DeleteTask handler returns the expected response when a valid request is received
func TestDeleteTask_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	taskService.On(testutils.FunctionName(t, ports.TaskService.Delete), mock.Anything, "user-123", "task-123").Return(nil).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/task-123", nil)
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestDeleteTask_ServiceError checks that DeleteTask handler returns an error when the service fails
func TestDeleteTask_ServiceError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	taskService := mocks.NewTaskService(t)
	taskService.On(testutils.FunctionName(t, ports.TaskService.Delete), mock.Anything, "user-123", "task-123").Return(errors.New("service failure")).Once()

	cfg := config.Config{}
	taskHandler := NewTaskHandler(context.Background(), cfg, taskService)
	SetTaskRoutes(r, taskHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/task-123", nil)
	claims := jwt.MapClaims{"user_id": "user-123"}
	req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, claims))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestGetUserID checks that the getUserID function handles scenarios as expected
func TestGetUserID(t *testing.T) {
	tests := []struct {
		name       string
		claims     interface{}
		wantUserID string
		wantErr    string
	}{
		{"ok", jwt.MapClaims{"user_id": "user-123"}, "user-123", ""},
		{"claims missing", nil, "", "claims not found in context"},
		{"claims wrong type", "not-a-map", "", "invalid claims type"},
		{"user_id missing", jwt.MapClaims{"other": "value"}, "", "user_id not found in claims"},
		{"user_id wrong type", jwt.MapClaims{"user_id": 12345}, "", "invalid user_id type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.claims != nil {
				req = req.WithContext(context.WithValue(req.Context(), middlewares.ClaimsKey, tt.claims))
			}

			userID, err := getUserID(req)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
				assert.Equal(t, "", userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUserID, userID)
			}
		})
	}
}
