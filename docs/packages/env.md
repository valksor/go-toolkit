# env

Environment variable expansion and layered environment management.

## Description

The `env` package provides environment variable expansion utilities and layered environment management.

**Key Features:**
- Expand environment variables in strings
- Layered environment management
- Priority-based variable resolution

## Installation

```go
import "github.com/valksor/go-toolkit/env"
```

## Usage

### Basic Expansion

```go
import "github.com/valksor/go-toolkit/env"

// Expand environment variables in a string
expanded := env.ExpandEnv("${HOME}/.config")
fmt.Println(expanded) // Prints: /home/user/.config

// Expand multiple variables
path := env.ExpandEnv("${HOME}/app/${ENVIRONMENT}/config.yaml")
fmt.Println(path) // Prints: /home/user/app/production/config.yaml
```

### Layered Environments

```go
// Create environment with multiple layers
layers := []map[string]string{
    {"DATABASE_URL": "postgres://localhost:5432/db"},
    {"DATABASE_URL": "postgres://prod.example.com:5432/db"},
}

e := env.New(layers...)

// Get value (last layer takes precedence)
fmt.Println(e.Get("DATABASE_URL"))
// Prints: postgres://prod.example.com:5432/db
```

## API Reference

### Types

- `Env` - Layered environment manager

### Functions

- `ExpandEnv(s string) string` - Expand environment variables in string
- `New(layers ...map[string]string) *Env` - Create new layered environment

### Methods

- `(e *Env) Get(key string) string` - Get value by key (last layer wins)
- `(e *Env) Set(key, value string)` - Set a value

## Common Patterns

### Configuration File Paths

```go
func GetConfigPath() string {
    return env.ExpandEnv("${XDG_CONFIG_HOME:-$HOME/.config}/myapp/config.yaml")
}
```

### Environment-Aware Configuration

```go
func GetDatabaseURL() string {
    layers := []map[string]string{
        {"DATABASE_URL": "postgres://localhost:5432/mydb"}, // Default
        GetEnvs(), // System environment
    }
    e := env.New(layers...)
    return e.Get("DATABASE_URL")
}
```

### Cross-Platform Paths

```go
// Works on Unix and Windows
cacheDir := env.ExpandEnv("${HOME}/.cache/app")
```

## See Also

- [envconfig](packages/envconfig.md) - Struct-based environment variable loading
- [cfg](packages/cfg.md) - Configuration file loading
