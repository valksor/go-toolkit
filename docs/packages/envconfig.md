# envconfig

Struct-based environment variable loading with validation.

## Description

The `envconfig` package provides generic environment variable handling and configuration loading utilities. It supports loading environment variables into structs using reflection.

**Key Features:**
- Parse .env files from byte arrays
- Load system environment variables
- Merge multiple environment sources with priority
- Fill struct fields from environment variables
- Nested struct support with dot notation
- String slice support (comma-separated)

## Installation

```go
import "github.com/valksor/go-toolkit/envconfig"
```

## Usage

### Basic Struct Loading

```go
type Config struct {
    Host string `required:"true"`
    Port int
}

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
```

### Nested Structs

```go
type Config struct {
    Host string
    Database struct {
        Host string
        Port int
    }
}

env := map[string]string{
    "host": "localhost",
    "database.host": "db-server",
    "database.port": "5432",
}

config := &Config{}
err := envconfig.FillStructFromEnv(
    "",
    reflect.ValueOf(config).Elem(),
    env,
)
```

### String Slices

```go
type Config struct {
    AllowedHosts []string
}

env := map[string]string{
    "allowedhosts": "localhost,example.com,test.com",
}

config := &Config{}
err := envconfig.FillStructFromEnv(
    "",
    reflect.ValueOf(config).Elem(),
    env,
)
// config.AllowedHosts = []string{"localhost", "example.com", "test.com"}
```

### .env File Parsing

```go
data := []byte(`
DB_HOST=localhost
DB_PORT=5432
# Comment
EMPTY_VAR=
`)

envVars := envconfig.ReadDotenvBytes(data)
// Returns: map[string]string{
//     "DB_HOST": "localhost",
//     "DB_PORT": "5432",
//     "EMPTY_VAR": "",
// }
```

### Multiple Environment Sources

```go
// Load from .env file
dotenvData, _ := os.ReadFile(".env")
envFromDotenv := envconfig.ReadDotenvBytes(dotenvData)

// Load system environment
envFromSystem := envconfig.GetEnvs()

// Merge with system environment taking precedence
merged := envconfig.MergeEnvMaps(envFromDotenv, envFromSystem)
```

### Mapstructure Tags

```go
type Config struct {
    Host string `mapstructure:"api_host"`
    Port int    `mapstructure:"api_port"`
}

env := map[string]string{
    "api_host": "api.example.com",
    "api_port": "443",
}

config := &Config{}
err := envconfig.FillStructFromEnv(
    "",
    reflect.ValueOf(config).Elem(),
    env,
)
```

## API Reference

### Functions

- `ReadDotenvBytes(data []byte) map[string]string` - Parse .env file content
- `GetEnvs() map[string]string` - Get all system environment variables
- `MergeEnvMaps(...map[string]string) map[string]string` - Merge multiple env maps
- `FillStructFromEnv(prefix string, value reflect.Value, env map[string]string) error` - Fill struct from env

## Common Patterns

### Configuration with Defaults

```go
type Config struct {
    Port int `env:"PORT" default:"8080"`
    Host string
}

func LoadConfig() (*Config, error) {
    config := &Config{Port: 8080} // Default value

    envMaps := []map[string]string{
        envconfig.ReadDotenvBytes(dotenvContent),
        envconfig.GetEnvs(),
    }
    merged := envconfig.MergeEnvMaps(envMaps...)

    err := envconfig.FillStructFromEnv(
        "",
        reflect.ValueOf(config).Elem(),
        merged,
    )

    return config, err
}
```

### .env File Loading

```go
func LoadConfigWithDotenv() (*Config, error) {
    config := &Config{}

    // Try to load .env file
    dotenvData, err := os.ReadFile(".env")
    if err == nil {
        envFromDotenv := envconfig.ReadDotenvBytes(dotenvData)
        envFromSystem := envconfig.GetEnvs()
        merged := envconfig.MergeEnvMaps(envFromDotenv, envFromSystem)

        err = envconfig.FillStructFromEnv(
            "",
            reflect.ValueOf(config).Elem(),
            merged,
        )
    }

    return config, err
}
```

### Priority Layers

```go
// Lowest to highest priority
envMaps := []map[string]string{
    mapFromDefaults,      // 1. Default values
    mapFromDotenv,        // 2. .env file
    mapFromSystem,        // 3. System environment
    mapFromCommandLine,   // 4. Command-line overrides
}
merged := envconfig.MergeEnvMaps(envMaps...)
```

## Key Mapping Rules

Environment variable keys are matched against struct field names using:
1. Field names are converted to lowercase
2. Nested structs use dot notation (e.g., "database.host")
3. Mapstructure tags override default field names

## Supported Types

- `string` - Direct assignment
- `[]string` - Comma-separated values
- Nested `struct` - Dot notation
- Anonymous embedded structs - Fields promoted to parent

## See Also

- [env](packages/env.md) - Environment variable expansion
- [cfg](packages/cfg.md) - Configuration file loading
