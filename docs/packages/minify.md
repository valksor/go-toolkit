# minify

JavaScript and CSS minification with content-based hashing.

## Description

The `minify` package provides JavaScript and CSS minification capabilities with content-based hashing for cache-busting. It supports both bundle-based and single-file workflows.

**Key Features:**
- Bundle configuration via JSON files
- Content-based hashing for cache-busting
- Automatic cleanup of old bundle versions
- Support for JavaScript and CSS minification
- Glob pattern support for flexible file selection

## Installation

```go
import "github.com/valksor/go-toolkit/minify"
```

## Usage

### Bundle-Based Workflow

Create a `bundles.json` file:

```json
{
  "bundles": [
    {
      "name": "app",
      "files": [
        "src/js/**/*.js",
        "vendor/**/*.js"
      ]
    },
    {
      "name": "styles",
      "files": [
        "src/css/**/*.css"
      ]
    }
  ]
}
```

Process bundles:

```go
config := minify.Config{
    BundlesFile: "bundles.json",
    OutputDir:   "./assets/static",
}
err := minify.ProcessBundles(config)
```

Output files:
- `app.a1b2c3d4.min.js`
- `styles.e5f6g7h8.min.css`

### Single File Workflow

```go
// Minify a single file with versioning
hash, err := minify.File("src/app.js", "assets/static/app.min.js")
// Creates: assets/static/app.{hash}.min.js
```

### Glob Pattern Support

```json
{
  "bundles": [
    {
      "name": "all-js",
      "files": [
        "src/**/*.js",
        "vendor/**/*.js"
      ]
    }
  ]
}
```

## API Reference

### Types

- `Config` - Configuration for bundle processing
- `Bundle` - Single bundle configuration
- `BundleConfig` - Structure of bundle configuration file

### Functions

- `ProcessBundles(Config) error` - Process all bundles defined in configuration
- `File(input, output string) (string, error)` - Minify a single file

## Common Patterns

### Build Process Integration

```go
// +build ignore

package main

func main() {
    config := minify.Config{
        BundlesFile: "bundles.json",
        OutputDir:   "./assets/static",
    }

    if err := minify.ProcessBundles(config); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Bundles processed successfully")
}
```

### Watch Mode Development

```go
// Re-minify on file changes
func watchAndRebuild() {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add("src/")

    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                minify.ProcessBundles(config)
            }
        }
    }
}
```

### Versioned Asset References

```go
// Get the hash of processed bundle
hash, _ := minify.File("src/app.js", "assets/static/app.min.js")

// Reference in templates
url := fmt.Sprintf("/assets/static/app.%s.min.js", hash)
```

## Bundle Configuration

### Bundle Options

```json
{
  "bundles": [
    {
      "name": "bundle-name",
      "files": [
        "pattern/**/*.js",
        "another-pattern/**/*.css"
      ]
    }
  ]
}
```

- `name` - Bundle identifier (used in output filename)
- `files` - List of glob patterns to include

### Content Hashing

The hash is generated from the minified content using xxHash, ensuring:
- Same content always produces the same hash
- Different content produces different hashes
- Efficient hash computation

## File Naming Convention

### Bundles
- Format: `{name}.{hash}.min.{ext}`
- Example: `app.a1b2c3d4.min.js`

### Single Files
- Format: `{basename}.{hash}.min.{ext}`
- Example: `app.e5f6g7h8.min.css`

## Cleanup

Old versions of bundles are automatically cleaned up during processing, keeping only the latest version for each bundle name.

## Dependencies

- [github.com/tdewolff/minify/v2](https://github.com/tdewolff/minify) - Minification engine
- [github.com/cespare/xxhash/v2](https://github.com/cespare/xxhash) - Fast hashing

## See Also

- [qtcwrap](packages/qtcwrap.md) - QuickTemplate compilation
