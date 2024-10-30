package tests

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Teardown(ctx context.Context, container *postgres.PostgresContainer) {
	if err := (*container).Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %v", err)
	}
}
