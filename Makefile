.PHONY: all build test clean run run-dev docker-build docker-run docker-down lint

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

# Run Docker container <hostPort>:<containerPort> - container port must be the same for the port defined in conf.json or default 8080
docker-run:
	@echo "Running Docker container..."
	@docker run --rm -p 5001:8080 \
		-v $(CURDIR)/config.json:/app/config.json:ro \
		-v $(CURDIR)/private.pem:/app/private.pem:ro \
		-v $(CURDIR)/public.pem:/app/public.pem:ro \
		$(BINARY_NAME):latest

# Stop and remove Docker container using docker-compose
docker-down:
	@echo "Stopping Docker container..."
	@docker-compose -f docker/docker-compose.yml down

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

# To run the service built locally
docker-compose-up:
	@docker-compose -f docker/docker-compose.yml up goauth-app -d
# @docker-compose -f docker/docker-compose.yml up app -d
# To run the service from Docker Hub
