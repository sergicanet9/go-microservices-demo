package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	handlersV1 "github.com/sergicanet9/go-microservices-demo/task-manager-api/app/handlers/v1"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/ports"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/core/services"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/infrastructure/mongo"
	"github.com/sergicanet9/scv-go-tools/v4/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	config   config.Config
	services svs
}

type svs struct {
	task ports.TaskService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api) {
	a.config = cfg

	db, err := infrastructure.ConnectMongoDB(ctx, a.config.DSN)
	if err != nil {
		observability.Logger().Fatal(err)
	}

	taskRepo, err := mongo.NewTaskRepository(ctx, db)
	if err != nil {
		observability.Logger().Fatal(err)
	}

	a.services.task = services.NewTaskService(a.config, taskRepo)

	return a
}

func (a *api) RunHTTP(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()

		router := mux.NewRouter()
		router.Use(middlewares.Logger("/swagger", "/docs.swagger.json", "/grpcui"))
		router.Use(middlewares.Recover)

		v1Router := router.PathPrefix("/task-manager-api/v1").Subrouter()

		healthHandler := handlersV1.NewHealthHandler(ctx, a.config)
		handlersV1.SetHealthRoutes(v1Router, healthHandler)

		taskHandler := handlersV1.NewTaskHandler(ctx, a.config, a.services.task)
		handlersV1.SetTaskRoutes(v1Router, taskHandler)

		v1Router.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.HTTPPort),
			Handler: router,
		}

		go shutdown(ctx, server)

		observability.Logger().Printf("Server listening on HTTP port %d", a.config.HTTPPort)
		return server.ListenAndServe()
	}
}

func shutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	observability.Logger().Printf("Shutting down HTTP server gracefully...")
	server.Shutdown(ctx)
}
