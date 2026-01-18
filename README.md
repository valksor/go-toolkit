# go-toolkit

A collection of Go utilities and packages for Valksor projects.

## Packages

### Configuration & Environment

- **cfg**: Configuration loading, saving, and merging utilities (YAML, JSON)
- **env**: Environment variable expansion and layered environment management
- **envconfig**: Struct-based environment variable loading with validation
- **project**: Project detection and context management

### CLI & Display

- **cli**: Cobra CLI helpers and common patterns
- **display**: Color and formatting utilities for terminal output
- **log**: Structured logging helpers

### Build & Template Tools

- **qtcwrap**: Wrapper for QuickTemplate (qtc) compiler
- **minify**: JavaScript and CSS minification with content-based hashing

### Utilities

- **cache**: Thread-safe in-memory TTL cache with automatic expiration
- **errors**: Error handling and wrapping utilities
- **output**: Output processing utilities including deduplicating writer
- **paths**: File path manipulation utilities
- **version**: Build version information helpers

## Installation

```bash
go get github.com/valksor/go-toolkit
```

## Usage Examples

### Environment Variable Expansion (env package)

```go
import "github.com/valksor/go-toolkit/env"

expanded := env.ExpandEnv("${HOME}/.config")
```

### Struct-Based Configuration (envconfig package)

```go
import "github.com/valksor/go-toolkit/envconfig"

type Config struct {
    Port     int    `required:"true"`
    Database string `required:"true"`
    Debug    bool   `env:"DEBUG"`
}

config := &Config{}
envMaps := []map[string]string{
    envconfig.ReadDotenvBytes(sharedEnv),
    envconfig.GetEnvs(),
}
merged := envconfig.MergeEnvMaps(envMaps...)
err := envconfig.FillStructFromEnv("", reflect.ValueOf(config).Elem(), merged)
```

### QuickTemplate Compilation (qtcwrap package)

```go
import "github.com/valksor/go-toolkit/qtcwrap"

config := qtcwrap.Config{
    Dir:              "templates",
    SkipLineComments: true,
    Ext:              ".qtpl",
}
qtcwrap.WithConfig(config)
```

### Asset Minification (minify package)

```go
import "github.com/valksor/go-toolkit/minify"

config := minify.Config{
    BundlesFile: "bundles.json",
    OutputDir:   "./assets/static",
}
minify.ProcessBundles(config)
```

### In-Memory Cache (cache package)

```go
import "github.com/valksor/go-toolkit/cache"
import "time"

// Create a new cache
c := cache.New()

// Store a value with 5-minute TTL
c.Set("key", data, 5*time.Minute)

// Retrieve a value
if val, ok := c.Get("key"); ok {
    // Use val (type assert to expected type)
}

// Start background cleanup scheduler (optional)
stop := c.StartCleanupScheduler(1 * time.Minute)
defer close(stop)
```

### Deduplicating Output Writer (output package)

```go
import "github.com/valksor/go-toolkit/output"

// Wrap any io.Writer to suppress consecutive duplicate lines
w := output.NewDeduplicatingWriter(os.Stdout)
w.Write([]byte("Processing...\n"))
w.Write([]byte("Processing...\n"))  // This line will be suppressed
w.Write([]byte("Done!\n"))          // This line will be written

// Remember to flush when done
w.Flush()
```

## Development

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Run quality checks (lint, format, security)
make quality

# Format code
make fmt

# Clean dependencies
make tidy
```

## Dependencies

See [go.mod](go.mod) for the complete list of dependencies.

## License

BSD 3-Clause License
