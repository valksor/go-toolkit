// Package cfg provides configuration loading, saving, and merging utilities.
//
// This package provides generic utilities for working with configuration files
// in YAML and JSON formats, with support for multi-layered configuration merging.
//
// Basic usage:
//
//	var config MyConfig
//	err := cfg.LoadYAML("config.yaml", &config)
//	err = cfg.SaveYAML("config.yaml", &config)
package cfg

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/valksor/go-toolkit/paths"
	"gopkg.in/yaml.v3"
)

// MergeMode defines how maps are merged.
type MergeMode string

const (
	// MergeModeOverlay merges override on top of base (default).
	MergeModeOverlay MergeMode = "overlay"
	// MergeModeReplace replaces base with override.
	MergeModeReplace MergeMode = "replace"
)

// LoadYAML loads a YAML file into the provided interface.
// The target must be a pointer to a struct or map.
func LoadYAML(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

// SaveYAML saves an interface to a YAML file.
// The source should be a struct or map.
func SaveYAML(path string, v interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// LoadJSON loads a JSON file into the provided interface.
// Note: This uses yaml.v3 which can unmarshal JSON as well.
func LoadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

// SaveJSON saves an interface to a JSON file.
// Note: This uses yaml.v3 which can marshal to JSON as well.
func SaveJSON(path string, v interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// PathExists checks if a path exists (file or directory) using os.Stat.
// Returns false if the path doesn't exist or if there's an error.
// For checking specifically files (not directories), use paths.FileExists.
func PathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// FileExists is an alias for PathExists for backward compatibility.
//
// Deprecated: Use PathExists for any path, or paths.FileExists for files only.
func FileExists(path string) bool {
	return PathExists(path)
}

// FindConfigInParents searches for a config file in the current directory
// and parent directories. Returns the path if found, empty string otherwise.
func FindConfigInParents(startDir string, filename string) string {
	dir := startDir

	for {
		candidate := filepath.Join(dir, filename)
		if PathExists(candidate) {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			return ""
		}

		dir = parent
	}
}

// FindConfigDirs searches for config directories in current and parent directories.
// Returns a list of directory paths where config files were found.
func FindConfigDirs(startDir string, configFilenames []string) []string {
	var found []string
	seen := make(map[string]bool)

	dir := startDir
	for {
		for _, filename := range configFilenames {
			candidate := filepath.Join(dir, filename)
			if PathExists(candidate) && !seen[dir] {
				found = append(found, dir)
				seen[dir] = true

				break
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return found
}

// MergeMaps merges two string maps based on the merge mode.
// If mode is MergeModeOverlay, override values are merged on top of base.
// If mode is MergeModeReplace, base is replaced with override.
func MergeMaps(base, override map[string]string, mode MergeMode) map[string]string {
	if len(override) == 0 {
		return CloneMap(base)
	}

	if len(base) == 0 || mode == MergeModeReplace {
		return CloneMap(override)
	}

	// Overlay mode: merge override on top of base
	result := CloneMap(base)
	for k, v := range override {
		result[k] = v
	}

	return result
}

// CloneMap creates a shallow copy of a string map.
func CloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}

	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}

	return result
}

// EnsureDir ensures a directory exists, creating it if necessary.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// ExpandPath expands ~ to the user's home directory.
// Note: This is a convenience wrapper. For more control, use the paths package.
func ExpandPath(path string) string {
	return paths.ExpandPath(path)
}

// LayeredConfig represents a multi-layered configuration.
// Layers are merged in order of precedence (highest first).
type LayeredConfig struct {
	// Local is the highest precedence layer (e.g., .tool/config.yaml).
	Local interface{}
	// Project is the middle layer (e.g., ~/.valksor/tool/projects.yaml).
	Project interface{}
	// Global is the lowest precedence layer (e.g., ~/.valksor/tool/config.yaml).
	Global interface{}
}

// MergeFunc is a function that merges multiple configuration layers.
// The function receives configs in order of precedence (highest first)
// and should return the merged configuration.
type MergeFunc func(layers []interface{}) interface{}

// LoadLayered loads configuration from multiple layers and merges them.
// The layers are loaded in order of precedence (highest first).
//
// Example:
//
//	merge := func(layers []interface{}) interface{} {
//	    // Custom merge logic
//	    return mergedConfig
//	}
//	merged, err := cfg.LoadLayered(layers, merge)
func LoadLayered(layers []Layer, merge MergeFunc) (*LayeredConfig, error) {
	result := &LayeredConfig{}

	// Load each layer in order
	for _, layer := range layers {
		var value interface{}

		// Skip optional layers that don't exist
		if layer.Optional && !PathExists(layer.Path) {
			continue
		}

		// Load based on layer type
		switch layer.Type {
		case LayerTypeYAML:
			var config interface{}
			if err := LoadYAML(layer.Path, &config); err != nil {
				return nil, err
			}
			value = config
		case LayerTypeJSON:
			var config interface{}
			if err := LoadJSON(layer.Path, &config); err != nil {
				return nil, err
			}
			value = config
		case LayerTypeCustom:
			if layer.Loader == nil {
				return nil, errors.New("custom loader required for LayerTypeCustom")
			}
			var config interface{}
			if err := layer.Loader.Load(layer.Path, &config); err != nil {
				return nil, err
			}
			value = config
		}

		// Assign to appropriate layer slot
		switch layer.Precedence {
		case PrecedenceLocal:
			result.Local = value
		case PrecedenceProject:
			result.Project = value
		case PrecedenceGlobal:
			result.Global = value
		}
	}

	// Apply merge function if provided
	if merge != nil {
		layers := []interface{}{result.Local, result.Project, result.Global}
		merged := merge(layers)
		result.Local = merged // Put merged result in highest precedence slot
	}

	return result, nil
}

// LayerType defines the type of configuration layer.
type LayerType int

const (
	LayerTypeYAML LayerType = iota
	LayerTypeJSON
	LayerTypeCustom
)

// LayerPrecedence defines the precedence level of a configuration layer.
type LayerPrecedence int

const (
	PrecedenceLocal LayerPrecedence = iota
	PrecedenceProject
	PrecedenceGlobal
)

// Layer represents a single configuration layer.
type Layer struct {
	// Path is the file path to load.
	Path string
	// Type is the type of configuration file.
	Type LayerType
	// Precedence is the precedence level of this layer.
	Precedence LayerPrecedence
	// Optional indicates whether the file must exist.
	Optional bool
	// Loader is a custom loader for LayerTypeCustom.
	Loader Loader
}

// Loader is a custom configuration loader interface.
type Loader interface {
	Load(path string, v interface{}) error
}

// LoadAllYAMLInDir loads all YAML files in a directory.
// Returns a map of filename (without extension) to loaded config.
func LoadAllYAMLInDir(dir string) (map[string]interface{}, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		var config interface{}
		if err := LoadYAML(path, &config); err != nil {
			return nil, err
		}

		name := entry.Name()[:len(entry.Name())-len(ext)]
		result[name] = config
	}

	return result, nil
}
