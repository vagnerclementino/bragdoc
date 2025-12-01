# 8. Incremental interface development CLI to TUI to Web

Date: 2025-12-01

## Status

Accepted

## Context

Bragdoc aims to provide three interfaces: CLI, TUI, and Web. We need to decide the development order and strategy.

Options:
1. **Parallel development**: Build all interfaces simultaneously
2. **User-driven**: Start with most requested interface
3. **Incremental**: Build interfaces in order of complexity
4. **MVP-first**: Build simplest interface, validate, then expand

Considerations:
- Time to market: Users need value quickly
- Validation: Need to validate product-market fit
- Architecture: Want to ensure business logic is interface-independent
- Resources: Limited development capacity
- Risk: Minimize risk of building wrong thing

## Decision

We will develop interfaces **incrementally** in this order:

**Phase 1: CLI Interface** (MVP - v1.0)
- Complete command-line interface with all core features
- Establishes business logic independent of interface
- Delivers immediate value to users
- Validates product-market fit
- Foundation for all future interfaces

**Phase 2: TUI Interface** (v1.1)
- Terminal UI with Bubbletea for interactive experience
- Reuses all business logic from Phase 1
- Improves user experience for interactive workflows
- Validates architecture's interface independence

**Phase 3: Web Interface** (v2.0)
- HTTP server with REST API and web UI
- Leverages battle-tested services and repositories
- Enables collaborative features and remote access
- Expands use cases beyond single-user CLI

**Phase 4: Feature Expansion** (v2.1+)
- Additional configuration formats (JSON, TOML)
- Full internationalization (Portuguese, Spanish, French)
- Advanced AI features
- Performance optimizations

## Consequences

**Positive:**
- **Fast MVP**: CLI delivers value in weeks, not months
- **Risk reduction**: Validate product before investing in complex interfaces
- **Architecture validation**: Proves business logic is truly interface-independent
- **Incremental learning**: Each phase informs the next
- **User feedback**: Early users guide feature priorities
- **Stable foundation**: Later interfaces built on tested business logic
- **Clear milestones**: Each phase is a shippable product

**Negative:**
- **Delayed features**: TUI and Web users must wait
- **Potential rework**: Early feedback might require changes
- **Feature parity**: Need to maintain consistency across interfaces

**Development Strategy:**

Phase 1 (CLI) focuses on:
- Core domain models (Brag, User, Tag)
- Business logic in services
- Repository interfaces and implementations
- Database with SQLite + SQLC
- Configuration management
- AI integration
- Document generation
- Complete CLI commands

Phase 2 (TUI) adds:
- Bubbletea models and views
- Interactive forms and navigation
- Real-time search and filtering
- Reuses all Phase 1 services

Phase 3 (Web) adds:
- HTTP server and routing
- REST API endpoints
- Web UI with server-side rendering
- Authentication and sessions
- Reuses all Phase 1 and 2 services

**Success Criteria:**
- Phase 1: 100+ active CLI users, positive feedback
- Phase 2: TUI requires zero changes to business logic
- Phase 3: Web interface reuses 90%+ of existing code

**Related Decisions:**
- See ADR-0007 for Clean Architecture enabling this strategy
- See ADR-0003 and ADR-0004 for v1 scope limitations
