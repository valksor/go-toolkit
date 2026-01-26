package license

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// goModule represents a Go module from `go list -json -m all`.
type goModule struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
	Dir     string `json:"Dir"`
}

// detectLicenses scans the Go module at rootDir and returns all dependency licenses.
func detectLicenses(ctx context.Context, rootDir string) ([]PackageLicense, error) {
	// Get all dependencies using go list
	modules, err := listModules(ctx, rootDir)
	if err != nil {
		return nil, fmt.Errorf("list modules: %w", err)
	}

	// Get Go module cache
	modCache, err := getModCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("get module cache: %w", err)
	}

	// Detect licenses for each module
	result := make([]PackageLicense, 0, len(modules))
	for _, mod := range modules {
		// Skip standard library
		if mod.Path == "" || strings.HasPrefix(mod.Path, "std/") {
			continue
		}

		licenseType, unknown := detectModuleLicense(mod, modCache)

		result = append(result, PackageLicense{
			Path:    mod.Path,
			License: licenseType,
			Unknown: unknown,
		})
	}

	// Sort by path for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].Path < result[j].Path
	})

	return result, nil
}

// listModules runs `go list -json -m all` and returns the modules.
func listModules(ctx context.Context, rootDir string) ([]goModule, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-json", "-m", "all")
	cmd.Dir = rootDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start go list: %w", err)
	}

	// Read stderr to check for errors
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Log or ignore stderr output
		}
	}()

	// Parse JSON output
	var modules []goModule
	decoder := json.NewDecoder(stdout)
	for decoder.More() {
		var mod goModule
		if err := decoder.Decode(&mod); err != nil {
			_ = cmd.Wait()

			return nil, fmt.Errorf("decode module: %w", err)
		}
		modules = append(modules, mod)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("go list: %w", err)
	}

	return modules, nil
}

// getModCache returns the Go module cache directory.
func getModCache(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "env", "GOMODCACHE")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get GOMODCACHE: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// detectModuleLicense detects the license for a module.
func detectModuleLicense(mod goModule, modCache string) (string, bool) {
	// If we have a Dir, use it directly
	var moduleDir string
	if mod.Dir != "" {
		moduleDir = mod.Dir
	} else {
		// Fallback: construct path from modCache
		// Format: modCache/path@version
		moduleDir = filepath.Join(modCache, mod.Path+"@"+mod.Version)
	}

	// Try to find license file in the module directory
	licensePath := findLicenseFile(moduleDir)
	if licensePath == "" {
		return "Unknown", true
	}

	// Read license file to detect SPDX identifier
	return identifyLicense(licensePath)
}

// findLicenseFile searches for a license file in the module directory.
func findLicenseFile(dir string) string {
	// Common license file names
	candidates := []string{
		"LICENSE",
		"LICENSE.md",
		"LICENSE.txt",
		"LICENCE",
		"LICENCE.md",
		"LICENCE.txt",
		"COPYING",
		"LICENSE-MIT",
		"LICENSE-APACHE",
		"LICENSE-BSD",
	}

	for _, name := range candidates {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try case-insensitive search in directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		name := strings.ToUpper(entry.Name())
		if strings.HasPrefix(name, "LICENSE") || strings.HasPrefix(name, "LICENCE") || name == "COPYING" {
			return filepath.Join(dir, entry.Name())
		}
	}

	return ""
}

// identifyLicense reads the license file and attempts to identify the SPDX license.
func identifyLicense(path string) (string, bool) {
	file, err := os.Open(path)
	if err != nil {
		return "Unknown", true
	}
	defer func() { _ = file.Close() }()

	// Read first few lines to look for SPDX identifier
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() && lineCount < 10 {
		line := scanner.Text()
		lineCount++

		// Look for SPDX identifier: SPDX-License-Identifier: ...
		if strings.Contains(line, "SPDX-License-Identifier:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				spdx := strings.TrimSpace(parts[1])
				if spdx != "" && spdx != "NONE" {
					return spdx, false
				}
			}
		}

		// Heuristic: check for common license text patterns
		if strings.Contains(line, "MIT License") {
			return "MIT", false
		}
		if strings.Contains(line, "Apache License") {
			return "Apache-2.0", false
		}
		if strings.Contains(line, "BSD License") {
			return "BSD-3-Clause", false
		}
	}

	// Infer from filename as fallback
	base := strings.ToUpper(filepath.Base(path))
	switch {
	case strings.Contains(base, "MIT"):
		return "MIT", false
	case strings.Contains(base, "APACHE"):
		return "Apache-2.0", false
	case strings.Contains(base, "BSD"):
		return "BSD-3-Clause", false
	case strings.Contains(base, "GPL"):
		return "GPL", false
	default:
		// Default to Apache-2.0 for generic LICENSE files
		return "Apache-2.0", false
	}
}
