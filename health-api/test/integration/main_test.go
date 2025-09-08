package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sergicanet9/go-microservices-demo/health-api/app/api"
	"github.com/sergicanet9/go-microservices-demo/health-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
)

const (
	contentType = "application/json"
)

// TestMain does the setup before running the tests and the teardown afterwards
func TestMain(m *testing.M) {
	// Runs the tests
	code := m.Run()

	os.Exit(code)
}

// New starts a testing instance of the API and returns its config
func New(t *testing.T) config.Config {
	t.Helper()

	cfg, err := testConfig(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	a := api.New(ctx, cfg)
	runHTTP := a.RunHTTP(ctx, cancel)

	go runHTTP()

	<-time.After(100 * time.Millisecond) // waiting time for letting the API start completely
	return cfg
}

func testConfig(t *testing.T) (c config.Config, err error) {
	c.Version = "Integration tests"
	c.Environment = "Integration tests"
	c.HTTPPort = testutils.FreePort(t)

	c.Timeout = utils.Duration{Duration: 30 * time.Second}

	return c, nil
}
