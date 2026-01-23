# Brag

## Available Commands

```bash
  add         Add a new brag entry
  edit        Edit an existing brag entry
  list        List brag entries
  remove      Remove brag entries
  show        Show detailed information about brag entries
```

### add

Add a new brag entry to document your professional achievements

```bash
Usage:
  bragdoc brag add [flags]

Flags:
  -c, --category string      Brag category (project|achievement|skill|leadership|innovation) (default "achievement")
  -d, --description string   Brag description (required)
  -h, --help                 help for add
      --tags strings         Comma-separated list of tags
  -t, --title string         Brag title (required)
```

### edit

Edit an existing brag entry by ID

```bash
Usage:
  bragdoc brag edit <id> [flags]

Flags:
  -c, --category string      New brag category (project|achievement|skill|leadership|innovation)
  -d, --description string   New brag description
  -h, --help                 help for edit
      --tags strings         New comma-separated list of tags (replaces existing tags)
  -t, --title string         New brag title
```

### list

List all your documented professional achievements with optional filters

```bash
Usage:
  bragdoc brag list [flags]

Flags:
  -c, --category string   Filter by category (project|achievement|skill|leadership|innovation)
  -f, --format string     Output format (table|json|yaml) (default "table")
  -h, --help              help for list
  -l, --limit int         Maximum number of results (default 50)
  -t, --tags strings      Filter by tags (comma-separated)
```

### remove

Remove one or more brag entries by ID.
Supports multiple IDs and ranges:
  - Single ID: bragdoc brag remove 1
  - Multiple IDs: bragdoc brag remove 1,2,3
  - Range: bragdoc brag remove 1-5
  - Combined: bragdoc brag remove 1,3,5-8

```bash
Usage:
  bragdoc brag remove <ids> [flags]

Flags:
  -f, --force   Skip confirmation prompt
  -h, --help    help for remove
```

### show

Show detailed information about one or more brag entries by ID.
Supports multiple IDs and ranges:
  - Single ID: bragdoc brag show 1
  - Multiple IDs: bragdoc brag show 1,2,3
  - Range: bragdoc brag show 1-5

```bash
Usage:
  bragdoc brag show <ids> [flags]

Flags:
  -h, --help   help for show
```
