package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
)

type healthHandler struct {
	ctx context.Context
	cfg config.Config
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(ctx context.Context, cfg config.Config) healthHandler {
	return healthHandler{
		ctx: ctx,
		cfg: cfg,
	}
}

// SetHealthRoutes creates health routes
func SetHealthRoutes(router *mux.Router, h healthHandler) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
}

const (
	taskManagerHealthURL    = "http://task-manager-api/v1/health"
	userManagementHealthURL = "http://user-management-api/v1/health"
)

// HealthStatus represents the status of a single service.
type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Global variable to store the final combined status.
var finalStatus = "OK"

// @Summary Health check
// @Description Returns basic runtime information of the API when the service is up
// @Tags Health
// @Success 200 "OK"
// @Router /health [get]
func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	results := make(chan HealthStatus, 2)

	// Set the global status back to OK for each new request.
	finalStatus = "OK"

	// Call the health check for Task Manager concurrently.
	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- checkServiceHealth(r.Context(), taskManagerHealthURL, "task-manager-api")
	}()

	// Call the health check for User Management concurrently.
	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- checkServiceHealth(r.Context(), userManagementHealthURL, "user-management-api")
	}()

	// Wait for both goroutines to finish.
	wg.Wait()
	close(results)

	// Collect the results.
	var healthStatuses []HealthStatus
	for res := range results {
		healthStatuses = append(healthStatuses, res)
	}

	w.Header().Set("Content-Type", "application/json")

	// Set the final HTTP status code.
	if finalStatus == "UNHEALTHY" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Write the JSON response.
	json.NewEncoder(w).Encode(healthStatuses)
}

// checkServiceHealth makes a GET request and returns a HealthStatus.
func checkServiceHealth(ctx context.Context, url string, serviceName string) HealthStatus {
	// Create an HTTP client with a 5-second timeout.
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		finalStatus = "UNHEALTHY"
		return HealthStatus{
			Service: serviceName,
			Status:  "UNHEALTHY",
			Message: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		finalStatus = "UNHEALTHY"
		return HealthStatus{
			Service: serviceName,
			Status:  "UNHEALTHY",
			Message: fmt.Sprintf("received status code: %d", resp.StatusCode),
		}
	}

	return HealthStatus{
		Service: serviceName,
		Status:  "OK",
	}
}
