.PHONY: build-server build-client build-all test clean dev-deps run-server run-client

# Build variables
BUILD_VERSION := 1.0.0
BUILD_DATE := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd'")
BUILD_COMMIT := $(shell git rev-parse --short HEAD 2>nul || echo N/A)

# Build targets
build-server:
	go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" -o bin/server.exe ./cmd/server

build-client:
	go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" -o bin/client.exe ./cmd/client

build-all: build-server build-client

# Test
test:
	go test ./...

# Clean
clean:
	rm -rf bin/

# Development
dev-deps:
	go mod download
	go mod tidy

# Database
db-migrate:
	# Add migration commands here when implemented

# Run
run-server: build-server
	./bin/server

run-client: build-client
	./bin/client

# Help
help:
	@echo "Available targets:"
	@echo "  build-server    - Build server binary"
	@echo "  build-client    - Build client binary"
	@echo "  build-all       - Build both server and client"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  dev-deps        - Download and tidy dependencies"
	@echo "  run-server      - Build and run server"
	@echo "  run-client      - Build and run client"
