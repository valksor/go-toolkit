# cli

Cobra CLI helpers and common patterns.

## Description

The `cli` package provides helpers and common patterns for building CLI applications using the Cobra framework.

**Key Features:**
- Common command patterns
- Reusable CLI utilities
- Cobra integration helpers

## Installation

```go
import "github.com/valksor/go-toolkit/cli"
```

## Usage

The `cli` package provides utilities for building CLI applications with Cobra. Please refer to the source code for detailed usage examples.

## API Reference

For full API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/valksor/go-toolkit/cli).

## Common Patterns

### Cobra Command Structure

```go
import (
    "github.com/spf13/cobra"
    "github.com/valksor/go-toolkit/cli"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My application",
    RunE:  runCommand,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### Using CLI Helpers

```go
// The cli package provides various helpers for common CLI tasks
// See package source code for specific helpers available
```

## See Also

- [cli/disambiguate](packages/disambiguate.md) - Symfony-style command shortcuts and prefix matching
- [Cobra Documentation](https://github.com/spf13/cobra)
- [display](packages/display.md) - Terminal color and formatting utilities
- [log](packages/log.md) - Structured logging helpers
