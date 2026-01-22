# cli/disambiguate

Symfony-style command shortcuts and prefix matching for Cobra CLI applications.

## Description

The `disambiguate` package enables shorthand command syntax like `c:v` → `config validate`, with case-insensitive prefix matching and optional interactive disambiguation.

**Key Features:**
- Colon notation for command shortcuts (`c:v` → `config validate`)
- Case-insensitive prefix matching
- Interactive disambiguation for ambiguous commands
- Non-interactive mode support for scripts

## Installation

```go
import "github.com/valksor/go-toolkit/cli/disambiguate"
```

## Usage

### Basic Integration

```go
import (
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/valksor/go-toolkit/cli/disambiguate"
)

func Execute() error {
    // Pre-process args to handle colon notation
    args := os.Args[1:]
    if len(args) > 0 && strings.Contains(args[0], ":") {
        resolved, matches, err := disambiguate.ResolveColonPath(rootCmd, args[0])
        if err == nil {
            if len(matches) == 0 {
                // Unambiguous match
                rootCmd.SetArgs(append(resolved, args[1:]...))
                return rootCmd.Execute()
            }
            // Ambiguous - handle selection
            if !disambiguate.IsInteractive() {
                return errors.New(disambiguate.FormatAmbiguousError(args[0], matches))
            }
            selected, err := disambiguate.SelectCommand(matches, args[0])
            if err != nil {
                return err
            }
            rootCmd.SetArgs(append(selected.Path, args[1:]...))
            return rootCmd.Execute()
        }
    }
    return rootCmd.Execute()
}
```

## Examples

```bash
# Full commands
myapp config validate
myapp agents list
myapp memory show

# Colon notation shortcuts
myapp c:v              # → config validate
myapp c:i              # → config init
myapp a:l              # → agents list
myapp m:s              # → memory show
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `ResolveColonPath(root, path)` | Resolves colon-separated path like `c:v` to `["config", "validate"]` |
| `FindPrefixMatches(parent, prefix)` | Finds all commands matching a case-insensitive prefix |
| `FindPrefixMatchesInPath(root, prefix)` | Searches entire command tree for matches |
| `FormatAmbiguousError(prefix, matches)` | Formats error message for ambiguous commands |
| `IsInteractive()` | Returns true if stdin is a TTY |
| `SelectCommand(matches, prefix)` | Prompts user to select from ambiguous matches |
| `FormatMatchNames(matches)` | Returns comma-separated command names |

### Types

| Type | Description |
|------|-------------|
| `CommandMatch` | Represents a command match with Command, MatchedName, and Path fields |

## Behavior

- **Case-insensitive**: `C:V` works the same as `c:v`
- **Prefix matching**: `conf:val` matches `config:validate`
- **Ambiguity handling**: Prompts for interactive selection when multiple commands match
- **Non-interactive**: In scripts, shows formatted error instead of prompting

## See Also

- [cli](packages/cli.md) - Cobra CLI helpers
- [Cobra Documentation](https://github.com/spf13/cobra)
