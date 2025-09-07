package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
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

// HealthStatus represents the status of a single service. # TODO move to models
type HealthStatus struct {
	ServiceURL string `json:"service_url"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

// TODO make it work with vscode launch.json
// @Summary Health check
// @Description Returns the status of all the microservices in the system
// @Tags Health
// @Success 200 "OK"
// @Router /health [get]
func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	servicesToCheck := strings.Split(h.cfg.URLs, ",")
	if len(servicesToCheck) == 0 {
		http.Error(w, "No services to check", http.StatusServiceUnavailable)
		return
	}

	var wg sync.WaitGroup
	results := make(chan HealthStatus, len(servicesToCheck))

	overallStatus := http.StatusOK

	for _, serviceURL := range servicesToCheck {
		if serviceURL == "" {
			continue
		}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			results <- h.checkServiceHealth(r.Context(), url, url)
		}(serviceURL)
	}

	wg.Wait()
	close(results)

	var healthStatuses []HealthStatus
	for res := range results {
		healthStatuses = append(healthStatuses, res)
		if res.Status == "UNHEALTHY" {
			overallStatus = http.StatusServiceUnavailable
		}
	}

	utils.SuccessResponse(w, overallStatus, healthStatuses)
}

func (h *healthHandler) checkServiceHealth(ctx context.Context, url string, serviceName string) HealthStatus {
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return HealthStatus{
			ServiceURL: serviceName,
			Status:     "UNHEALTHY",
			Message:    err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthStatus{
			ServiceURL: serviceName,
			Status:     "UNHEALTHY",
			Message:    fmt.Sprintf("received status code: %d", resp.StatusCode),
		}
	}

	return HealthStatus{
		ServiceURL: serviceName,
		Status:     "OK",
	}
}
