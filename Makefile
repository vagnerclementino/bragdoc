.PHONY: test test-unit test-integration test-coverage update-golden clean build smoke

# Run all tests
test: build-with-coverage
	@rm -rf .coverdata
	@mkdir -p .coverdata
	@go test ./...

# Run only unit tests (fast)
test-unit:
	@go test -short ./...

# Build binary with coverage instrumentation
build-with-coverage:
	@echo "Building test binary with coverage..."
	@go build -cover -o bragdoc-test -ldflags "-X 'github.com/vagnerclementino/bragdoc/internal/command.Version=0.1.0' -X 'github.com/vagnerclementino/bragdoc/internal/command.Build=test'" ./cmd/cli
	@rm -f bragdoc-test

# Check test coverage
test-coverage: test
	@echo "\n📊 Coverage Report:"
	@go tool covdata percent -i=.coverdata

# Update golden files
update-golden:
	@echo "Updating golden files..."
	@go test ./cmd/cli -update

# Run smoke tests
smoke: build
	@echo "🔥 Running smoke tests..."
	@./smoke.sh

# Clean build artifacts and coverage data
clean:
	@rm -rf .coverdata
	@rm -f bragdoc-test
	@rm -f bragdoc
	@echo "✨ Cleaned build artifacts"

# Build production binary
build:
	@go build -o bragdoc ./cmd/cli
	@echo "✅ Built bragdoc binary"
