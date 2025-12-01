# 2. Use SQLite as database

Date: 2025-12-01

## Status

Accepted

## Context

Bragdoc needs a reliable, lightweight database solution for storing user data, brags, tags, and their relationships. The application is a CLI tool that should work seamlessly on user machines without requiring external database servers or complex setup procedures.

Key requirements:
- Zero-configuration deployment
- Single-file database for easy backup and portability
- ACID compliance for data integrity
- Cross-platform support (macOS, Linux, Windows)
- Sufficient performance for personal use (thousands of records)
- No external dependencies or server processes

## Decision

We will use SQLite as the database engine for Bragdoc.

SQLite will be embedded directly into the application binary, providing:
- A single `.db` file stored in `~/.bragdoc/bragdoc.db`
- Full SQL support with ACID transactions
- Native Go support via `database/sql` and SQLite drivers
- Automatic migrations on application startup
- Built-in backup capabilities

## Consequences

**Positive:**
- **Zero setup**: Users don't need to install or configure a database server
- **Portability**: The entire database is a single file that can be easily backed up or moved
- **Reliability**: SQLite is battle-tested and provides ACID guarantees
- **Performance**: More than sufficient for personal use cases (millions of rows)
- **Simplicity**: Standard SQL interface, well-documented, mature ecosystem
- **Cross-platform**: Works identically on macOS, Linux, and Windows

**Negative:**
- **Concurrency limitations**: Not suitable for high-concurrency scenarios (not a concern for CLI tool)
- **No network access**: Cannot be accessed remotely (acceptable for personal tool)
- **Single writer**: Only one process can write at a time (not an issue for single-user CLI)

**Risks and Mitigations:**
- **Data corruption**: Mitigated by implementing automatic backups before write operations
- **Migration failures**: Mitigated by version tracking and rollback capabilities
- **File locking**: Mitigated by proper transaction handling and connection management
