.PHONY: all build test clean generate wire sqlc prisma dev test test-unit test-integration test-coverage

# Go related variables
BINARY_NAME=mallbots-api
MAIN_PACKAGE=./cmd/api

# Tools installation
install-tools:
	@echo "Installing required tools..."
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/cosmtrek/air@latest

# Development server with hot reload
dev:
	air

# Build the application
build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build files
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf modules/token/infrastructure/query/gen/
	rm -rf **/wire_gen.go

# Generate all
generate: wire sqlc prisma-generate

# Wire dependency injection
wire:
	@echo "Generating wire..."
	wire ./...

# Generate SQLC
sqlc:
	@echo "Generating SQLC..."
	sqlc generate

# Prisma operations
prisma-init:
	npx prisma init

prisma-format:
	npx prisma format

prisma-generate:
	npx prisma generate

prisma-db-pull:
	npx prisma db pull

prisma-db-push:
	npx prisma db push

prisma-migrate-dev:
	npx prisma migrate dev

prisma-migrate-reset:
	npx prisma migrate reset

prisma-studio:
	npx prisma studio

# Generate SQL from Prisma schema
prisma-gen-sql:
	npx prisma migrate diff --from-empty --to-schema-datamodel=./prisma/schema.prisma --script > ./schema.gen.sql

# Combine generate
prisma-all: prisma-format prisma-generate prisma-gen-sql


# Docker compose operations
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Linting and formatting
lint:
	golangci-lint run

fmt:
	go fmt ./...

# Dev environment setup
setup: install-tools
	@echo "Setting up development environment..."
	cp .env.example .env
	go mod download
	go mod tidy

# Run all tests
test:
	go test   ./...

# Run tests with coverage
test-coverage:
	go test  -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Help
help:
	@echo "Available commands:"
	@echo "Development:"
	@echo "  make setup              - Install tools and setup development environment"
	@echo "  make dev               - Run development server with hot reload"
	@echo "  make build             - Build the application"
	@echo "  make fmt               - Format code"
	@echo "  make lint              - Run linter"
	@echo ""
	@echo "Testing:"
	@echo "  make test              - Run all tests with race detection"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo ""
	@echo "Code Generation:"
	@echo "  make generate          - Generate wire, sqlc and prisma"
	@echo "  make wire              - Generate wire dependency injection"
	@echo "  make sqlc              - Generate SQLC"
	@echo ""
	@echo "Database Operations:"
	@echo "  make prisma-init       - Initialize Prisma"
	@echo "  make prisma-generate   - Generate Prisma client"
	@echo "  make prisma-db-pull    - Pull database schema"
	@echo "  make prisma-db-push    - Push schema to database"
	@echo "  make prisma-migrate-dev - Create new migration"
	@echo "  make prisma-studio     - Open Prisma Studio"
