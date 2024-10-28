# Declare phony targets
.PHONY: dev test help

# Default target
.DEFAULT_GOAL := help
dev:
	go run main.go

test:
	go test -v -race -cover ./${module}...

# Help target
help:
	@echo "Available targets:"
	@echo "  dev   - Run the application in development mode"
	@echo "  test  - Run tests with verbose output and race detection"
