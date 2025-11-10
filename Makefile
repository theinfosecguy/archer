.PHONY: build clean install test run docker-build docker-run lint fmt help install-tools imports check-imports

# Variables
BINARY_NAME=archer
MAIN_PATH=./cmd/archer
BUILD_DIR=./build
TEMPLATES_DIR=templates
GOIMPORTS=$(shell go env GOPATH)/bin/goimports

# Build flags
LDFLAGS=-ldflags="-w -s"
GOFLAGS=-trimpath

# Default target
all: build

## build: Build the binary
build: check-imports
	@echo "Building $(BINARY_NAME)..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(GOFLAGS) $(LDFLAGS) $(MAIN_PATH)
	@echo "Install complete"

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests complete"

## run: Run the application (use: make run ARGS="list")
run: build
	@./$(BINARY_NAME) $(ARGS)

## install-tools: Install required development tools
install-tools:
	@echo "Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed successfully"

## imports: Run goimports to organize imports
imports:
	@echo "Running goimports..."
	@if [ ! -f $(GOIMPORTS) ]; then \
		echo "goimports not found. Installing..."; \
		$(MAKE) install-tools; \
	fi
	@$(GOIMPORTS) -w -local github.com/theinfosecguy/archer .
	@echo "Imports organized"

## check-imports: Check if imports are properly formatted
check-imports:
	@if [ ! -f $(GOIMPORTS) ]; then \
		echo "goimports not found. Installing..."; \
		$(MAKE) install-tools; \
	fi
	@if [ -n "$$($(GOIMPORTS) -l -local github.com/theinfosecguy/archer .)" ]; then \
		echo "The following files have import issues:"; \
		$(GOIMPORTS) -l -local github.com/theinfosecguy/archer .; \
		echo "Run 'make imports' to fix them"; \
		exit 1; \
	fi

## fmt: Format Go code
fmt: imports
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

## lint: Run linter
lint:
	@echo "Running linter..."
	@go vet ./...
	@echo "Lint complete"

## tidy: Tidy go.mod
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Tidy complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "Docker build complete"

## docker-run: Run Docker container (use: make docker-run ARGS="list")
docker-run:
	@docker run --rm $(BINARY_NAME):latest $(ARGS)

## release-linux: Build Linux binary
release-linux:
	@echo "Building Linux binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

## release-darwin: Build macOS binary
release-darwin:
	@echo "Building macOS binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "macOS builds complete"

## release-windows: Build Windows binary
release-windows:
	@echo "Building Windows binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

## release-all: Build binaries for all platforms
release-all: release-linux release-darwin release-windows
	@echo "All release builds complete"

## help: Display this help message
help:
	@echo "Archer - Secret Validator"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
