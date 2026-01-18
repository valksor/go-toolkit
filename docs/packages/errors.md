# errors

Error handling and wrapping utilities.

## Description

The `errors` package provides error handling utilities for Go applications.

**Key Features:**
- Error wrapping
- Error formatting
- Common error patterns

## Installation

```go
import "github.com/valksor/go-toolkit/errors"
```

## Usage

The `errors` package provides utilities for error handling. Please refer to the source code for detailed usage examples.

## API Reference

For full API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/valksor/go-toolkit/errors).

## Common Patterns

### Error Wrapping

```go
import (
    "fmt"
    "github.com/valksor/go-toolkit/errors"
)

func DoSomething() error {
    if err := someOperation(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
```

### Error Creation

```go
// The errors package provides various utilities for error handling
// See package source code for specific functions available
```

## See Also

- Go standard library `errors` and `fmt` packages
