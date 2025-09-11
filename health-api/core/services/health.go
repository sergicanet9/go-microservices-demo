package services

import (
	"context"
	"fmt"
	"sync"

	commonPorts "github.com/sergicanet9/go-microservices-demo/common/clients/ports"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
)

// healthService adapter of an health service
type healthService struct {
	config               config.Config
	taskManagerClient    commonPorts.TaskManagerV1HTTPClient
	userManagementClient commonPorts.UserManagementV1GRPCClient
}

// NewHealthService creates a new health service
func NewHealthService(cfg config.Config, taskManagerClient commonPorts.TaskManagerV1HTTPClient, userManagementClient commonPorts.UserManagementV1GRPCClient) ports.HealthService {
	return &healthService{
		config:               cfg,
		taskManagerClient:    taskManagerClient,
		userManagementClient: userManagementClient,
	}
}

// HealthCheck all services concurrently
func (h *healthService) HealthCheck(ctx context.Context) ([]models.HealthResp, error) {
	var wg sync.WaitGroup
	results := make(chan models.HealthResp, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- models.HealthResp{
			Service: "health-api (self)",
			Status:  "OK",
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.userManagementClient.Health(ctx)
		if err != nil {
			results <- models.HealthResp{Service: "user-management-api", Status: "UNHEALTHY", Error: err.Error()}
			return
		}
		results <- models.HealthResp{Service: "user-management-api", Status: "OK"}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.taskManagerClient.Health(ctx)
		if err != nil {
			results <- models.HealthResp{Service: "task-manager-api", Status: "UNHEALTHY", Error: err.Error()}
			return
		}
		results <- models.HealthResp{Service: "task-manager-api", Status: "OK"}
	}()

	wg.Wait()
	close(results)

	var healthResps []models.HealthResp
	var serviceUnavailable bool
	for res := range results {
		if res.Status == "UNHEALTHY" {
			serviceUnavailable = true
		}
		healthResps = append(healthResps, res)
	}

	if serviceUnavailable {
		return healthResps, wrappers.NewServiceUnavailableErr(fmt.Errorf("service unavailable"))
	}

	return healthResps, nil
}
