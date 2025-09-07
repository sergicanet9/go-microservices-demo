package ports

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/health-api/core/models"
)

// HealthService interface
type HealthService interface {
	HealthCheck(ctx context.Context) ([]models.HealthResp, error)
}
