# 7. Adopt Clean Architecture

Date: 2025-12-01

## Status

Accepted

## Context

Bragdoc will evolve from a simple CLI to include TUI and Web interfaces. We need an architecture that:

- Separates business logic from interface concerns
- Allows adding new interfaces without changing core logic
- Makes testing easier by isolating dependencies
- Provides clear boundaries between layers
- Scales from simple to complex features

Without proper architecture, we risk:
- Business logic scattered across UI code
- Difficulty testing without UI
- Code duplication across interfaces
- Tight coupling making changes risky

## Decision

We will adopt **Clean Architecture** (also known as Hexagonal Architecture or Ports and Adapters).

**Layer Structure:**

1. **Domain Layer** (`internal/{domain}/`)
   - Entities with business rules
   - Pure Go structs, no external dependencies
   - Example: `Brag`, `User`, `Tag` entities

2. **Use Case Layer** (`internal/{domain}/service.go`)
   - Business logic and workflows
   - UseCase interfaces define operations
   - Service structs implement UseCases
   - Depends only on domain and repository interfaces

3. **Repository Layer** (`internal/{domain}/repository.go`)
   - Data access interfaces
   - Implemented by infrastructure layer
   - Example: `BragRepository`, `UserRepository`

4. **Infrastructure Layer** (`internal/database/`, `internal/config/`)
   - Concrete implementations of repositories
   - Database connections, file I/O
   - External service integrations

5. **Interface Layer** (`cmd/`, `internal/command/`)
   - CLI commands, TUI screens, Web handlers
   - Depends on use cases, not implementations
   - Multiple interfaces share same business logic

**Dependency Rule:**
- Dependencies point inward: Interface → UseCase → Domain
- Inner layers never depend on outer layers
- Use dependency injection for flexibility

## Consequences

**Positive:**
- **Testability**: Business logic testable without database or UI
- **Flexibility**: Easy to add new interfaces (CLI, TUI, Web)
- **Maintainability**: Clear separation of concerns
- **Reusability**: Same business logic across all interfaces
- **Independence**: Can change database, UI, or external services without affecting business logic
- **Team scalability**: Different team members can work on different layers

**Negative:**
- **Initial complexity**: More files and interfaces than simple approach
- **Learning curve**: Team needs to understand architecture principles
- **Boilerplate**: More code for simple operations

**Example Structure:**
```
internal/
├── brag/
│   ├── entity.go          # Domain: Brag entity
│   ├── service.go         # UseCase: Business logic
│   └── repository.go      # Interface: Data access contract
├── database/
│   └── brag_repository.go # Infrastructure: Repository implementation
cmd/
├── cli/
│   └── commands/
│       └── brag.go        # Interface: CLI commands
└── api/
    └── handlers/
        └── brag.go        # Interface: Web handlers
```

**Validation:**
- Business validations in service layer
- Structural validations in entity methods
- Input validations in interface layer

**Related Decisions:**
- See ADR-0008 for incremental interface development strategy
