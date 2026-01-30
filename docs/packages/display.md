# display

Color formatting, spinners, and terminal output utilities.

## Description

The `display` package provides utilities for adding color, formatting, and animated spinners to terminal output.

**Key Features:**
- Semantic color functions (Success, Error, Warning, Info)
- Animated CLI spinners for long-running operations
- NO_COLOR environment variable support
- Cross-platform compatibility
- Thread-safe spinner operations

## Installation

```go
import "github.com/valksor/go-toolkit/display"
```

## Usage

### Color Functions

```go
// Initialize colors (usually at startup)
display.InitColors(noColorFlag)

// Semantic colors
fmt.Println(display.Success("Operation completed"))
fmt.Println(display.Error("Something went wrong"))
fmt.Println(display.Warning("Be careful"))
fmt.Println(display.Info("FYI"))
fmt.Println(display.Muted("Secondary info"))
fmt.Println(display.Bold("Important"))
fmt.Println(display.Cyan("code"))

// Prefixed messages
fmt.Println(display.SuccessMsg("File saved"))      // ✓ File saved
fmt.Println(display.ErrorMsg("File not found"))    // ✗ File not found
fmt.Println(display.WarningMsg("Disk almost full")) // ⚠ Disk almost full
fmt.Println(display.InfoMsg("Processing..."))      // → Processing...
```

### Spinners

```go
// Create and start a spinner
spinner := display.NewSpinner("Loading data...")
spinner.Start()

// Do work...
time.Sleep(2 * time.Second)

// Stop with success message
spinner.StopWithSuccess("Data loaded!")

// Or stop with error/warning
spinner.StopWithError("Failed to load data")
spinner.StopWithWarning("Loaded with warnings")

// Update message while running
spinner.UpdateMessage("Processing item 5 of 10...")
```

### Check Color Status

```go
if display.ColorsEnabled() {
    // Colors are enabled
}

// Manually control colors (useful for testing)
display.SetColorsEnabled(false)
```

## API Reference

### Color Functions

- `InitColors(noColor bool)` - Initialize color system
- `ColorsEnabled() bool` - Check if colors are enabled
- `SetColorsEnabled(enabled bool)` - Manually control colors

### Semantic Colors

- `Success(text string) string` - Green text
- `Error(text string) string` - Red text
- `Warning(text string) string` - Yellow text
- `Info(text string) string` - Blue text
- `Muted(text string) string` - Gray text
- `Bold(text string) string` - Bold text
- `Cyan(text string) string` - Cyan text

### Prefix Functions

- `SuccessPrefix() string` - Returns green checkmark (✓)
- `ErrorPrefix() string` - Returns red X (✗)
- `WarningPrefix() string` - Returns yellow warning (⚠)
- `InfoPrefix() string` - Returns blue arrow (→)

### Message Functions

- `SuccessMsg(format string, args ...any) string` - Checkmark + message
- `ErrorMsg(format string, args ...any) string` - X + red message
- `WarningMsg(format string, args ...any) string` - Warning + yellow message
- `InfoMsg(format string, args ...any) string` - Arrow + message

### Spinner Type

- `NewSpinner(message string) *Spinner` - Create a new spinner
- `(s *Spinner) Start()` - Begin animation
- `(s *Spinner) Stop()` - Stop and clear line
- `(s *Spinner) StopWithSuccess(message string)` - Stop with success message
- `(s *Spinner) StopWithError(message string)` - Stop with error message
- `(s *Spinner) StopWithWarning(message string)` - Stop with warning message
- `(s *Spinner) UpdateMessage(message string)` - Change message while running
- `(s *Spinner) SetWriter(w io.Writer)` - Set custom output (for testing)

## Common Patterns

### Long-Running Operations

```go
func processFiles(files []string) error {
    spinner := display.NewSpinner("Processing files...")
    spinner.Start()

    for i, file := range files {
        spinner.UpdateMessage(fmt.Sprintf("Processing %d/%d: %s", i+1, len(files), file))

        if err := processFile(file); err != nil {
            spinner.StopWithError(fmt.Sprintf("Failed on %s: %v", file, err))
            return err
        }
    }

    spinner.StopWithSuccess(fmt.Sprintf("Processed %d files", len(files)))
    return nil
}
```

### Respecting NO_COLOR

The package automatically respects the `NO_COLOR` environment variable:

```bash
NO_COLOR=1 myapp  # Disables all colors
```

### Testing with Custom Writer

```go
func TestSpinner(t *testing.T) {
    var buf bytes.Buffer
    spinner := display.NewSpinner("Testing...")
    spinner.SetWriter(&buf)

    spinner.Start()
    time.Sleep(10 * time.Millisecond)
    spinner.StopWithSuccess("Done!")

    output := buf.String()
    if !strings.Contains(output, "Done!") {
        t.Error("Expected success message")
    }
}
```

## Important Notes

### Thread Safety

The Spinner type is thread-safe. You can call `UpdateMessage()` from any goroutine while the spinner is running.

### Terminal Detection

When colors are disabled (via `--no-color` flag or `NO_COLOR` env), the spinner falls back to static output without animation.

## See Also

- [cli](packages/cli.md) - Cobra CLI helpers
- [log](packages/log.md) - Structured logging
- [chart](packages/chart.md) - ASCII chart rendering
