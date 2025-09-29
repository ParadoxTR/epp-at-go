# EPP-AT-GO Makefile

.PHONY: help build test lint clean fmt vet deps check-deps install-tools

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the project
build: ## Build the project
	@echo "Building EPP-AT-GO..."
	go build ./...

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint: check-deps ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning..."
	go clean ./...
	rm -f coverage.out coverage.html

# Download dependencies
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Check if required tools are installed
check-deps: ## Check if required development tools are installed
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run 'make install-tools' first."; exit 1; }

# Install development tools
install-tools: ## Install required development tools
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	else \
		echo "golangci-lint is already installed"; \
	fi

# Run all checks (format, vet, lint, test)
check: fmt vet lint test ## Run all code quality checks

# Example usage
example: ## Run the basic usage example
	@echo "Running basic usage example..."
	cd examples && go run basic_usage.go

# Generate documentation
docs: ## Generate and serve documentation
	@echo "Generating documentation..."
	godoc -http=:6060 -play
	@echo "Documentation server started at http://localhost:6060"

# Security scan
security: ## Run security scan with gosec
	@if command -v gosec >/dev/null 2>&1; then \
		echo "Running security scan..."; \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Benchmark tests
bench: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Update dependencies
update-deps: ## Update all dependencies to latest versions
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Validate go.mod
validate-mod: ## Validate go.mod file
	@echo "Validating go.mod..."
	go mod verify

# Pre-commit hook (run before committing)
pre-commit: fmt vet lint test ## Run pre-commit checks
	@echo "Pre-commit checks completed successfully!"