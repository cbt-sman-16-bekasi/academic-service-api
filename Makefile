.PHONY: build run dev test clean docker swagger migrate-module

# Binary name
BINARY=academic-api
MAIN_PATH=./src/main/cmd/api

# Build the application
build:
	@echo "Building..."
	@go build -o bin/$(BINARY) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running..."
	@./bin/$(BINARY)

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development server..."
	@go run $(MAIN_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Build docker image
docker:
	@echo "Building Docker image..."
	@docker build -t $(BINARY):latest .

# Run docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(BINARY):latest

# Generate swagger docs
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g src/main/cmd/api/main.go -o docs

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Lint code
lint:
	@echo "Linting..."
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting..."
	@go fmt ./...

# Create new module structure
# Usage: make new-module name=modulename
new-module:
	@echo "Creating module: $(name)"
	@mkdir -p src/main/internal/modules/$(name)/entity
	@mkdir -p src/main/internal/modules/$(name)/dto
	@mkdir -p src/main/internal/modules/$(name)/repository
	@mkdir -p src/main/internal/modules/$(name)/service
	@mkdir -p src/main/internal/modules/$(name)/handler
	@echo "package $(name)\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc RegisterRoutes(router *gin.RouterGroup) {\n\t// TODO: Implement routes\n}" > src/main/internal/modules/$(name)/routes.go
	@echo "Module $(name) created successfully!"

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Build and run the application"
	@echo "  make dev          - Run in development mode"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker       - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make deps         - Download dependencies"
	@echo "  make lint         - Lint code"
	@echo "  make fmt          - Format code"
	@echo "  make new-module   - Create new module (usage: make new-module name=modulename)"

