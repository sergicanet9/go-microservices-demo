package v1

import (
	"context"
	"fmt"

	"github.com/sergicanet9/go-microservices-demo/common/clients/ports"
	"github.com/sergicanet9/go-microservices-demo/common/proto/usermanagementapi/v1/gen/go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcClient struct {
	healthClient pb.HealthServiceClient
	userClient   pb.UserServiceClient
	conn         *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client for User Management API v1
func NewGRPCClient(ctx context.Context, target string) (ports.UserManagementV1GRPCClient, error) {
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

// Health calls HealthCheck
func (c *grpcClient) Health(ctx context.Context) error {
	_, err := c.healthClient.HealthCheck(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("HealthCheck call failed: %w", err)
	}
	return nil
}

// Exists calls GetByID and returns if the player exists
func (c *grpcClient) Exists(ctx context.Context, token, userID string) (bool, error) {
	md := metadata.Pairs("authorization", token)
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := c.userClient.GetByID(mdCtx, &pb.GetUserByIDRequest{Id: userID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return false, nil
		}
		return false, err

	}
	if resp.Id == userID {
		return true, nil
	}
	return false, fmt.Errorf("unexpected GetByID response: %v", resp)
}
