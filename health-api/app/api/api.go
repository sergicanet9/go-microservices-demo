package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	taskManagerClient "github.com/sergicanet9/go-microservices-demo/common/clients/taskmanagerapi/v1"
	userManagementClient "github.com/sergicanet9/go-microservices-demo/common/clients/usermanagementapi/v1"
	handlersV1 "github.com/sergicanet9/go-microservices-demo/health-api/app/handlers/v1"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/ports"
	"github.com/sergicanet9/go-microservices-demo/health-api/core/services"
	"github.com/sergicanet9/scv-go-tools/v4/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	config   config.Config
	services svs
}

type svs struct {
	health ports.HealthService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api) {
	a.config = cfg

	taskManagerClient := taskManagerClient.NewHTTPClient(cfg.TaskManagerURL)

	userManagementClient, err := userManagementClient.NewGRPCClient(ctx, cfg.UserManagementTarget)
	if err != nil {
		observability.Logger().Fatal(err)
	}

	a.services.health = services.NewHealthService(a.config, taskManagerClient, userManagementClient)

	return a
}

func (a *api) RunHTTP(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()

		router := mux.NewRouter()
		router.Use(middlewares.Logger("/swagger", "/docs.swagger.json", "/grpcui"))
		router.Use(middlewares.Recover)

		v1Router := router.PathPrefix("/health-api/v1").Subrouter()

		healthHandler := handlersV1.NewHealthHandler(ctx, a.config, a.services.health)
		handlersV1.SetHealthRoutes(v1Router, healthHandler)

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
