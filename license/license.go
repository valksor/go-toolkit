// Package license provides license information and detection utilities.
//
// It can:
//   - Return an embedded project license text
//   - Scan Go module dependencies and detect their SPDX license types
//
// The license detection uses github.com/google/go-licenses for accurate
// SPDX license identification.
package license

import (
	"context"
	"fmt"
	"path/filepath"
)

// PackageLicense represents a Go package/module with its detected license.
type PackageLicense struct {
	// Path is the Go module path (e.g., "github.com/spf13/cobra").
	Path string
	// License is the SPDX license ID (e.g., "MIT", "Apache-2.0", "BSD-3-Clause").
	License string
	// Unknown is true if the license type couldn't be determined.
	Unknown bool
}

// GetDependencyLicenses returns all dependency licenses for the Go module
// at rootDir. It uses SPDX license detection via go-licenses.
//
// The rootDir should contain a go.mod file. The function scans all transitive
// dependencies and attempts to identify their licenses.
//
// Returns a slice of PackageLicense, sorted by path.
func GetDependencyLicenses(ctx context.Context, rootDir string) ([]PackageLicense, error) {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolve absolute path: %w", err)
	}

	libs, err := detectLicenses(ctx, absRoot)
	if err != nil {
		return nil, fmt.Errorf("detect licenses: %w", err)
	}

	return libs, nil
}
