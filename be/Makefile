.PHONY: lint lint-fix build run clean test

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

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...