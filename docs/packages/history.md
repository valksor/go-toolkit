# history

Generic JSON-persisted history tracking with automatic pruning.

## Description

The `history` package provides a generic `History[T]` type for tracking attempts, runs, or any historical data that should be persisted across sessions. It automatically prunes old entries to prevent unbounded growth.

**Key Features:**
- Generic type support (Go 1.18+)
- JSON persistence to disk
- Automatic pruning (keeps last N entries)
- Thread-safe operations
- Zero external dependencies

## Installation

```go
import "github.com/valksor/go-toolkit/history"
```

## Usage

### Basic Usage

```go
// Define your entry type
type BuildAttempt struct {
    Timestamp time.Time `json:"timestamp"`
    Success   bool      `json:"success"`
    Output    string    `json:"output"`
}

// Create a history manager
h := history.New[BuildAttempt]("/path/to/workspace", "build_history.json")

// Save an entry
attempt := BuildAttempt{
    Timestamp: time.Now(),
    Success:   true,
    Output:    "Build successful",
}
if err := h.Save(attempt); err != nil {
    log.Fatal(err)
}

// Load all entries
attempts, err := h.Load()
if err != nil {
    log.Fatal(err)
}

// Get the most recent entry
last, ok, err := h.Last()
if err != nil {
    log.Fatal(err)
}
if ok {
    fmt.Printf("Last build: %v\n", last.Timestamp)
}
```

### Custom Max Entries

```go
// Keep only last 5 entries (default is 10)
h := history.New[BuildAttempt](
    "/path/to/workspace",
    "build_history.json",
    history.WithMaxEntries(5),
)
```

### Clear History

```go
// Remove all history
if err := h.Clear(); err != nil {
    log.Fatal(err)
}
```

## API Reference

### Types

- `History[T]` - Generic history manager for type T (must be JSON-serializable)
- `Option` - Functional option for configuring History

### Constants

- `DefaultMaxEntries` - Default number of entries to keep (10)

### Functions

- `New[T](dir, filename string, opts ...Option) *History[T]` - Creates a new History manager
- `WithMaxEntries(n int) Option` - Sets maximum entries to keep

### Methods

- `(h *History[T]) Path() string` - Returns the full path to the history file
- `(h *History[T]) Load() ([]T, error)` - Loads all entries from disk
- `(h *History[T]) Save(entry T) error` - Appends an entry and persists to disk
- `(h *History[T]) SaveAll(entries []T) error` - Replaces all entries
- `(h *History[T]) Clear() error` - Removes all history
- `(h *History[T]) Len() (int, error)` - Returns the number of entries
- `(h *History[T]) Last() (T, bool, error)` - Returns the most recent entry

## Common Patterns

### Tracking Command Attempts

```go
type CommandAttempt struct {
    Timestamp time.Time `json:"timestamp"`
    Command   string    `json:"command"`
    ExitCode  int       `json:"exit_code"`
    IsDryRun  bool      `json:"dry_run"`
}

h := history.New[CommandAttempt]("~/.myapp", "command_history.json")

// Record an attempt
h.Save(CommandAttempt{
    Timestamp: time.Now(),
    Command:   "deploy",
    ExitCode:  0,
    IsDryRun:  false,
})
```

### Dry Run to Real Execution Pattern

```go
// First, dry run
h.Save(Attempt{IsDryRun: true, ...})

// Then, real execution
h.Save(Attempt{IsDryRun: false, ...})

// Check if last attempt was a dry run
last, ok, _ := h.Last()
if ok && last.IsDryRun {
    fmt.Println("Last attempt was a dry run")
}
```

## Important Notes

### File Location

The history file is stored at `{dir}/{filename}`. The directory is created automatically if it doesn't exist.

### Automatic Pruning

When `Save()` is called and the number of entries exceeds `maxEntries`, the oldest entries are automatically removed. Only the most recent `maxEntries` are kept.

### JSON Serialization

The entry type `T` must be JSON-serializable. All fields should have appropriate `json` tags for consistent serialization.

## See Also

- [cache](packages/cache.md) - In-memory caching with TTL
- [cfg](packages/cfg.md) - Configuration management
