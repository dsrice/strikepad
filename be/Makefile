.PHONY: lint lint-fix build run clean test dev debug

# Lint the code
lint:
	golangci-lint run

# Lint and fix auto-fixable issues
lint-fix:
	golangci-lint run --fix

# Build the application
build:
	go build -o bin/strikepad-backend ./main.go

# Run the application
run:
	go run main.go

# Run with hot reloading using Air
dev:
	air

# Debug the application with Delve
debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient

# Debug with specific breakpoint
debug-break:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient --continue --accept-multiclient main.go

# Clean build artifacts
clean:
	rm -rf bin/ tmp/

# Run tests
test:
	CGO_ENABLED=1 go test -v ./...

# Run tests with coverage
test-coverage:
	CGO_ENABLED=1 go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in watch mode (requires entr)
test-watch:
	find . -name "*.go" | entr -c CGO_ENABLED=1 go test -v ./...

# Run specific test
test-run:
	@read -p "Enter test name pattern: " pattern; \
	CGO_ENABLED=1 go test -v -run $$pattern ./...

# Run benchmark tests
test-bench:
	CGO_ENABLED=1 go test -v -bench=. ./...