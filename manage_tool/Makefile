.PHONY: lint lint-fix lint-verbose lint-ci build run clean dev air

# Lint the code
lint:
	golangci-lint run

# Lint and fix auto-fixable issues
lint-fix:
	golangci-lint run --fix

# Lint with verbose output
lint-verbose:
	golangci-lint run --verbose

# Lint and generate JSON report (same as CI)
lint-ci:
	golangci-lint run --timeout=10m --out-format=json:golangci-lint-report.json,line-number
	@if [ -f golangci-lint-report.json ]; then \
		echo "Lint report generated: golangci-lint-report.json"; \
		cat golangci-lint-report.json; \
	fi

# Build the application
build:
	go build -o bin/manage-tool ./main.go

# Run the application
run:
	go run main.go

# Run with hot reloading using Air
dev:
	air

# Alias for dev (for consistency)
air:
	air

# Clean build artifacts
clean:
	rm -rf bin/ tmp/ *.exe golangci-lint-report.json

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	gofmt -s -w .
	goimports -w .

# Verify dependencies
verify:
	go mod verify

# Check for security vulnerabilities
security:
	gosec ./...

# Run all checks (lint + security + build)
check: lint security build

# Setup development environment
setup:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go mod download