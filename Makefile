.PHONY: all build test clean run run-dev docker-build docker-run lint

# Variables
BINARY_NAME=auth
MAIN_PATH=cmd/server/main.go
CONFIG_PATH=config.json

# Default target
all: clean build

# Build the application
build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Run the application
run:
	@echo "Running application..."
	@go run $(MAIN_PATH)

# Run in development mode
run-dev:
	@echo "Running in development mode..."
	@go run $(MAIN_PATH) -mode=dev

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest -f docker/Dockerfile .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Generate mocks for testing
generate-mocks:
	@echo "Generating mocks..."
	@go generate ./...

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@go run $(MAIN_PATH) migrate up

# Run database migrations down
migrate-down:
	@echo "Rolling back database migrations..."
	@go run $(MAIN_PATH) migrate down

# Health check
health:
	@echo "Checking health..."
	@curl -X GET http://localhost:8080/api/v1/auth/health -w "\n"