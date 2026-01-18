# output

Output processing utilities including deduplicating writer.

## Description

The `output` package provides the `DeduplicatingWriter`, which suppresses consecutive duplicate lines.

**Key Features:**
- Suppress consecutive duplicate lines
- Thread-safe operation
- Buffered partial line handling
- Auto-flush for multi-line writes

## Installation

```go
import "github.com/valksor/go-toolkit/output"
```

## Usage

### Basic Usage

```go
import "github.com/valksor/go-toolkit/output"

// Wrap any io.Writer
w := output.NewDeduplicatingWriter(os.Stdout)

// Write lines - consecutive duplicates are suppressed
w.Write([]byte("Processing...\n"))
w.Write([]byte("Processing...\n")) // This line will be suppressed
w.Write([]byte("Done!\n"))          // This line will be written

// Remember to flush when done
w.Flush()
```

### With Log Output

```go
// Suppress duplicate log messages
logWriter := output.NewDeduplicatingWriter(logFile)
logger := log.New(logWriter, "", 0)

logger.Println("Checking service...")
logger.Println("Checking service...") // Suppressed
logger.Println("Service is up")        // Written
logWriter.Flush()
```

### Reset State

```go
w := output.NewDeduplicatingWriter(os.Stdout)

w.Write([]byte("Message\n"))
w.Write([]byte("Message\n")) // Suppressed

w.Reset() // Clear state

w.Write([]byte("Message\n")) // Written (state was reset)
```

## API Reference

### Types

- `DeduplicatingWriter` - Wrapper that suppresses consecutive duplicate lines

### Functions

- `NewDeduplicatingWriter(w io.Writer) *DeduplicatingWriter` - Create new deduplicating writer

### Methods

- `(d *DeduplicatingWriter) Write(p []byte) (int, error)` - Write with deduplication
- `(d *DeduplicatingWriter) Flush() error` - Flush any buffered content
- `(d *DeduplicatingWriter) Reset()` - Clear deduplication state

## Common Patterns

### Polling Status Messages

```go
func pollService() {
    w := output.NewDeduplicatingWriter(os.Stdout)
    defer w.Flush()

    for {
        status := checkService()
        fmt.Fprintf(w, "Service status: %s\n", status)
        time.Sleep(5 * time.Second)
    }
}
// Only status changes are printed
```

### Loop Progress Reporting

```go
func processItems(items []Item) {
    w := output.NewDeduplicatingWriter(os.Stdout)
    defer w.Flush()

    for i, item := range items {
        fmt.Fprintf(w, "Processing %d/%d\n", i+1, len(items))
        processItem(item)
    }
    fmt.Fprintf(w, "Done!\n")
}
```

### HTTP Request Logging

```go
type loggingResponseWriter struct {
    http.ResponseWriter
    output *output.DeduplicatingWriter
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
    // Log request body (suppressing duplicates)
    l.output.Write(b)
    return l.ResponseWriter.Write(b)
}
```

## Important Notes

### Thread Safety

`DeduplicatingWriter` is thread-safe and can be used concurrently from multiple goroutines.

### Flush Required

Always call `Flush()` when done writing to ensure no buffered content is lost. This is especially important for partial lines (lines without newlines).

### Partial Line Handling

The writer buffers partial lines until a newline is received. Multi-line writes are auto-flushed after processing complete lines.

### Comparison Logic

Lines are compared without trailing newlines for deduplication purposes, but newlines are preserved when writing.

## Behavior Examples

```go
w := output.NewDeduplicatingWriter(os.Stdout)

w.Write([]byte("Line 1\n"))  // Written
w.Write([]byte("Line 1\n"))  // Suppressed (duplicate)
w.Write([]byte("Line 2\n"))  // Written (different)
w.Write([]byte("Line 1\n"))  // Written (new sequence)

w.Flush()
```

Output:
```
Line 1
Line 2
Line 1
```

## See Also

- Go standard library `io` and `bufio` packages
