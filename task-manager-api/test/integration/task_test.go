package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/entities"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestCreateTask_Ok checks that CreateTask endpoint returns the expected response when everything goes as expected
func TestCreateTask_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testTask := getNewTestTask()

	// Act
	body := models.CreateTaskReq(testTask)
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://:%d/task-manager-api/v1/tasks", cfg.HTTPPort)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusCreated, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response models.CreateTaskResp
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.NotEmpty(t, response.ID)
	createdTask, err := findTask(response.ID, cfg)
	if err != nil {
		t.Fatalf("unexpected error while finding the created task: %s", err)
	}
	assert.Equal(t, testTask.UserID, createdTask.UserID)
	assert.Equal(t, testTask.Title, createdTask.Title)
	assert.Equal(t, testTask.Description, createdTask.Description)
	assert.NotNil(t, createdTask.CreatedAt)
}

// TestGetTasks_Ok checks that GetTasks endpoint returns the expected response when everything goes as expected
func TestGetTasks_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)

	// Act
	url := fmt.Sprintf("http://:%d/task-manager-api/v1/tasks", cfg.HTTPPort)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response []models.GetTaskResp
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.NotEmpty(t, response)
}

// TestDeleteTask_Ok checks that DeleteTask endpoint returns the expected response when everything goes as expected
func TestDeleteUser_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testTask := getNewTestTask()
	err := insertTask(&testTask, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	url := fmt.Sprintf("http://:%d/task-manager-api/v1/tasks/%s", cfg.HTTPPort, testTask.ID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	_, err = findTask(testTask.ID, cfg)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

// HELP FUNCTIONS

func getNewTestTask() entities.Task {
	return entities.Task{
		UserID:      tokenUserID,
		Title:       "test-title",
		Description: "test-description",
	}
}

func insertTask(t *entities.Task, cfg config.Config) error {
	now := time.Now().UTC()
	db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
	if err != nil {
		return err
	}
	t.CreatedAt = now
	result, err := db.Collection(entities.EntityNameTask).InsertOne(context.Background(), t)
	t.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func findTask(ID string, cfg config.Config) (entities.Task, error) {
	db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
	if err != nil {
		return entities.Task{}, err
	}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return entities.Task{}, err
	}

	var u entities.Task
	err = db.Collection(entities.EntityNameTask).FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
	return u, err

}
