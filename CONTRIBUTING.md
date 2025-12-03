# Contributing to Bragdoc

Thank you for your interest in contributing to Bragdoc! We welcome contributions from the community and are grateful for any help you can provide.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Smoke Testing](#smoke-testing)
- [Pull Request Process](#pull-request-process)
- [Architecture Decision Records](#architecture-decision-records)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/bragdoc.git
   cd bragdoc
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/vagnerclementino/bragdoc.git
   ```

## Development Setup

### Prerequisites

- **Go 1.21.1** or higher
- **Make** (for build automation)
- **Git** (for version control)

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify installation
make test
```

### Build the Project

```bash
# Build the binary
make build

# Run the binary
./bragdoc --help
```

## How to Contribute

### Reporting Bugs

If you find a bug, please [open an issue](https://github.com/vagnerclementino/bragdoc/issues) with:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Your environment (OS, Go version, etc.)
- Any relevant logs or screenshots

### Suggesting Features

We welcome feature suggestions! Please [open an issue](https://github.com/vagnerclementino/bragdoc/issues) with:

- A clear description of the feature
- The problem it solves
- Potential implementation approach (optional)
- Any relevant examples or mockups

### Submitting Changes

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our [coding standards](#coding-standards)

3. **Write tests** for your changes

4. **Run tests** to ensure everything works:
   ```bash
   make test
   make smoke
   ```

5. **Commit your changes** with clear, descriptive messages:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request** on GitHub

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names

### Project Structure

Bragdoc follows Clean Architecture principles:

```
├── cmd/                    # Application entry points
│   └── cli/               # CLI application
├── internal/              # Private application code
│   ├── command/           # CLI commands (organized by entity)
│   ├── domain/            # Domain entities (pure data structures)
│   ├── service/           # Business logic and validations
│   ├── repository/        # Data access layer
│   └── database/          # Database setup and migrations
├── config/                # Configuration management
└── docs/adr/             # Architecture Decision Records
```

### Key Principles

1. **Domain Entities**: Pure data structures with no business logic
2. **Service Layer**: Contains all business validations and logic
3. **Repository Layer**: Handles data persistence only
4. **Command Layer**: CLI interface and user interaction
5. **Dependency Injection**: Services receive dependencies via constructors

### Code Quality

Before submitting, ensure your code passes:

```bash
# Format code
make fmt

# Run linter
make lint

# Run static analysis
make vet

# Check imports
make imports

# Run all quality checks
make quality
```

## Testing

Bragdoc uses a combination of unit tests and integration tests with golden files.

### Running Tests

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run only unit tests (fast)
make test-unit

# Update golden files (when CLI output changes)
make update-golden
```

### Test Structure

- **Unit Tests**: Test individual functions and components in isolation
- **Integration Tests**: Test the CLI application end-to-end using golden files

### Integration Testing with Golden Files

Integration tests build a test binary and run it like a real user would:

```go
func TestCLIVersion(t *testing.T) {
    output, err := runBinary([]string{"version"}, nil)
    if err != nil {
        t.Fatal(err)
    }

    if *update {
        writeFixture(t, "version.golden", output)
    }

    actual := string(output)
    expected := loadFixture(t, "version.golden")

    if !reflect.DeepEqual(actual, expected) {
        t.Fatalf("actual = %s, expected = %s", actual, expected)
    }
}
```

Golden files are stored in `testdata/golden/` and contain expected CLI outputs.

### Adding New Integration Tests

1. Add test function in `cmd/cli/main_test.go`
2. Run with `-update` flag to generate golden file:
   ```bash
   go test ./cmd/cli -update
   ```
3. Review the generated golden file
4. Run tests normally to verify:
   ```bash
   make test
   ```

### Test Coverage Goals

- **CLI Package**: Minimum 60% coverage
- **Services**: Minimum 80% coverage
- **Commands**: Minimum 65% coverage
- **Database**: Minimum 45% coverage

Current coverage can be checked with:
```bash
make test-coverage
```

### Writing Tests

- Use table-driven tests when appropriate
- Mock external dependencies using interfaces
- Test both success and error cases
- Use descriptive test names
- Keep golden files small and focused

Example unit test:

```go
func TestBragService_Create_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockBragRepository)
    service := service.NewBragService(mockRepo)
    
    // Act
    result, err := service.Create(ctx, brag)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Updating Golden Files

When you intentionally change CLI output:

```bash
make update-golden
git diff testdata/golden/  # Review changes
make test                  # Verify tests pass
```

## Smoke Testing

Smoke tests verify core functionality works as expected. These tests catch major issues quickly.

### Running Smoke Tests

```bash
# Run smoke tests for your platform
make smoke

# Run for specific platform
TARGET_OS=darwin/amd64 make smoke
TARGET_OS=linux/amd64 make smoke
```

### Smoke Test Coverage

The smoke test suite covers:

1. **Help Commands** - `--help`, `version`
2. **Initialization** - `init` command
3. **Brag Management** - Create, list, edit, remove
4. **Filtering** - By category and tags
5. **Output Formats** - Table, JSON, YAML
6. **Error Handling** - Validation errors

### Adding Smoke Tests

When adding new features, update `smoke.sh`:

```bash
# Step X: Test new feature
echo -e "${BLUE}Step X: Testing new feature...${NC}"
if ./${BINARY_NAME} new-command 2>&1 | grep -q "expected output"; then
    print_result "new feature test" 0
else
    print_result "new feature test" 1 "Failed to execute new feature"
fi
```

### Troubleshooting Smoke Tests

If tests fail:

```bash
# Clean up test artifacts
rm -rf ./.bragdoc
rm -f ./bragdoc

# Run again
make smoke
```

## Pull Request Process

1. **Update documentation** if you're changing functionality
2. **Add tests** for new features or bug fixes
3. **Run all tests** and ensure they pass
4. **Update the README** if needed
5. **Reference any related issues** in your PR description
6. **Request review** from maintainers

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Smoke tests passing
- [ ] Documentation updated
- [ ] Commit messages are clear and descriptive
- [ ] No merge conflicts with main branch

### Review Process

- Maintainers will review your PR within a few days
- Address any feedback or requested changes
- Once approved, a maintainer will merge your PR

## Architecture Decision Records

We use ADRs to document significant architectural decisions. When making architectural changes:

1. Create a new ADR in `docs/adr/`
2. Use the template from existing ADRs
3. Number it sequentially (e.g., `0010-your-decision.md`)
4. Include:
   - Context and problem statement
   - Decision drivers
   - Considered options
   - Decision outcome
   - Consequences

Example:

```markdown
# 10. Your Decision Title

Date: 2024-01-15

## Status

Accepted

## Context

Describe the context and problem...

## Decision

Describe the decision...

## Consequences

Describe the consequences...
```

## Development Commands

### Build Commands

```bash
make build          # Build the binary
make run            # Build and run
make clean          # Remove build artifacts
make install        # Install to /usr/local/bin
```

### Quality Commands

```bash
make test           # Run tests with coverage
make test-race      # Test for race conditions
make lint           # Run golangci-lint
make fmt            # Format code
make vet            # Run go vet
make imports        # Run goimports
make quality        # Run all quality checks
```

### Packaging Commands

```bash
make package        # Create distribution packages
make tidy           # Clean up dependencies
```

## Getting Help

- **Documentation**: Check the [docs/](docs/) directory
- **Issues**: Search [existing issues](https://github.com/vagnerclementino/bragdoc/issues)
- **Discussions**: Start a [discussion](https://github.com/vagnerclementino/bragdoc/discussions)

## Recognition

Contributors will be recognized in:

- The project README
- Release notes
- GitHub contributors page

Thank you for contributing to Bragdoc! 🎉
