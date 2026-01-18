# project

Project detection and context management.

## Description

The `project` package allows detecting project context based on directory patterns and local configuration directories. It's useful for tools that need to understand which project the user is currently working in.

**Key Features:**
- Project detection from local config directories
- Project registry for known projects
- Context-aware project resolution
- Multiple detection sources (local, registry, explicit, auto-detect)

## Installation

```go
import "github.com/valksor/go-toolkit/project"
```

## Usage

### Basic Project Detection

```go
cfg := &paths.Config{
    Vendor:   "valksor",
    ToolName: "mytool",
    LocalDir: ".mytool",
}

detector := project.NewDetector(cfg, ".mytool")

// Detect project from current working directory
ctx, err := detector.DetectFromCwd()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Project: %s\n", ctx.Name)
fmt.Printf("Directory: %s\n", ctx.Directory)
```

### Project Context

```go
type Context struct {
    Name           string           // Project name
    Directory      string           // Project root directory
    LocalConfigDir string           // Path to local config directory
    Config         interface{}      // Parsed configuration
    Source         DetectionSource  // How project was detected
}
```

### Detection Sources

- `SourceLocal` - Detected from local config directory
- `SourceRegistry` - Matched from project registry
- `SourceExplicit` - Explicitly specified via flag
- `SourceAutoDetect` - Auto-detected from directory name
- `SourceNone` - No project context detected

## API Reference

### Types

- `Context` - Resolved project context
- `DetectionSource` - How the project was detected
- `Detector` - Project detector
- `PathResolver` - Interface for resolving config paths
- `ConfigLoader` - Function for loading project config

### Functions

- `NewDetector(cfg PathResolver, localDir string) *Detector` - Create new detector

### Methods

- `(d *Detector) DetectFromCwd() (*Context, error)` - Detect project from current directory
- `(d *Detector) Detect(dir string) (*Context, error)` - Detect project from specific directory

## Common Patterns

### CLI Integration

```go
var projectFlag string

rootCmd := &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg := &paths.Config{
            Vendor:   "valksor",
            ToolName: "mytool",
            LocalDir: ".mytool",
        }

        detector := project.NewDetector(cfg, ".mytool")

        var ctx *project.Context
        var err error

        if projectFlag != "" {
            ctx, err = detector.Detect(projectFlag)
        } else {
            ctx, err = detector.DetectFromCwd()
        }

        if err != nil {
            return err
        }

        return runCommand(ctx)
    },
}
```

### Configuration Loading

```go
detector := project.NewDetector(cfg, ".mytool")

// Register config loader
loader := func(path string) (interface{}, error) {
    var config ProjectConfig
    data, _ := os.ReadFile(path)
    yaml.Unmarshal(data, &config)
    return config, nil
}

detector.SetConfigLoader(loader)

ctx, _ := detector.DetectFromCwd()
if ctx.Config != nil {
    config := ctx.Config.(ProjectConfig)
    // Use config...
}
```

### Project Registry

```go
// Register known projects
registry := project.NewRegistry()
registry.Add("myproject", "/path/to/myproject")
registry.AddByPattern("/workspace/projects/*")

detector.SetRegistry(registry)

ctx, _ := detector.DetectFromCwd()
```

## Detection Order

The detector searches for projects in the following order:

1. **Local Config Directory** - Looks for `.mytool` (or configured) directory
2. **Project Registry** - Checks if current directory matches known projects
3. **Explicit Flag** - Uses project specified via command-line flag
4. **Auto-Detect** - Infers project name from directory name

If no detection method succeeds, returns `SourceNone`.

## Common Use Cases

### Multi-Project Tools

Tools that operate on multiple projects can use project detection to:

- Determine which project the user is working in
- Load project-specific configuration
- Provide context-aware commands
- Validate project setup

### Workspace Management

```go
func ListProjects() error {
    detector := project.NewDetector(cfg, ".mytool")

    // Walk through directory tree and find all projects
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if info.IsDir() {
            ctx, err := detector.Detect(path)
            if err == nil && ctx.Source != project.SourceNone {
                fmt.Printf("Found project: %s at %s\n", ctx.Name, ctx.Directory)
            }
        }
        return nil
    })

    return nil
}
```

### Project Validation

```go
func ValidateProject() error {
    detector := project.NewDetector(cfg, ".mytool")
    ctx, err := detector.DetectFromCwd()

    if err != nil {
        return fmt.Errorf("not in a valid project: %w", err)
    }

    if ctx.Source == project.SourceNone {
        return errors.New("no project detected")
    }

    // Validate project structure...
    return nil
}
```

## See Also

- [paths](packages/paths.md) - Path resolution utilities
- [cfg](packages/cfg.md) - Configuration loading
