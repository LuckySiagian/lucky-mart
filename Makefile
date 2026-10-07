.PHONY: run build test lint swagger docker-up docker-down clean help

# Default target
.DEFAULT_GOAL := help

APP_NAME=backend-sembako
MAIN_PATH=./cmd/api

## run: Run the application locally
run:
	@echo "Starting application..."
	@go run $(MAIN_PATH)/main.go

## build: Build application binary
build:
	@echo "Building binary..."
	@go build -o bin/$(APP_NAME) $(MAIN_PATH)/main.go
	@echo "Binary created at bin/$(APP_NAME)"

## test: Run unit tests
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## swagger: Generate Swagger OpenAPI documentation
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs

## docker-up: Start Docker containers (Backend, Postgres, PgAdmin)
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d --build

## docker-down: Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down -v

## clean: Clean build binaries and temporary files
clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out

## help: Display available commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
