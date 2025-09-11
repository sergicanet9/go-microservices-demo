package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
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

// @Summary Health check
// @Description Returns basic runtime information of the API when the service is up
// @Tags Health
// @Success 200 "OK"
// @Router /health [get]
func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Version", h.cfg.Version)
	r.Header.Add("Environment", h.cfg.Environment)
	r.Header.Add("HTTPPort", strconv.Itoa(h.cfg.HTTPPort))

	dsnValue := "***FILTERED***"
	if h.cfg.Environment == "local" {
		dsnValue = h.cfg.DSN
	}
	r.Header.Add("DSN", dsnValue)

	utils.SuccessResponse(w, http.StatusOK, "OK")
}
