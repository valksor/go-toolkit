// Package providerconfig provides configuration types for providers.
package providerconfig

import "context"

// Config holds provider configuration.
type Config struct {
	options map[string]any
}

// NewConfig creates a new config.
func NewConfig() Config {
	return Config{options: make(map[string]any)}
}

// Set sets an option and returns the config for chaining.
func (c Config) Set(key string, value any) Config {
	c.options[key] = value
	return c
}

// Get gets an option.
func (c Config) Get(key string) any {
	return c.options[key]
}

// GetString gets a string option.
func (c Config) GetString(key string) string {
	if v, ok := c.options[key].(string); ok {
		return v
	}
	return ""
}

// GetBool gets a bool option.
func (c Config) GetBool(key string) bool {
	if v, ok := c.options[key].(bool); ok {
		return v
	}
	return false
}

// GetInt gets an int option.
func (c Config) GetInt(key string) int {
	if v, ok := c.options[key].(int); ok {
		return v
	}
	if v, ok := c.options[key].(int64); ok {
		return int(v)
	}
	return 0
}

// GetIntWithDefault gets an int option with a default value.
func (c Config) GetIntWithDefault(key string, defaultVal int) int {
	if v, ok := c.options[key].(int); ok {
		return v
	}
	if v, ok := c.options[key].(int64); ok {
		return int(v)
	}
	return defaultVal
}

// GetStringWithDefault gets a string option with a default value.
func (c Config) GetStringWithDefault(key, defaultVal string) string {
	if v, ok := c.options[key].(string); ok {
		return v
	}
	return defaultVal
}

// GetStringSlice gets a string slice option.
func (c Config) GetStringSlice(key string) []string {
	if v, ok := c.options[key].([]string); ok {
		return v
	}
	return nil
}

// Has checks if a key exists.
func (c Config) Has(key string) bool {
	_, ok := c.options[key]
	return ok
}

// Keys returns all keys.
func (c Config) Keys() []string {
	keys := make([]string, 0, len(c.options))
	for k := range c.options {
		keys = append(keys, k)
	}
	return keys
}

// ConfigReader reads config from a source.
type ConfigReader interface {
	Read(ctx context.Context) (Config, error)
}

// ConfigWriter writes config to a source.
type ConfigWriter interface {
	Write(ctx context.Context, cfg Config) error
}

// ConfigSource locates and manages config files.
type ConfigSource interface {
	Path() string
	Exists() bool
}

// ConfigManager combines read and write capabilities.
// Providers implement this so external code (CLI/Web UI) can configure them.
type ConfigManager interface {
	ConfigReader
	ConfigWriter
	ConfigSource
}
