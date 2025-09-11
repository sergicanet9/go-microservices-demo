package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sergicanet9/go-microservices-demo/common/clients/ports"
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
