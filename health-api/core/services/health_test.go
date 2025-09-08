package services

import (
	"testing"

	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/stretchr/testify/assert"
)

// TestNewHealthService_Ok checks that NewHealthService returns a new service instance
func TestNewHealthService_Ok(t *testing.T) {
	// Arrange
	cfg := config.Config{}
	cfg.TaskManagerURL = "http://testing.com"
	cfg.UserManagementURL = "http://testing.com"

	// Act
	service := NewHealthService(cfg)

	// Assert
	assert.NotEmpty(t, service)
}

// TODO test
