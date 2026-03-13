# 3. YAML only for configuration in v1

Date: 2025-12-01

## Status

Accepted

## Context

Bragdoc requires a configuration file format for storing user preferences, AI settings, prompts, and other application settings. While the Viper library we're using supports multiple formats (YAML, JSON, TOML), we need to decide which format(s) to support in the first version.

Considerations:
- Time to market: Supporting multiple formats adds complexity
- User experience: Most users are comfortable with YAML
- Maintainability: Fewer formats mean simpler testing and documentation
- Future flexibility: Architecture should allow adding formats later

## Decision

Version 1 (v1) of Bragdoc will support **YAML only** for configuration files.

The configuration file will be located at `~/.bragdoc/config.yaml`.

Future versions may add support for JSON and TOML formats based on user demand.

## Consequences

**Positive:**
- **Faster MVP delivery**: Reduces scope and complexity for initial release
- **Simpler documentation**: Only need to document one format
- **Easier testing**: Fewer test cases and edge cases to handle
- **Human-readable**: YAML is widely understood and easy to edit manually
- **Good defaults**: YAML supports comments, multi-line strings, and is less verbose than JSON

**Negative:**
- **Limited choice**: Users who prefer JSON or TOML must wait for future versions
- **YAML quirks**: YAML has some parsing edge cases (indentation sensitivity, type coercion)

**Future Work:**
- v2 can add JSON support with minimal changes (Viper already supports it)
- v3 can add TOML support if requested
- The `bragdoc init` command will eventually support a `--format` flag

**Technical Notes:**
- Architecture uses Viper, which abstracts format handling
- Adding new formats later requires minimal code changes
- Configuration struct remains format-agnostic
