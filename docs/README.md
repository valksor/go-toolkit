# Valksor Go Toolkit

A collection of Go utilities and packages for Valksor projects.

## Overview

go-toolkit provides a comprehensive set of reusable Go packages designed to accelerate development across Valksor projects. Each package is focused, well-tested, and ready for production use.

## Installation

```bash
go get github.com/valksor/go-toolkit
```

## Packages

### Configuration & Environment

- **[cfg](packages/cfg.md)** - Configuration loading, saving, and merging utilities (YAML, JSON)
- **[env](packages/env.md)** - Environment variable expansion and layered environment management
- **[envconfig](packages/envconfig.md)** - Struct-based environment variable loading with validation
- **[project](packages/project.md)** - Project detection and context management

### CLI & Display

- **[cli](packages/cli.md)** - Cobra CLI helpers and common patterns
- **[display](packages/display.md)** - Color and formatting utilities for terminal output
- **[log](packages/log.md)** - Structured logging helpers

### Build & Template Tools

- **[qtcwrap](packages/qtcwrap.md)** - Wrapper for QuickTemplate (qtc) compiler
- **[minify](packages/minify.md)** - JavaScript and CSS minification with content-based hashing

### Utilities

- **[cache](packages/cache.md)** - Thread-safe in-memory TTL cache with automatic expiration
- **[errors](packages/errors.md)** - Error handling and wrapping utilities
- **[output](packages/output.md)** - Output processing utilities including deduplicating writer
- **[paths](packages/paths.md)** - File path manipulation utilities
- **[version](packages/version.md)** - Build version information helpers

## Quick Links

- [Quick Start Guide](quickstart.md) - Get started in 5 minutes
- [Package Documentation](#packages) - Detailed documentation for each package
- [Contributing](legal/contributing.md) - Learn how to contribute

## Popular Packages

### envconfig

The `envconfig` package provides struct-based environment variable loading with validation:

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

### cache

The `cache` package provides a thread-safe in-memory TTL cache:

```go
import "github.com/valksor/go-toolkit/cache"
import "time"

c := cache.New()
c.Set("key", data, 5*time.Minute)

if val, ok := c.Get("key"); ok {
    // Use val (type assert to expected type)
}
```

### minify

The `minify` package handles JavaScript and CSS minification:

```go
import "github.com/valksor/go-toolkit/minify"

config := minify.Config{
    BundlesFile: "bundles.json",
    OutputDir:   "./assets/static",
}
minify.ProcessBundles(config)
```

## Documentation

- [Browse all packages](#packages)
- [Quick Start Guide](quickstart.md)
- [Contributing](legal/contributing.md)
- [License](legal/license.md)

## Links

- [GitHub Repository](https://github.com/valksor/go-toolkit)
- [pkg.go.dev](https://pkg.go.dev/github.com/valksor/go-toolkit)
- [Report Issues](https://github.com/valksor/go-toolkit/issues)
