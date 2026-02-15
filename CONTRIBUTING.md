# Contributing to go-toolkit

Thank you for your interest in contributing to go-toolkit! This is a shared Go library used across Valksor projects.

## Development Setup

### Prerequisites

- **Go 1.25+** - Required for building
- **Git** - Required for version control
- **make** - For build automation

### Getting Started

```bash
git clone https://github.com/valksor/go-toolkit.git
cd go-toolkit

make deps
make test
```

### Running Tests

```bash
make test        # Run all tests
make coverage    # Generate coverage
make quality     # Run linters
```

## Code Style

- **Imports**: stdlib → third-party (alphabetical)
- **Naming**: PascalCase exported, camelCase unexported
- **Errors**: `fmt.Errorf("pkg: %w", err)`
- **Modern Go (1.25+)**: Use `slices`, `maps`, `log/slog`, `context.Context`

## Testing

- Use table-driven tests
- Target 80%+ coverage
- Place tests next to source (`foo_test.go`)

## Pull Request Process

1. **Format**: `make fmt`
2. **Lint**: `make quality`
3. **Test**: `make test`
4. Fork → branch → commit → push → PR

## What Belongs in go-toolkit?

- Generic utilities with no project-specific dependencies
- Cross-project shared functionality
- Well-tested, stable APIs

## License

BSD 3-Clause License. See [LICENSE](LICENSE).
