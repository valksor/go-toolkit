# slug

Convert text to URL-safe slugs.

## Description

The `slug` package provides utilities for converting text to URL-safe slugs suitable for branch names, URLs, and identifiers.

**Key Features:**
- Unicode normalization and diacritic removal
- Lowercasing
- Replaces spaces/underscores with hyphens
- Removes non-alphanumeric characters
- Smart truncation at word boundaries
- Collapses multiple consecutive hyphens

## Installation

```go
import "github.com/valksor/go-toolkit/slug"
```

## Usage

### Basic Slug Generation

```go
import "github.com/valksor/go-toolkit/slug"

s := slug.Slugify("Add user authentication", 0)
// Returns: "add-user-authentication"
```

### With Length Limit

```go
s := slug.Slugify("Add user authentication for the admin dashboard panel", 30)
// Returns: "add-user-authentication-for"
```

### Unicode Support

```go
s := slug.Slugify("Créer un compte utilisateur", 50)
// Returns: "creer-un-compte-utilisateur"

s := slug.Slugify("用户认证", 20)
// Returns: "user-ren-zheng"
```

### Existing Slugs (Pass-through)

```go
s := slug.Slugify("already-a-slug", 0)
// Returns: "already-a-slug"
```

### Title Case to Slug

```go
s := slug.Slugify("Fix Bug #123: Memory Leak in Parser", 50)
// Returns: "fix-bug-123-memory-leak-in-parser"
```

## API Reference

### Functions

#### Slugify
```go
func Slugify(title string, maxLen int) string
```

Converts a title to a URL-safe slug suitable for branch names.

**Parameters:**
- `title` - The input string to convert
- `maxLen` - Maximum length (0 or negative = no truncation)

**Returns:** URL-safe slug string

**Behavior:**
1. Normalizes unicode and removes diacritics
2. Lowercases the string
3. Replaces spaces and underscores with hyphens
4. Removes non-alphanumeric characters (except hyphens)
5. Collapses multiple consecutive hyphens
6. Trims leading/trailing hyphens
7. Truncates to `maxLen` at word boundary if needed

## Common Patterns

### Git Branch Names

```go
func CreateBranchName(title string) string {
    // Git branch names typically limited to ~50 chars
    return slug.Slugify(title, 50)
}

branchName := CreateBranchName("FEATURE: Add OAuth2 login support")
// Returns: "feature-add-oauth2-login-support"
```

### URL Path Generation

```go
func GenerateURLPath(articleTitle string) string {
    // No length limit for URLs
    return slug.Slugify(articleTitle, 0)
}

path := GenerateURLPath("How to Build REST APIs with Go")
// Returns: "how-to-build-rest-apis-with-go"
```

### Issue/Task Identifiers

```go
func TaskToBranch(taskType, taskID, title string) string {
    prefix := slug.Slugify(taskType, 20)
    idPart := slug.Slugify(taskID, 20)
    titlePart := slug.Slugify(title, 40)
    return fmt.Sprintf("%s-%s-%s", prefix, idPart, titlePart)
}

branch := TaskToBranch("feature", "PROJ-123", "Add dark mode support")
// Returns: "feature-proj-123-add-dark-mode-support"
```

### Consistent Identifiers

```go
// For creating consistent file names from titles
func FileNameFromTitle(title string) string {
    return slug.Slugify(title, 100) + ".md"
}

file := FileNameFromTitle("My Document: Introduction (2024)")
// Returns: "my-document-introduction-2024.md"
```

## Truncation Behavior

When `maxLen` is specified, truncation prefers word boundaries:

```go
// Long title gets truncated at hyphen (word boundary)
slug.Slugify("this-is-a-very-long-title-that-needs-truncation", 20)
// Returns: "this-is-a-very-long"

// If no hyphen exists, truncates directly and removes trailing hyphen
slug.Slugify("superlongwordwithoutanyhyphens", 10)
// Returns: "superlongwo"
```

## See Also

- [paths](paths.md) - File path manipulation utilities
- [env](env.md) - Environment variable expansion
