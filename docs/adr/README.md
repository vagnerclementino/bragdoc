# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the Bragdoc project.

## What are ADRs?

ADRs document important architectural decisions made during the development of the project. Each ADR describes the context, the decision made, and the consequences of that decision.

## Index

- [ADR-0001](0001-record-architecture-decisions.md) - Record architecture decisions
- [ADR-0002](0002-use-sqlite-as-database.md) - Use SQLite as database
- [ADR-0003](0003-yaml-only-for-configuration-in-v1.md) - YAML only for configuration in v1
- [ADR-0004](0004-english-only-interface-for-v1.md) - English only interface for v1
- [ADR-0005](0005-language-defined-at-document-generation-time.md) - Language defined at document generation time
- [ADR-0006](0006-use-sqlc-for-type-safe-queries.md) - Use SQLC for type-safe queries
- [ADR-0007](0007-adopt-clean-architecture.md) - Adopt Clean Architecture
- [ADR-0008](0008-incremental-interface-development-cli-to-tui-to-web.md) - Incremental interface development (CLI → TUI → Web)
- [ADR-0009](0009-testify-as-testing-framework.md) - Testify as testing framework

## Key Decisions Summary

### Database & Persistence
- **SQLite** for embedded database (ADR-0002)
- **SQLC** for type-safe query generation (ADR-0006)

### Configuration & Localization
- **YAML only** in v1, JSON/TOML in future versions (ADR-0003)
- **English interface** in v1, i18n in future versions (ADR-0004)
- **Language at generation time** for document content (ADR-0005)

### Architecture & Development
- **Clean Architecture** for separation of concerns (ADR-0007)
- **Incremental development**: CLI → TUI → Web (ADR-0008)
- **Testify** for testing framework (ADR-0009)

## Creating New ADRs

To create a new ADR:

```bash
cd docs/adr
adr new "Title of the decision"
```

Then edit the generated file to fill in the Context, Decision, and Consequences sections.

## ADR Lifecycle

- **Proposed**: Decision is being considered
- **Accepted**: Decision has been made and is being implemented
- **Deprecated**: Decision is no longer relevant
- **Superseded**: Decision has been replaced by a newer ADR
