// Package minify provides JavaScript and CSS minification capabilities with content-based
// hashing for cache-busting. It supports both bundle-based and single-file workflows,
// with automatic versioning and cleanup of old files.
//
// The package offers:
// - Bundle configuration via JSON files
// - Content-based hashing for cache-busting
// - Automatic cleanup of old bundle versions
// - Support for JavaScript and CSS minification
// - Glob pattern support for flexible file selection
// - Single file minification with versioning
package minify

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

// MIME type constants.
const (
	mimeJS  = "application/javascript"
	mimeCSS = "text/css"
)

// Filename format constants.
const (
	bundleFilenameFormat = "%s.%s.min.js"
	hashTruncateLength   = 8
)

// computeContentHash generates a truncated base36 hash string from content.
func computeContentHash(content string) string {
	return computeContentHashBytes([]byte(content))
}

// computeContentHashBytes generates a truncated base36 hash string from byte content.
func computeContentHashBytes(content []byte) string {
	hash := xxhash.Sum64(content)
	hashStr := strconv.FormatUint(hash, 36)
	if len(hashStr) < hashTruncateLength {
		return hashStr
	}

	return hashStr[:hashTruncateLength]
}

// Config represents the configuration for bundle processing.
// It specifies the location of the bundle configuration file and
// the output directory for processed bundles.
type Config struct {
	// BundlesFile is the path to the JSON file containing bundle definitions.
	// This file should contain a BundleConfig structure with bundle definitions.
	BundlesFile string

	// OutputDir is the directory where minified bundle files will be written.
	// The directory will be created if it doesn't exist.
	OutputDir string
}

// Bundle represents a single bundle configuration.
// Each bundle has a name and a list of file patterns to include.
type Bundle struct {
	// Name is the bundle identifier, used in the output filename.
	// The final filename will be: {Name}.{Hash}.min.js
	Name string `json:"name"`

	// Files is a list of file patterns (glob patterns) to include in the bundle.
	// Patterns are processed in order and all matching files are combined.
	Files []string `json:"files"`
}

// BundleConfig represents the structure of the bundle configuration file.
// It contains an array of bundle definitions.
type BundleConfig struct {
	// Bundles is the list of bundle configurations to process.
	Bundles []Bundle `json:"bundles"`
}

// ProcessBundles processes all bundles defined in the configuration file.
// It reads the bundle configuration, processes each bundle by combining
// and minifying the specified files, and writes the output to the configured
// output directory with content-based hashing.
//
// The function performs the following steps for each bundle:
// 1. Loads the bundle configuration from the JSON file
// 2. Processes each bundle by combining files and minifying JavaScript content
// 3. Generates a content-based hash for cache-busting
// 4. Writes the minified output to the specified directory
//
// Parameters:
//   - config: Configuration specifying the bundles file and output directory
//
// Returns:
//   - error: Any error that occurred during processing, or nil on success
//
// Example:
//
//	config := Config{
//		BundlesFile: "bundles.json",
//		OutputDir:   "./assets/static",
//	}
//	err := ProcessBundles(config)
//	if err != nil {
//		log.Fatalf("Failed to process bundles: %v", err)
//	}
func ProcessBundles(config Config) error {
	bundleConfig, err := loadBundleConfig(config.BundlesFile)
	if err != nil {
		return fmt.Errorf("failed to load bundle config: %w", err)
	}

	m := minify.New()
	m.AddFunc(mimeJS, js.Minify)

	for _, bundle := range bundleConfig.Bundles {
		if err := processBundle(m, bundle, config.OutputDir); err != nil {
			return fmt.Errorf("failed to process bundle %s: %w", bundle.Name, err)
		}
	}

	return nil
}

// loadBundleConfig loads and parses the bundle configuration file.
// It reads the JSON file and unmarshals it into a BundleConfig structure.
//
// Parameters:
//   - bundlesFile: Path to the JSON bundle configuration file
//
// Returns:
//   - *BundleConfig: The parsed bundle configuration
//   - Error: Any error that occurred during loading or parsing
func loadBundleConfig(bundlesFile string) (*BundleConfig, error) {
	data, err := os.ReadFile(bundlesFile) // #nosec G304 -- bundlesFile is validated by caller
	if err != nil {
		return nil, fmt.Errorf("failed to read bundle config file: %w", err)
	}

	var config BundleConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle config: %w", err)
	}

	return &config, nil
}

// processBundle processes a single bundle by combining files, minifying content,
// and writing the output with a content-based hash.
//
// The function:
// 1. Expands glob patterns to find all matching files
// 2. Reads and combines all file contents
// 3. Minifies the combined JavaScript content
// 4. Generates a content-based hash
// 5. Writes the minified output with the hash in the filename
//
// Parameters:
//   - minifier: The minifier instance to use for processing
//   - bundle: The bundle configuration to process
//   - outputDir: The directory where the output file should be written
//
// Returns:
//   - error: Any error that occurred during processing
func processBundle(minifier *minify.M, bundle Bundle, outputDir string) error {
	allFiles, err := collectBundleFiles(bundle.Files)
	if err != nil {
		return err
	}

	combinedContent, err := readAndCombineFiles(allFiles)
	if err != nil {
		return err
	}

	// Minify the combined content
	minified, err := minifier.String(mimeJS, combinedContent)
	if err != nil {
		return fmt.Errorf("failed to minify bundle %s: %w", bundle.Name, err)
	}

	// Generate content-based hash for cache-busting
	hashStr := computeContentHash(minified)
	filename := fmt.Sprintf(bundleFilenameFormat, bundle.Name, hashStr)
	outputPath := filepath.Join(outputDir, filename)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the minified content to the output file
	if err := os.WriteFile(outputPath, []byte(minified), 0o600); err != nil {
		return fmt.Errorf("failed to write minified file: %w", err)
	}

	return nil
}

// GetBundleHash calculates the content-based hash for a specific bundle.
// This is useful for generating cache-busting URLs or checking if a bundle
// needs to be regenerated.
//
// The function:
// 1. Loads the bundle configuration
// 2. Finds the specified bundle
// 3. Collects all files matching the bundle's patterns
// 4. Combines and minifies the content
// 5. Generates a hash from the minified content
//
// Parameters:
//   - bundleName: The name of the bundle to hash
//   - bundlesFile: Path to the bundle configuration file
//
// Returns:
//   - String: The 8-character hash string in base36 format
//   - Error: Any error that occurred during processing
//
// Example:
//
//	hash, err := GetBundleHash("base", "bundles.json")
//	if err != nil {
//		log.Fatalf("Failed to get bundle hash: %v", err)
//	}
//	fmt.Printf("Bundle hash: %s\n", hash) // Output: "a1b2c3d4"
func GetBundleHash(bundleName, bundlesFile string) (string, error) {
	bundleConfig, err := loadBundleConfig(bundlesFile)
	if err != nil {
		return "", fmt.Errorf("failed to load bundle config: %w", err)
	}

	targetBundle, err := findBundle(bundleConfig.Bundles, bundleName)
	if err != nil {
		return "", err
	}

	allFiles, err := collectBundleFiles(targetBundle.Files)
	if err != nil {
		return "", err
	}

	combinedContent, err := readAndCombineFiles(allFiles)
	if err != nil {
		return "", err
	}

	minified, err := contentMinify(combinedContent)
	if err != nil {
		return "", fmt.Errorf("failed to minify bundle %s: %w", bundleName, err)
	}

	return computeContentHash(minified), nil
}

// findBundle searches for a bundle by name in the list of bundles.
// It returns a pointer to the bundle if found, or an error if not found.
//
// Parameters:
//   - bundles: The list of bundles to search in
//   - bundleName: The name of the bundle to find
//
// Returns:
//   - *Bundle: Pointer to the found bundle
//   - error: Error if the bundle is not found
func findBundle(bundles []Bundle, bundleName string) (*Bundle, error) {
	for _, bundle := range bundles {
		if bundle.Name == bundleName {
			return &bundle, nil
		}
	}

	return nil, fmt.Errorf("bundle %s not found", bundleName)
}

// collectBundleFiles expands glob patterns to collect all matching files.
// It processes each pattern and returns a list of all files that match.
//
// Parameters:
//   - patterns: List of glob patterns to expand
//
// Returns:
//   - []string: List of all matching file paths
//   - Error: Any error that occurred during glob expansion
func collectBundleFiles(patterns []string) ([]string, error) {
	var allFiles []string
	for _, file := range patterns {
		matches, err := filepath.Glob(file)
		if err != nil {
			return nil, fmt.Errorf("failed to glob pattern %s: %w", file, err)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("no files found for pattern: %s", file)
		}
		allFiles = append(allFiles, matches...)
	}

	return allFiles, nil
}

// readAndCombineFiles reads all specified files and combines their content.
// Each file's content is separated by a newline in the combined output.
//
// Parameters:
//   - files: List of file paths to read and combine
//
// Returns:
//   - string: Combined content of all files
//   - Error: Any error that occurred during file reading
func readAndCombineFiles(files []string) (string, error) {
	var combinedContent strings.Builder
	for _, file := range files {
		content, err := os.ReadFile(file) // #nosec G304 -- file paths are from glob matches
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", file, err)
		}
		_, _ = combinedContent.Write(content)
		_, _ = combinedContent.WriteString("\n")
	}

	return combinedContent.String(), nil
}

// contentMinify minifies JavaScript content using the tdewolff/minify library.
// It creates a new minifier instance and processes the content as JavaScript.
//
// Parameters:
//   - content: The JavaScript content to minify
//
// Returns:
//   - string: The minified JavaScript content
//   - Error: Any error that occurred during minification
func contentMinify(content string) (string, error) {
	minifier := minify.New()
	minifier.AddFunc(mimeJS, js.Minify)
	result, err := minifier.String(mimeJS, content)
	if err != nil {
		return "", fmt.Errorf("failed to minify content: %w", err)
	}

	return result, nil
}

// GetBundleFilename generates the complete filename for a bundle including its hash.
// This is useful for generating URLs or checking if a bundle file exists.
//
// The filename format is: {bundleName}.{hash}.min.js
//
// Parameters:
//   - bundleName: The name of the bundle
//   - bundlesFile: Path to the bundle configuration file
//
// Returns:
//   - string: Complete filename with hash (e.g., "base.a1b2c3d4.min.js")
//   - error: Any error that occurred during processing
//
// Example:
//
//	filename, err := GetBundleFilename("base", "bundles.json")
//	if err != nil {
//		log.Fatalf("Failed to get bundle filename: %v", err)
//	}
//	fmt.Printf("Bundle filename: %s\n", filename) // Output: "base.a1b2c3d4.min.js"
func GetBundleFilename(bundleName, bundlesFile string) (string, error) {
	hash, err := GetBundleHash(bundleName, bundlesFile)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(bundleFilenameFormat, bundleName, hash), nil
}

// BundleExists checks if a bundle file already exists in the output directory.
// This is useful for determining if a bundle needs to be regenerated.
//
// Parameters:
//   - bundleName: The name of the bundle to check
//   - bundlesFile: Path to the bundle configuration file
//   - outputDir: The directory where bundle files are stored
//
// Returns:
//   - bool: True if the bundle file exists, false otherwise
//   - Error: Any error that occurred during the check
//
// Example:
//
//	exists, err := BundleExists("base", "bundles.json", "./assets/static")
//	if err != nil {
//		log.Fatalf("Failed to check bundle existence: %v", err)
//	}
//	if !exists {
//		// Bundle needs to be generated
//	}
func BundleExists(bundleName, bundlesFile, outputDir string) (bool, error) {
	filename, err := GetBundleFilename(bundleName, bundlesFile)
	if err != nil {
		return false, err
	}

	outputPath := filepath.Join(outputDir, filename)
	_, err = os.Stat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to stat bundle file: %w", err)
	}

	return true, nil
}

// CleanOldBundles removes old versions of a bundle, keeping only the current version.
// This helps prevent disk space issues by removing outdated bundle files.
//
// The function:
// 1. Determines the current bundle filename
// 2. Finds all files matching the bundle pattern
// 3. Removes files that don't match the current version
//
// Parameters:
//   - bundleName: The name of the bundle to clean
//   - bundlesFile: Path to the bundle configuration file
//   - outputDir: The directory containing bundle files
//
// Returns:
//   - error: Any error that occurred during cleanup
//
// Example:
//
//	err := CleanOldBundles("base", "bundles.json", "./assets/static")
//	if err != nil {
//		log.Printf("Warning: Failed to clean old bundles: %v", err)
//	}
func CleanOldBundles(bundleName, bundlesFile, outputDir string) error {
	currentFilename, err := GetBundleFilename(bundleName, bundlesFile)
	if err != nil {
		return err
	}

	pattern := filepath.Join(outputDir, bundleName+".*.min.js")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob old bundles: %w", err)
	}

	currentPath := filepath.Join(outputDir, currentFilename)
	for _, match := range matches {
		if match != currentPath {
			if err := os.Remove(match); err != nil {
				return fmt.Errorf("failed to remove old bundle %s: %w", match, err)
			}
		}
	}

	return nil
}

// fileContentMinify minifies file content based on the specified file type.
// It supports both CSS and JavaScript minification.
//
// Parameters:
//   - fileType: The type of file to minify ("css" or "js")
//   - content: The content to minify as a byte array
//
// Returns:
//   - []byte: The minified content
//   - error: Any error that occurred during minification
func fileContentMinify(fileType string, content []byte) ([]byte, error) {
	var minified []byte
	var err error
	switch fileType {
	case "css":
		m := minify.New()
		m.AddFunc(mimeCSS, css.Minify)
		minified, err = m.Bytes(mimeCSS, content)
	case "js":
		m := minify.New()
		m.AddFunc(mimeJS, js.Minify)
		minified, err = m.Bytes(mimeJS, content)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to minify %s content: %w", fileType, err)
	}

	return minified, nil
}

// generateMinifiedFilename generates a versioned filename for a minified file.
// The filename includes a content-based hash for cache-busting.
//
// Parameters:
//   - inputPath: The path to the original file
//   - fileType: The type of file ("css" or "js")
//   - minified: The minified content (used for hash generation)
//
// Returns:
//   - string: The generated filename with hash
func generateMinifiedFilename(inputPath, fileType string, minified []byte) string {
	hashStr := computeContentHashBytes(minified)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	if fileType == "js" {
		return fmt.Sprintf(bundleFilenameFormat, baseName, hashStr)
	}

	return fmt.Sprintf("%s.%s.%s", baseName, hashStr, fileType)
}

// cleanupOldMinifiedFiles removes old versions of minified files.
// This helps prevent accumulation of old file versions.
//
// Parameters:
//   - inputPath: The original file path
//   - outputPath: The current output file path
//   - outputDir: The output directory
//   - baseName: The base name of the file
//   - fileType: The file type extension
func cleanupOldMinifiedFiles(inputPath, outputPath, outputDir, baseName, fileType string) {
	pattern := filepath.Join(outputDir, baseName+".*."+fileType)
	matches, _ := filepath.Glob(pattern)
	for _, oldFile := range matches {
		if oldFile != inputPath && oldFile != outputPath {
			_ = os.Remove(oldFile)
		}
	}
}

// AndVersionFile minifies a file and creates a versioned copy with content-based hashing.
// This is useful for individual file processing outside the bundle system.
//
// The function:
// 1. Reads the input file
// 2. Minifies the content based on file type
// 3. Generates a content-based hash
// 4. Creates a versioned filename
// 5. Writes the minified content if it doesn't already exist
// 6. Cleans up old versions of the file
//
// Parameters:
//   - inputPath: Path to the input file to minify
//   - outputDir: Directory where the minified file should be written
//   - fileType: Type of file to minify ("css" or "js")
//
// Returns:
//   - string: The filename of the created minified file
//   - Error: Any error that occurred during processing
//
// Example:
//
//	filename, err := AndVersionFile("assets/css/main.css", "public/css", "css")
//	if err != nil {
//		log.Fatalf("Failed to minify CSS: %v", err)
//	}
//	fmt.Printf("Minified CSS: %s\n", filename) // Output: "main.a1b2c3d4.css"
func AndVersionFile(inputPath, outputDir string, fileType string) (string, error) {
	// Read the input file
	content, err := os.ReadFile(inputPath) // #nosec G304 -- inputPath is validated elsewhere
	if err != nil {
		return "", fmt.Errorf("failed to read %s file: %w", fileType, err)
	}

	// Minify the content
	minified, err := fileContentMinify(fileType, content)
	if err != nil {
		return "", fmt.Errorf("failed to minify %s: %w", fileType, err)
	}

	// Generate versioned filename
	filename := generateMinifiedFilename(inputPath, fileType, minified)
	outputPath := filepath.Join(outputDir, filename)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		return filename, nil
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the minified content
	if err := os.WriteFile(outputPath, minified, 0o600); err != nil {
		return "", fmt.Errorf("failed to write minified %s file: %w", fileType, err)
	}

	// Clean up old versions
	cleanupOldMinifiedFiles(inputPath, outputPath, outputDir, baseName, fileType)

	return filename, nil
}

// AndVersionCSS is a convenience function for minifying CSS files.
// It's a wrapper around AndVersionFile specifically for CSS files.
//
// Parameters:
//   - inputPath: Path to the input CSS file
//   - outputDir: Directory where the minified CSS file should be written
//
// Returns:
//   - string: The filename of the created minified CSS file
//   - error: Any error that occurred during processing
//
// Example:
//
//	filename, err := AndVersionCSS("assets/css/main.css", "public/css")
//	if err != nil {
//		log.Fatalf("Failed to minify CSS: %v", err)
//	}
//	fmt.Printf("Minified CSS: %s\n", filename) // Output: "main.a1b2c3d4.css"
func AndVersionCSS(inputPath, outputDir string) (string, error) {
	return AndVersionFile(inputPath, outputDir, "css")
}
