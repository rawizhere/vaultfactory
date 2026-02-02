.PHONY: help build test lint clean deps

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build both server and client
	@echo "Building server and client"
	@go build -o bin/vaultfactory-server ./cmd/server
	@go build -o bin/vaultfactory-client ./cmd/client

build-server: ## Build server only
	@echo "Building server"
	@go build -o bin/vaultfactory-server ./cmd/server

build-client: ## Build client only
	@echo "Building client"
	@go build -o bin/vaultfactory-client ./cmd/client

test: ## Run tests
	@echo "Running tests"
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## Run linters
	@echo "Running linters"
	@golangci-lint run
	@staticcheck ./...
	@go vet ./...

lint-fix: ## Run linters with auto-fix
	@echo "Running linters with auto-fix"
	@golangci-lint run --fix

fmt: ## Format code
	@echo "Formatting code"
	@go fmt ./...

fmt-check: ## Check if code is formatted
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Files not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi

security: ## Run security checks
	@echo "Running security checks"
	@govulncheck ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts"
	@rm -rf bin/
	@rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies"
	@go mod download
	@go mod verify

deps-update: ## Update dependencies
	@echo "Updating dependencies"
	@go get -u ./...
	@go mod tidy

dev-setup: deps ## Setup development environment
	@echo "Installing development tools"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest

db-migrate: ## Run database migrations
	@echo "Running database migrations"
	@migrate -path scripts/migrations -database "postgres://vaultfactory:password@localhost:5432/vaultfactory?sslmode=disable" up

db-rollback: ## Rollback database migrations
	@echo "Rolling back database migrations"
	@migrate -path scripts/migrations -database "postgres://vaultfactory:password@localhost:5432/vaultfactory?sslmode=disable" down

docker-build: ## Build Docker images
	@echo "Building Docker images"
	@docker build -t vaultfactory-server -f Dockerfile.server .
	@docker build -t vaultfactory-client -f Dockerfile.client .

docker-run: ## Run with Docker Compose
	@echo "Starting services with Docker Compose"
	@docker-compose up -d

docker-stop: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services"
	@docker-compose down

ci: fmt-check lint test ## Run CI checks locally

release: clean test lint ## Prepare for release