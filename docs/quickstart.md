# Quick Start

Get up and running with go-toolkit in 5 minutes.

## Installation

Install the toolkit using Go:

```bash
go get github.com/valksor/go-toolkit
```

## Your First Usage

### Environment Variable Expansion

The `env` package makes it easy to expand environment variables in strings:

```go
package main

import (
    "fmt"
    "github.com/valksor/go-toolkit/env"
)

func main() {
    // Expand environment variables in a string
    expanded := env.ExpandEnv("${HOME}/.config")
    fmt.Println(expanded) // Prints: /home/user/.config

    // Use layered environments
    layers := []map[string]string{
        {"DATABASE_URL": "postgres://localhost:5432/db"},
        {"DATABASE_URL": "postgres://prod.example.com:5432/db"},
    }
    e := env.New(layers...)
    fmt.Println(e.Get("DATABASE_URL")) // Uses last layer
}
```

### Struct-Based Configuration

The `envconfig` package loads environment variables into structs with validation:

```go
package main

import (
    "fmt"
    "reflect"
    "github.com/valksor/go-toolkit/envconfig"
)

type Config struct {
    Port     int    `required:"true"`
    Host     string `required:"true"`
    Debug    bool   `env:"DEBUG"`
    LogLevel string `env:"LOG_LEVEL" default:"info"`
}

func main() {
    config := &Config{}

    // Get environment variables
    envMaps := []map[string]string{
        envconfig.GetEnvs(),
    }
    merged := envconfig.MergeEnvMaps(envMaps...)

    // Fill struct from environment
    err := envconfig.FillStructFromEnv(
        "",
        reflect.ValueOf(config).Elem(),
        merged,
    )
    if err != nil {
        panic(err)
    }

    fmt.Printf("Config: %+v\n", config)
}
```

### In-Memory Cache

The `cache` package provides a thread-safe TTL cache:

```go
package main

import (
    "fmt"
    "time"
    "github.com/valksor/go-toolkit/cache"
)

func main() {
    // Create a new cache
    c := cache.New()

    // Store a value with 5-minute TTL
    c.Set("user:123", map[string]string{"name": "Alice"}, 5*time.Minute)

    // Retrieve a value
    if val, ok := c.Get("user:123"); ok {
        user := val.(map[string]string)
        fmt.Println("User:", user["name"])
    }

    // Start background cleanup (optional)
    stop := c.StartCleanupScheduler(1 * time.Minute)
    defer close(stop)
}
```

### Configuration Loading

The `cfg` package handles YAML and JSON configuration:

```go
package main

import (
    "fmt"
    "github.com/valksor/go-toolkit/cfg"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        URL string `yaml:"url"`
    } `yaml:"database"`
}

func main() {
    var config Config

    // Load from YAML file
    err := cfg.Load("config.yaml", &config)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)

    // Save configuration
    err = cfg.Save("config.output.yaml", &config)
    if err != nil {
        panic(err)
    }
}
```

### Structured Logging

The `log` package provides structured logging helpers:

```go
package main

import (
    "context"
    "log/slog"
    "github.com/valksor/go-toolkit/log"
)

func main() {
    // Create a logger
    logger := slog.New(log.NewHandler(slog.LevelInfo))

    // Log messages
    logger.Info("Application started",
        "version", "1.0.0",
        "environment", "production",
    )

    // With context
    ctx := context.Background()
    logger.InfoContext(ctx, "Processing request",
        "request_id", "abc-123",
        "user_id", "user-456",
    )
}
```

## Next Steps

- [Browse all packages](README.md#packages)
- [Read detailed documentation](README.md#packages)
- [Contributing](legal/contributing.md)

## Common Use Cases

### CLI Applications

Use the `cli` package with Cobra for building CLI tools:

```go
import "github.com/valksor/go-toolkit/cli"

// cli provides common patterns and helpers for Cobra commands
```

### Template Compilation

Use `qtcwrap` for QuickTemplate compilation:

```go
import "github.com/valksor/go-toolkit/qtcwrap"

config := qtcwrap.Config{
    Dir:              "templates",
    SkipLineComments: true,
    Ext:              ".qtpl",
}
qtcwrap.WithConfig(config)
```

### Asset Minification

Use `minify` for JavaScript and CSS assets:

```go
import "github.com/valksor/go-toolkit/minify"

config := minify.Config{
    BundlesFile: "bundles.json",
    OutputDir:   "./assets/static",
}
minify.ProcessBundles(config)
```

## Need Help?

- [Documentation](README.md)
- [Report an Issue](https://github.com/valksor/go-toolkit/issues)
- [Package Documentation](README.md#packages)
