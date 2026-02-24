# DOC

## Available Commands

```bash
  generate    Generate a brag document
```

### generate

Generate a professional achievement document from your brags.

You can filter which brags to include using IDs, categories, or tags.
If no filters are specified, all brags will be included.

Examples:
  # Generate markdown document with all brags
  bragdoc doc generate

  # Generate and save to file
  bragdoc doc generate --output achievements.md

  # Generate with specific brags
  bragdoc doc generate --brags 1,2,3

  # Generate with brags in specific categories
  bragdoc doc generate --category project,leadership

  # Generate with brags having specific tags
  bragdoc doc generate --tags promotion,review

```bash
Usage:
  bragdoc doc generate [flags]

Flags:
  -b, --brags strings      Specific brag IDs to include (comma-separated)
  -c, --category strings   Include only these categories (comma-separated)
      --enhance-with-ai    Enhance descriptions using AI (not yet implemented)
  -f, --format string      Document format (markdown|pdf|docx) - MVP supports only markdown (default "markdown")
  -h, --help               help for generate
  -o, --output string      Output file path (if not specified, prints to stdout)
  -t, --tags strings       Include only brags with these tags (comma-separated)
      --template string    Document template (default|executive|technical) - MVP supports only default (default "default")
```
