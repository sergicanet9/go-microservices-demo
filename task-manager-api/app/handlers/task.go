package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/models"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v4/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
)

type taskHandler struct {
	ctx context.Context
	cfg config.Config
	svc ports.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(ctx context.Context, cfg config.Config, svc ports.TaskService) taskHandler {
	return taskHandler{
		ctx: ctx,
		cfg: cfg,
		svc: svc,
	}
}

// SetTaskRoutes creates task routes
func SetTaskRoutes(router *mux.Router, t taskHandler) {
	router.HandleFunc("/tasks", t.createTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", t.getTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", t.deleteTask).Methods(http.MethodDelete)
}

// @Summary Create task
// @Description Creates a new task for the logged in user
// @Tags Tasks
// @Security Bearer
// @Param user body models.CreateTaskReq true
// @Success 201 {object} models.CreateTaskResp
// @Router /task [post]
func (t *taskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(t.ctx, t.cfg.Timeout.Duration)
	defer cancel()

	userID, err := getUserID(r)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	var createTaskReq models.CreateTaskReq
	err = json.Unmarshal(body, &createTaskReq)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	createTaskResp, err := t.svc.Create(ctx, userID, createTaskReq)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.SuccessResponse(w, http.StatusCreated, createTaskResp)
}

// @Summary Get tasks
// @Description Gets all tasks for the logged in user
// @Tags Tasks
// @Security Bearer
// @Success 200 {array} models.GetTaskResp
// @Router /task [get]
func (t *taskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(t.ctx, t.cfg.Timeout.Duration)
	defer cancel()

	userID, err := getUserID(r)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	getTasksResp, err := t.svc.GetByUserID(ctx, userID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.SuccessResponse(w, http.StatusOK, getTasksResp)
}

// @Summary Delete task
// @Description Deletes a tasks for the logged in user
// @Tags Tasks
// @Security Bearer
// @Param id path string true
// @Success 200
// @Router /task [delete]
func (t *taskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(t.ctx, t.cfg.Timeout.Duration)
	defer cancel()

	userID, err := getUserID(r)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	var params = mux.Vars(r)
	err = t.svc.Delete(ctx, userID, params["id"])
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.SuccessResponse(w, http.StatusOK, nil)
}

func getUserID(r *http.Request) (string, error) {
	claimsValue := r.Context().Value(middlewares.ClaimsKey)
	if claimsValue == nil {
		return "", fmt.Errorf("claims not found in context")
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims type")
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return "", fmt.Errorf("user_id not found in claims")
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		return "", fmt.Errorf("invalid user_id type")
	}

	return userID, nil
}
