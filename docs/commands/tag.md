# TAG

## Available Commands

```bash
  add         Add a new tag
  list        List all tags
  remove      Remove a tag
```

### add

Create a new tag for organizing your brags

```bash
Usage:
  bragdoc tag add [flags]

Flags:
  -h, --help          help for add
  -n, --name string   Tag name (required)
```

### list

List all tags created by the user

```bash
Usage:
  bragdoc tag list [flags]

Flags:
  -f, --format string   Output format (table|json|yaml) (default "table")
  -h, --help            help for list
```

### remove

Remove a tag by ID. This will also remove all associations with brags.

```bash
Usage:
  bragdoc tag remove <id> [flags]

Flags:
  -f, --force   Skip confirmation prompt
  -h, --help    help for remove
```
