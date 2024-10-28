package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDBContainer struct {
	Id         string
	ConnString string
}

func Setup(t *testing.T) (testcontainers.Container, string, error) {
	container, err := SetupPostgresContainer(t)
	if err != nil {
		return nil, "", err
	}

	connString, err := GetConnString(t, container)
	if err != nil {
		container.Terminate(context.Background())
		return nil, "", err
	}

	return container, connString, nil
}

func GetConnString(t *testing.T, postgresC testcontainers.Container) (string, error) {
	ctx := context.Background()
	host, err := postgresC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
		return "", err
	}

	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
		return "", err
	}

	return fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port()), nil
}

func SetupPostgresContainer(t *testing.T) (testcontainers.Container, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
		return nil, err
	}

	return postgresC, nil
}
