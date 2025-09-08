package entities

import "time"

const EntityNameTask = "tasks"

type Task struct {
	ID          string    `bson:"_id,omitempty"`
	UserID      string    `bson:"userID"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
}
