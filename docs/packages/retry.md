# retry

Retry operations with exponential backoff and jitter.

## Description

The `retry` package provides utilities for retrying operations with configurable exponential backoff, jitter to prevent thundering herd problems, and context cancellation support.

**Key Features:**
- Configurable retry attempts, delays, and exponential backoff
- Jitter support to prevent thundering herd
- Context cancellation support
- Custom retryable error detection
- Auto-detects temporary errors (implementing `Temporary() bool`)
- Manual retry control via `RetryContext`

## Installation

```go
import "github.com/valksor/go-toolkit/retry"
```

## Usage

### Basic Retry

```go
import "github.com/valksor/go-toolkit/retry"
import "context"

config := retry.DefaultConfig()
err := config.Do(ctx, func() error {
    return callExternalAPI()
})
```

### Custom Configuration

```go
config := retry.Config{
    MaxAttempts:     5,
    BaseDelay:       100 * time.Millisecond,
    MaxDelay:        10 * time.Second,
    ExponentialBase: 2.0,
    Jitter:          true,
}
err := config.Do(ctx, func() error {
    return callExternalAPI()
})
```

### Custom Retryable Errors

```go
config := retry.Config{
    MaxAttempts: 3,
    IsRetryableFunc: func(err error) bool {
        // Only retry on specific error types
        return errors.Is(err, ErrTemporary) || errors.Is(err, ErrRateLimit)
    },
}
err := config.Do(ctx, func() error {
    return callExternalAPI()
})
```

### With Context-Aware Function

```go
config := retry.DefaultConfig()
err := config.DoWithContext(ctx, func(ctx context.Context) error {
    // Function receives context for cancellation
    return callAPIWithContext(ctx)
})
```

### Manual Retry Control

```go
config := retry.DefaultConfig()
rc := retry.NewRetryContext(config)

for rc.ShouldContinue() {
    err := doOperation()
    if err == nil {
        break
    }

    if !rc.HandleError(err) {
        return err // Not retryable
    }

    if err := rc.Delay(ctx); err != nil {
        return err // Context cancelled
    }
}
```

## API Reference

### Types

#### Config
```go
type Config struct {
    MaxAttempts     int             // Maximum number of retry attempts (default: 3)
    BaseDelay       time.Duration   // Base delay between retries (default: 1s)
    MaxDelay        time.Duration   // Maximum delay (default: 60s)
    ExponentialBase float64         // Multiplier for exponential backoff (default: 2.0)
    Jitter          bool            // Add randomness to delay (default: true)
    IsRetryableFunc IsRetryableFunc // Custom retryable check (optional)
}
```

#### RetryContext
```go
type RetryContext struct {
    // Contains unexported fields for retry state management
}
```

### Functions

- `DefaultConfig() Config` - Returns a config with sensible defaults
- `NewRetryContext(config Config) *RetryContext` - Creates a new retry context for manual control

### Methods

#### Config
- `(c Config) Do(ctx context.Context, fn func() error) error` - Execute function with retry
- `(c Config) DoWithContext(ctx context.Context, fn func(context.Context) error) error` - Execute with context-aware function
- `(c Config) IsRetryable(err error) bool` - Check if error should trigger retry
- `(c Config) CalculateDelay(attempt int) time.Duration` - Calculate delay for attempt number

#### RetryContext
- `(r *RetryContext) ShouldContinue() bool` - Returns true if more attempts allowed
- `(r *RetryContext) HandleError(err error) bool` - Record error and returns true if should retry
- `(r *RetryContext) Delay(ctx context.Context) error` - Wait for appropriate delay
- `(r *RetryContext) LastError() error` - Returns most recent error
- `(r *RetryContext) AttemptCount() int` - Returns number of attempts made

## Common Patterns

### HTTP Client with Retry

```go
func fetchWithRetry(ctx context.Context, url string) (*http.Response, error) {
    config := retry.Config{
        MaxAttempts: 3,
        BaseDelay:   100 * time.Millisecond,
        MaxDelay:    5 * time.Second,
    }

    var resp *http.Response
    err := config.Do(ctx, func() error {
        var err error
        resp, err = http.Get(url)
        if err != nil {
            return err // Network errors are retryable
        }
        if resp.StatusCode >= 500 {
            resp.Body.Close()
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }
        return nil
    })

    return resp, err
}
```

### Database Query Retry

```go
func queryWithRetry(ctx context.Context, db *sql.DB, query string) error {
    config := retry.DefaultConfig()
    config.IsRetryableFunc = func(err error) bool {
        // Only retry on connection errors
        if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) {
            return true
        }
        return false
    }

    return config.Do(ctx, func() error {
        _, err := db.ExecContext(ctx, query)
        return err
    })
}
```

### Rate Limit Handling

```go
func callWithRateLimit(ctx context.Context) error {
    config := retry.Config{
        MaxAttempts:     10,
        BaseDelay:       1 * time.Second,
        MaxDelay:        60 * time.Second,
        ExponentialBase: 2.0,
    }
    config.IsRetryableFunc = func(err error) bool {
        return errors.Is(err, ErrRateLimit)
    }

    return config.Do(ctx, func() error {
        return callAPI()
    })
}
```

## See Also

- [cache](cache.md) - In-memory caching for reducing retry frequency
- [errors](errors.md) - Error handling utilities
