package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
)

type healthHandler struct {
	ctx context.Context
	cfg config.Config
	svc ports.HealthService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(ctx context.Context, cfg config.Config, svc ports.HealthService) healthHandler {
	return healthHandler{
		ctx: ctx,
		cfg: cfg,
		svc: svc,
	}
}

// SetHealthRoutes creates health routes
func SetHealthRoutes(router *mux.Router, h healthHandler) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
}

// @Summary Health check
// @Description Returns the status of all the microservices in the system
// @Tags Health
// @Success 200 "OK"
// @Router /health [get]
func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(h.ctx, h.cfg.Timeout.Duration)
	defer cancel()

	response, err := h.svc.HealthCheck(ctx)
	if err != nil {
		if errors.Is(err, wrappers.ServiceUnavailableErr) {
			utils.SuccessResponse(w, http.StatusServiceUnavailable, response)
			return
		}
		utils.ErrorResponse(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}
