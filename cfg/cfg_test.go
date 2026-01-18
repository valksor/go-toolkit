package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create test YAML file
	content := `
key1: value1
key2: value2
nested:
  item1: "test"
`
	err := os.WriteFile(configPath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load the config
	var result map[string]interface{}
	err = LoadYAML(configPath, &result)
	if err != nil {
		t.Fatalf("LoadYAML() error = %v", err)
	}

	if result["key1"] != "value1" {
		t.Errorf("key1 = %v, want 'value1'", result["key1"])
	}
}

func TestSaveYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save config
	config := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	err := SaveYAML(configPath, &config)
	if err != nil {
		t.Fatalf("SaveYAML() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveYAML() did not create file")
	}

	// Load and verify
	var loaded map[string]interface{}
	err = LoadYAML(configPath, &loaded)
	if err != nil {
		t.Fatalf("Failed to load saved file: %v", err)
	}

	if loaded["key1"] != "value1" {
		t.Errorf("Loaded key1 = %v, want 'value1'", loaded["key1"])
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	nonExistingFile := filepath.Join(tmpDir, "nonexistent.txt")

	// Create existing file
	err := os.WriteFile(existingFile, []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !FileExists(existingFile) {
		t.Error("FileExists() returned false for existing file")
	}

	if FileExists(nonExistingFile) {
		t.Error("FileExists() returned true for non-existing file")
	}
}

func TestFindConfigInParents(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "level1", "level2", "level3")
	err := os.MkdirAll(nestedDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create config file at level1
	configPath := filepath.Join(tmpDir, "level1", "config.yaml")
	err = os.WriteFile(configPath, []byte("test: true"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Should find config from level3
	found := FindConfigInParents(nestedDir, "config.yaml")
	if found != configPath {
		t.Errorf("FindConfigInParents() = %v, want %v", found, configPath)
	}

	// Should not find non-existent file
	found = FindConfigInParents(nestedDir, "nonexistent.yaml")
	if found != "" {
		t.Errorf("FindConfigInParents() = %v, want empty string", found)
	}
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]string
		override map[string]string
		mode     MergeMode
		expected map[string]string
	}{
		{
			name:     "overlay mode",
			base:     map[string]string{"key1": "value1", "key2": "value2"},
			override: map[string]string{"key2": "newvalue", "key3": "value3"},
			mode:     MergeModeOverlay,
			expected: map[string]string{"key1": "value1", "key2": "newvalue", "key3": "value3"},
		},
		{
			name:     "replace mode",
			base:     map[string]string{"key1": "value1", "key2": "value2"},
			override: map[string]string{"key3": "value3"},
			mode:     MergeModeReplace,
			expected: map[string]string{"key3": "value3"},
		},
		{
			name:     "nil base",
			base:     nil,
			override: map[string]string{"key1": "value1"},
			mode:     MergeModeOverlay,
			expected: map[string]string{"key1": "value1"},
		},
		{
			name:     "nil override",
			base:     map[string]string{"key1": "value1"},
			override: nil,
			mode:     MergeModeOverlay,
			expected: map[string]string{"key1": "value1"},
		},
		{
			name:     "both nil",
			base:     nil,
			override: nil,
			mode:     MergeModeOverlay,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMaps(tt.base, tt.override, tt.mode)

			if len(result) != len(tt.expected) {
				t.Errorf("MergeMaps() length = %d, want %d", len(result), len(tt.expected))
			}

			for k, expectedVal := range tt.expected {
				if result[k] != expectedVal {
					t.Errorf("MergeMaps()[%q] = %q, want %q", k, result[k], expectedVal)
				}
			}
		})
	}
}

func TestCloneMap(t *testing.T) {
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	cloned := CloneMap(original)

	// Verify values match
	for k, v := range original {
		if cloned[k] != v {
			t.Errorf("CloneMap()[%q] = %q, want %q", k, cloned[k], v)
		}
	}

	// Modify clone shouldn't affect original
	cloned["key1"] = "modified"
	if original["key1"] == "modified" {
		t.Error("Modifying clone affected original")
	}

	// Test nil map
	cloned = CloneMap(nil)
	if cloned != nil {
		t.Error("CloneMap(nil) should return nil, not empty map")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()

	newDir := filepath.Join(tmpDir, "new", "nested", "dir")
	err := EnsureDir(newDir)
	if err != nil {
		t.Fatalf("EnsureDir() error = %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(newDir)
	if err != nil {
		t.Error("EnsureDir() did not create directory")
	}

	if !info.IsDir() {
		t.Error("EnsureDir() created a file, not a directory")
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "tilde alone",
			input:    "~",
			contains: home,
		},
		{
			name:     "tilde with path",
			input:    "~/config",
			contains: filepath.Join(home, "config"),
		},
		{
			name:     "absolute path",
			input:    "/absolute/path",
			contains: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			contains: "relative/path",
		},
		{
			name:     "empty string",
			input:    "",
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.contains {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.contains)
			}
		})
	}
}

func TestLoadLayered(t *testing.T) {
	tmpDir := t.TempDir()

	// Create global config
	globalPath := filepath.Join(tmpDir, "global.yaml")
	globalContent := `key1: global1
key2: global2`
	err := os.WriteFile(globalPath, []byte(globalContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Create local config
	localPath := filepath.Join(tmpDir, "local.yaml")
	localContent := `key2: local2
key3: local3`
	err = os.WriteFile(localPath, []byte(localContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create local config: %v", err)
	}

	// Define merge function
	merge := func(layers []interface{}) interface{} {
		// Simple merge: combine all maps
		result := make(map[string]interface{})
		for _, layer := range layers {
			if m, ok := layer.(map[string]interface{}); ok {
				for k, v := range m {
					result[k] = v
				}
			}
		}

		return result
	}

	// Load layers
	layers := []Layer{
		{
			Path:       globalPath,
			Type:       LayerTypeYAML,
			Precedence: PrecedenceGlobal,
			Optional:   false,
		},
		{
			Path:       localPath,
			Type:       LayerTypeYAML,
			Precedence: PrecedenceLocal,
			Optional:   false,
		},
	}

	result, err := LoadLayered(layers, merge)
	if err != nil {
		t.Fatalf("LoadLayered() error = %v", err)
	}

	if result.Local == nil {
		t.Error("LoadLayered() did not merge layers")
	}
}

func TestSaveYAMLCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")

	config := map[string]string{"key": "value"}
	err := SaveYAML(configPath, &config)
	if err != nil {
		t.Fatalf("SaveYAML() error = %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(filepath.Dir(configPath)); os.IsNotExist(err) {
		t.Error("SaveYAML() did not create parent directory")
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveYAML() did not create file")
	}
}

// mockLoader implements the Loader interface for testing.
type mockLoader struct {
	loadFunc func(path string, v interface{}) error
}

func (m *mockLoader) Load(path string, v interface{}) error {
	return m.loadFunc(path, v)
}

func TestCustomLoader(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom.config")

	// Create file
	err := os.WriteFile(configPath, []byte("custom content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create mock loader
	loader := &mockLoader{
		loadFunc: func(path string, v interface{}) error {
			// Just set a value for testing
			if m, ok := v.(*map[string]interface{}); ok {
				if *m == nil {
					*m = make(map[string]interface{})
				}
				(*m)["loaded"] = true
			}

			return nil
		},
	}

	// Load with custom loader
	var result map[string]interface{}
	err = loader.Load(configPath, &result)
	if err != nil {
		t.Fatalf("Custom loader error = %v", err)
	}

	if loaded, ok := result["loaded"].(bool); !ok || !loaded {
		t.Error("Custom loader was not called")
	}
}
