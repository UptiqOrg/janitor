package tests

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Setup(ctx context.Context) (*postgres.PostgresContainer, string, error) {
	container, err := SetupPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("Error setting up postgers container: %s", err)
		return nil, "", err
	}

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}

	return container, connString, nil
}

func SetupPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		postgres.WithInitScripts(filepath.Join("init.sql")),
		testcontainers.WithImage("postgres:15.3"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}
