BINARY_NAME=zgo
SERVER_NAME=server
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build build-server install clean test cover lint generate wire docs mock dev air help

# Default target
all: lint test build

# Build the CLI tool (auto-runs wire first)
build: wire
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) cmd/zgo/main.go

# Build the server
build-server:
	@echo "Building $(SERVER_NAME)..."
	go build $(LDFLAGS) -o $(SERVER_NAME) cmd/server/main.go

# Build all
build-all: build build-server

# Install CLI to $GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/zgo

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME) $(SERVER_NAME)
	rm -f coverage.txt coverage.html

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test ./... -race -v

# Run tests with coverage
cover:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Fix lint issues automatically
lint-fix:
	@echo "Fixing lint issues..."
	golangci-lint run --fix ./...

# Generate code (Wire, etc.)
generate:
	@echo "Generating code..."
	go generate ./...

# Run Wire dependency injection
wire:
	@echo "Running Wire..."
	@which wire > /dev/null || (echo "Installing wire..." && go install github.com/google/wire/cmd/wire@latest)
	cd internal/wiring && wire

# Generate API documentation (Swagger)
docs:
	@echo "Generating API documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/server/main.go -o docs/swagger

# Generate mocks
mock:
	@echo "Generating mocks..."
	@which mockgen > /dev/null || (echo "Installing mockgen..." && go install go.uber.org/mock/mockgen@latest)
	go generate ./...

# Run development server
dev:
	go run cmd/server/main.go

# Run with Air (hot reload)
air:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Run server
server:
	go run cmd/server/main.go

# Database migrations
migrate:
	./$(BINARY_NAME) migrate

migrate-fresh:
	./$(BINARY_NAME) migrate:fresh

migrate-rollback:
	./$(BINARY_NAME) migrate:rollback

# Database seeding
seed:
	./$(BINARY_NAME) db:seed

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Update dependencies
update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Check for vulnerabilities
vuln:
	@echo "Checking for vulnerabilities..."
	@which govulncheck > /dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)
	govulncheck ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  make build        - Build the CLI tool"
	@echo "  make build-server - Build the server"
	@echo "  make build-all    - Build all binaries"
	@echo "  make install      - Install CLI to GOPATH/bin"
	@echo "  make test         - Run tests"
	@echo "  make test-race    - Run tests with race detection"
	@echo "  make cover        - Run tests with coverage report"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make lint-fix     - Fix lint issues automatically"
	@echo "  make generate     - Run go generate"
	@echo "  make wire         - Run Wire DI generator"
	@echo "  make docs         - Generate Swagger documentation"
	@echo "  make mock         - Generate mocks"
	@echo "  make dev          - Run development server"
	@echo "  make air          - Run with hot reload (Air)"
	@echo "  make migrate      - Run database migrations"
	@echo "  make vuln         - Check for vulnerabilities"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make help         - Show this help"
