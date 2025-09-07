package services

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/stretchr/testify/assert"
)

// TestNewHealthService_Ok checks that NewHealthService returns a new service instance
func TestNewHealthService_Ok(t *testing.T) {
	// Arrange
	cfg := config.Config{
		URLs: "http://testing.com",
	}

	// Act
	service := NewHealthService(cfg)

	// Assert
	assert.NotEmpty(t, service)
}

// TestHealthCheck_NoURLs checks that the service returns an error when no URLs are provided
func TestHealthCheck_NoURLs(t *testing.T) {
	// Arrange
	cfg := config.Config{}
	service := NewHealthService(cfg)

	// Act
	resps, err := service.HealthCheck(context.Background())

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "no URLs provided", err.Error())
	assert.Nil(t, resps)
}
