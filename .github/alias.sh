#!/bin/bash
set -euo pipefail

# Check unnecessary aliased imports where no conflict exists
# Flags alias != basename and where basename is not already imported
# Exception: when basename equals the current package name, alias is NOT flagged
# (since using the same name would conflict with the current package)

find . -type f -name '*.go' \
  -not -name '*.qtpl.go' \
  -not -path '*/vendor/*' \
  | while read -r file; do
    gawk '
      BEGIN {
        inblock = 0
        line_count = 0
        current_package = ""
      }

      # Store all lines of the file and detect current package
      {
        lines[line_count++] = $0
        # Detect package declaration
        if (match($0, /^package\s+([a-zA-Z_][a-zA-Z0-9_]*)/, pkg)) {
          current_package = pkg[1]
        }
      }

      # First pass: collect all identifiers in use (aliases and base names)
      /^\s*import\s*\(/ { inblock = 1; next }
      inblock && /^\s*\)/ { inblock = 0; next }

      inblock && match($0, /^\s*"([^"]+)"\s*$/, m) {
        split(m[1], parts, "/")
        base = parts[length(parts)]
        used[base] = 1
        next
      }

      inblock && match($0, /^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+"([^"]+)"/, m) {
        alias = m[1]
        used[alias] = 1
        next
      }

      match($0, /^\s*import\s+"([^"]+)"\s*$/, m) {
        split(m[1], parts, "/")
        base = parts[length(parts)]
        used[base] = 1
        next
      }

      match($0, /^\s*import\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+"([^"]+)"\s*$/, m) {
        alias = m[1]
        used[alias] = 1
        next
      }

      END {
        inblock = 0
        for (i = 0; i < line_count; i++) {
          line = lines[i]

          if (line ~ /^\s*import\s*\(/) { inblock = 1; continue }
          if (inblock && line ~ /^\s*\)/) { inblock = 0; continue }

          if (inblock && match(line, /^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+"([^"]+)"/, m)) {
            alias = m[1]
            path = m[2]
            split(path, parts, "/")
            base = parts[length(parts)]
            # Skip if alias is underscore, dot, equals base, or base is already used
            # Also skip if base equals current package name (alias is necessary to avoid conflict)
            if (alias != "_" && alias != "." && alias != base && !(base in used) && base != current_package) {
              print FILENAME ":" i+1 ":" line
            }
            continue
          }

          if (!inblock && match(line, /^\s*import\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+"([^"]+)"/, m)) {
            alias = m[1]
            path = m[2]
            split(path, parts, "/")
            base = parts[length(parts)]
            # Same check for single-line import
            if (alias != "_" && alias != "." && alias != base && !(base in used) && base != current_package) {
              print FILENAME ":" i+1 ":" line
            }
          }
        }
      }
    ' "$file"
  done
