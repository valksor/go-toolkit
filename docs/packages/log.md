# log

Structured logging helpers using Go's `log/slog`.

## Description

The `log` package provides structured logging using Go's standard `log/slog` package. It supports multiple output formats (text/JSON), configurable levels, and contextual logging.

**Key Features:**
- Structured logging with `log/slog`
- Text and JSON output formats
- Configurable log levels
- Context-aware logging
- Thread-safe

## Installation

```go
import "github.com/valksor/go-toolkit/log"
```

## Usage

### Basic Logging

```go
import (
    "log/slog"
    "github.com/valksor/go-toolkit/log"
)

// Default configuration (text format, info level)
log.Info("application started", "version", "1.0.0")
log.Debug("processing request", "id", "abc-123")
log.Warn("high memory usage", "usage", "90%")
log.Error("request failed", "error", err)
```

### Configuration

```go
// Configure with custom options
log.Configure(log.Options{
    Output: os.Stderr,      // Output writer
    Level:  log.LevelDebug, // Log level
    JSON:    true,          // Use JSON format
})
```

### Context-Aware Logging

```go
ctx := context.Background()

// Log with context
log.InfoContext(ctx, "processing request",
    "request_id", "abc-123",
    "user_id", "user-456",
)
```

### Different Log Levels

```go
// Debug level - detailed diagnostics
log.Debug("database query", "query", "SELECT * FROM users")

// Info level - general informational messages
log.Info("user logged in", "user", "alice")

// Warn level - warning messages
log.Warn("deprecated API used", "endpoint", "/v1/users")

// Error level - error messages
log.Error("operation failed", "error", err)
```

## API Reference

### Types

- `Level = slog.Level` - Logging level type
- `Options` - Logger configuration options

### Constants

- `LevelDebug = slog.LevelDebug` - Debug level
- `LevelInfo = slog.LevelInfo` - Info level
- `LevelWarn = slog.LevelWarn` - Warning level
- `LevelError = slog.LevelError` - Error level

### Functions

- `Configure(Options)` - Configure the global logger
- `Debug(msg string, args ...any)` - Log at debug level
- `Info(msg string, args ...any)` - Log at info level
- `Warn(msg string, args ...any)` - Log at warning level
- `Error(msg string, args ...any)` - Log at error level
- `DebugContext(ctx context.Context, msg string, args ...any)` - Log at debug level with context
- `InfoContext(ctx context.Context, msg string, args ...any)` - Log at info level with context
- `WarnContext(ctx context.Context, msg string, args ...any)` - Log at warning level with context
- `ErrorContext(ctx context.Context, msg string, args ...any)` - Log at error level with context
- `NewHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler` - Create a new handler

## Common Patterns

### Request Logging

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    requestID := r.Header.Get("X-Request-ID")

    log.InfoContext(r.Context(), "handling request",
        "method", r.Method,
        "path", r.URL.Path,
        "request_id", requestID,
    )

    // Handle request...
}
```

### Error Logging

```go
func processData(data []byte) error {
    if len(data) == 0 {
        log.Error("empty data received", "source", "api")
        return fmt.Errorf("empty data")
    }

    log.Debug("processing data", "size", len(data))
    // Process data...
    return nil
}
```

### Environment-Based Logging

```go
func init() {
    level := log.LevelInfo
    if os.Getenv("DEBUG") != "" {
        level = log.LevelDebug
    }

    log.Configure(log.Options{
        Level: level,
        JSON:  os.Getenv("ENVIRONMENT") == "production",
    })
}
```

### JSON Logging for Production

```go
if os.Getenv("ENVIRONMENT") == "production" {
    log.Configure(log.Options{
        Level: log.LevelInfo,
        JSON:  true,
    })
}
```

## Output Formats

### Text Format (Default)

```
time=2024-01-15T12:00:00.000Z level=INFO msg="application started" version=1.0.0
```

### JSON Format

```json
{"time":"2024-01-15T12:00:00.000Z","level":"INFO","msg":"application started","version":"1.0.0"}
```

## See Also

- [Go `log/slog` documentation](https://pkg.go.dev/log/slog)
- [display](packages/display.md) - Terminal color and formatting
