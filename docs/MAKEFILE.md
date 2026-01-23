# Makefile Guide

This document describes all available Make targets and how to use them effectively.

## Quick Reference

```bash
make help           # Show all available targets
make build          # Build the application
make test           # Run tests
make quality        # Run all quality checks
```

## Target Categories

### Application Targets

Build and run the application:

```bash
# Build the binary
make build
# Output: ./bragdoc

# Build and run immediately
make run

# Install to system path
sudo make install
# Installs to: /usr/local/bin/bragdoc

# Create distribution packages
make package
# Creates: bragdoc.zip and bragdoc.tar.gz

# Create and push release tag
make release VERSION=v1.0.0
# Creates tag and triggers GitHub Actions

# Clean all build artifacts
make clean
# Removes: binary, packages, coverage files
```

### Quality Targets

Ensure code quality:

```bash
# Run all tests with coverage
make test
# Generates: coverage.txt

# Check for race conditions
make test-race

# Run linter (auto-installs golangci-lint if needed)
make lint

# Format code
make fmt

# Static analysis
make vet

# Organize imports
make imports

# Run ALL quality checks
make quality
# Runs: test, test-race, fmt, vet, imports, lint

# Run smoke tests
make smoke
# Executes: ./smoke.sh
```

### Helper Targets

Development utilities:

```bash
# Generate SQLC code
make generate

# Clean up go.mod
make tidy

# Update golden test files
make update-golden
# Use when CLI output intentionally changes

# Show help
make help
```

## Common Workflows

### Starting Development

```bash
# Clean slate
make clean

# Generate code
make generate

# Verify everything works
make test
```

### Before Committing

```bash
# Run all quality checks
make quality

# If quality passes, run smoke tests
make smoke

# Build final binary
make build
```

### Fixing Lint Issues

```bash
# Run individual checks
make fmt        # Fix formatting
make imports    # Fix imports
make vet        # Check for issues
make lint       # Final validation
```

### Updating Tests

```bash
# When CLI output changes intentionally
make update-golden

# Verify changes
git diff testdata/golden/

# Run tests to confirm
make test
```

### Creating a Release

```bash
# Build and package
make clean
make build
make package

# Verify binaries
unzip -l bragdoc.zip
tar -tzf bragdoc.tar.gz

# Create release (pushes tag and triggers CI/CD)
make release VERSION=v1.0.0
```

## Configuration

### Makefile Variables

Defined in `Makefile.vars`:

```makefile
BINARY_NAME := bragdoc
GOPATH      := $(shell go env GOPATH)
VERSION     := $(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD       := $(shell git rev-parse --short HEAD 2> /dev/null || echo "undefined")
LDFLAGS     := "-X 'github.com/vagnerclementino/bragdoc/internal/command.Version=$(VERSION)' -X 'github.com/vagnerclementino/bragdoc/internal/command.Build=$(BUILD)'"
GOOS        := darwin
GOARCH      := amd64
```

### Customizing Build

```bash
# Build for different platform
GOOS=linux GOARCH=amd64 make build

# Override version
VERSION=1.0.0 make build
```

## Troubleshooting

### Lint Fails

```bash
# Install/update golangci-lint manually
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

# Run lint again
make lint
```

### Tests Fail

```bash
# Clean and rebuild
make clean
make generate
make test

# Check specific package
go test ./internal/service -v
```

### Build Fails

```bash
# Verify dependencies
make tidy

# Check Go version
go version  # Should be 1.23+

# Clean and rebuild
make clean
make build
```

### Golden Files Out of Sync

```bash
# Update golden files
make update-golden

# Review changes
git diff testdata/golden/

# Commit if intentional
git add testdata/golden/
git commit -m "chore: update golden files"
```

## Advanced Usage

### Silent Mode

All targets run in silent mode by default (`.SILENT:` in Makefile).

### Parallel Execution

```bash
# Run multiple targets (not recommended for quality)
make clean build test
```

### Dependency Chain

Some targets have dependencies:

- `build` → `generate`
- `run` → `build`
- `install` → `build`
- `package` → `build`
- `quality` → `test`, `test-race`, `fmt`, `vet`, `imports`, `lint`

### Ignoring Unknown Commands

The Makefile ignores unknown commands to prevent errors:

```makefile
%:
    @:
```

## Integration with CI/CD

### GitHub Actions

The Makefile targets are used in CI/CD pipelines:

**Quality Pipeline** (`.github/workflows/quality.yml`):
```yaml
- name: Run tests
  run: go test ./... -v -coverprofile=coverage.txt -covermode=atomic
```

**Release Pipeline** (`.github/workflows/release.yml`):
```yaml
- name: Building binary
  run: go build -o bragdoc-macos-intel -ldflags "..." ./cmd/cli
```

### Local CI Simulation

```bash
# Simulate quality pipeline
make clean
make tidy
make quality

# Simulate release build
make clean
make build
make package
```

## Best Practices

1. **Always run `make help`** first to see available targets
2. **Run `make quality`** before committing
3. **Use `make clean`** when switching branches
4. **Run `make generate`** after modifying SQL queries
5. **Update golden files** only when output changes are intentional
6. **Run `make smoke`** before creating PRs
7. **Use `make tidy`** to keep dependencies clean

## See Also

- [CONTRIBUTING.md](CONTRIBUTING.md) - Full contribution guide
- [README.md](README.md) - Project overview
- [GETTING_STARTED.md](GETTING_STARTED.md) - User guide
