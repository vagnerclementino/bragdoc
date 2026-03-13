# 6. Use SQLC for type-safe queries

Date: 2025-12-01

## Status

Accepted

## Context

When working with SQL databases in Go, there are several approaches for writing queries:

1. **Raw SQL with database/sql**: Maximum control but no type safety, prone to errors
2. **ORM (GORM, Ent)**: High-level abstraction but adds complexity and magic
3. **Query builders**: Programmatic query construction but still runtime errors
4. **SQLC**: Generates type-safe Go code from SQL queries

Requirements:
- Type safety at compile time
- Performance (no reflection overhead)
- Maintainability (clear SQL, clear Go code)
- Simplicity (minimal magic, easy to debug)
- SQLite compatibility

## Decision

We will use **SQLC** for generating type-safe database access code.

Implementation approach:
- Write SQL queries in `.sql` files in `internal/database/queries/`
- Write schema migrations in `internal/database/migrations/`
- SQLC generates Go code with type-safe functions
- Generated code uses standard `database/sql` interfaces
- No runtime reflection or ORM overhead

## Consequences

**Positive:**
- **Compile-time safety**: Type errors caught at compile time, not runtime
- **Performance**: No reflection, direct SQL execution
- **Clarity**: SQL is SQL, Go is Go - no DSL to learn
- **Maintainability**: Easy to understand generated code
- **Tooling**: Great IDE support for both SQL and generated Go code
- **Testability**: Easy to test with real database or mocks
- **Migration friendly**: Schema changes automatically update generated code

**Negative:**
- **Build step**: Requires running `sqlc generate` when queries change
- **Learning curve**: Team needs to learn SQLC conventions
- **Generated code**: Need to commit generated code or generate in CI

**Workflow:**
```bash
# 1. Write SQL query
# internal/database/queries/brags.sql
-- name: GetBrag :one
SELECT * FROM brags WHERE id = ? LIMIT 1;

# 2. Generate Go code
sqlc generate

# 3. Use type-safe function
brag, err := queries.GetBrag(ctx, id)
```

**Configuration:**
```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "sqlite"
    queries: "internal/database/queries/"
    schema: "internal/database/migrations/"
    gen:
      go:
        package: "queries"
        out: "internal/database/queries"
        emit_json_tags: true
        emit_interface: true
```

**Alternatives Considered:**
- **GORM**: Too much magic, harder to debug, performance overhead
- **Raw SQL**: No type safety, error-prone
- **Ent**: Too complex for our needs, opinionated schema management
