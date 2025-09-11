package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealth_Ok checks that the HTTP client handles a successful call to the Health endpoint as expected
func TestHealth_Ok(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/health", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.Health(context.Background())

	// Assert
	assert.NoError(t, err)
}

// TestHealth_HTTPError checks that the HTTP client handles an unsuccessful call to the Health endpoint as expected
func TestHealth_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.Health(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected http status code")
}

// TestCreateTask_Ok checks that the HTTP client handles a successful call to the Create endpoint as expected
func TestCreateTask_Ok(t *testing.T) {
	// Arrange
	expectedResp := models.CreateTaskResp{ID: "123"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(expectedResp)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	resp, err := client.CreateTask(context.Background(), models.CreateTaskReq{})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

// TestCreateTask_HTTPError checks that the HTTP client handles an unsuccessful call to the Create endpoint as expected
func TestCreateTask_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	_, err := client.CreateTask(context.Background(), models.CreateTaskReq{})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected http status code")
}

// TestGetTasks_Ok checks that the HTTP client handles a successful call to the GetTasks endpoint as expected
func TestGetTasks_Ok(t *testing.T) {
	// Arrange
	expectedResp := []models.GetTaskResp{{ID: "123", Title: "Test task"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResp)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	resp, err := client.GetTasks(context.Background(), "Bearer test-token")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

// TestGetTasks_HTTPError checks that the HTTP client handles an unsuccessful call to the GetTasks endpoint as expected
func TestGetTasks_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	_, err := client.GetTasks(context.Background(), "Bearer test-token")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected http status code")
}

// TestDeleteTask_Ok checks that the HTTP client handles a successful call to the DeleteTask endpoint as expected
func TestDeleteTask_Ok(t *testing.T) {
	// Arrange
	taskID := "123"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, fmt.Sprintf("/tasks/%s", taskID), r.URL.Path)
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.DeleteTask(context.Background(), "Bearer test-token", taskID)

	// Assert
	assert.NoError(t, err)
}

// TestDeleteTask_HTTPError checks that the HTTP client handles an unsuccessful call to the DeleteTask endpoint as expected
func TestDeleteTask_HTTPError(t *testing.T) {
	// Arrange
	taskID := "123"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Act
	err := client.DeleteTask(context.Background(), "Bearer test-token", taskID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected http status code")
}
