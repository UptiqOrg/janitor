package tests

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
)

func Teardown(container *testcontainers.Container) {
	ctx := context.Background()
	if err := (*container).Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %v", err)
	}
}
