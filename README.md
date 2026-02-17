# Valksor Toolkit

A collection of Go utilities and packages for Valksor projects.

## What is go-toolkit?

go-toolkit provides a comprehensive set of reusable Go packages designed to accelerate development across Valksor projects. Each package is focused, well-tested, and ready for production use.

## When to use it

You're building a Go project and need:
- Configuration management (YAML/JSON, environment variables)
- CLI helpers (Cobra-based patterns, spinners, ASCII charts)
- Terminal display formatting with colors
- Build utilities (template compilation, asset minification)
- Event-driven communication between components
- JSON-RPC 2.0 protocol support for plugins
- Validation result tracking with severity levels
- Common utilities (caching, error handling, output processing)
- Retry logic with exponential backoff
- URL-safe slug generation
- Filename parsing and template expansion
- History tracking with auto-pruning

## Installation

```bash
go get github.com/valksor/go-toolkit
```

## Packages

### Configuration & Environment

| Package | Description | Docs |
|---------|-------------|------|
| **cfg** | Configuration loading, saving, and merging (YAML, JSON) | [docs](https://valksor.com/docs/go-toolkit/#/packages/cfg) |
| **env** | Environment variable expansion and layered environments | [docs](https://valksor.com/docs/go-toolkit/#/packages/env) |
| **envconfig** | Struct-based environment variable loading with validation | [docs](https://valksor.com/docs/go-toolkit/#/packages/envconfig) |
| **project** | Project detection and context management | [docs](https://valksor.com/docs/go-toolkit/#/packages/project) |

### CLI & Display

| Package | Description | Docs |
|---------|-------------|------|
| **cli** | Cobra CLI helpers and common patterns | [docs](https://valksor.com/docs/go-toolkit/#/packages/cli) |
| **cli/disambiguate** | Symfony-style command shortcuts (`c:v` → `config validate`) | [docs](https://valksor.com/docs/go-toolkit/#/packages/disambiguate) |
| **display** | Color formatting, spinners, and terminal output utilities | [docs](https://valksor.com/docs/go-toolkit/#/packages/display) |
| **log** | Structured logging helpers | [docs](https://valksor.com/docs/go-toolkit/#/packages/log) |

### Build & Template Tools

| Package | Description | Docs |
|---------|-------------|------|
| **qtcwrap** | Wrapper for QuickTemplate (qtc) compiler | [docs](https://valksor.com/docs/go-toolkit/#/packages/qtcwrap) |
| **minify** | JavaScript and CSS minification with content-based hashing | [docs](https://valksor.com/docs/go-toolkit/#/packages/minify) |

### Utilities

| Package | Description | Docs |
|---------|-------------|------|
| **cache** | Thread-safe in-memory TTL cache with automatic expiration | [docs](https://valksor.com/docs/go-toolkit/#/packages/cache) |
| **chart** | ASCII chart rendering (bar, line, pie) for terminal output | [docs](https://valksor.com/docs/go-toolkit/#/packages/chart) |
| **errors** | Error handling and wrapping utilities | [docs](https://valksor.com/docs/go-toolkit/#/packages/errors) |
| **eventbus** | Generic pub/sub event system with async publishing | [docs](https://valksor.com/docs/go-toolkit/#/packages/eventbus) |
| **history** | Generic JSON-persisted history tracking with auto-pruning | [docs](https://valksor.com/docs/go-toolkit/#/packages/history) |
| **naming** | Filename parsing, template expansion, and ticket ID extraction | [docs](https://valksor.com/docs/go-toolkit/#/packages/naming) |
| **output** | Output processing utilities including deduplicating writer | [docs](https://valksor.com/docs/go-toolkit/#/packages/output) |
| **paths** | File path manipulation utilities | [docs](https://valksor.com/docs/go-toolkit/#/packages/paths) |
| **retry** | Retry operations with exponential backoff and jitter | [docs](https://valksor.com/docs/go-toolkit/#/packages/retry) |
| **slug** | Convert text to URL-safe slugs | [docs](https://valksor.com/docs/go-toolkit/#/packages/slug) |
| **validate** | Generic validation result system with severity levels | [docs](https://valksor.com/docs/go-toolkit/#/packages/validate) |
| **version** | Build version information helpers | [docs](https://valksor.com/docs/go-toolkit/#/packages/version) |

## Documentation

Full documentation available at [valksor.com/docs/go-toolkit](https://valksor.com/docs/go-toolkit)

## Development

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Run quality checks (lint, format, security)
make quality

# Format code
make fmt

# Clean dependencies
make tidy
```

## License

BSD 3-Clause License
