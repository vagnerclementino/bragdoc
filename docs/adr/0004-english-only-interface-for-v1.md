# 4. English only interface for v1

Date: 2025-12-01

## Status

Accepted

## Context

Bragdoc is being developed with potential for international users. While full internationalization (i18n) is desirable, it adds significant complexity to the initial development:

- Message catalogs for multiple languages
- Locale-specific formatting (dates, numbers)
- Testing across different languages
- Documentation in multiple languages
- Maintenance overhead for translations

The question is whether to implement i18n from the start or focus on core functionality first.

## Decision

Version 1 (v1) of Bragdoc will have an **English-only interface**.

This includes:
- CLI command names and descriptions
- Help text and error messages
- Log messages
- Documentation

The codebase will be structured to support i18n in future versions:
- Use an `internal/i18n` package with translation-ready architecture
- Keep user-facing strings separate from logic
- Design with localization in mind

## Consequences

**Positive:**
- **Faster time to market**: Focus on core features rather than translations
- **Simpler initial development**: No need for message catalogs or locale handling
- **Easier testing**: Single language reduces test matrix
- **Clearer MVP validation**: Validate product-market fit before investing in translations
- **Better foundation**: Can implement i18n properly after understanding user needs

**Negative:**
- **Limited audience**: Non-English speakers may find the tool harder to use
- **Potential rework**: Some strings may need refactoring for proper i18n later

**Future Work:**
- v2+ will add internationalization support
- Priority languages based on user demand: Portuguese (pt-BR), Spanish (es), French (fr)
- The `internal/i18n` package will provide translation functions
- Configuration will include `language` and `locale` settings

**Important Note:**
This decision applies only to the **interface** (commands, messages, help text). The **content** of brags and generated documents can be in any language - see ADR-0005 for details on document language handling.
