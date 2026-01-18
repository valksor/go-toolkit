# version

Build version information helpers.

## Description

The `version` package provides build version information. Version variables are set via ldflags at build time.

**Key Features:**
- Build-time version information
- Formatted version output
- Git commit tracking
- Build timestamp

## Installation

```go
import "github.com/valksor/go-toolkit/version"
```

## Usage

### Setting Version at Build Time

```bash
# Build with version information
go build -ldflags \
  "-X github.com/valksor/go-toolkit/version.Version=1.0.0 \
   -X github.com/valksor/go-toolkit/version.Commit=abc123 \
   -X github.com/valksor/go-toolkit/version.BuildTime=2024-01-15T12:00:00Z"
```

### Display Version Information

```go
import "github.com/valksor/go-toolkit/version"

// Full version info
fmt.Println(version.Info("myapp"))
```

Output:
```
myapp 1.0.0
  Commit: abc123
  Built:  2024-01-15T12:00:00Z
  Go:     go1.21.5
```

### Short Version

```go
// Just the version string
fmt.Println(version.Short())
// Output: 1.0.0
```

### CLI Version Command

```go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(version.Info("myapp"))
    },
}
```

## API Reference

### Variables

- `Version string` - Application version (set via ldflags)
- `Commit string` - Git commit hash (set via ldflags)
- `BuildTime string` - Build timestamp (set via ldflags)

### Functions

- `Info(appName string) string` - Formatted version information
- `Short() string` - Just the version string
- `Set(v, c, bt string)` - Set version information (for testing)

## Common Patterns

### Makefile Integration

```makefile
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/valksor/go-toolkit/version.Version=$(VERSION) \
           -X github.com/valksor/go-toolkit/version.Commit=$(COMMIT) \
           -X github.com/valksor/go-toolkit/version.BuildTime=$(BUILDTIME)"

build:
    go build $(LDFLAGS) -o myapp
```

### GoReleaser Configuration

```yaml
builds:
  - env:
      - CGO_ENABLED=0
    ldflags:
      - -X github.com/valksor/go-toolkit/version.Version={{.Version}}
      - -X github.com/valksor/go-toolkit/version.Commit={{.Commit}}
      - -X github.com/valksor/go-toolkit/version.BuildTime={{.Date}}
    goos:
      - linux
      - darwin
      - windows
```

### GitHub Actions

```yaml
- name: Set version variables
  run: |
    echo "VERSION=$(git describe --tags --always --dirty)" >> $GITHUB_ENV
    echo "COMMIT=$(git rev-parse --short HEAD)" >> $GITHUB_ENV
    echo "BUILDTIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")" >> $GITHUB_ENV

- name: Build
  run: |
    go build -ldflags \
      "-X github.com/valksor/go-toolkit/version.Version=${{ env.VERSION }} \
       -X github.com/valksor/go-toolkit/version.Commit=${{ env.COMMIT }} \
       -X github.com/valksor/go-toolkit/version.BuildTime=${{ env.BUILDTIME }}"
```

### Development Builds

For development builds without setting ldflags, the defaults are:
- `Version = "dev"`
- `Commit = "none"`
- `BuildTime = "unknown"`

## Testing

```go
func TestVersion(t *testing.T) {
    // Set version for testing
    version.Set("1.0.0", "abc123", "2024-01-15T12:00:00Z")
    defer version.Set("dev", "none", "unknown") // Reset

    info := version.Info("testapp")
    // Assert on info...
}
```

## Default Values

If not set via ldflags:
- `Version`: "dev"
- `Commit`: "none"
- `BuildTime`: "unknown"

## See Also

- Go build `ldflags` documentation
