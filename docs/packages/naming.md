# naming

Filename parsing, template expansion, and ticket ID extraction.

## Description

The `naming` package provides utilities for extracting information from filenames and working with naming conventions common in project management tools.

**Key Features:**
- Task type extraction from filenames (feature, fix, docs, etc.)
- Ticket ID parsing (JIRA-123, FEATURE-456, etc.)
- Template expansion for branch names
- Git branch name sanitization
- Type alias normalization (bug → fix, feat → feature)

## Installation

```go
import "github.com/valksor/go-toolkit/naming"
```

## Usage

### Extract Task Type from Filename

```go
// From ticket pattern
taskType := naming.TaskTypeFromFilename("FEATURE-123.md")
// Returns: "feature"

// From type prefix
taskType = naming.TaskTypeFromFilename("fix-login-bug.md")
// Returns: "fix"

// Aliases are normalized
taskType = naming.TaskTypeFromFilename("BUG-456.md")
// Returns: "fix" (bug is aliased to fix)

// Unknown patterns default to "task"
taskType = naming.TaskTypeFromFilename("my-custom-file.md")
// Returns: "task"
```

### Extract Key from Filename

```go
key := naming.KeyFromFilename("FEATURE-123.md")
// Returns: "FEATURE-123"

key = naming.KeyFromFilename("/path/to/fix-login-bug.md")
// Returns: "fix-login-bug"
```

### Parse Ticket ID

```go
ticketID, ok := naming.ParseTicketID("JIRA-456-extra-stuff")
// Returns: "JIRA-456", true

ticketID, ok = naming.ParseTicketID("my-task")
// Returns: "", false
```

### Template Expansion

```go
vars := naming.TemplateVars{
    Key:    "FEATURE-123",
    TaskID: "a1b2c3d4",
    Type:   "feature",
    Slug:   "add-user-auth",
    Title:  "Add user authentication",
}

// Expand branch name pattern
branch := naming.ExpandTemplate("{type}/{key}--{slug}", vars)
// Returns: "feature/FEATURE-123--add-user-auth"

// Expand commit prefix
prefix := naming.ExpandTemplate("[{key}]", vars)
// Returns: "[FEATURE-123]"
```

### Validate Pattern

```go
unknown := naming.ValidatePattern("{type}/{unknown_var}")
// Returns: []string{"unknown_var"}

unknown = naming.ValidatePattern("{type}/{key}")
// Returns: nil (all variables are valid)
```

### Clean Branch Name

```go
branch := naming.CleanBranchName("feature//name---test/")
// Returns: "feature/name--test"
```

## API Reference

### Types

- `TemplateVars` - Variables for template expansion

### Variables

- `TypePrefixes` - Known type prefixes for task classification

### Functions

- `TaskTypeFromFilename(filename string) string` - Extracts task type from filename
- `KeyFromFilename(filename string) string` - Extracts key from filename
- `KeyFromDirectory(dirPath string) string` - Extracts key from directory path
- `ParseTicketID(s string) (string, bool)` - Extracts ticket ID from string
- `ExpandTemplate(pattern string, vars TemplateVars) string` - Expands a pattern with variables
- `ValidatePattern(pattern string) []string` - Returns unknown variables in pattern
- `CleanBranchName(name string) string` - Sanitizes a git branch name

### TemplateVars Fields

| Field | Description | Example |
|-------|-------------|---------|
| `Key` | External key | "FEATURE-123" |
| `TaskID` | Internal task ID | "a1b2c3d4" |
| `Type` | Task type | "feature" |
| `Slug` | Slugified title | "add-user-auth" |
| `Title` | Original title | "Add user authentication" |

### Supported Template Variables

- `{key}` - External key
- `{task_id}` - Internal task ID
- `{type}` - Task type
- `{slug}` - Slugified title
- `{title}` - Original title

## Common Patterns

### Recognized Type Prefixes

The following prefixes are recognized (case-insensitive):

| Prefix | Canonical Type |
|--------|---------------|
| feature, feat | feature |
| fix, bug, bugfix, hotfix | fix |
| docs, doc | docs |
| refactor, refact | refactor |
| perf, performance | perf |
| test, tests | test |
| chore | chore |
| style | style |
| ci | ci |
| build | build |
| task | task |

### Ticket ID Pattern

Ticket IDs must match the pattern `^[A-Z]+-\d+`:
- Starts with uppercase letters
- Followed by a hyphen
- Ends with digits

Examples: `JIRA-123`, `FEATURE-1`, `ABC-9999`

## See Also

- [slug](packages/slug.md) - URL-safe slug generation
- [paths](packages/paths.md) - File path utilities
