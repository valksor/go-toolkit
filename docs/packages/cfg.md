# cfg

Configuration loading, saving, and merging utilities (YAML, JSON).

## Description

The `cfg` package provides generic utilities for working with configuration files in YAML and JSON formats, with support for multi-layered configuration merging.

**Key Features:**
- Load and save YAML/JSON files
- Multi-layered configuration with precedence
- Merge configurations with different modes
- Find config files in parent directories
- Layered configuration support

## Installation

```go
import "github.com/valksor/go-toolkit/cfg"
```

## Usage

### Basic Loading and Saving

```go
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        URL string `yaml:"url"`
    } `yaml:"database"`
}

var config Config

// Load from YAML file
err := cfg.LoadYAML("config.yaml", &config)
if err != nil {
    panic(err)
}

// Save to YAML file
err = cfg.SaveYAML("config.output.yaml", &config)
```

### Finding Configuration Files

```go
// Search for config file in current and parent directories
path := cfg.FindConfigInParents(".", "config.yaml")
if path != "" {
    fmt.Println("Found config at:", path)
}

// Find all config directories
dirs := cfg.FindConfigDirs(".", []string{".git", "go.mod"})
fmt.Println("Config dirs:", dirs)
```

### Merging Configurations

```go
base := map[string]string{
    "host": "localhost",
    "port": "8080",
}
override := map[string]string{
    "host": "prod.example.com",
}

// Overlay mode (default): merges override on top of base
merged := cfg.MergeMaps(base, override, cfg.MergeModeOverlay)
// Result: {"host": "prod.example.com", "port": "8080"}

// Replace mode: replaces base with override
merged = cfg.MergeMaps(base, override, cfg.MergeModeReplace)
// Result: {"host": "prod.example.com"}
```

### Layered Configuration

```go
layers := []cfg.Layer{
    {
        Path:        "./config.local.yaml",
        Type:        cfg.LayerTypeYAML,
        Precedence:  cfg.PrecedenceLocal,
        Optional:    true,
    },
    {
        Path:        "~/.config/app/config.yaml",
        Type:        cfg.LayerTypeYAML,
        Precedence:  cfg.PrecedenceGlobal,
    },
}

merged, err := cfg.LoadLayered(layers, func(layers []interface{}) interface{} {
    // Custom merge logic
    result := layers[0] // Start with local
    // Merge other layers...
    return result
})
```

### Utility Functions

```go
// Check if file exists
exists := cfg.FileExists("config.yaml")

// Ensure directory exists
err := cfg.EnsureDir("./config/dir")

// Expand ~ to home directory
expanded := cfg.ExpandPath("~/config.yaml")

// Load all YAML files in a directory
configs, err := cfg.LoadAllYAMLInDir("./config")
```

## API Reference

### Types

- `MergeMode` - How maps are merged (`MergeModeOverlay`, `MergeModeReplace`)
- `LayeredConfig` - Multi-layered configuration
- `Layer` - Single configuration layer
- `LayerType` - Type of configuration layer (`YAML`, `JSON`, `Custom`)
- `LayerPrecedence` - Precedence level (`Local`, `Project`, `Global`)
- `MergeFunc` - Function that merges configuration layers
- `Loader` - Custom configuration loader interface

### Functions

#### Loading/Saving
- `LoadYAML(path string, v interface{}) error` - Load YAML file
- `SaveYAML(path string, v interface{}) error` - Save YAML file
- `LoadJSON(path string, v interface{}) error` - Load JSON file
- `SaveJSON(path string, v interface{}) error` - Save JSON file
- `LoadAllYAMLInDir(dir string) (map[string]interface{}, error)` - Load all YAML files in directory

#### Finding Files
- `FileExists(path string) bool` - Check if file exists
- `FindConfigInParents(startDir string, filename string) string` - Find config in parent directories
- `FindConfigDirs(startDir string, configFilenames []string) []string` - Find config directories

#### Merging
- `MergeMaps(base, override map[string]string, mode MergeMode)` - Merge two maps
- `CloneMap(map[string]string) map[string]string` - Clone a map
- `LoadLayered(layers []Layer, merge MergeFunc)` - Load layered configuration

#### Utilities
- `EnsureDir(path string) error` - Ensure directory exists
- `ExpandPath(path string) string` - Expand ~ to home directory

## Common Patterns

### Configuration with Defaults

```go
func LoadConfig() (*Config, error) {
    var config Config

    // Start with defaults
    config.Server.Host = "localhost"
    config.Server.Port = 8080

    // Override with file if exists
    if cfg.FileExists("config.yaml") {
        if err := cfg.LoadYAML("config.yaml", &config); err != nil {
            return nil, err
        }
    }

    return &config, nil
}
```

### Multi-Environment Configuration

```go
func LoadConfig(env string) (*Config, error) {
    var config Config

    // Load base config
    if err := cfg.LoadYAML("config.base.yaml", &config); err != nil {
        return nil, err
    }

    // Overlay environment-specific config
    envFile := fmt.Sprintf("config.%s.yaml", env)
    if cfg.FileExists(envFile) {
        var envConfig Config
        if err := cfg.LoadYAML(envFile, &envConfig); err != nil {
            return nil, err
        }
        // Merge configs...
    }

    return &config, nil
}
```

### User-Local Configuration

```go
func LoadConfig() (*Config, error) {
    var config Config

    // Load system-wide config
    if cfg.FileExists("/etc/app/config.yaml") {
        cfg.LoadYAML("/etc/app/config.yaml", &config)
    }

    // Overlay user config
    userConfig := filepath.Join(os.Getenv("HOME"), ".config/app/config.yaml")
    if cfg.FileExists(userConfig) {
        cfg.LoadYAML(userConfig, &config)
    }

    // Overlay local config
    if cfg.FileExists("./config.local.yaml") {
        cfg.LoadYAML("./config.local.yaml", &config)
    }

    return &config, nil
}
```
