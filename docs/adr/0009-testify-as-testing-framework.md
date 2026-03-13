# 9. Testify as testing framework

Date: 2025-12-01

## Status

Accepted

## Context

Go has built-in testing support with the `testing` package, but additional frameworks can improve test readability, reduce boilerplate, and provide better assertions and mocking capabilities.

Options for testing in Go:
1. **Standard library only**: Minimal dependencies but verbose assertions
2. **Testify**: Popular framework with assertions, mocks, and suites
3. **Ginkgo/Gomega**: BDD-style testing with rich matchers
4. **GoConvey**: Web UI and BDD-style assertions

Requirements:
- Clear, readable test assertions
- Easy mocking for interfaces
- Table-driven test support
- Minimal learning curve
- Good community support
- Compatible with standard Go tooling

## Decision

We will use **Testify** as our testing framework.

Specifically, we'll use:
- `testify/assert`: Readable assertions
- `testify/require`: Assertions that stop test on failure
- `testify/mock`: Interface mocking
- `testify/suite`: Test suite organization (when needed)

## Consequences

**Positive:**
- **Readable assertions**: `assert.Equal(t, expected, actual)` vs `if expected != actual { t.Errorf(...) }`
- **Better error messages**: Testify provides detailed failure output
- **Easy mocking**: Generate mocks from interfaces automatically
- **Less boilerplate**: Reduces repetitive test code
- **Industry standard**: Most popular Go testing library (20k+ GitHub stars)
- **Good documentation**: Extensive examples and community resources
- **Standard compatible**: Works with `go test` and all standard tooling

**Negative:**
- **External dependency**: Adds one more dependency to project
- **Learning curve**: Team needs to learn Testify API (minimal)
- **Mock generation**: Requires mockery tool or manual mock creation

**Example Usage:**

```go
// Without Testify
func TestCreateBrag(t *testing.T) {
    brag, err := service.Create(ctx, input)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if brag.Title != "Expected Title" {
        t.Errorf("expected title %q, got %q", "Expected Title", brag.Title)
    }
}

// With Testify
func TestCreateBrag(t *testing.T) {
    brag, err := service.Create(ctx, input)
    require.NoError(t, err)
    assert.Equal(t, "Expected Title", brag.Title)
}
```

**Mocking Example:**

```go
// Define mock
type MockBragRepository struct {
    mock.Mock
}

func (m *MockBragRepository) Insert(ctx context.Context, brag *Brag) (*Brag, error) {
    args := m.Called(ctx, brag)
    return args.Get(0).(*Brag), args.Error(1)
}

// Use in test
func TestBragService(t *testing.T) {
    mockRepo := new(MockBragRepository)
    mockRepo.On("Insert", mock.Anything, mock.Anything).Return(&Brag{ID: 1}, nil)
    
    service := NewService(mockRepo)
    result, err := service.Create(ctx, input)
    
    assert.NoError(t, err)
    assert.Equal(t, int64(1), result.ID)
    mockRepo.AssertExpectations(t)
}
```

**Testing Strategy:**
- Unit tests: Test services with mocked repositories
- Integration tests: Test repositories with real SQLite database
- Table-driven tests: Use Testify with Go's table-driven pattern
- Coverage target: 80%+ for domain and service layers

**Installation:**
```bash
go get github.com/stretchr/testify
```

**Alternatives Considered:**
- **Standard library only**: Too verbose, harder to read
- **Ginkgo/Gomega**: More complex, BDD style not needed
- **GoConvey**: Web UI unnecessary, less popular
