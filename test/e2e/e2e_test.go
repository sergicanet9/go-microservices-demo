package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	healthAPIAddr         = "http://localhost:82/health-api/v1"
	taskManagerAPIAddr    = "http://localhost:82/task-manager-api/v1"
	userManagementAPIAddr = "http://localhost:82/user-management-api/v1"
)

var (
	userEmail string
	userID    string
	userToken string
	taskID    string
	client    = &http.Client{Timeout: 6 * time.Second}
)

func TestMain(m *testing.M) {
	maxAttempts := 20
	if !servicesReady(maxAttempts) {
		fmt.Println("Error: services are not running. Run 'make up' first.")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func servicesReady(maxAttempts int) bool {
	for i := 1; i <= maxAttempts; i++ {
		fmt.Printf("Waiting for services to be ready... (attempt %d/%d)\n", i, maxAttempts)

		resp, err := client.Get(healthAPIAddr + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("services ready")
			return true
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Printf("Health check failed after %d attempts.\n", maxAttempts)
	return false
}

// TestE2EWorkflow runs the entire end-to-end flow across task-manager-api and user-management-api
func TestE2EWorkflow(t *testing.T) {
	t.Run("Create a user", testCreateUser)
	t.Run("Login with the user", testLoginUser)
	t.Run("Get user by ID", testGetUserByID)
	t.Run("Create a new task", testCreateTask)
	t.Run("Get user's tasks", testGetTasks)
	t.Run("Delete the created task", testDeleteTask)
	t.Run("Delete the user", testDeleteUser)
}

func testCreateUser(t *testing.T) {
	userEmail = fmt.Sprintf("e2e-test-%d@test.com", time.Now().UnixNano())

	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": "e2e-password",
		"claimIds": []int{0},
	}
	body, _ := json.Marshal(reqBody)

	resp, err := client.Post(userManagementAPIAddr+"/users", "application/json", bytes.NewReader(body))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	userID = respBody["id"]
	assert.NotEmpty(t, userID)
}

func testLoginUser(t *testing.T) {
	reqBody := map[string]interface{}{
		"email":    userEmail,
		"password": "e2e-password",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := client.Post(userManagementAPIAddr+"/users/login", "application/json", bytes.NewReader(body))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	userToken = respBody["token"]
	assert.NotEmpty(t, userToken)
}

func testGetUserByID(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", userManagementAPIAddr, userID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, userID, respBody["id"])
}

func testCreateTask(t *testing.T) {
	reqBody := map[string]interface{}{
		"title":       "E2E Test Task",
		"description": "This is a task created by an E2E test.",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", taskManagerAPIAddr+"/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	taskID = respBody["id"]
	assert.NotEmpty(t, taskID)
}

func testGetTasks(t *testing.T) {
	req, _ := http.NewRequest("GET", taskManagerAPIAddr+"/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NotEmpty(t, respBody)
	assert.Equal(t, taskID, respBody[0]["id"])
}

func testDeleteTask(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/tasks/%s", taskManagerAPIAddr, taskID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testDeleteUser(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", userManagementAPIAddr, userID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
