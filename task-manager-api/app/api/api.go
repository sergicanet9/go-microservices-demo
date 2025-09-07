package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/app/handlers"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	config   config.Config
	services svs
}

type svs struct {
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api) {
	a.config = cfg
	return a
}

func (a *api) RunHTTP(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()

		router := mux.NewRouter()

		healthHandler := handlers.NewHealthHandler(ctx, a.config)
		handlers.SetHealthRoutes(router, healthHandler)

		router.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

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
