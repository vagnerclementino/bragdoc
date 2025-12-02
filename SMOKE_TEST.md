# Smoke Test Documentation

## Overview

The smoke test suite (`smoke.sh`) performs basic functional testing of the Bragdoc CLI to verify core functionality works as expected. These tests are designed to catch major issues quickly before more comprehensive testing.

## Running Smoke Tests

### Local Development

Run smoke tests for your current platform:

```bash
make smoke
```

### CI/CD Pipeline

Specify target OS using the `TARGET_OS` environment variable:

```bash
# macOS Intel
TARGET_OS=darwin/amd64 make smoke

# macOS ARM (Apple Silicon)
TARGET_OS=darwin/arm64 make smoke

# Linux AMD64
TARGET_OS=linux/amd64 make smoke

# Linux ARM64
TARGET_OS=linux/arm64 make smoke
```

### Direct Script Execution

You can also run the script directly:

```bash
./smoke.sh [target_os]

# Examples:
./smoke.sh darwin/arm64
./smoke.sh linux/amd64
```

## Test Coverage

The smoke test suite covers the following scenarios:

### 1. Help Commands (No Initialization Required)
- `bragdoc --help`
- `bragdoc version`
- `bragdoc brag --help`

### 2. Initialization Requirement
- Verifies commands fail appropriately without initialization
- Tests `brag list` and `brag add` rejection

### 3. Initialization
- `bragdoc init` with required parameters
- Verifies configuration and database creation

### 4. Brag Creation
- Add brag with tags
- Add brag with different categories (achievement, leadership, innovation)
- Verify successful creation

### 5. Brag Listing
- List in table format (default)
- List in JSON format
- List in YAML format

### 6. Brag Filtering
- Filter by category
- Filter by tags

### 7. Brag Display
- Show single brag by ID
- Show multiple brags (comma-separated IDs)
- Show range of brags (ID range)

### 8. Brag Editing
- Edit brag title
- Verify edit persistence

### 9. Brag Removal
- Remove brag with `--force` flag
- Verify removal

### 10. Error Handling
- Validation error for short title
- Validation error for short description
- Validation error for invalid category

## Test Results

The script provides colored output:
- ✓ Green checkmark for passed tests
- ✗ Red X for failed tests
- Summary with total, passed, and failed counts

## Exit Codes

- `0`: All tests passed
- `1`: One or more tests failed

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Smoke Tests

on: [push, pull_request]

jobs:
  smoke-test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
        include:
          - os: ubuntu-latest
            target: linux/amd64
          - os: macos-latest
            target: darwin/arm64
    
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Run Smoke Tests
        env:
          TARGET_OS: ${{ matrix.target }}
        run: make smoke
```

### GitLab CI Example

```yaml
smoke-test:
  stage: test
  parallel:
    matrix:
      - TARGET_OS: [darwin/amd64, darwin/arm64, linux/amd64, linux/arm64]
  script:
    - make smoke
  artifacts:
    when: on_failure
    paths:
      - bragdoc
```

## Troubleshooting

### Tests Fail After Previous Run

The script automatically cleans up, but if you encounter issues:

```bash
# Clean up manually
rm -rf ./.bragdoc
rm -f ./bragdoc

# Run again
make smoke
```

### Cross-Compilation Issues

If cross-compilation fails, ensure you have the necessary toolchain:

```bash
# macOS: Install Xcode Command Line Tools
xcode-select --install

# Linux: Install build-essential
sudo apt-get install build-essential
```

### CGO Requirements

Bragdoc uses SQLite which requires CGO. Ensure `CGO_ENABLED=1` when building.

## Adding New Tests

To add new smoke tests:

1. Add a new test step in `smoke.sh`
2. Use the `run_test` or `print_result` functions
3. Follow the existing pattern for consistency
4. Update this documentation

Example:

```bash
# Step X: Test new feature
echo -e "${BLUE}Step X: Testing new feature...${NC}"
if ./${BINARY_NAME} new-command 2>&1 | grep -q "expected output"; then
    print_result "new feature test" 0
else
    print_result "new feature test" 1 "Failed to execute new feature"
fi
echo ""
```

## Maintenance

- Review and update tests when adding new features
- Keep test execution time reasonable (< 2 minutes)
- Ensure tests are idempotent and don't affect system state
- Use temporary directories for all test data
