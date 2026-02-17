# CLAUDE.md

# IT IS YEAR 2026 !!! Please use 2026 in web searches!!!

Guidance for Claude Code when working with go-toolkit.

## Project Overview

go-toolkit is a **Go utility library** providing reusable packages for Valksor projects. Each package is focused, independent, and production-ready: configuration, CLI helpers, display formatting, event bus, caching, retry, validation, and more.

**Key constraint**: All packages must be generic with no project-specific dependencies (no mehrhof, assern, etc.).

---

## Critical Rules

### 1. Tests & Docs Required

Every feature MUST include:

| Requirement | Location | Target |
|-------------|----------|--------|
| Unit tests | `*_test.go` next to source | 80%+ coverage |
| Package docs | `doc.go` or package comment | Usage + examples |

Write tests FIRST (TDD). Use table-driven tests. Run `make test` before committing code.

### 2. Quality Checks by Scope

Run checks **only for code you changed**:

| Changed | Command |
|---------|---------|
| `*.go` | `make quality && make test` |
| `docs/**`, `*.md` | None |

If tests fail, fix them first. No exceptions for "not my code."

### 3. Use Make Commands

Always use `make` commands, not direct `go` commands:

| Operation | Command              |
|-----------|----------------------|
| Test      | `make test`          |
| Race      | `make race`          |
| Coverage  | `make coverage`      |
| Quality   | `make quality`       |
| Format    | `make fmt`           |
| Tidy      | `make tidy`          |

`make quality` runs: golangci-lint, gofmt, goimports, gofumpt, govulncheck, check-alias.

### 4. No nolint Abuse

**`//nolint` is a LAST RESORT.** Never disable linters in `.golangci.yml`.

**Acceptable**:
- `//nolint:unparam // Required by interface`
- `//nolint:errcheck // String builder WriteString won't fail`

**Never acceptable**:
- `//nolint:errcheck` without justification
- `//nolint:gosec` (fix the security issue)
- `//nolint:all` (never suppress all linters)

Always: specify linter name, include justification, place on specific line.

### 5. File Size < 500 Lines

Keep all Go files under 500 lines. Split by feature or responsibility:

```go
// Split cache.go (800 lines) into:
cache_core.go      // Core cache logic
cache_options.go   // Option types and defaults
cache_expiry.go    // Expiration handling
```

**Exceptions**: Generated code, single-responsibility modules.

### 6. Git Command Policy

All git commands are classified into three tiers. **No exceptions, no force flags, no overrides.** Tier 2 and 3 commands are never used autonomously — the agent must have explicit user instruction before running any write operation on the repository.

#### Tier 1 — Always Allowed

Safe read-only commands, available anytime:

`git status`, `git diff`, `git log`, `git show`, `git blame`, `git grep`, `git branch` (read-only), `git remote -v` (read-only), `git fetch`, `git reflog`, `git shortlog`, `git describe`, `git checkout`, `git switch`, `git restore`

#### Tier 2 — User-Requested Only

**Only use when the user explicitly asks.** Never run these commands autonomously — not for convenience, not as part of a workflow, not "to be helpful." If the task seems to need one of these commands but the user hasn't asked, ask first.

`git add`, `git commit`, `git rm`, `git mv`, `git apply`, `git am`

#### Tier 3 — Always Blocked

**NEVER use these commands.** No time window, no override, no exceptions:

`git push`, `git pull`, `git merge`, `git rebase`, `git reset`, `git revert`, `git cherry-pick`, `git tag`, `git stash` (all subcommands), `git worktree` (all subcommands), `git clean`, `git bisect`, `git notes`, `git submodule` (write operations)

Do not suggest, recommend, or implement workflows that rely on any Tier 3 command. If a task seems to need one, use a Tier 1 or Tier 2 alternative, or ask the user to perform the operation manually.

**⛔ `git worktree` — absolute prohibition.** No `git worktree add`, `remove`, `list`, `prune`, or any other worktree subcommand. Do not suggest, recommend, or implement any workflow that involves worktrees. No force flag, no override, no exceptions — ever. If a task seems to benefit from worktrees, use separate clones or branches instead.

---

## Code Style

- **Imports**: stdlib → third-party → local (alphabetical within groups)
- **Naming**: PascalCase exported, camelCase unexported
- **Errors**: `fmt.Errorf("prefix: %w", err)`; `errors.Join(errs...)`
- **Logging**: `log/slog`
- **Formatting**: `make fmt` (gofmt, goimports, gofumpt)
- **Quality**: `make quality`

### Modern Go (1.25+)

- Use `slices.Contains()`, `maps.Clone()` instead of manual loops
- Use `wg.Go(func() { ... })` instead of `wg.Add(1); go func() { defer wg.Done() }()`
- Always pass `context.Context` for cancelable operations

---

## Testing

- Run: `make test`
- Coverage: `make coverage` (output: `.coverage/coverage.html`)
- Style: Table-driven with `tests := []struct{...}{...}`
- Target: 80%+ coverage
- Race detector: `make race`

---

## See Also

- [README.md](README.md) - Package overview, installation
- [Documentation](https://valksor.com/docs/go-toolkit) - Full package docs
