package mongo

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/entities"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// taskRepository adapter of a task repository for mongo.
type taskRepository struct {
	infrastructure.MongoRepository
}

// NewTaskRepository creates a task repository for mongo
func NewTaskRepository(ctx context.Context, db *mongo.Database) (ports.TaskRepository, error) {
	r := &taskRepository{
		infrastructure.MongoRepository{
			DB:         db,
			Collection: db.Collection(entities.EntityNameTask),
			Target:     entities.Task{},
		},
	}

	_, err := r.Collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	)
	return r, err
}
