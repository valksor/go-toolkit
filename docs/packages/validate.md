# validate

Generic validation result system with severity levels and multiple output formats.

## Overview

The `validate` package provides a flexible way to track and report validation findings. It supports:

- Three severity levels: error, warning, info
- Findings with codes, messages, paths, and fix suggestions
- Multiple output formats (text, JSON)
- Result merging for combining multiple validations

## Installation

```bash
go get github.com/valksor/go-toolkit/validate
```

## Usage

### Creating a Validation Result

```go
result := validate.NewResult()

// Add findings
result.AddError("INVALID_PORT", "Port must be between 1-65535", "server.port", "config.yaml")
result.AddWarning("DEPRECATED", "This field is deprecated", "legacy.enabled", "config.yaml")
result.AddInfo("INFO", "Using default value", "timeout", "config.yaml")

// Check if valid
if !result.Valid {
    fmt.Println(result.Format("text"))
}
```

### Adding Findings with Suggestions

```go
result.AddErrorWithSuggestion(
    "INVALID_VALUE",
    "Value must be positive",
    "worker.count",
    "config.yaml",
    "Change to a positive integer",
)
```

### Merging Results

```go
configResult := validateConfig(config)
securityResult := validateSecurity(config)

// Combine results
configResult.Merge(securityResult)

if !configResult.Valid {
    fmt.Println(configResult.Format("text"))
}
```

### Output Formats

#### Text Format (default)

```
config.yaml:
  ERROR [INVALID_PORT] server.port: Port must be between 1-65535
    Suggestion: Change to a value between 1 and 65535
  WARNING [DEPRECATED] legacy.enabled: This field is deprecated
    Suggestion: Migrate to the new configuration format

Summary: 1 error(s), 1 warning(s)
Configuration is INVALID
```

#### JSON Format

```json
{
  "findings": [
    {
      "severity": "error",
      "code": "INVALID_PORT",
      "message": "Port must be between 1-65535",
      "path": "server.port",
      "file": "config.yaml",
      "suggestion": "Change to a value between 1 and 65535"
    }
  ],
  "errors": 1,
  "warnings": 1,
  "valid": false
}
```

## API Reference

### Types

#### Severity

```go
type Severity string

const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)
```

#### Finding

```go
type Finding struct {
    Severity   Severity // error, warning, or info
    Code       string   // Machine-readable error code
    Message    string   // Human-readable message
    Path       string   // Config path (e.g., "server.port")
    File       string   // Source file
    Suggestion string   // How to fix the issue
}
```

#### Result

```go
type Result struct {
    Findings []Finding // All findings
    Errors   int       // Total errors
    Warnings int       // Total warnings
    Valid    bool      // true if no errors
}
```

### Functions

#### NewResult

Creates an empty validation result.

```go
func NewResult() *Result
```

#### AddError

Adds an error finding. Sets `Valid` to false.

```go
func (r *Result) AddError(code, message, path, file string)
```

#### AddWarning

Adds a warning finding. Does not affect validity.

```go
func (r *Result) AddWarning(code, message, path, file string)
```

#### AddInfo

Adds an informational finding. Does not affect validity.

```go
func (r *Result) AddInfo(code, message, path, file string)
```

#### Format

Returns the result in the specified format ("json" or "text").

```go
func (r *Result) Format(format string) string
```

#### Merge

Combines another result into this one.

```go
func (r *Result) Merge(other *Result)
```

## Best Practices

1. **Use specific error codes**: Create codes that uniquely identify the issue type
2. **Provide actionable suggestions**: Help users fix the issue
3. **Include file paths**: Make it easy to locate the problem
4. **Use warnings for deprecations**: Warnings don't fail validation but should be addressed
5. **Group related validations**: Use `Merge()` to combine results from different validators

## Example: Config Validator

```go
func validateConfig(cfg *Config) *validate.Result {
    result := validate.NewResult()

    // Validate server port
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        result.AddError(
            "INVALID_PORT",
            fmt.Sprintf("Port %d is out of range", cfg.Server.Port),
            "server.port",
            cfg.File,
        )
    }

    // Validate database URL
    if cfg.Database.URL == "" {
        result.AddError(
            "MISSING_URL",
            "Database URL is required",
            "database.url",
            cfg.File,
        )
    }

    // Check for deprecated fields
    if cfg.Legacy.Enabled {
        result.AddWarning(
            "DEPRECATED",
            "legacy.enabled is deprecated, use features.v2 instead",
            "legacy.enabled",
            cfg.File,
        )
    }

    return result
}
```
