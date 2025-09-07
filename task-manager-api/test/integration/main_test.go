package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/app/api"
	"github.com/sergicanet9/go-microservices-demo/task-manager-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
)

const (
	contentType        = "application/json"
	mongoDBName        = "test-db"
	mongoUser          = "mongo"
	mongoPassword      = "test"
	mongoContainerPort = "27017/tcp"
	mongoDSNEnv        = "mongoDSN"
	jwtSecret          = "eaeBbXUxks"
	nonExpiryToken     = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.cCKM32os5ROKxeE3IiDWoOyRew9T8puzPUKurPhrDug"
)

// TestMain does the setup before running the tests and the teardown afterwards
func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("could not connect to docker: %s", err)
	}

	mongoResource := setupMongo(pool)

	// Runs the tests
	code := m.Run()

	// When it´s done, kill and remove the containers
	if err = pool.Purge(mongoResource); err != nil {
		log.Panicf("could not purge resource: %s", err)
	}
	os.Unsetenv(mongoDSNEnv)

	os.Exit(code)
}

func setupMongo(pool *dockertest.Pool) *dockertest.Resource {
	// creates filekey
	_, filePath, _, _ := runtime.Caller(0)
	fileKey, err := os.CreateTemp(path.Dir(filePath), "")
	if err != nil {
		log.Panic(err)
	}
	defer os.Remove(fileKey.Name())

	bytes := []byte(`secret123`)
	err = os.WriteFile(fileKey.Name(), bytes, 0644)
	if err != nil {
		log.Panic(err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			fmt.Sprintf("MONGO_INITDB_DATABASE=%s", mongoDBName),
			fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", mongoUser),
			fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", mongoPassword),
		},
		Mounts: []string{fmt.Sprintf("%s:/auth/file.key", fileKey.Name())},
		Cmd:    []string{"--keyFile", "/auth/file.key", "--replSet", "rs0", "--bind_ip_all"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	exitCode, err := resource.Exec([]string{"/bin/sh", "-c", "chown mongodb:mongodb /auth/file.key"}, dockertest.ExecOptions{})
	if err != nil {
		log.Panicf("failure executing command in the resource: %s", err)
	}
	if exitCode != 0 {
		log.Panicf("failure executing command in the resource, exit code was %d", exitCode)
	}

	dsn := fmt.Sprintf("mongodb://%s:%s@localhost:%s/%s?authSource=admin&connect=direct", mongoUser, mongoPassword, resource.GetPort(mongoContainerPort), mongoDBName)
	os.Setenv(mongoDSNEnv, dsn)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		_, err = infrastructure.ConnectMongoDB(context.Background(), dsn)
		return err
	})
	if err != nil {
		log.Panicf("Could not connect to docker: %s", err)
	}

	exitCode, err = resource.Exec([]string{"/bin/sh", "-c", fmt.Sprintf("echo 'rs.initiate().ok' | mongosh -u %s -p %s --quiet", mongoUser, mongoPassword)}, dockertest.ExecOptions{})
	if err != nil {
		log.Panicf("failure executing command in the resource: %s", err)
	}
	if exitCode != 0 {
		log.Panicf("failure executing command in the resource, exit code was %d", exitCode)
	}

	return resource
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
	c.DSN = os.Getenv(mongoDSNEnv)
	c.JWTSecret = jwtSecret

	c.Timeout = utils.Duration{Duration: 30 * time.Second}

	return c, nil
}
