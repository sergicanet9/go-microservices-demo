package mongo

import (
	"context"
	"testing"

	"github.com/sergicanet9/scv-go-tools/v4/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestNewTaskRepository_Ok checks that NewTaskRepository creates a new taskRepository struct
func TestNewTaskRepository_Ok(t *testing.T) {
	mt := mocks.NewMongoDB(t)

	mt.Run("", func(mt *mtest.T) {
		// Arrange
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		repo, err := NewTaskRepository(context.Background(), mt.DB)

		// Assert
		assert.NotEmpty(t, repo)
		assert.Nil(t, err)
	})
}
