package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sergicanet9/go-microservices-demo/common/clients/ports"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
)

type httpClient struct {
	baseURL string
	client  *http.Client
}

// NewHTTPClient creates a new HTTP client for Task Manager API v1
func NewHTTPClient(baseURL string) ports.TaskManagerV1HTTPClient {
	return &httpClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Health call
func (c *httpClient) Health(ctx context.Context) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/health", c.baseURL), nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status code: %s", httpResp.Status)
	}

	return
}

// TODO unused - remove? and also tests...

// CreateTask call
func (c *httpClient) CreateTask(ctx context.Context, task models.CreateTaskReq) (resp models.CreateTaskResp, err error) {
	body, err := json.Marshal(task)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/tasks", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		return resp, fmt.Errorf("unexpected http status code: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetTasks call
func (c *httpClient) GetTasks(ctx context.Context, token string) (resp []models.GetTaskResp, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/tasks", c.baseURL), nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	httpResp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected http status code: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteTask call
func (c *httpClient) DeleteTask(ctx context.Context, token string, taskID string) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/tasks/%s", c.baseURL, taskID), nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	httpResp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status code: %s", httpResp.Status)
	}

	return
}
