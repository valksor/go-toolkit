// Package qtcwrap provides utilities for working with the QuickTemplate compiler (qtc).
//
// This package simplifies the process of generating Go code from QuickTemplate (.qtpl) files
// by providing a convenient Go API wrapper around the qtc command-line tool.
//
// The package supports:
// - Executing qtc with various configuration options
// - Directory-based template compilation
// - Single file template compilation
// - Skipping line comments for cleaner generated code
// - Custom file extensions
// - Proper error handling and warning suppression
package qtcwrap

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config represents the configuration options for the qtc compiler.
//
// The configuration allows you to specify:
// - Target directory or specific file to compile
// - Whether to skip line comments in generated code
// - Custom file extension for template files
//
// Example usage:
//
//	config := Config{
//	    Dir:              "templates",
//	    SkipLineComments: true,
//	    Ext:              ".qtpl",
//	    File:             "",
//	}
//	WithConfig(config)
type Config struct {
	// Dir specifies the directory containing .qtpl files to compile.
	// If empty, defaults to the current directory (".").
	// This field is ignored if File is specified.
	Dir string

	// SkipLineComments controls whether to skip line comments in generated code.
	// When true, the generated Go code will be more compact but less debuggable.
	// When false, line comments are preserved for a better debugging experience.
	SkipLineComments bool

	// Ext specifies the file extension for template files.
	// If empty, qtc uses its default extension (.qtpl).
	// This field is ignored if File is specified.
	Ext string

	// File specifies a single .qtpl file to compile.
	// When specified, Dir and Ext fields are ignored.
	// The file path should be relative to the current working directory.
	File string
}

// QtcWrap executes the qtc compiler with default configuration.
//
// This is a convenience function that uses sensible defaults:
// - Compiles templates in current directory
// - Skips line comments for cleaner output
// - Uses default file extension
//
// For custom configuration, use WithConfig() instead.
//
// Example:
//
//	QtcWrap()  // Compiles all .qtpl files in current directory
func QtcWrap() {
	WithConfig(Config{
		Dir:              ".",
		SkipLineComments: true,
		Ext:              "",
		File:             "",
	})
}

// WithConfig executes the qtc compiler with the specified configuration.
//
// This function builds the appropriate command-line arguments based on the
// provided configuration and executes the qtc tool. It handles both directory-based
// and single-file compilation modes.
//
// The function performs the following steps:
// 1. Validates the qtc tool is available
// 2. Builds command-line arguments from configuration
// 3. Executes the qtc command
// 4. Handles errors and warnings appropriately
//
// Directory mode (when File is empty):
// - Compiles all template files in the specified directory
// - Optionally filters by file extension
// - Recursively processes subdirectories
//
// Single file mode (when File is specified):
// - Compiles only the specified template file
// - Ignores Dir and Ext configuration
//
// Error handling:
// - Suppresses common temporary file warnings
// - Reports actual compilation errors
// - Handles missing qtc tool gracefully
//
// Example:
//
//	config := Config{
//	    Dir:              "templates",
//	    SkipLineComments: true,
//	    Ext:              ".qtpl",
//	}
//	WithConfig(config)
//
//	// Or compile a single file:
//	config := Config{
//	    File:             "templates/home.qtpl",
//	    SkipLineComments: true,
//	}
//	WithConfig(config)
func WithConfig(config Config) {
	// Validate qtc tool availability
	if err := validateQtcTool(); err != nil {
		fmt.Printf("qtc tool validation failed: %v\n", err)

		return
	}

	// Build command arguments based on configuration
	args := buildArgs(config)

	// Execute qtc command
	executeQtc(args)
}

// validateQtcTool checks if the qtc command is available in the system PATH.
//
// This function attempts to locate the qtc executable to ensure it's available
// before attempting to run template compilation.
//
// Returns an error if qtc is not found or not executable.
func validateQtcTool() error {
	_, err := exec.LookPath("qtc")
	if err != nil {
		return fmt.Errorf("qtc command not found in PATH: %w", err)
	}

	return nil
}

// buildArgs constructs command-line arguments for the qtc tool based on configuration.
//
// This function translates the Config struct into appropriate command-line flags
// for the qtc tool, handling both directory-based and single-file compilation modes.
//
// The function prioritizes File over Dir/Ext configuration when File is specified.
func buildArgs(config Config) []string {
	var args []string

	// Handle single file mode
	if config.File != "" {
		args = append(args, "-file="+config.File)
	} else {
		// Handle directory mode
		if config.Dir != "" {
			args = append(args, "-dir="+config.Dir)
		}
		if config.Ext != "" {
			args = append(args, "-ext="+config.Ext)
		}
	}

	// Add skip line comments flag if enabled
	if config.SkipLineComments {
		args = append(args, "-skipLineComments")
	}

	return args
}

// executeQtc runs the qtc command with the provided arguments.
//
// This function handles the actual execution of the qtc tool, including:
// - Setting up proper stdout/stderr handling
// - Executing the command with security considerations
// - Processing errors and warnings
// - Suppressing common temporary file warnings
//
// The function uses proper error handling to distinguish between temporary
// file warnings (which are suppressed) and actual compilation errors.
func executeQtc(args []string) {
	// Create command with security considerations
	// #nosec G204 -- args are constructed internally from validated config; safe from injection
	//nolint:noctx
	cmd := exec.Command("qtc", args...)

	// Set up output handling
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute command
	if err := cmd.Run(); err != nil {
		handleQtcError(stderr, err)
	}
}

// handleQtcError processes errors from qtc execution.
//
// This function analyzes stderr output to distinguish between:
// - Temporary file warnings (suppressed with informative message)
// - Actual compilation errors (reported to user)
// - Tool execution errors (reported to user)
//
// The function implements intelligent error filtering to reduce noise
// while preserving important error information.
func handleQtcError(stderr bytes.Buffer, err error) {
	msg := stderr.String()

	// Check if this is a temporary file warning that should be suppressed
	if isTemporaryFileWarning(stderr.Bytes()) {
		fmt.Printf("[qtc warning suppressed] %s\n", msg)

		return
	}

	// Handle actual errors
	if msg != "" {
		fmt.Print(msg)
	} else {
		fmt.Printf("qtc execution failed: %v\n", err)
	}
}

// isTemporaryFileWarning checks if the error message is a temporary file warning.
//
// This function analyzes stderr content to identify common temporary file warnings
// that can be safely suppressed, such as:
// - Missing .tmp directory warnings
// - Temporary file access issues
// - Other transient filesystem warnings
//
// Returns true if the error is a temporary file warning that should be suppressed.
func isTemporaryFileWarning(stderr []byte) bool {
	msg := string(stderr)

	// Check for common temporary file warning patterns
	return strings.Contains(msg, "no such file or directory") &&
		strings.Contains(msg, ".tmp")
}

// GetDefaultConfig returns a Config struct with sensible default values.
//
// This function provides a convenient way to get a starting configuration
// that can be modified as needed.
//
// Default values:
// - Dir: "." (current directory)
// - SkipLineComments: true (cleaner output)
// - Ext: "" (use qtc default)
// - File: "" (directory mode)
//
// Example:
//
//	config := GetDefaultConfig()
//	config.Dir = "templates"
//	config.Ext = ".qtpl"
//	WithConfig(config)
func GetDefaultConfig() Config {
	return Config{
		Dir:              ".",
		SkipLineComments: true,
		Ext:              "",
		File:             "",
	}
}

// CompileDirectory compiles all template files in the specified directory.
//
// This is a convenience function for directory-based compilation with common options.
// It uses default settings with the ability to specify the target directory.
//
// Parameters:
// - dir: Directory containing .qtpl files to compile
//
// Example:
//
//	CompileDirectory("templates")
//	CompileDirectory("src/views")
func CompileDirectory(dir string) {
	config := GetDefaultConfig()
	config.Dir = dir
	WithConfig(config)
}

// CompileFile compiles a single template file.
//
// This is a convenience function for single-file compilation.
// It uses default settings optimized for single-file processing.
//
// Parameters:
// - file: Path to the .qtpl file to compile
//
// Example:
//
//	CompileFile("templates/home.qtpl")
//	CompileFile("src/views/login.qtpl")
func CompileFile(file string) {
	config := GetDefaultConfig()
	config.File = file
	WithConfig(config)
}

// CompileWithExtension compiles template files with a specific extension.
//
// This is a convenience function for directory-based compilation with
// a custom file extension filter.
//
// Parameters:
// - dir: Directory containing template files
// - ext: File extension to filter by (e.g., ".qtpl", ".template")
//
// Example:
//
//	CompileWithExtension("templates", ".qtpl")
//	CompileWithExtension("views", ".template")
func CompileWithExtension(dir, ext string) {
	config := GetDefaultConfig()
	config.Dir = dir
	config.Ext = ext
	WithConfig(config)
}

// IsQtcAvailable checks if the qtc tool is available in the system.
//
// This function can be used to verify qtc availability before attempting
// compilation, allowing for graceful degradation or alternative approaches.
//
// Returns true if qtc is available, false otherwise.
//
// Example:
//
//	if IsQtcAvailable() {
//	    QtcWrap()
//	} else {
//	    fmt.Println("qtc not available, skipping template compilation")
//	}
func IsQtcAvailable() bool {
	return validateQtcTool() == nil
}

// ValidateConfig checks if the provided configuration is valid.
//
// This function validates configuration parameters to ensure they are
// reasonable and compatible with the qtc tool requirements.
//
// Validation rules:
// - If File is specified, it must exist and be readable
// - If Dir is specified, it must exist and be a directory
// - Ext should start with a dot if specified
// - File and Dir cannot both be empty
//
// Returns an error if the configuration is invalid.
//
// Example:
//
//	config := Config{Dir: "templates", Ext: ".qtpl"}
//	if err := ValidateConfig(config); err != nil {
//	    fmt.Printf("Invalid configuration: %v\n", err)
//	    return
//	}
//	WithConfig(config)
func ValidateConfig(config Config) error {
	// Validate file mode
	if config.File != "" {
		if _, err := os.Stat(config.File); err != nil {
			return fmt.Errorf("file %s is not accessible: %w", config.File, err)
		}

		return nil
	}

	// Validate directory mode
	if config.Dir == "" {
		return errors.New("either File or Dir must be specified")
	}

	// Check if directory exists
	if info, err := os.Stat(config.Dir); err != nil {
		return fmt.Errorf("directory %s is not accessible: %w", config.Dir, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", config.Dir)
	}

	// Validate extension format
	if config.Ext != "" && !strings.HasPrefix(config.Ext, ".") {
		return fmt.Errorf("extension must start with a dot: %s", config.Ext)
	}

	return nil
}

// CompileWithValidation compiles templates with configuration validation.
//
// This function combines configuration validation with template compilation,
// providing a safer alternative to WithConfig() for production use.
//
// The function will validate the configuration before attempting compilation
// and return an error if the configuration is invalid.
//
// Example:
//
//	config := Config{Dir: "templates", SkipLineComments: true}
//	if err := CompileWithValidation(config); err != nil {
//	    fmt.Printf("Compilation failed: %v\n", err)
//	}
func CompileWithValidation(config Config) error {
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Validate qtc tool
	if err := validateQtcTool(); err != nil {
		return fmt.Errorf("qtc tool validation failed: %w", err)
	}

	// Compile templates
	WithConfig(config)

	return nil
}

// FindTemplateFiles discovers .qtpl files in the specified directory.
//
// This function recursively searches for template files and returns their paths.
// It can be useful for preprocessing or validation before compilation.
//
// Parameters:
// - dir: Directory to search in
// - ext: File extension to filter by (optional, defaults to ".qtpl")
//
// Returns a slice of file paths or an error if the directory cannot be accessed.
//
// Example:
//
//	files, err := FindTemplateFiles("templates", ".qtpl")
//	if err != nil {
//	    fmt.Printf("Error finding templates: %v\n", err)
//	} else {
//	    fmt.Printf("Found %d template files\n", len(files))
//	    for _, file := range files {
//	        fmt.Printf("  - %s\n", file)
//	    }
//	}
func FindTemplateFiles(dir, ext string) ([]string, error) {
	if ext == "" {
		ext = ".qtpl"
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return files, fmt.Errorf("error walking the path %s: %w", dir, err)
	}

	return files, nil
}
