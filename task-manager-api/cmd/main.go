package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/app/api"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/app/async"
	_ "github.com/sergicanet9/go-microservices-demo/task-manager-api/app/docs"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
)

// @title Task Manager API
// @version v1
// @BasePath /task-manager-api/v1
// @tag.description Powered by scv-go-tools
// @tag.docs.url https://github.com/sergicanet9/scv-go-tools
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	var opts struct {
		Version     string `long:"ver" description:"Version" required:"true"`
		Environment string `long:"env" description:"Environment" choice:"local" choice:"prod" required:"true"`
		HTTPPort    int    `long:"hport" description:"Running HTTP port" required:"true"`
		DSN         string `long:"dsn" description:"Database DSN" required:"true"`
		JWTSecret   string `long:"jsecret" description:"Secret used to sign and validate JWT tokens" required:"true"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		observability.Logger().Fatal(fmt.Errorf("provided flags not valid: %s, %w", args, err))
	}

	cfg, err := config.ReadConfig(opts.Version, opts.Environment, opts.HTTPPort, opts.DSN, opts.JWTSecret, "config")
	if err != nil {
		observability.Logger().Fatal(fmt.Errorf("cannot parse config file for env %s: %w", opts.Environment, err))
	}

	observability.Logger().Printf("Version: %s", cfg.Version)
	observability.Logger().Printf("Environment: %s", cfg.Environment)

	var g multierror.Group
	ctx, cancel := context.WithCancel(context.Background())

	a := api.New(ctx, cfg)
	g.Go(a.RunHTTP(ctx, cancel))

	if cfg.Async.Run {
		async := async.New(cfg)
		g.Go(async.Run(ctx, cancel))
	}

	<-ctx.Done()
	observability.Logger().Printf("context canceled, the application will terminate...")

	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()

	select {
	case <-done:
		observability.Logger().Printf("application terminated gracefully")
	case <-time.After(10 * time.Second):
		observability.Logger().Fatalf("some processes did not terminate gracefully, application termination forced")
	}
}
