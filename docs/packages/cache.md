# cache

Thread-safe in-memory cache with TTL (time-to-live) support.

## Description

The `cache` package provides a simple, thread-safe in-memory cache with automatic expiration based on TTL. It's designed for caching provider API responses to reduce rate limit usage, but can be used for any caching needs.

**Key Features:**
- Thread-safe: All methods are safe for concurrent use
- TTL-based expiration with lazy expiration
- Optional background cleanup scheduler
- Runtime enable/disable functionality
- Zero external dependencies

## Installation

```go
import "github.com/valksor/go-toolkit/cache"
```

## Usage

### Basic Usage

```go
// Create a new cache
c := cache.New()

// Store a value with 5-minute TTL
c.Set("key", data, 5*time.Minute)

// Retrieve a value
if val, ok := c.Get("key"); ok {
    // Use val (type assert to expected type)
    data := val.(*MyType)
}

// Delete a specific key
c.Delete("key")

// Clear all entries
c.Clear()

// Get cache size
size := c.Size()
```

### Background Cleanup

```go
// Start background cleanup scheduler (optional)
stop := c.StartCleanupScheduler(1 * time.Minute)
defer close(stop) // Stop scheduler when done
```

### Enable/Disable Cache

```go
// Disable caching at runtime (all Get operations will return miss)
c.Disable()

// Re-enable caching
c.Enable()

// Check if cache is enabled
if c.Enabled() {
    // Cache is enabled
}
```

## API Reference

### Types

- `Cache` - Thread-safe in-memory cache with TTL support

### Functions

- `New() *Cache` - Creates a new Cache instance

### Methods

- `(c *Cache) Get(key string) (any, bool)` - Retrieves a value by key
- `(c *Cache) Set(key string, data any, ttl time.Duration)` - Stores a value with TTL
- `(c *Cache) Delete(key string)` - Removes a value from the cache
- `(c *Cache) Clear()` - Removes all entries from the cache
- `(c *Cache) Size() int` - Returns the number of entries
- `(c *Cache) Cleanup()` - Removes all expired entries
- `(c *Cache) Enable()` - Enables the cache
- `(c *Cache) Disable()` - Disables the cache
- `(c *Cache) Enabled() bool` - Returns true if cache is enabled
- `(c *Cache) StartCleanupScheduler(interval time.Duration) chan struct{}` - Starts periodic cleanup

## Common Patterns

### Predefined TTL Constants

```go
import "github.com/valksor/go-toolkit/cache"

ttl := cache.DefaultIssueTTL    // 5 minutes
ttl = cache.DefaultCommentsTTL  // 1 minute
ttl = cache.DefaultMetadataTTL  // 30 minutes
ttl = cache.DefaultDatabaseTTL  // 1 hour
ttl = cache.DefaultPluginTTL    // 10 minutes
```

### Caching API Responses

```go
func fetchUser(id string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", id)

    // Check cache first
    if val, ok := cache.Get(cacheKey); ok {
        return val.(*User), nil
    }

    // Fetch from API
    user, err := api.GetUser(id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    cache.Set(cacheKey, user, cache.DefaultIssueTTL)
    return user, nil
}
```

### With Cleanup Scheduler

```go
func main() {
    c := cache.New()

    // Start cleanup scheduler (runs in background)
    stop := c.StartCleanupScheduler(1 * time.Minute)
    defer close(stop)

    // Use cache...
}
```

## Important Notes

### Immutability Warning

The `Get()` method returns a reference to the stored value without copying. **Do not modify the returned value directly**, as it will corrupt the cache. For mutable types (slices, maps, structs), make a copy before modifying:

```go
if val, ok := c.Get("key"); ok {
    // WRONG: This modifies the cached value
    data := val.(*MyType)
    data.Field = "new value"

    // RIGHT: Make a copy first
    data := val.(*MyType)
    copy := *data
    copy.Field = "new value"
}
```

Alternatively, store only immutable values in the cache.

### Lazy Expiration

The cache uses lazy expiration: expired entries are not immediately removed but return a cache miss when accessed. The background cleanup scheduler handles actual deletion to avoid lock contention.
