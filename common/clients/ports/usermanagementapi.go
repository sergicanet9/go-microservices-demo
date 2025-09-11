package ports

import "github.com/sergicanet9/go-microservices-demo/common/proto/usermanagementapi/v1/gen/go/pb"

// UserManagementAPIV1GRPCClient interface for a User Management API v1 gRPC Client
type UserManagementAPIV1GRPCClient interface {
	Close() error
	User() pb.UserServiceClient
	Health() pb.HealthServiceClient
}
