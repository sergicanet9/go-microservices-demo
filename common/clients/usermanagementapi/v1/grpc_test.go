package v1

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/sergicanet9/go-microservices-demo/common/clients/ports"
	"github.com/sergicanet9/go-microservices-demo/common/proto/usermanagementapi/v1/gen/go/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestGRPCClient_Ok checks that a new gRPC Client can be created and closed
func TestGRPCClient_Ok(t *testing.T) {
	// Arrange
	ctx := context.Background()
	target := "test-target"

	// Act
	client, err := NewGRPCClient(ctx, target)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)

	if client != nil {
		err = client.Close()
		assert.NoError(t, err)
	}
}

// TestGRPCClient checks that the client handles scenarios as expected
func TestGRPCClient(t *testing.T) {
	serverAddr, grpcServer, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer grpcServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewGRPCClient(ctx, serverAddr)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	tests := []struct {
		name    string
		runTest func(t *testing.T, client ports.UserManagementV1GRPCClient)
	}{
		{
			name: "HealthCheck Ok",
			runTest: func(t *testing.T, client ports.UserManagementV1GRPCClient) {
				err := client.Health(context.Background())
				assert.NoError(t, err, "Health check should not return an error")
			},
		},
		{
			name: "Exists - user exists",
			runTest: func(t *testing.T, client ports.UserManagementV1GRPCClient) {
				exists, err := client.Exists(context.Background(), "Bearer test-token", "test-id")
				assert.NoError(t, err)
				assert.True(t, exists)
			},
		},
		{
			name: "Exists - user does not exist",
			runTest: func(t *testing.T, client ports.UserManagementV1GRPCClient) {
				exists, err := client.Exists(context.Background(), "Bearer test-token", "non-existent-id")
				assert.NoError(t, err)
				assert.False(t, exists)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.runTest(t, client)
		})
	}
}

// HELP FUNCTIONS
type mockUserManagementServer struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedHealthServiceServer
}

func (s *mockUserManagementServer) GetByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	expectedToken := "Bearer test-token"
	if authHeader[0] != expectedToken {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token")
	}

	if req.Id == "test-id" {
		return &pb.GetUserResponse{Id: "test-id"}, nil
	}
	return nil, status.Errorf(codes.NotFound, "user not found")
}

func (s *mockUserManagementServer) HealthCheck(context.Context, *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{}, nil
}

func newTestServer() (string, *grpc.Server, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &mockUserManagementServer{})
	pb.RegisterHealthServiceServer(grpcServer, &mockUserManagementServer{})

	go func() {
		grpcServer.Serve(lis)
	}()

	return lis.Addr().String(), grpcServer, nil
}
