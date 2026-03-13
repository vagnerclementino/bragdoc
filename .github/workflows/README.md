# GitHub Actions Workflows

## Quality Pipeline

The `quality.yml` workflow ensures code quality and functionality through multiple stages of testing and validation.

### Trigger Events

- **Pull Requests** targeting `main` branch only

### Pipeline Stages

#### 1. Quality Checks
**Purpose**: Validate code quality and run unit tests

**Steps**:
- Check Go module consistency (`go mod tidy`)
- Run `golangci-lint` for code quality (warnings only, doesn't fail pipeline)
- Execute unit tests with race detection
- Generate and upload coverage reports to Codecov

**Runs on**: `ubuntu-latest`

#### 2. Smoke Tests
**Purpose**: Verify core functionality works on different platforms

**Steps**:
- Run smoke test suite on multiple OS platforms
- Test basic CLI operations (init, add, list, edit, remove)
- Verify error handling and validation

**Matrix**:
- Linux (AMD64)
- macOS (ARM64)



#### 3. Integration Tests
**Purpose**: Run integration and E2E tests

**Steps**:
- Execute integration tests
- Run end-to-end test suite

**Runs on**: `ubuntu-latest`

#### 4. Quality Summary
**Purpose**: Aggregate results and provide final status

**Steps**:
- Check all previous job results
- Report overall pipeline status
- Fail if any stage failed

### Concurrency

The pipeline uses concurrency groups to cancel in-progress runs when new commits are pushed to the same branch.

```yaml
concurrency:
  group: quality-${{ github.ref }}
  cancel-in-progress: true
```

### Dependencies Between Jobs

```
quality-checks
    ├── smoke-tests
    ├── integration-tests
    └── quality-summary
```

### Environment Variables

- `TARGET_OS`: Target operating system for smoke tests (e.g., `linux/amd64`, `darwin/arm64`)
- `CGO_ENABLED`: Set to `1` for SQLite support
- `GOOS`: Target operating system for build
- `GOARCH`: Target architecture for build

### Artifacts

- **Coverage Reports**: Uploaded to Codecov after unit tests

### Status Badges

Add these badges to your README.md:

```markdown
[![Quality Pipeline](https://github.com/vagnerclementino/bragdoc/actions/workflows/quality.yml/badge.svg)](https://github.com/vagnerclementino/bragdoc/actions/workflows/quality.yml)
```

### Local Testing

Before pushing, you can run the same checks locally:

```bash
# Check go mod
go mod tidy
git diff --exit-code go.mod go.sum

# Run linter
golangci-lint run --timeout 5m

# Run unit tests
go test ./... -v -race -coverprofile=coverage.txt

# Run smoke tests
make smoke

# Run integration tests
go test ./... -v -run Integration
```

### Troubleshooting

#### Linting Issues

The `golangci-lint` step is configured to not fail the pipeline (`continue-on-error: true` and `--issues-exit-code=0`). This means:
- Linting issues will be reported as warnings
- The pipeline will continue even if there are linting issues
- You should still review and fix linting warnings before merging

To run linting locally and fix issues:
```bash
golangci-lint run --timeout 5m --fix
```

#### Smoke Tests Fail on macOS

If smoke tests fail on macOS runners, check:
1. CGO is enabled
2. Xcode Command Line Tools are available
3. Binary has execute permissions

#### Tests Fail on Specific Platform

If tests fail on a specific platform:
1. Review the platform-specific logs in GitHub Actions
2. Test locally on that platform if possible
3. Check for platform-specific issues (CGO, dependencies, etc.)

#### Coverage Upload Fails

If Codecov upload fails:
1. Check if `CODECOV_TOKEN` secret is set (optional for public repos)
2. Verify coverage file is generated correctly
3. Check Codecov service status

### Optimization Tips

1. **Cache Go modules**: The workflow uses Go cache to speed up builds
2. **Parallel execution**: Jobs run in parallel when possible
3. **Fail-fast disabled**: Build matrix continues even if one platform fails
4. **Concurrency control**: Cancels outdated runs automatically

### Adding New Checks

To add new quality checks:

1. Add a new job or step in `quality.yml`
2. Update dependencies in the `needs` field
3. Update this documentation
4. Test locally before pushing

Example:

```yaml
security-scan:
  name: Security Scan
  needs: quality-checks
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Run Gosec
      uses: securego/gosec@master
      with:
        args: ./...
```

### Maintenance

- Review and update Go version regularly
- Keep action versions up to date
- Monitor pipeline execution time
- Adjust timeout values if needed
- Review artifact retention policies
