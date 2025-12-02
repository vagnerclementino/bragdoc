# Bragdoc Document Generation Guide

## Quick Start

Generate a professional achievement document from your brags:

```bash
# Generate and display in terminal
bragdoc doc generate

# Save to a file
bragdoc doc generate --output my-achievements.md
```

## Command Reference

### Basic Usage

```bash
bragdoc doc generate [flags]
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | stdout |
| `--format` | `-f` | Document format (markdown\|pdf\|docx) | markdown |
| `--brags` | `-b` | Specific brag IDs (comma-separated) | all |
| `--category` | `-c` | Filter by categories (comma-separated) | all |
| `--tags` | `-t` | Filter by tags (comma-separated) | all |
| `--template` | | Template to use (default\|executive\|technical) | default |
| `--enhance-with-ai` | | Enhance with AI (not yet implemented) | false |

### MVP Limitations

- **Format**: Only Markdown is currently supported (PDF and Docx coming soon)
- **Template**: Only "default" template is available
- **AI Enhancement**: Not yet implemented (coming in next phase)

## Examples

### Generate All Brags

```bash
# Display in terminal
bragdoc doc generate

# Save to file
bragdoc doc generate --output achievements.md
```

### Filter by Specific Brags

```bash
# Single brag
bragdoc doc generate --brags 1 --output brag-1.md

# Multiple brags
bragdoc doc generate --brags 1,2,3 --output selected-brags.md
```

### Filter by Category

```bash
# Single category
bragdoc doc generate --category leadership --output leadership.md

# Multiple categories
bragdoc doc generate --category project,achievement --output projects-achievements.md
```

Available categories:
- `project`
- `achievement`
- `skill`
- `leadership`
- `innovation`

### Filter by Tags

```bash
# Single tag
bragdoc doc generate --tags promotion --output promotion-doc.md

# Multiple tags
bragdoc doc generate --tags review,important --output review-doc.md
```

### Combine Filters

```bash
# Category + Tags
bragdoc doc generate \
  --category achievement \
  --tags important,promotion \
  --output key-achievements.md

# Specific brags + Output
bragdoc doc generate \
  --brags 1,5,7 \
  --output selected.md
```

## Document Structure

The generated Markdown document includes:

1. **Header**
   - Your name
   - Job title (if configured)
   - Company (if configured)
   - Generation date

2. **Summary**
   - Total number of brags
   - Categories covered

3. **Achievements by Category**
   - Organized by category (achievement, leadership, project, etc.)
   - Each brag includes:
     - Title
     - Description
     - Tags
     - Creation date

4. **Footer**
   - Information about Bragdoc CLI

## Tips

### For Performance Reviews

```bash
# Generate achievements from the last year with specific tags
bragdoc doc generate \
  --tags review,2024 \
  --output performance-review-2024.md
```

### For Promotion Packets

```bash
# Focus on leadership and major projects
bragdoc doc generate \
  --category leadership,project \
  --tags promotion,impact \
  --output promotion-packet.md
```

### For Resume Updates

```bash
# Get all achievements to update resume
bragdoc doc generate \
  --category achievement,skill \
  --output resume-achievements.md
```

## Troubleshooting

### No brags found

If you see "no brags found matching the criteria":
- Check your filters are correct
- Verify you have brags in the database: `bragdoc brag list`
- Try without filters first: `bragdoc doc generate`

### Format not supported

If you see "format not yet supported":
- Use `--format markdown` (or omit the flag)
- PDF and Docx support coming in future releases

### AI enhancement not available

If you see "not yet implemented":
- Remove the `--enhance-with-ai` flag
- AI enhancement will be available in the next release

## Next Steps

After generating your document:

1. **Review the content** - Make sure all important achievements are included
2. **Edit if needed** - The Markdown file can be edited in any text editor
3. **Convert format** - Use tools like Pandoc to convert to PDF or Word if needed
4. **Share** - Use the document for reviews, promotions, or job applications

## Future Features

Coming soon:
- PDF and Docx export
- AI-powered content enhancement
- Multiple templates (executive, technical)
- Custom templates
- Language translation
- Collaborative features

## Help

For more information:
```bash
bragdoc doc --help
bragdoc doc generate --help
```
