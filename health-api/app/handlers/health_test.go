package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/ports"
	"github.com/sergicanet9/go-microservices-demo/health-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHealthCheck_Ok checks that healthCheck handler does not return an error when the service does not fail
func TestHealthCheck_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()
	healthService := mocks.NewHealthService(t)
	expectedResponse := []models.HealthResp{
		{
			ServiceURL: "http://test.com/health",
			Status:     "HEALTHY",
		},
	}
	healthService.On(testutils.FunctionName(t, ports.HealthService.HealthCheck), mock.Anything).Return(expectedResponse, nil)

	cfg := config.Config{}
	healthHandler := NewHealthHandler(context.Background(), cfg, healthService)
	SetHealthRoutes(r, healthHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}

	var response []models.HealthResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}

	assert.ElementsMatch(t, expectedResponse, response)
}

// TestHealthCheck_ServiceUnavailable checks that healthCheck handler returns a ServiceUnavailable status code when the service returns a service unavailable error
func TestHealthCheck_ServiceUnavailable(t *testing.T) {
	// Arrange
	r := mux.NewRouter()
	healthService := mocks.NewHealthService(t)
	expectedResponse := []models.HealthResp{
		{
			ServiceURL: "http://test.com/health",
			Status:     "UNHEALTHY",
		},
	}
	healthService.On(testutils.FunctionName(t, ports.HealthService.HealthCheck), mock.Anything).Return(expectedResponse, wrappers.NewServiceUnavailableErr(fmt.Errorf("test-error")))

	cfg := config.Config{}
	healthHandler := NewHealthHandler(context.Background(), cfg, healthService)
	SetHealthRoutes(r, healthHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/health"
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusServiceUnavailable, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}

	var response []models.HealthResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}

	assert.ElementsMatch(t, expectedResponse, response)
}

// TestHealthCheck_ServiceError checks that healthCheck handler returns an error response when the service fails
func TestHealthCheck_ServiceError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()
	healthService := mocks.NewHealthService(t)
	expectedError := "service-error"
	healthService.On(testutils.FunctionName(t, ports.HealthService.HealthCheck), mock.Anything).Return([]models.HealthResp{}, errors.New(expectedError))

	cfg := config.Config{}
	healthHandler := NewHealthHandler(context.Background(), cfg, healthService)
	SetHealthRoutes(r, healthHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/health"
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}
