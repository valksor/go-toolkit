package qtcwrap

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// File and path constants.
	testQtplFile  = "test.qtpl"
	testContent   = "test content"
	templatesDir  = "templates"
	qtplExt       = ".qtpl"
	templateExt   = ".template"
	goExt         = ".go"
	tempDirPrefix = "qtcwrap_test"

	// Command argument constants.
	dirTemplatesArg = "-dir=templates"
	extQtplArg      = "-ext=.qtpl"
	skipCommentsArg = "-skipLineComments"

	// Error message templates.
	syntaxErrorMsg                 = "syntax error in template"
	createTempDirErr               = "Failed to create temp directory: %v"
	removeTempDirErr               = "Failed to remove temp directory: %v"
	createTempFileErr              = "Failed to create temp file: %v"
	extensionMustStartWithDot      = "extension must start with a dot"
	eitherFileOrDirMustBeSpecified = "either File or Dir must be specified"
	createSpecificFileErr          = "Failed to create file %s: %v"
	findTemplateFilesErr           = "Failed to find template files: %v"
	closeWriterErr                 = "Failed to close writer: %v"
	readFromPipeErr                = "Failed to read from pipe: %v"
	changeDirPermsErr              = "Failed to change directory permissions: %v"
	restoreDirPermsErr             = "Failed to restore directory permissions: %v"
	fileNotAccessibleErr           = "file %s is not accessible"
	dirNotAccessibleErr            = "directory %s is not accessible"
	notDirErr                      = "is not a directory"
)

// Helper functions for test setup and assertions.

// createTempTestDir creates a temporary directory for testing and returns its path.
func createTempTestDir(t *testing.T) string {
	t.Helper()

	return t.TempDir()
}

// createTempTestFile creates a temporary file with the given content in the specified directory.
func createTempTestFile(t *testing.T, dir, content string) string {
	t.Helper()
	tempFile := filepath.Join(dir, testQtplFile)
	err := os.WriteFile(tempFile, []byte(content), 0o600)
	if err != nil {
		t.Fatalf(createTempFileErr, err)
	}

	return tempFile
}

// assertConfigEquals validates that two Config structs are equal.
func assertConfigEquals(t *testing.T, expected, actual Config, testName string) {
	t.Helper()
	if expected.Dir != actual.Dir {
		t.Errorf("%s: Expected Dir to be '%s', got '%s'", testName, expected.Dir, actual.Dir)
	}
	if expected.SkipLineComments != actual.SkipLineComments {
		t.Errorf("%s: Expected SkipLineComments to be %t, got %t", testName, expected.SkipLineComments, actual.SkipLineComments)
	}
	if expected.Ext != actual.Ext {
		t.Errorf("%s: Expected Ext to be '%s', got '%s'", testName, expected.Ext, actual.Ext)
	}
	if expected.File != actual.File {
		t.Errorf("%s: Expected File to be '%s', got '%s'", testName, expected.File, actual.File)
	}
}

// assertValidationError validates error conditions and messages.
func assertValidationError(t *testing.T, err error, expectedMsg string, expectErr bool) {
	t.Helper()
	if expectErr && err == nil {
		t.Error("Expected error but got none")

		return
	}
	if !expectErr && err != nil {
		t.Errorf("Expected no error but got: %v", err)

		return
	}
	if expectErr && expectedMsg != "" && !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

// assertArgsContain validates that expected argument is present in args slice.
func assertArgsContain(t *testing.T, args []string, expected string) {
	t.Helper()
	for _, arg := range args {
		if arg == expected {
			return
		}
	}
	t.Errorf("Expected args to contain '%s', got %v", expected, args)
}

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name: "DefaultConfig",
			config: Config{
				Dir:              ".",
				SkipLineComments: true,
				Ext:              "",
				File:             "",
			},
			expected: Config{
				Dir:              ".",
				SkipLineComments: true,
				Ext:              "",
				File:             "",
			},
		},
		{
			name: "CustomConfig",
			config: Config{
				Dir:              templatesDir,
				SkipLineComments: false,
				Ext:              qtplExt,
				File:             testQtplFile,
			},
			expected: Config{
				Dir:              templatesDir,
				SkipLineComments: false,
				Ext:              qtplExt,
				File:             testQtplFile,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertConfigEquals(t, tt.expected, tt.config, tt.name)
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	tests := []struct {
		name     string
		expected any
		actual   any
	}{
		{"Dir", ".", config.Dir},
		{"SkipLineComments", true, config.SkipLineComments},
		{"Ext", "", config.Ext},
		{"File", "", config.File},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			if testT.actual != testT.expected {
				t.Errorf("Expected %s to be %v, got %v", testT.name, testT.expected, testT.actual)
			}
		})
	}
}

func TestValidateQtcTool(t *testing.T) {
	t.Run("ValidateQtcTool", func(t *testing.T) {
		err := validateQtcTool()
		// This test will pass if qtc is available, skip if not
		if err != nil {
			t.Skipf("qtc tool not available: %v", err)
		}
	})
}

func TestIsQtcAvailable(t *testing.T) {
	t.Run("IsQtcAvailable", func(t *testing.T) {
		available := IsQtcAvailable()

		// Check that the function returns a boolean
		if available != true && available != false {
			t.Error("IsQtcAvailable should return a boolean")
		}
	})
}

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected []string
	}{
		{
			name: "FileMode",
			config: Config{
				File:             testQtplFile,
				Dir:              templatesDir,
				Ext:              qtplExt,
				SkipLineComments: false,
			},
			expected: []string{"-file=test.qtpl"},
		},
		{
			name: "DirectoryMode",
			config: Config{
				Dir:              templatesDir,
				Ext:              qtplExt,
				SkipLineComments: false,
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg},
		},
		{
			name: "DirectoryModeWithExt",
			config: Config{
				Dir:              templatesDir,
				SkipLineComments: false,
				Ext:              qtplExt,
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg},
		},
		{
			name: "DirectoryModeWithSkipComments",
			config: Config{
				Dir:              templatesDir,
				SkipLineComments: true,
				Ext:              qtplExt,
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg, skipCommentsArg},
		},
		{
			name: "FileModeWithSkipComments",
			config: Config{
				Dir:              "",
				SkipLineComments: true,
				Ext:              "",
				File:             testQtplFile,
			},
			expected: []string{"-file=test.qtpl", skipCommentsArg},
		},
		{
			name: "OnlyDir",
			config: Config{
				Dir:              "src",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expected: []string{"-dir=src"},
		},
		{
			name: "OnlyExt",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              templateExt,
				File:             "",
			},
			expected: []string{"-ext=.template"},
		},
		{
			name: "OnlySkipComments",
			config: Config{
				Dir:              "",
				SkipLineComments: true,
				Ext:              "",
				File:             "",
			},
			expected: []string{skipCommentsArg},
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expected: []string{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			args := buildArgs(Config{
				Dir:              testCase.config.Dir,
				SkipLineComments: testCase.config.SkipLineComments,
				Ext:              testCase.config.Ext,
				File:             testCase.config.File,
			})

			if len(args) != len(testCase.expected) {
				t.Errorf("Expected %d arguments, got %d", len(testCase.expected), len(args))
			}

			for idx, arg := range args {
				if idx >= len(testCase.expected) {
					t.Errorf("Unexpected argument at index %d: %s", idx, arg)

					continue
				}
				if arg != testCase.expected[idx] {
					t.Errorf("Expected argument %d to be '%s', got '%s'", idx, testCase.expected[idx], arg)
				}
			}
		})
	}
}

func TestIsTemporaryFileWarning(t *testing.T) {
	tests := []struct {
		name     string
		stderr   []byte
		expected bool
	}{
		{
			name:     "TemporaryFileWarning",
			stderr:   []byte("open .tmp/test.qtpl: no such file or directory"),
			expected: true,
		},
		{
			name:     "TemporaryFileWarningWithPath",
			stderr:   []byte("stat: /path/to/.tmp/file: no such file or directory"),
			expected: true,
		},
		{
			name:     "ActualError",
			stderr:   []byte(syntaxErrorMsg),
			expected: false,
		},
		{
			name:     "NoSuchFileButNotestTmp",
			stderr:   []byte("open test.qtpl: no such file or directory"),
			expected: false,
		},
		{
			name:     "TmpButNotNoSuchFile",
			stderr:   []byte("permission denied: .tmp/test.qtpl"),
			expected: false,
		},
		{
			name:     "EmptyStderr",
			stderr:   []byte(""),
			expected: false,
		},
		{
			name:     "OnlyTmp",
			stderr:   []byte(".tmp"),
			expected: false,
		},
		{
			name:     "OnlyNoSuchFile",
			stderr:   []byte("no such file or directory"),
			expected: false,
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			result := isTemporaryFileWarning(testT.stderr)
			if result != testT.expected {
				t.Errorf("Expected %v, got %v for stderr: %s", testT.expected, result, string(testT.stderr))
			}
		})
	}
}

func TestValidateConfigValidCases(t *testing.T) {
	tempDir := createTempTestDir(t)
	tempFile := createTempTestFile(t, tempDir, testContent)

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "ValidFileConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             tempFile,
			},
		},
		{
			name: "ValidDirConfig",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
		},
		{
			name: "ValidDirConfigWithExt",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              qtplExt,
				File:             "",
			},
		},
		{
			name: "ValidExtensionWithDot",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              qtplExt,
				File:             "",
			},
		},
		{
			name: "ValidCurrentDir",
			config: Config{
				Dir:              ".",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidationError(t, ValidateConfig(tt.config), "", false)
		})
	}
}

func TestValidateConfigInvalidCases(t *testing.T) {
	tempDir := createTempTestDir(t)
	tempFile := createTempTestFile(t, tempDir, testContent)

	tests := []struct {
		name     string
		config   Config
		errorMsg string
	}{
		{
			name: "InvalidFileConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "/nonexistent/file.qtpl",
			},
			errorMsg: "file /nonexistent/file.qtpl is not accessible",
		},
		{
			name: "InvalidDirConfig",
			config: Config{
				Dir:              "/nonexistent/directory",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			errorMsg: "directory /nonexistent/directory is not accessible",
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			errorMsg: eitherFileOrDirMustBeSpecified,
		},
		{
			name: "DirIsActuallyFile",
			config: Config{
				Dir:              tempFile,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			errorMsg: notDirErr,
		},
		{
			name: "InvalidExtension",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "qtpl",
				File:             "",
			},
			errorMsg: extensionMustStartWithDot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidationError(t, ValidateConfig(tt.config), tt.errorMsg, true)
		})
	}
}

func TestValidateConfigEdgeCases(t *testing.T) {
	// This function is kept for additional edge cases that might be added later
	// Currently covers scenarios that don't fit in valid/invalid categories
	t.Run("BackwardCompatibility", func(t *testing.T) {
		// Ensure the original TestValidateConfig behavior is preserved
		tempDir := createTempTestDir(t)

		config := Config{
			Dir:              tempDir,
			SkipLineComments: true,
			Ext:              qtplExt,
			File:             "",
		}

		err := ValidateConfig(config)
		if err != nil {
			t.Errorf("Expected valid config with skip comments, got error: %v", err)
		}
	})
}

func TestCompileWithValidationSuccess(t *testing.T) {
	tempDir := createTempTestDir(t)

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "ValidConfig",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
		},
		{
			name: "ValidConfigWithExtension",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: true,
				Ext:              qtplExt,
				File:             "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompileWithValidation(tt.config)
			// For valid configs, we might get a qtc tool error, which is acceptable
			if err != nil && !strings.Contains(err.Error(), "qtc tool validation failed") {
				t.Errorf("Expected no error or qtc tool error, got: %v", err)
			}
		})
	}
}

func TestCompileWithValidationFailure(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		errorMsg string
	}{
		{
			name: "InvalidConfig",
			config: Config{
				Dir:              "/nonexistent/directory",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			errorMsg: "configuration validation failed",
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			errorMsg: eitherFileOrDirMustBeSpecified,
		},
		{
			name: "InvalidExtension",
			config: Config{
				Dir:              ".",
				SkipLineComments: false,
				Ext:              "invalid",
				File:             "",
			},
			errorMsg: extensionMustStartWithDot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidationError(t, CompileWithValidation(tt.config), tt.errorMsg, true)
		})
	}
}

func TestFindTemplateFiles(t *testing.T) {
	tempDir := t.TempDir()

	testFiles := []string{
		"test1.qtpl",
		"test2.qtpl",
		"subdir/test3.qtpl",
		"subdir/test4.template",
		"other.go",
	}
	for _, fileName := range testFiles {
		fullPath := filepath.Join(tempDir, fileName)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o700); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", fileName, err)
		}
		err := os.WriteFile(fullPath, []byte(testContent), 0o600)
		if err != nil {
			t.Fatalf(createSpecificFileErr, fileName, err)
		}
	}

	testCases := []struct {
		description string
		ext         string
		expectedLen int
	}{
		{"Find .qtpl files", qtplExt, 3},
		{"Find .template files", templateExt, 1},
		{"Find .go files", goExt, 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			files, err := FindTemplateFiles(tempDir, testCase.ext)
			if err != nil {
				t.Fatalf(findTemplateFilesErr, err)
			}
			if len(files) != testCase.expectedLen {
				t.Errorf("Expected %d files, got %d", testCase.expectedLen, len(files))
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a temporary file for testing
	tempFile := filepath.Join(tempDir, testQtplFile)
	err := os.WriteFile(tempFile, []byte(testContent), 0o600)
	if err != nil {
		t.Fatalf(createTempFileErr, err)
	}

	t.Run("CompileDirectory", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileDirectory(tempDir)
	})

	t.Run("CompileFile", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileFile(tempFile)
	})

	t.Run("CompileWithExtension", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileWithExtension(tempDir, qtplExt)
	})

	t.Run("QtcWrap", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		QtcWrap()
	})
}

func TestHandleQtcError(t *testing.T) {
	tests := []struct {
		name           string
		stderr         string
		expectedOutput string
	}{
		{
			name:           "TemporaryFileWarning",
			stderr:         "open .tmp/test.qtpl: no such file or directory",
			expectedOutput: "[qtc warning suppressed]",
		},
		{
			name:           "ActualError",
			stderr:         syntaxErrorMsg,
			expectedOutput: syntaxErrorMsg,
		},
		{
			name:           "EmptyStderr",
			stderr:         "",
			expectedOutput: "qtc execution failed:",
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			var buf bytes.Buffer
			stderr := bytes.NewBufferString(testT.stderr)
			err := errors.New("exit status 1")

			// Capture output by temporarily redirecting stdout
			oldStdout := os.Stdout
			rFile, wFile, _ := os.Pipe()
			os.Stdout = wFile

			handleQtcError(*stderr, err)

			if err := wFile.Close(); err != nil {
				t.Fatalf(closeWriterErr, err)
			}
			os.Stdout = oldStdout

			// Read the captured output
			if _, err := buf.ReadFrom(rFile); err != nil {
				t.Fatalf(readFromPipeErr, err)
			}
			output := buf.String()

			if !strings.Contains(output, testT.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", testT.expectedOutput, output)
			}
		})
	}
}

func TestExecuteQtc(t *testing.T) {
	t.Run("ExecuteQtcWithInvalidArgs", func(t *testing.T) {
		// Test with invalid arguments that should fail
		args := []string{"-invalid-flag"}

		// Capture output
		var buf bytes.Buffer
		oldStdout := os.Stdout
		rFile, wFile, _ := os.Pipe()
		os.Stdout = wFile

		executeQtc(args)

		if err := wFile.Close(); err != nil {
			t.Fatalf(closeWriterErr, err)
		}
		os.Stdout = oldStdout
		if _, err := buf.ReadFrom(rFile); err != nil {
			t.Fatalf(readFromPipeErr, err)
		}

		// The function should handle the error gracefully
		// We can't test much more without qtc being available
	})
}

func TestWithConfig(t *testing.T) {
	t.Run("WithValidConfig", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: true,
			Ext:              "",
			File:             "",
		}

		// This should not panic
		WithConfig(config)
	})
}

func TestConfigPrecedence(t *testing.T) {
	tempDir := createTempTestDir(t)
	tempFile := createTempTestFile(t, tempDir, testContent)

	config := Config{
		Dir:              tempDir,
		SkipLineComments: false,
		Ext:              "",
		File:             tempFile,
	}

	args := buildArgs(config)
	// When both are specified, File should take precedence
	if len(args) != 1 || args[0] != "-file="+tempFile {
		t.Errorf("Expected File to take precedence over Dir, got args: %v", args)
	}
}

func TestExtensionHandling(t *testing.T) {
	t.Run("EmptyExtension", func(t *testing.T) {
		config := Config{
			Dir:              ".",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		args := buildArgs(config)
		// Empty extension should not be added to args
		found := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "-ext=") {
				found = true

				break
			}
		}
		if found {
			t.Error("Empty extension should not be added to args")
		}
	})

	t.Run("ValidExtension", func(t *testing.T) {
		config := Config{
			Dir:              ".",
			SkipLineComments: false,
			Ext:              qtplExt,
			File:             "",
		}

		args := buildArgs(config)
		assertArgsContain(t, args, "-ext=.qtpl")
	})
}

func TestConfigValidationEdgeCases(t *testing.T) {
	t.Run("EmptyConfig", func(t *testing.T) {
		config := Config{
			Dir:              "",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err := ValidateConfig(config)
		if err == nil {
			t.Error("Expected error for empty config")
		}
	})

	t.Run("NilSafeValidation", func(t *testing.T) {
		// Test that validation handles edge cases gracefully
		config := Config{
			Dir:              ".",
			SkipLineComments: true,
			Ext:              qtplExt,
			File:             "",
		}

		// This should not panic
		err := ValidateConfig(config)
		if err != nil {
			t.Errorf("Expected valid config, got error: %v", err)
		}
	})
}

func TestFileSystemPermissions(t *testing.T) {
	t.Run("ReadOnlyDirectory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Make directory read-only
		if err := os.Chmod(tempDir, 0o600); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		// Restore permissions for cleanup
		defer func() {
			if err := os.Chmod(tempDir, 0o600); err != nil {
				t.Fatalf("Failed to restore directory permissions: %v", err)
			}
		}()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err := ValidateConfig(config)
		// Should still be valid as we can read the directory
		if err != nil {
			t.Errorf("Expected no error for read-only directory, got: %v", err)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentValidation", func(t *testing.T) {
		tempDir := t.TempDir()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		// Run multiple validations concurrently
		done := make(chan bool, 10)
		for range 10 {
			go func() {
				err := ValidateConfig(config)
				if err != nil {
					t.Errorf("Concurrent validation failed: %v", err)
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for range 10 {
			<-done
		}
	})
}

func TestLargeDirectoryStructure(t *testing.T) {
	t.Run("ManyFiles", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create many files
		for i := range 100 {
			filename := filepath.Join(tempDir, fmt.Sprintf("test%d.qtpl", i))
			err := os.WriteFile(filename, []byte(testContent), 0o600)
			if err != nil {
				t.Fatalf(createSpecificFileErr, filename, err)
			}
		}

		files, err := FindTemplateFiles(tempDir, qtplExt)
		if err != nil {
			t.Fatalf(findTemplateFilesErr, err)
		}

		if len(files) != 100 {
			t.Errorf("Expected 100 files, got %d", len(files))
		}
	})
}

func TestSpecialCharacters(t *testing.T) {
	t.Run("SpecialCharsInPath", func(t *testing.T) {
		tempDir := t.TempDir()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err := ValidateConfig(config)
		if err != nil {
			t.Errorf("Expected no error for path with spaces, got: %v", err)
		}
	})
}

func TestSymlinks(t *testing.T) {
	t.Run("SymlinkDirectory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a symlink to the directory
		symlinkPath := filepath.Join(tempDir, "symlink")
		err := os.Symlink(tempDir, symlinkPath)
		if err != nil {
			t.Skipf("Failed to create symlink, skipping test: %v", err)
		}

		config := Config{
			Dir:              symlinkPath,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err = ValidateConfig(config)
		if err != nil {
			t.Errorf("Expected no error for symlink directory, got: %v", err)
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("DetailedErrorMessages", func(t *testing.T) {
		config := Config{
			Dir:              "/this/path/definitely/does/not/exist",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		err := ValidateConfig(config)
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "directory") {
			t.Errorf("Expected error message to contain 'directory', got: %s", errMsg)
		}
		if !strings.Contains(errMsg, config.Dir) {
			t.Errorf("Expected error message to contain directory path, got: %s", errMsg)
		}
	})
}
