package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	commonMocks "github.com/sergicanet9/go-microservices-demo/common/test/mocks"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHealthCheck checks that HealthCheck handles scenarios as expected.
func TestHealthCheck(t *testing.T) {
	// Arrange
	tests := []struct {
		name               string
		userManagementErr  error
		taskManagerErr     error
		expectedError      error
		expectedUserStatus string
		expectedTaskStatus string
	}{
		{
			name:               "All services are healthy",
			userManagementErr:  nil,
			taskManagerErr:     nil,
			expectedError:      nil,
			expectedUserStatus: "OK",
			expectedTaskStatus: "OK",
		},
		{
			name:               "User Management service is unhealthy",
			userManagementErr:  errors.New("connection error"),
			taskManagerErr:     nil,
			expectedError:      wrappers.NewServiceUnavailableErr(fmt.Errorf("service unavailable")),
			expectedUserStatus: "UNHEALTHY",
			expectedTaskStatus: "OK",
		},
		{
			name:               "Task Manager service is unhealthy",
			userManagementErr:  nil,
			taskManagerErr:     errors.New("HTTP 500"),
			expectedError:      wrappers.NewServiceUnavailableErr(fmt.Errorf("service unavailable")),
			expectedUserStatus: "OK",
			expectedTaskStatus: "UNHEALTHY",
		},
		{
			name:               "Both services are unhealthy",
			userManagementErr:  errors.New("connection error"),
			taskManagerErr:     errors.New("HTTP 500"),
			expectedError:      wrappers.NewServiceUnavailableErr(fmt.Errorf("service unavailable")),
			expectedUserStatus: "UNHEALTHY",
			expectedTaskStatus: "UNHEALTHY",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userManagementClientMock := commonMocks.NewUserManagementV1GRPCClient(t)
			taskManagerClientMock := commonMocks.NewTaskManagerV1HTTPClient(t)

			userManagementClientMock.On("Health", mock.Anything).Return(tc.userManagementErr).Once()
			taskManagerClientMock.On("Health", mock.Anything).Return(tc.taskManagerErr).Once()

			service := &healthService{
				config:               config.Config{},
				taskManagerClient:    taskManagerClientMock,
				userManagementClient: userManagementClientMock,
			}

			// Act
			resp, err := service.HealthCheck(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.expectedError, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Len(t, resp, 3)

			// Verificar el estado de cada servicio de forma independiente
			var selfFound, userFound, taskFound bool
			for _, r := range resp {
				switch r.Service {
				case "health-api (self)":
					assert.Equal(t, "OK", r.Status)
					selfFound = true
				case "user-management-api":
					assert.Equal(t, tc.expectedUserStatus, r.Status)
					userFound = true
				case "task-manager-api":
					assert.Equal(t, tc.expectedTaskStatus, r.Status)
					taskFound = true
				}
			}
			assert.True(t, selfFound, "self health status not found in response")
			assert.True(t, userFound, "user-management-api health status not found in response")
			assert.True(t, taskFound, "task-manager-api health status not found in response")
		})
	}
}
