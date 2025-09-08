package handlers

import (
	"context"

	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
)

type taskHandler struct {
	ctx context.Context
	cfg config.Config
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(ctx context.Context, cfg config.Config) taskHandler {
	return taskHandler{
		ctx: ctx,
		cfg: cfg,
	}
}

// // SetTaskRoutes creates task routes
// func SetTaskRoutes(router *mux.Router, t taskHandler) {
// 	router.HandleFunc("/tasks", t.createTask).Methods(http.MethodPost)
// 	router.HandleFunc("/tasks", t.getTasks).Methods(http.MethodGet)
// 	router.HandleFunc("/tasks/{id}", t.deleteTask).Methods(http.MethodDelete)
// }
