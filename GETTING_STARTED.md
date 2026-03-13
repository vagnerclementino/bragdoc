# Getting Started with Bragdoc

Welcome to Bragdoc! This guide will help you get up and running quickly so you can start documenting your professional achievements.

## Table of Contents

- [What is Bragdoc?](#what-is-bragdoc)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Basic Usage](#basic-usage)
- [Advanced Features](#advanced-features)
- [Configuration](#configuration)
- [Tips and Best Practices](#tips-and-best-practices)
- [Troubleshooting](#troubleshooting)

## What is Bragdoc?

Bragdoc is a command-line tool that helps you track and document your professional achievements. It's designed to make it easy to:

- Record your accomplishments as they happen
- Organize achievements by category and tags
- Generate professional documents for performance reviews
- Export your achievements in multiple formats

Inspired by the concept of "Brag Documents" popularized by [Julia Evans](https://jvns.ca/blog/brag-documents/), Bragdoc helps you maintain a comprehensive record of your professional growth.

## Installation

### Prerequisites

- Go 1.21.1 or higher (for building from source)
- macOS or Linux operating system

### Building from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/vagnerclementino/bragdoc.git
   cd bragdoc
   ```

2. **Build the binary**:
   ```bash
   make build
   ```

3. **Install (optional)**:
   ```bash
   make install
   ```
   This installs the binary to `/usr/local/bin/bragdoc`

4. **Verify installation**:
   ```bash
   bragdoc version
   ```

## Quick Start

### 1. Initialize Bragdoc

Before using Bragdoc, you need to initialize it with your information:

```bash
bragdoc init \
  --name "Your Name" \
  --email "your.email@example.com" \
  --job-title "Software Engineer" \
  --company "Your Company" \
  --language en-US
```

This creates:
- Configuration file at `~/.bragdoc/config.yaml`
- SQLite database at `~/.bragdoc/bragdoc.db`

### 2. Add Your First Achievement

```bash
bragdoc brag add \
  --title "Implemented User Authentication" \
  --description "Designed and implemented a secure authentication system using OAuth 2.0, reducing login time by 50% and improving security" \
  --category achievement \
  --tags "security,backend,performance"
```

### 3. List Your Achievements

```bash
bragdoc brag list
```

Output:
```
ID  TITLE                              CATEGORY     TAGS                        CREATED
--  -----                              --------     ----                        -------
1   Implemented User Authentication    achievement  security, backend, performance  2024-01-15
```

Congratulations! You've created your first brag entry. 🎉

## Basic Usage

### Managing Brags

#### Add a Brag

```bash
bragdoc brag add \
  --title "Your Achievement Title" \
  --description "Detailed description of what you accomplished" \
  --category achievement \
  --tags "tag1,tag2,tag3"
```

**Categories:**
- `project` - Project-related achievements
- `achievement` - General accomplishments
- `skill` - New skills learned
- `leadership` - Leadership activities
- `innovation` - Innovative solutions

#### List All Brags

```bash
# Table format (default)
bragdoc brag list

# JSON format
bragdoc brag list --format json

# YAML format
bragdoc brag list --format yaml
```

#### Filter Brags

```bash
# By category
bragdoc brag list --category leadership

# By tags
bragdoc brag list --tags "backend,performance"
```

#### Show Specific Brags

```bash
# Single brag
bragdoc brag show --id 1

# Multiple brags
bragdoc brag show --id 1,2,3

# Range of brags
bragdoc brag show --id 1-5
```

#### Edit a Brag

```bash
bragdoc brag edit \
  --id 1 \
  --title "Updated Title" \
  --description "Updated description"
```

#### Remove a Brag

```bash
bragdoc brag remove --id 1 --force
```

### Managing Tags

#### Add a Tag

```bash
bragdoc tag add --name "kubernetes"
```

#### List All Tags

```bash
bragdoc tag list
```

#### Remove a Tag

```bash
bragdoc tag remove --id 1 --force
```

### Generating Documents

Generate a professional document from your brags:

```bash
bragdoc doc generate \
  --output my-achievements.md \
  --format markdown
```

## Advanced Features

### Output Formats

Bragdoc supports multiple output formats:

```bash
# Markdown (default)
bragdoc doc generate --format markdown

# JSON
bragdoc brag list --format json

# YAML
bragdoc brag list --format yaml
```

### Filtering and Searching

Combine filters for precise results:

```bash
# Leadership achievements with specific tags
bragdoc brag list --category leadership --tags "team,mentoring"

# All project-related brags
bragdoc brag list --category project
```

### Batch Operations

Show multiple brags at once:

```bash
# Specific IDs
bragdoc brag show --id 1,3,5,7

# Range
bragdoc brag show --id 1-10

# Combine with format
bragdoc brag show --id 1-5 --format json
```

## Configuration

### Configuration File

Located at `~/.bragdoc/config.yaml`:

```yaml
database:
  path: /Users/yourname/.bragdoc/bragdoc.db

user:
  id: 1
  name: Your Name
  email: your.email@example.com
  job_title: Software Engineer
  company: Your Company
  language: en-US
```

### Custom Database Location

You can specify a custom database path in the config:

```yaml
database:
  path: ~/Documents/my-brags/bragdoc.db
```

Bragdoc will:
- Expand `~` to your home directory
- Create parent directories if they don't exist
- Validate the path is accessible

### Supported Languages

Currently supported:
- `en-US` - English (United States)
- `pt-BR` - Portuguese (Brazil)

## Tips and Best Practices

### 1. Record Achievements Regularly

Don't wait for performance review time! Add achievements as they happen:

```bash
# Set up a weekly reminder to add brags
# Add this to your calendar or task manager
```

### 2. Be Specific and Quantifiable

Good:
```bash
bragdoc brag add \
  --title "Optimized Database Queries" \
  --description "Reduced query execution time by 75% (from 4s to 1s) by adding proper indexes and optimizing N+1 queries, improving user experience for 10,000+ daily users"
```

Not as good:
```bash
bragdoc brag add \
  --title "Made things faster" \
  --description "Improved performance"
```

### 3. Use Consistent Tags

Create a tagging system and stick to it:

```bash
# Technical areas
--tags "backend,frontend,devops,database"

# Skills
--tags "leadership,mentoring,communication"

# Impact
--tags "performance,security,cost-savings"
```

### 4. Categorize Appropriately

- **project**: Major projects or initiatives
- **achievement**: Specific accomplishments
- **skill**: New technologies or skills learned
- **leadership**: Team leadership, mentoring
- **innovation**: Creative solutions or new approaches

### 5. Include Context

Always include:
- What you did
- Why it mattered
- The impact (quantified if possible)
- Technologies or skills used

### 6. Regular Reviews

Review your brags quarterly:

```bash
# Export all brags
bragdoc doc generate --output quarterly-review.md

# Review and update as needed
```

## Troubleshooting

### Command Not Found

If you get `command not found: bragdoc`:

```bash
# Check if binary exists
ls -la ./bragdoc

# Add to PATH or use full path
export PATH=$PATH:/path/to/bragdoc

# Or install globally
make install
```

### Database Errors

If you encounter database errors:

```bash
# Check database location
cat ~/.bragdoc/config.yaml

# Verify database exists
ls -la ~/.bragdoc/bragdoc.db

# Reinitialize if needed (WARNING: This will delete existing data)
rm -rf ~/.bragdoc
bragdoc init --name "Your Name" --email "your@email.com"
```

### Validation Errors

Common validation errors and fixes:

**Title too short:**
```bash
# Error: brag title must be at least 5 characters
# Fix: Use a more descriptive title
--title "Implemented Feature X"  # Good
--title "Fix"                     # Too short
```

**Description too short:**
```bash
# Error: brag description must be at least 20 characters
# Fix: Provide more detail
--description "Implemented a comprehensive authentication system..."  # Good
--description "Fixed bug"                                             # Too short
```

**Invalid category:**
```bash
# Error: invalid category
# Fix: Use one of the valid categories
--category achievement  # Valid
--category random      # Invalid
```

### Permission Errors

If you get permission errors:

```bash
# Check directory permissions
ls -la ~/.bragdoc

# Fix permissions if needed
chmod 755 ~/.bragdoc
chmod 644 ~/.bragdoc/config.yaml
chmod 644 ~/.bragdoc/bragdoc.db
```

## Next Steps

Now that you're familiar with the basics:

1. **Start documenting**: Add your recent achievements
2. **Explore features**: Try different output formats and filters
3. **Establish a routine**: Set reminders to add brags regularly
4. **Generate reports**: Create documents for performance reviews
5. **Contribute**: Check out [CONTRIBUTING.md](CONTRIBUTING.md) to help improve Bragdoc

## Getting Help

- **Documentation**: See the [README](README.md) for more information
- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/vagnerclementino/bragdoc/issues)
- **Examples**: Check the [smoke test](smoke.sh) for usage examples

Happy bragging! 🚀
