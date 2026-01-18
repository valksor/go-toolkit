# paths

File path manipulation utilities.

## Description

The `paths` package provides path resolution utilities for Valksor tools. It handles XDG directory support, home directory expansion, and provides a configurable way to manage global and local configuration paths.

**Key Features:**
- Global and local config path resolution
- Home directory expansion
- XDG directory support
- Local directory discovery
- File and directory existence checks

## Installation

```go
import "github.com/valksor/go-toolkit/paths"
```

## Usage

### Basic Configuration

```go
cfg := &paths.Config{
    Vendor:   "valksor",
    ToolName: "mytool",
    LocalDir: ".mytool",
}

// Get global config path: ~/.valksor/mytool/config.yaml
globalPath, err := cfg.GlobalConfigPath()
```

### Global Paths

```go
cfg := &paths.Config{
    Vendor:   "valksor",
    ToolName: "mytool",
    LocalDir: ".mytool",
}

// Global config directory: ~/.valksor/mytool/
globalDir, err := cfg.GlobalDir()

// Global config file: ~/.valksor/mytool/config.yaml
globalConfig, err := cfg.GlobalConfigPath()

// Global file path: ~/.valksor/mytool/mcp.json
globalFile, err := cfg.GlobalFilePath("mcp.json")

// Ensure global directory exists
dir, err := cfg.EnsureGlobalDir()
```

### Local Paths

```go
cfg := &paths.Config{
    Vendor:   "valksor",
    ToolName: "mytool",
    LocalDir: ".mytool",
}

// Find local config directory (searches upward from current dir)
localDir := cfg.FindLocalConfigDir(".")

// Local config file: /path/to/project/.mytool/config.yaml
localConfig := cfg.LocalConfigPath(localDir)

// Local file path: /path/to/project/.mytool/data.json
localFile := cfg.LocalFilePath(localDir, "data.json")

// Ensure local directory exists
dir, err := cfg.EnsureLocalDir("/path/to/project")
```

### Path Expansion

```go
// Expand ~ to home directory
expanded := paths.ExpandPath("~/config.yaml")
// Returns: /home/user/config.yaml

// Already expanded path is returned as-is
same := paths.ExpandPath("/etc/config.yaml")
// Returns: /etc/config.yaml
```

### File/Directory Checks

```go
// Check if file exists
exists := paths.FileExists("/path/to/file.txt")

// Check if directory exists
exists := paths.DirExists("/path/to/dir")
```

## API Reference

### Types

- `Config` - Path resolution configuration

### Functions

- `ExpandPath(path string) string` - Expand ~ to home directory
- `FileExists(path string) bool` - Check if file exists
- `DirExists(path string) bool` - Check if directory exists

### Methods

- `(c *Config) GlobalDir() (string, error)` - Get global config directory
- `(c *Config) GlobalConfigPath() (string, error)` - Get global config file path
- `(c *Config) GlobalFilePath(filename string) (string, error)` - Get global file path
- `(c *Config) FindLocalConfigDir(startDir string) string` - Find local config directory
- `(c *Config) LocalConfigPath(localDir string) string` - Get local config file path
- `(c *Config) LocalFilePath(localDir, filename string) string` - Get local file path
- `(c *Config) EnsureGlobalDir() (string, error)` - Create global directory if needed
- `(c *Config) EnsureLocalDir(baseDir string) (string, error)` - Create local directory if needed

## Common Patterns

### Multi-Level Configuration Loading

```go
func LoadConfig() (*Config, error) {
    cfg := &paths.Config{
        Vendor:   "valksor",
        ToolName: "mytool",
        LocalDir: ".mytool",
    }

    // Load global config
    globalPath, _ := cfg.GlobalConfigPath()
    if paths.FileExists(globalPath) {
        loadYAML(globalPath, &config)
    }

    // Load local config (overrides global)
    localDir := cfg.FindLocalConfigDir(".")
    if localDir != "" {
        localPath := cfg.LocalConfigPath(localDir)
        if paths.FileExists(localPath) {
            loadYAML(localPath, &config)
        }
    }

    return &config, nil
}
```

### Initialization

```go
func Init() error {
    cfg := &paths.Config{
        Vendor:   "valksor",
        ToolName: "mytool",
        LocalDir: ".mytool",
    }

    // Ensure global directory exists
    globalDir, err := cfg.EnsureGlobalDir()
    if err != nil {
        return fmt.Errorf("failed to create global dir: %w", err)
    }

    // Create default global config
    globalConfig := cfg.GlobalConfigPath()
    if !paths.FileExists(globalConfig) {
        createDefaultConfig(globalConfig)
    }

    return nil
}
```

### Find Project Root

```go
func FindProjectRoot() (string, error) {
    cfg := &paths.Config{
        Vendor:   "valksor",
        ToolName: "mytool",
        LocalDir: ".mytool",
    }

    localDir := cfg.FindLocalConfigDir(".")
    if localDir == "" {
        return "", errors.New("not in a project directory")
    }

    // localDir is the path to .mytool directory
    // Parent directory is the project root
    return filepath.Dir(localDir), nil
}
```

## Testing Support

The package provides testing utilities:

```go
func TestPathResolution(t *testing.T) {
    // Set custom home directory for testing
    restore := paths.SetHomeDirForTesting("/test/home")
    defer restore()

    cfg := &paths.Config{
        Vendor:   "valksor",
        ToolName: "mytool",
        LocalDir: ".mytool",
    }

    globalDir, _ := cfg.GlobalDir()
    // Returns: /test/home/valksor/mytool
}
```

## Path Examples

Given configuration:
```go
cfg := &paths.Config{
    Vendor:   "valksor",
    ToolName: "mytool",
    LocalDir: ".mytool",
}
```

- `GlobalDir()` → `~/.valksor/mytool/`
- `GlobalConfigPath()` → `~/.valksor/mytool/config.yaml`
- `GlobalFilePath("data.json")` → `~/.valksor/mytool/data.json`
- `FindLocalConfigDir("/project/subdir")` → `/project/.mytool` (if exists)
- `LocalConfigPath("/project/.mytool")` → `/project/.mytool/config.yaml`

## See Also

- [cfg](packages/cfg.md) - Configuration loading utilities
- [env](packages/env.md) - Environment variable expansion
