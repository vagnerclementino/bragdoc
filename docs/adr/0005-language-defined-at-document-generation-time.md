# 5. Language defined at document generation time

Date: 2025-12-01

## Status

Accepted

## Context

Users may want to generate brag documents in different languages for different purposes (e.g., English for international companies, Portuguese for local companies). We need to decide when and how language is determined for document content.

Two main approaches:
1. **Language per brag**: Each brag has a fixed language, documents inherit from brags
2. **Language at generation time**: Brags are language-agnostic, language chosen when generating document

Considerations:
- Flexibility: Users may need the same achievements in multiple languages
- Data model complexity: Storing language per brag adds fields and constraints
- AI capabilities: Modern AI can translate and adapt content effectively
- User workflow: How do users typically work with multilingual content?

## Decision

**Language is defined at document generation time**, not stored with individual brags.

Implementation:
- Brags are stored language-agnostic (users write in whatever language they prefer)
- Default language is configured in `config.yaml` (`user.language` field)
- `bragdoc doc generate` command accepts `--language` flag to override default
- AI translates and adapts content to target language during generation
- Same brags can be used to generate documents in different languages

## Consequences

**Positive:**
- **Maximum flexibility**: Same brag can appear in documents in different languages
- **Simpler data model**: No language field needed in brags table
- **Easier data entry**: Users write brags naturally without worrying about language
- **Better reusability**: One set of brags serves all language needs
- **AI-powered translation**: Leverages AI for natural, context-aware translation
- **Less duplication**: No need to maintain parallel brags in different languages

**Negative:**
- **Translation quality**: Depends on AI quality (mitigated by using good AI models)
- **No language mixing**: Can't have some brags in English and others in Portuguese in same document (acceptable trade-off)
- **AI dependency**: Document generation requires AI for non-native language output

**Example Workflow:**
```bash
# User writes brags in their preferred language (e.g., Portuguese)
bragdoc brag add --title "Implementei sistema de cache" --description "..."

# Generate document in English for international company
bragdoc doc generate --language en --output resume-en.pdf

# Generate document in Portuguese for local company
bragdoc doc generate --language pt --output resume-pt.pdf
```

**Configuration:**
```yaml
user:
  language: "en"  # Default language for document generation
```

**Related Decisions:**
- See ADR-0004 for interface language (English only in v1)
- Interface language and document language are independent concerns
