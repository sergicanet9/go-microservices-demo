package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
)

// userService adapter of an health service
type healthService struct {
	config config.Config
}

// NewHealthService creates a new health service
func NewHealthService(cfg config.Config) ports.HealthService {
	return &healthService{
		config: cfg,
	}
}

// HealthCheck all services
func (h *healthService) HealthCheck(ctx context.Context) ([]models.HealthResp, error) {
	if h.config.URLs == "" {
		return nil, errors.New("no URLs provided")
	}

	urls := strings.Split(h.config.URLs, ",")

	var wg sync.WaitGroup
	results := make(chan models.HealthResp, len(urls)+1)

	results <- models.HealthResp{
		ServiceURL: "self",
		Status:     "OK",
	}

	for _, serviceURL := range urls {
		if serviceURL == "" {
			continue
		}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			results <- h.checkServiceHealth(ctx, url)
		}(serviceURL)
	}

	wg.Wait()
	close(results)

	var healthResps []models.HealthResp
	var err error
	for res := range results {
		if res.Error != "" {
			err = wrappers.NewServiceUnavailableErr(fmt.Errorf("service unavailable"))
		}
		healthResps = append(healthResps, res)
	}

	return healthResps, err
}

func (h *healthService) checkServiceHealth(ctx context.Context, url string) models.HealthResp {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return models.HealthResp{
			ServiceURL: url,
			Status:     "UNHEALTHY",
			Error:      err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.HealthResp{
			ServiceURL: url,
			Status:     "UNHEALTHY",
			Error:      fmt.Sprintf("received status code: %d", resp.StatusCode),
		}
	}

	return models.HealthResp{
		ServiceURL: url,
		Status:     "OK",
	}
}
