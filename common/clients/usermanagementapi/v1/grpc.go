package v1

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/common/clients/ports"
	"github.com/sergicanet9/go-microservices-demo/common/proto/usermanagementapi/v1/gen/go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	healthClient pb.HealthServiceClient
	userClient   pb.UserServiceClient
	conn         *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client for User Management API v1
func NewGRPCClient(ctx context.Context, target string) (ports.UserManagementAPIV1GRPCClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	healthClient := pb.NewHealthServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)

	return &grpcClient{
		healthClient: healthClient,
		userClient:   userClient,
		conn:         conn,
	}, nil
}

// Close closes the gRPC connection
func (c *grpcClient) Close() error {
	return c.conn.Close()
}

// User returns the gRPC User client
func (c *grpcClient) User() pb.UserServiceClient {
	return c.userClient
}

// Health returns the gRPC Health client
func (c *grpcClient) Health() pb.HealthServiceClient {
	return c.healthClient
}
