.PHONY: lint lint-fix lint-verbose lint-ci build run clean test dev debug

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
	CGO_ENABLED=1 go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Generate coverage report only
coverage-report:
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# View coverage in browser
coverage-view: coverage-report
	@echo "Opening coverage report in browser..."
	@if command -v xdg-open > /dev/null; then \
		xdg-open coverage.html; \
	elif command -v open > /dev/null; then \
		open coverage.html; \
	elif command -v start > /dev/null; then \
		start coverage.html; \
	else \
		echo "Please open coverage.html manually in your browser"; \
	fi

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

# Database migration commands
migrate-status:
	./atlas.exe migrate status --env dev

migrate-apply:
	./atlas.exe migrate apply --env dev

migrate-diff:
	@read -p "Enter migration name: " name; \
	./atlas.exe migrate diff $$name --env dev

migrate-validate:
	./atlas.exe migrate validate --env dev

migrate-down:
	@read -p "Enter number of migrations to rollback: " count; \
	./atlas.exe migrate down $$count --env dev