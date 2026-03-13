# 10. Repository Dependency Inversion

Date: 2025-12-03

## Status

Accepted

## Context

Following Clean Architecture principles (ADR-0007), we need to ensure proper separation between the repository interfaces and their implementations. The Dependency Inversion Principle states that:

- High-level modules should not depend on low-level modules. Both should depend on abstractions.
- Abstractions should not depend on details. Details should depend on abstractions.

Initially, repository implementations were mixed with their interfaces in the `internal/repository` package, violating this principle. This created tight coupling between the domain layer and the database implementation details.

## Decision

We will separate repository interfaces from their implementations:

**Repository Interfaces** (`internal/repository/`)
- Define contracts for data access
- No implementation details
- No database-specific code
- Pure Go interfaces

**Repository Implementations** (`internal/database/`)
- SQLite-specific implementations
- Database connection management
- Query execution details
- Data mapping logic

**Structure:**
```text
internal/
├── repository/
│   ├── user.go          # UserRepository interface
│   ├── brag.go          # BragRepository interface
│   └── tag.go           # TagRepository interface
└── database/
    ├── sqlite.go        # SQLite connection wrapper
    ├── user.go          # SQLite UserRepository implementation
    ├── brag.go          # SQLite BragRepository implementation
    └── tag.go           # SQLite TagRepository implementation
```

**Key Components:**

1. **SQLiteDB Wrapper** (`database/sqlite.go`)
   - Wraps `sql.DB` and `queries.Queries`
   - Provides transaction support
   - Centralizes database access

2. **Repository Implementations** (`database/*_repository.go`)
   - Implement repository interfaces
   - Handle database-specific operations
   - Convert between database and domain models

3. **Repository Interfaces** (`repository/*.go`)
   - Define data access contracts
   - No implementation details
   - Used by service layer

## Consequences

**Positive:**
- **Testability**: Easy to mock repositories for testing
- **Flexibility**: Can swap database implementations without changing business logic
- **Clarity**: Clear separation between interface and implementation
- **Maintainability**: Changes to database logic don't affect interface contracts
- **Scalability**: Easy to add new database implementations (PostgreSQL, MySQL, etc.)

**Negative:**
- **More files**: Separation creates additional files
- **Indirection**: One more layer to navigate
- **Initial setup**: Requires more upfront design

**Migration Path:**
1. Create `internal/database/sqlite.go` wrapper
2. Move implementations to `internal/database/*_repository.go`
3. Keep only interfaces in `internal/repository/*.go`
4. Update service layer to use new constructors
5. Update tests to use new structure

**Example Usage:**
```go
// In main.go or initialization code
db := database.Open(dbPath)
sqliteDB := database.NewSQLiteDB(db)

// Create repositories
userRepo := database.NewUserRepository(sqliteDB)
bragRepo := database.NewBragRepository(sqliteDB, userRepo)
tagRepo := database.NewTagRepository(sqliteDB)

// Create services (depend on interfaces, not implementations)
userService := service.NewUserService(userRepo)
bragService := service.NewBragService(bragRepo, userRepo)
tagService := service.NewTagService(tagRepo)
```

**Testing:**
```go
// Mock repository for testing
type mockUserRepository struct {
    repository.UserRepository
    // ... mock methods
}

// Test service with mock
func TestUserService(t *testing.T) {
    mockRepo := &mockUserRepository{}
    service := service.NewUserService(mockRepo)
    // ... test service logic
}
```

## Future Enhancements

### Turso Database Integration

The dependency inversion pattern enables easy integration with **Turso Database**, a distributed SQLite database built on libSQL:

**Benefits:**
- **SQLite compatibility**: Works with existing schema and SQLC queries
- **Edge deployment**: Database at the edge for low latency
- **Multi-region replication**: Automatic data replication across regions
- **Embedded replicas**: Local-first architecture with automatic sync
- **Same codebase**: Reuse SQLite repository implementations

**Implementation approach:**
```go
// internal/database/turso.go
type TursoDB struct {
    db      *sql.DB
    queries *queries.Queries
}

func NewTursoDB(url, authToken string) (*TursoDB, error) {
    db, err := sql.Open("libsql", url+"?authToken="+authToken)
    if err != nil {
        return nil, err
    }
    return &TursoDB{db: db, queries: queries.New(db)}, nil
}

// Reuse existing SQLite implementations
func NewTursoUserRepository(db *TursoDB) repository.UserRepository {
    return &sqliteUserRepository{db: db}
}
```

**Configuration:**
```bash
# Local SQLite (default)
bragdoc init

# Turso Database (future)
export TURSO_URL="libsql://your-database.turso.io"
export TURSO_AUTH_TOKEN="your-token"
bragdoc init --database turso
```

This architecture decision makes Turso integration straightforward without changing business logic or service layer code.

## Related Decisions

- ADR-0007: Adopt Clean Architecture
- ADR-0002: Use SQLite as database
- ADR-0006: Use SQLC for type-safe queries
