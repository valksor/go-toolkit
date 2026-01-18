# qtcwrap

Wrapper for QuickTemplate (qtc) compiler.

## Description

The `qtcwrap` package provides utilities for working with the QuickTemplate compiler (qtc). It simplifies the process of generating Go code from QuickTemplate (.qtpl) files by providing a convenient Go API wrapper.

**Key Features:**
- Directory-based template compilation
- Single file template compilation
- Skip line comments for cleaner generated code
- Custom file extensions
- Proper error handling

## Installation

```go
import "github.com/valksor/go-toolkit/qtcwrap"
```

## Usage

### Directory Compilation

```go
import "github.com/valksor/go-toolkit/qtcwrap"

// Compile all .qtpl files in a directory
config := qtcwrap.Config{
    Dir:              "templates",
    SkipLineComments: true,
    Ext:              ".qtpl",
}
err := qtcwrap.WithConfig(config)
```

### Single File Compilation

```go
// Compile a single .qtpl file
config := qtcwrap.Config{
    File:             "templates/index.qtpl",
    SkipLineComments: false,
}
err := qtcwrap.WithConfig(config)
```

### Default Options

```go
// Use default configuration (current dir, .qtpl extension)
err := qtcwrap.Default()
```

## API Reference

### Types

- `Config` - Configuration options for the qtc compiler

### Fields

- `Dir` - Directory containing .qtpl files (default: ".")
- `File` - Specific file to compile (overrides Dir)
- `SkipLineComments` - Skip line comments in generated code (default: false)
- `Ext` - File extension for templates (default: ".qtpl")

### Functions

- `WithConfig(Config) error` - Run qtc with specified configuration
- `Default() error` - Run qtc with default configuration
- `WithDir(dir string) error` - Run qtc on directory with defaults

## Common Patterns

### Build Integration

```go
//go:generate go run github.com/valksor/go-toolkit/cmd/qtcwrap
package templates

// Or use in Makefile
// compile-templates:
//     go run github.com/valksor/go-toolkit/cmd/qtcwrap -dir=templates
```

### Pre-commit Hook

```go
// Automatically compile templates before commit
func preCommitHook() error {
    config := qtcwrap.Config{
        Dir:              "templates",
        SkipLineComments: true,
    }

    if err := qtcwrap.WithConfig(config); err != nil {
        return fmt.Errorf("template compilation failed: %w", err)
    }

    return nil
}
```

### Watch Mode

```go
func watchTemplates() {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add("templates")

    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".qtpl") {
                config := qtcwrap.Config{File: event.Name}
                if err := qtcwrap.WithConfig(config); err != nil {
                    log.Printf("Compilation error: %v", err)
                }
            }
        }
    }
}
```

### Production Build

```go
// Skip line comments for production (smaller binaries)
config := qtcwrap.Config{
    Dir:              "templates",
    SkipLineComments: true,
}

// Or keep for development (better stack traces)
config := qtcwrap.Config{
    Dir:              "templates",
    SkipLineComments: false,
}
```

## Generated File Location

The qtc compiler generates Go files in the same directory as the .qtpl files:

```
templates/
  index.qtpl
  index.qtpl.go    # Generated
  layout.qtpl
  layout.qtpl.go   # Generated
```

## SkipLineComments Option

### When `SkipLineComments: true`:
- Generated code is more compact
- Smaller file sizes
- Harder to debug (no line correlation)

### When `SkipLineComments: false`:
- Generated code preserves line numbers
- Better stack traces
- Easier debugging
- Larger file sizes

## Error Handling

The package wraps qtc errors with additional context:

```go
err := qtcwrap.WithConfig(config)
if err != nil {
    // Error includes command output and context
    fmt.Fprintf(os.Stderr, "Template compilation failed: %v\n", err)
    os.Exit(1)
}
```

## Dependencies

- [QuickTemplate](https://github.com/valksor/qtc) - Template compiler

## See Also

- [minify](packages/minify.md) - Asset minification
- QuickTemplate [documentation](https://github.com/valksor/qtc)
