package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotNil(t, client.User())
	assert.NotNil(t, client.Health())

	if client != nil {
		err = client.Close()
		assert.NoError(t, err)
	}
}
