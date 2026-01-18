package minify

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test file path constants.
const (
	testFile1Path = "src/file1.js"
	testFile2Path = "src/file2.js"
	testCSSPath   = "src/main.css"
	testJSPattern = "src/*.js"
)

// Test file name constants.
const (
	testFile1Name        = "file1.js"
	testFile2Name        = "file2.js"
	testBundleConfigFile = "bundles.json"
)

// File extension constants.
const (
	minJSExtension = ".min.js"
)

// Error message constants.
const (
	errFailedToRemoveTempDir    = "Failed to remove temp directory: %v"
	errFailedToProcessBundles   = "Failed to process bundles: %v"
	errFailedToVersionCSSFile   = "Failed to version CSS file: %v"
	errMsgNonExistentFile       = "Expected error for non-existent file"
	errMsgNonExistentBundle     = "Expected error for non-existent bundle"
	errMsgMinifiedSmaller       = "Expected minified content to be smaller"
	errMsgFilenameStartWithMain = "Expected filename to start with 'main.', got %s"
	errMsgFilenameEndWithMinJS  = "Expected filename to end with '.min.js', got %s"
	errMsgFilenameEndWithCSS    = "Expected filename to end with '.css', got %s"
	errMsgOutputFileExist       = "Expected output file to exist: %v"
)

// Test helper functions

func createTempDir(tb testing.TB) string {
	tb.Helper()

	return tb.TempDir()
}

func createTempFile(tb testing.TB, dir, filename, content string) string {
	tb.Helper()
	path := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		tb.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		tb.Fatalf("Failed to write file: %v", err)
	}

	return path
}

func createBundleConfig(tb testing.TB, dir string, bundles []Bundle) string {
	tb.Helper()
	config := BundleConfig{Bundles: bundles}
	data, err := json.Marshal(config)
	if err != nil {
		tb.Fatalf("Failed to marshal bundle config: %v", err)
	}

	return createTempFile(tb, dir, testBundleConfigFile, string(data))
}

// Additional helper functions for reducing test complexity

func setupTestEnvironment(tb testing.TB) string {
	tb.Helper()
	tempDir := createTempDir(tb)
	tb.Cleanup(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			tb.Logf(errFailedToRemoveTempDir, err)
		}
	})

	return tempDir
}

func createTestFiles(tb testing.TB, tempDir string, fileSpecs map[string]string) {
	tb.Helper()
	for path, content := range fileSpecs {
		createTempFile(tb, tempDir, path, content)
	}
}

func createTestBundleEnvironment(tb testing.TB, bundleName string, files []string) (string, string) {
	tb.Helper()
	tempDir := setupTestEnvironment(tb)

	// Create source files
	fileSpecs := make(map[string]string)
	for _, file := range files {
		if strings.HasSuffix(file, ".js") {
			fileSpecs[file] = sampleJS
		} else if strings.HasSuffix(file, ".css") {
			fileSpecs[file] = sampleCSS
		}
	}
	createTestFiles(tb, tempDir, fileSpecs)

	// Create bundle config
	bundles := []Bundle{{Name: bundleName, Files: []string{filepath.Join(tempDir, testJSPattern)}}}
	configFile := createBundleConfig(tb, tempDir, bundles)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")

	return configFile, outputDir
}

func assertFilenameFormat(tb testing.TB, filename, prefix, suffix string, expectedParts int) {
	tb.Helper()
	if !strings.HasPrefix(filename, prefix) {
		tb.Errorf("Expected filename to start with '%s', got %s", prefix, filename)
	}

	if !strings.HasSuffix(filename, suffix) {
		tb.Errorf("Expected filename to end with '%s', got %s", suffix, filename)
	}

	parts := strings.Split(filename, ".")
	if len(parts) != expectedParts {
		tb.Errorf("Expected filename to have %d parts, got %d: %s", expectedParts, len(parts), filename)
	}
}

func assertHashProperties(tb testing.TB, hash string, expectedLength int) {
	tb.Helper()
	if len(hash) != expectedLength {
		tb.Errorf("Expected hash length to be %d, got %d", expectedLength, len(hash))
	}
}

func assertHashConsistency(tb testing.TB, getHashFunc func(string, string) (string, error), configFile, bundleName string) {
	tb.Helper()
	// Test consistency - same content should produce same hash
	hash1, err := getHashFunc(bundleName, configFile)
	if err != nil {
		tb.Fatalf("Failed to get hash first time: %v", err)
	}

	hash2, err := getHashFunc(bundleName, configFile)
	if err != nil {
		tb.Fatalf("Failed to get hash second time: %v", err)
	}

	if hash1 != hash2 {
		tb.Errorf("Expected consistent hash, got %s and %s", hash1, hash2)
	}
}

func assertBundleExistence(tb testing.TB, exists, shouldExist bool, context string) {
	tb.Helper()
	if shouldExist && !exists {
		tb.Errorf("Expected bundle to exist %s", context)
	}
	if !shouldExist && exists {
		tb.Errorf("Expected bundle to not exist %s", context)
	}
}

// Test sample JavaScript and CSS content.
const (
	sampleJS = `
		// Sample JavaScript file
		function hello() {
			console.log("Hello, World!");
		}
		
		var x = 1;
		var y = 2;
		var z = x + y;
		
		hello();
	`

	sampleCSS = `
		/* Sample CSS file */
		body {
			margin: 0;
			padding: 0;
			font-family: Arial, sans-serif;
		}
		
		.container {
			max-width: 1200px;
			margin: 0 auto;
		}
		
		h1 {
			color: #333;
			font-size: 2em;
		}
	`

	sampleJS2 = `
		// Another JavaScript file
		function goodbye() {
			console.log("Goodbye!");
		}
		
		goodbye();
	`
)

// Tests for Config type

func TestConfig(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := Config{
			BundlesFile: testBundleConfigFile,
			OutputDir:   "./output",
		}

		if config.BundlesFile != testBundleConfigFile {
			t.Errorf("Expected BundlesFile to be 'bundles.json', got %s", config.BundlesFile)
		}

		if config.OutputDir != "./output" {
			t.Errorf("Expected OutputDir to be './output', got %s", config.OutputDir)
		}
	})
}

// Tests for Bundle type

func TestBundle(t *testing.T) {
	t.Run("ValidBundle", func(t *testing.T) {
		bundle := Bundle{
			Name:  "test",
			Files: []string{testFile1Name, testFile2Name},
		}

		if bundle.Name != "test" {
			t.Errorf("Expected Name to be 'test', got %s", bundle.Name)
		}

		if len(bundle.Files) != 2 {
			t.Errorf("Expected 2 files, got %d", len(bundle.Files))
		}
	})
}

// Tests for BundleConfig type

func TestBundleConfig(t *testing.T) {
	t.Run("ValidBundleConfig", func(t *testing.T) {
		config := BundleConfig{
			Bundles: []Bundle{
				{Name: "bundle1", Files: []string{testFile1Name}},
				{Name: "bundle2", Files: []string{testFile2Name}},
			},
		}

		if len(config.Bundles) != 2 {
			t.Errorf("Expected 2 bundles, got %d", len(config.Bundles))
		}
	})
}

// Tests for loadBundleConfig function

func TestLoadBundleConfig(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		tempDir := setupTestEnvironment(t)

		bundles := []Bundle{
			{Name: "test", Files: []string{"*.js"}},
		}
		configFile := createBundleConfig(t, tempDir, bundles)

		config, err := loadBundleConfig(configFile)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if len(config.Bundles) != 1 {
			t.Errorf("Expected 1 bundle, got %d", len(config.Bundles))
		}

		if config.Bundles[0].Name != "test" {
			t.Errorf("Expected bundle name 'test', got %s", config.Bundles[0].Name)
		}
	})

	t.Run("FileNotFound", func(t *testing.T) {
		_, err := loadBundleConfig("nonexistent.json")
		if err == nil {
			t.Error(errMsgNonExistentFile)
		}

		if !strings.Contains(err.Error(), "failed to read bundle config file") {
			t.Errorf("Expected 'failed to read bundle config file' error, got: %v", err)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		tempDir := setupTestEnvironment(t)

		configFile := createTempFile(t, tempDir, "invalid.json", "{invalid json")

		_, err := loadBundleConfig(configFile)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}

		if !strings.Contains(err.Error(), "failed to unmarshal bundle config") {
			t.Errorf("Expected 'failed to read bundle config' error, got: %v", err)
		}
	})
}

// Tests for findBundle function

func TestFindBundle(t *testing.T) {
	bundles := []Bundle{
		{Name: "bundle1", Files: []string{"file1.js"}},
		{Name: "bundle2", Files: []string{"file2.js"}},
		{Name: "bundle3", Files: []string{"file3.js"}},
	}

	t.Run("BundleFound", func(t *testing.T) {
		bundle, err := findBundle(bundles, "bundle2")
		if err != nil {
			t.Fatalf("Failed to find bundle: %v", err)
		}

		if bundle.Name != "bundle2" {
			t.Errorf("Expected bundle name 'bundle2', got %s", bundle.Name)
		}
	})

	t.Run("BundleNotFound", func(t *testing.T) {
		_, err := findBundle(bundles, "nonexistent")
		if err == nil {
			t.Error(errMsgNonExistentBundle)
		}

		if !strings.Contains(err.Error(), "bundle nonexistent not found") {
			t.Errorf("Expected 'bundle nonexistent not found' error, got: %v", err)
		}
	})
}

// Tests for collectBundleFiles function

func TestCollectBundleFiles(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create test files
	createTempFile(t, tempDir, testFile1Name, sampleJS)
	createTempFile(t, tempDir, testFile2Name, sampleJS2)
	createTempFile(t, tempDir, "subdir/file3.js", sampleJS)

	t.Run("ValidPatterns", func(t *testing.T) {
		patterns := []string{
			filepath.Join(tempDir, "*.js"),
			filepath.Join(tempDir, "subdir/*.js"),
		}

		files, err := collectBundleFiles(patterns)
		if err != nil {
			t.Fatalf("Failed to collect files: %v", err)
		}

		if len(files) != 3 {
			t.Errorf("Expected 3 files, got %d", len(files))
		}
	})

	t.Run("NoFilesFound", func(t *testing.T) {
		patterns := []string{filepath.Join(tempDir, "*.nonexistent")}

		_, err := collectBundleFiles(patterns)
		if err == nil {
			t.Error("Expected error for no files found")
		}

		if !strings.Contains(err.Error(), "no files found for pattern") {
			t.Errorf("Expected 'no files found for pattern' error, got: %v", err)
		}
	})

	t.Run("InvalidPattern", func(t *testing.T) {
		patterns := []string{"["}

		_, err := collectBundleFiles(patterns)
		if err == nil {
			t.Error("Expected error for invalid pattern")
		}

		if !strings.Contains(err.Error(), "failed to glob pattern") {
			t.Errorf("Expected 'failed to glob pattern' error, got: %v", err)
		}
	})
}

// Tests for readAndCombineFiles function

func TestReadAndCombineFiles(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	file1 := createTempFile(t, tempDir, testFile1Name, "content1")
	file2 := createTempFile(t, tempDir, testFile2Name, "content2")

	t.Run("ValidFiles", func(t *testing.T) {
		files := []string{file1, file2}

		content, err := readAndCombineFiles(files)
		if err != nil {
			t.Fatalf("Failed to read and combine files: %v", err)
		}

		expected := "content1\ncontent2\n"
		if content != expected {
			t.Errorf("Expected content %q, got %q", expected, content)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		files := []string{"nonexistent.js"}

		_, err := readAndCombineFiles(files)
		if err == nil {
			t.Error(errMsgNonExistentFile)
		}

		if !strings.Contains(err.Error(), "failed to read file") {
			t.Errorf("Expected 'failed to read file' error, got: %v", err)
		}
	})
}

// Tests for contentMinify function

func TestContentMinify(t *testing.T) {
	t.Run("ValidJavaScript", func(t *testing.T) {
		content := `
			function hello() {
				console.log("Hello, World!");
			}
			hello();
		`

		minified, err := contentMinify(content)
		if err != nil {
			t.Fatalf("Failed to minify content: %v", err)
		}

		if len(minified) >= len(content) {
			t.Error(errMsgMinifiedSmaller)
		}

		// Check that the minified content contains the essential parts
		if !strings.Contains(minified, "function hello()") {
			t.Error("Expected minified content to contain function definition")
		}
	})

	t.Run("InvalidJavaScript", func(t *testing.T) {
		content := "function unclosed() {"

		_, err := contentMinify(content)
		if err == nil {
			t.Error("Expected error for invalid JavaScript")
		}

		if !strings.Contains(err.Error(), "failed to minify content") {
			t.Errorf("Expected 'failed to minify content' error, got: %v", err)
		}
	})
}

// Tests for fileContentMinify function

func TestFileContentMinify(t *testing.T) {
	t.Run("ValidCSS", func(t *testing.T) {
		content := []byte(`
			body {
				margin: 0;
				padding: 0;
			}
		`)

		minified, err := fileContentMinify("css", content)
		if err != nil {
			t.Fatalf("Failed to minify CSS: %v", err)
		}

		if len(minified) >= len(content) {
			t.Error(errMsgMinifiedSmaller)
		}

		minifiedStr := string(minified)
		if !strings.Contains(minifiedStr, "body{margin:0;padding:0}") {
			t.Error("Expected minified CSS to be compressed")
		}
	})

	t.Run("ValidJavaScript", func(t *testing.T) {
		content := []byte(`
			function hello() {
				console.log("Hello");
			}
		`)

		minified, err := fileContentMinify("js", content)
		if err != nil {
			t.Fatalf("Failed to minify JavaScript: %v", err)
		}

		if len(minified) >= len(content) {
			t.Error(errMsgMinifiedSmaller)
		}
	})

	t.Run("UnsupportedFileType", func(t *testing.T) {
		content := []byte("test content")

		_, err := fileContentMinify("txt", content)
		if err == nil {
			t.Error("Expected error for unsupported file type")
		}

		if !strings.Contains(err.Error(), "unsupported file type") {
			t.Errorf("Expected 'unsupported file type' error, got: %v", err)
		}
	})
}

// Tests for generateMinifiedFilename function

func TestGenerateMinifiedFilename(t *testing.T) {
	content := []byte("test content")

	tests := []struct {
		name          string
		inputPath     string
		fileType      string
		expectedName  string
		expectedExt   string
		expectedParts int
	}{
		{
			name:          "JavaScriptFile",
			inputPath:     "/path/to/main.js",
			fileType:      "js",
			expectedName:  "main",
			expectedExt:   minJSExtension,
			expectedParts: 4,
		},
		{
			name:          "CSSFile",
			inputPath:     "/path/to/style.css",
			fileType:      "css",
			expectedName:  "style",
			expectedExt:   ".css",
			expectedParts: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := generateMinifiedFilename(test.inputPath, test.fileType, content)

			assertFilenameFormat(t, filename, test.expectedName, test.expectedExt, test.expectedParts)

			parts := strings.Split(filename, ".")
			assertHashProperties(t, parts[1], 8)
		})
	}
}

// Tests for ProcessBundles function

func TestProcessBundlesValidBundles(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source files
	createTempFile(t, tempDir, testFile1Path, sampleJS)
	createTempFile(t, tempDir, testFile2Path, sampleJS2)

	// Create bundle config
	bundles := []Bundle{
		{
			Name:  "test",
			Files: []string{filepath.Join(tempDir, testJSPattern)},
		},
	}
	configFile := createBundleConfig(t, tempDir, bundles)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")

	config := Config{
		BundlesFile: configFile,
		OutputDir:   outputDir,
	}

	err := ProcessBundles(config)
	if err != nil {
		t.Fatalf(errFailedToProcessBundles, err)
	}

	// Check that output file was created
	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 output file, got %d", len(files))
	}

	// Check filename format
	filename := files[0].Name()
	if !strings.HasPrefix(filename, "test.") {
		t.Errorf("Expected filename to start with 'test.', got %s", filename)
	}

	if !strings.HasSuffix(filename, minJSExtension) {
		t.Errorf(errMsgFilenameEndWithMinJS, filename)
	}
}

func TestProcessBundlesInvalidConfigFile(t *testing.T) {
	config := Config{
		BundlesFile: "nonexistent.json",
		OutputDir:   "/tmp",
	}

	err := ProcessBundles(config)
	if err == nil {
		t.Fatal("Expected error for invalid config file")
	}

	if !strings.Contains(err.Error(), "failed to load bundle config") {
		t.Errorf("Expected 'failed to load bundle config' error, got: %v", err)
	}
}

func TestProcessBundlesNoFilesFound(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	bundles := []Bundle{
		{
			Name:  "test",
			Files: []string{filepath.Join(tempDir, "*.nonexistent")},
		},
	}
	configFile := createBundleConfig(t, tempDir, bundles)

	config := Config{
		BundlesFile: configFile,
		OutputDir:   filepath.Join(tempDir, "output"),
	}

	err := ProcessBundles(config)
	if err == nil {
		t.Fatal("Expected error for no files found")
	}

	if !strings.Contains(err.Error(), "failed to process bundle") {
		t.Errorf("Expected 'failed to process bundle' error, got: %v", err)
	}
}

// Tests for GetBundleHash function

func TestGetBundleHash(t *testing.T) {
	t.Run("ValidBundle", func(t *testing.T) {
		configFile, _ := createTestBundleEnvironment(t, "test", []string{testFile1Path, testFile2Path})

		hash, err := GetBundleHash("test", configFile)
		if err != nil {
			t.Fatalf("Failed to get bundle hash: %v", err)
		}

		assertHashProperties(t, hash, 8)
		assertHashConsistency(t, GetBundleHash, configFile, "test")
	})

	t.Run("BundleNotFound", func(t *testing.T) {
		tempDir := setupTestEnvironment(t)

		bundles := []Bundle{
			{Name: "other", Files: []string{"*.js"}},
		}
		configFile := createBundleConfig(t, tempDir, bundles)

		_, err := GetBundleHash("nonexistent", configFile)
		if err == nil {
			t.Error(errMsgNonExistentBundle)
		}

		if !strings.Contains(err.Error(), "bundle nonexistent not found") {
			t.Errorf("Expected 'bundle nonexistent not found' error, got: %v", err)
		}
	})
}

// Tests for GetBundleFilename function

func TestGetBundleFilename(t *testing.T) {
	t.Run("ValidBundle", func(t *testing.T) {
		tempDir := createTempDir(t)
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Logf(errFailedToRemoveTempDir, err)
			}
		}()

		// Create source files
		createTempFile(t, tempDir, testFile1Path, sampleJS)

		// Create bundle config
		bundles := []Bundle{
			{
				Name:  "test",
				Files: []string{filepath.Join(tempDir, testJSPattern)},
			},
		}
		configFile := createBundleConfig(t, tempDir, bundles)

		filename, err := GetBundleFilename("test", configFile)
		if err != nil {
			t.Fatalf("Failed to get bundle filename: %v", err)
		}

		if !strings.HasPrefix(filename, "test.") {
			t.Errorf("Expected filename to start with 'test.', got %s", filename)
		}

		if !strings.HasSuffix(filename, minJSExtension) {
			t.Errorf(errMsgFilenameEndWithMinJS, filename)
		}

		// Check format: name.hash.min.js
		parts := strings.Split(filename, ".")
		if len(parts) != 4 {
			t.Errorf("Expected filename to have 4 parts, got %d: %s", len(parts), filename)
		}

		if len(parts[1]) != 8 {
			t.Errorf("Expected hash to be 8 characters, got %d", len(parts[1]))
		}
	})
}

// Tests for BundleExists function

func TestBundleExists(t *testing.T) {
	t.Run("BundleExists", func(t *testing.T) {
		configFile, outputDir := createTestBundleEnvironment(t, "test", []string{testFile1Path})

		// Initially bundle doesn't exist
		exists, err := BundleExists("test", configFile, outputDir)
		if err != nil {
			t.Fatalf("Failed to check bundle existence: %v", err)
		}

		assertBundleExistence(t, exists, false, "initially")

		// Process bundle
		config := Config{
			BundlesFile: configFile,
			OutputDir:   outputDir,
		}

		err = ProcessBundles(config)
		if err != nil {
			t.Fatalf(errFailedToProcessBundles, err)
		}

		// Now bundle should exist
		exists, err = BundleExists("test", configFile, outputDir)
		if err != nil {
			t.Fatalf("Failed to check bundle existence after processing: %v", err)
		}

		assertBundleExistence(t, exists, true, "after processing")
	})

	t.Run("BundleNotFound", func(t *testing.T) {
		tempDir := setupTestEnvironment(t)

		bundles := []Bundle{
			{Name: "other", Files: []string{"*.js"}},
		}
		configFile := createBundleConfig(t, tempDir, bundles)

		_, err := BundleExists("nonexistent", configFile, tempDir)
		if err == nil {
			t.Error(errMsgNonExistentBundle)
		}
	})
}

// Tests for CleanOldBundles function

func TestCleanOldBundles(t *testing.T) {
	t.Run("CleanOldVersions", func(t *testing.T) {
		tempDir := createTempDir(t)
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Logf(errFailedToRemoveTempDir, err)
			}
		}()

		// Create source files
		createTempFile(t, tempDir, testFile1Path, sampleJS)

		// Create bundle config
		bundles := []Bundle{
			{
				Name:  "test",
				Files: []string{filepath.Join(tempDir, testJSPattern)},
			},
		}
		configFile := createBundleConfig(t, tempDir, bundles)

		outputDir := filepath.Join(tempDir, "output")

		// Create some old bundle files
		createTempFile(t, outputDir, "test.old1hash.min.js", "old content 1")
		createTempFile(t, outputDir, "test.old2hash.min.js", "old content 2")

		// Process current bundle
		config := Config{
			BundlesFile: configFile,
			OutputDir:   outputDir,
		}

		err := ProcessBundles(config)
		if err != nil {
			t.Fatalf(errFailedToProcessBundles, err)
		}

		// Clean old bundles
		err = CleanOldBundles("test", configFile, outputDir)
		if err != nil {
			t.Fatalf("Failed to clean old bundles: %v", err)
		}

		// Check that only current bundle remains
		files, err := os.ReadDir(outputDir)
		if err != nil {
			t.Fatalf("Failed to read output directory: %v", err)
		}

		if len(files) != 1 {
			t.Errorf("Expected 1 file after cleanup, got %d", len(files))
		}

		// Check that remaining file is the current bundle
		filename := files[0].Name()
		if !strings.HasPrefix(filename, "test.") || !strings.HasSuffix(filename, minJSExtension) {
			t.Errorf("Expected current bundle file, got %s", filename)
		}
	})
}

// Tests for AndVersionFile function

func TestAndVersionFileValidCSSFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source CSS file
	inputFile := createTempFile(t, tempDir, testCSSPath, sampleCSS)
	outputDir := filepath.Join(tempDir, "output")

	filename, err := AndVersionFile(inputFile, outputDir, "css")
	if err != nil {
		t.Fatalf(errFailedToVersionCSSFile, err)
	}

	if !strings.HasPrefix(filename, "main.") {
		t.Errorf(errMsgFilenameStartWithMain, filename)
	}

	if !strings.HasSuffix(filename, ".css") {
		t.Errorf(errMsgFilenameEndWithCSS, filename)
	}

	// Check that file was created
	outputPath := filepath.Join(outputDir, filename)
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf(errMsgOutputFileExist, err)
	}

	// Check that content is minified
	content, err := os.ReadFile(outputPath) // #nosec G304 -- outputPath is in test temp directory
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) >= len(sampleCSS) {
		t.Error(errMsgMinifiedSmaller)
	}
}

func TestAndVersionFileValidJSFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source JS file
	inputFile := createTempFile(t, tempDir, "src/main.js", sampleJS)
	outputDir := filepath.Join(tempDir, "output")

	filename, err := AndVersionFile(inputFile, outputDir, "js")
	if err != nil {
		t.Fatalf("Failed to version JS file: %v", err)
	}

	if !strings.HasPrefix(filename, "main.") {
		t.Errorf(errMsgFilenameStartWithMain, filename)
	}

	if !strings.HasSuffix(filename, minJSExtension) {
		t.Errorf(errMsgFilenameEndWithMinJS, filename)
	}

	// Check that file was created
	outputPath := filepath.Join(outputDir, filename)
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf(errMsgOutputFileExist, err)
	}
}

func TestAndVersionFileUnsupportedFileType(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	inputFile := createTempFile(t, tempDir, "src/main.txt", "text content")
	outputDir := filepath.Join(tempDir, "output")

	_, err := AndVersionFile(inputFile, outputDir, "txt")
	if err == nil {
		t.Fatal("Expected error for unsupported file type")
	}

	if !strings.Contains(err.Error(), "unsupported file type") {
		t.Errorf("Expected 'unsupported file type' error, got: %v", err)
	}
}

func TestAndVersionFileNonExistentFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	inputFile := filepath.Join(tempDir, "nonexistent.css")
	outputDir := filepath.Join(tempDir, "output")

	_, err := AndVersionFile(inputFile, outputDir, "css")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read css file") {
		t.Errorf("Expected 'failed to read css file' error, got: %v", err)
	}
}

func TestAndVersionFileExistingFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source CSS file
	inputFile := createTempFile(t, tempDir, testCSSPath, sampleCSS)
	outputDir := filepath.Join(tempDir, "output")

	// First call should create the file
	filename1, err := AndVersionFile(inputFile, outputDir, "css")
	if err != nil {
		t.Fatalf(errFailedToVersionCSSFile, err)
	}

	// Second call should return existing filename without recreating
	filename2, err := AndVersionFile(inputFile, outputDir, "css")
	if err != nil {
		t.Fatalf("Failed to version CSS file again: %v", err)
	}

	if filename1 != filename2 {
		t.Errorf("Expected same filename for same content, got %s and %s", filename1, filename2)
	}
}

// Tests for AndVersionCSS function

func TestAndVersionCSSValidCSSFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source CSS file
	inputFile := createTempFile(t, tempDir, "src/styles.css", sampleCSS)
	outputDir := filepath.Join(tempDir, "output")

	filename, err := AndVersionCSS(inputFile, outputDir)
	if err != nil {
		t.Fatalf(errFailedToVersionCSSFile, err)
	}

	if !strings.HasPrefix(filename, "styles.") {
		t.Errorf("Expected filename to start with 'styles.', got %s", filename)
	}

	if !strings.HasSuffix(filename, ".css") {
		t.Errorf(errMsgFilenameEndWithCSS, filename)
	}

	// Check that file was created
	outputPath := filepath.Join(outputDir, filename)
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf(errMsgOutputFileExist, err)
	}

	// Check that content is minified
	content, err := os.ReadFile(outputPath) // #nosec G304 -- outputPath is in test temp directory
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) >= len(sampleCSS) {
		t.Error(errMsgMinifiedSmaller)
	}
}

func TestAndVersionCSSNonExistentFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	inputFile := filepath.Join(tempDir, "nonexistent.css")
	outputDir := filepath.Join(tempDir, "output")

	_, err := AndVersionCSS(inputFile, outputDir)
	if err == nil {
		t.Error(errMsgNonExistentFile)
	} else if !strings.Contains(err.Error(), "failed to read css file") {
		t.Errorf("Expected 'failed to read css file' error, got: %v", err)
	}
}

// Benchmark tests

func BenchmarkProcessBundles(b *testing.B) {
	tempDir := createTempDir(b)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			b.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source files
	createTempFile(b, tempDir, testFile1Path, sampleJS)
	createTempFile(b, tempDir, testFile2Path, sampleJS2)

	// Create bundle config
	bundles := []Bundle{
		{
			Name:  "test",
			Files: []string{filepath.Join(tempDir, testJSPattern)},
		},
	}
	configFile := createBundleConfig(b, tempDir, bundles)

	config := Config{
		BundlesFile: configFile,
		OutputDir:   filepath.Join(tempDir, "output"),
	}

	b.ResetTimer()
	for range b.N {
		err := ProcessBundles(config)
		if err != nil {
			b.Fatalf("Failed to process bundles: %v", err)
		}
	}
}

func BenchmarkGetBundleHash(b *testing.B) {
	tempDir := createTempDir(b)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			b.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	// Create source files
	createTempFile(b, tempDir, testFile1Path, sampleJS)
	createTempFile(b, tempDir, testFile2Path, sampleJS2)

	// Create bundle config
	bundles := []Bundle{
		{
			Name:  "test",
			Files: []string{filepath.Join(tempDir, testJSPattern)},
		},
	}
	configFile := createBundleConfig(b, tempDir, bundles)

	b.ResetTimer()
	for range b.N {
		_, err := GetBundleHash("test", configFile)
		if err != nil {
			b.Fatalf("Failed to get bundle hash: %v", err)
		}
	}
}

func BenchmarkAndVersionFile(b *testing.B) {
	tempDir := createTempDir(b)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			b.Logf(errFailedToRemoveTempDir, err)
		}
	}()

	inputFile := createTempFile(b, tempDir, testCSSPath, sampleCSS)
	outputDir := filepath.Join(tempDir, "output")

	b.ResetTimer()
	for i := range b.N {
		// Use different output directory for each iteration to avoid cache
		currentOutputDir := filepath.Join(outputDir, fmt.Sprintf("iter%d", i))
		_, err := AndVersionFile(inputFile, currentOutputDir, "css")
		if err != nil {
			b.Fatalf("Failed to version file: %v", err)
		}
	}
}
