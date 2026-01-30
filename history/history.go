// Package history provides generic JSON-persisted history tracking for CLI tools.
//
// This package implements a common pattern for tracking attempts, runs, or any
// other historical data that should be persisted across sessions with automatic
// pruning to prevent unbounded growth.
//
// Basic usage:
//
//	type BuildAttempt struct {
//	    Timestamp time.Time `json:"timestamp"`
//	    Success   bool      `json:"success"`
//	    Output    string    `json:"output"`
//	}
//
//	h := history.New[BuildAttempt]("/path/to/workspace", "build_history.json")
//	attempts, _ := h.Load()
//	h.Save(BuildAttempt{Timestamp: time.Now(), Success: true})
package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DefaultMaxEntries is the default number of entries to keep.
const DefaultMaxEntries = 10

// History manages persistence of historical entries of type T.
// T must be JSON-serializable.
type History[T any] struct {
	dir        string
	filename   string
	maxEntries int
}

// Option configures a History instance.
type Option func(*options)

type options struct {
	maxEntries int
}

// WithMaxEntries sets the maximum number of entries to keep.
// Defaults to DefaultMaxEntries (10).
func WithMaxEntries(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.maxEntries = n
		}
	}
}

// New creates a new History manager for storing entries of type T.
// The dir parameter is the directory where the history file will be stored.
// The filename parameter is the name of the JSON file.
func New[T any](dir, filename string, opts ...Option) *History[T] {
	o := &options{
		maxEntries: DefaultMaxEntries,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &History[T]{
		dir:        dir,
		filename:   filename,
		maxEntries: o.maxEntries,
	}
}

// Path returns the full path to the history file.
func (h *History[T]) Path() string {
	return filepath.Join(h.dir, h.filename)
}

// Load loads all entries from disk.
// Returns an empty slice if the file doesn't exist.
func (h *History[T]) Load() ([]T, error) {
	data, err := os.ReadFile(h.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}

		return nil, err
	}

	var entries []T
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// Save appends an entry and persists to disk, keeping only the last maxEntries.
func (h *History[T]) Save(entry T) error {
	entries, err := h.Load()
	if err != nil {
		return err
	}

	entries = append(entries, entry)
	if len(entries) > h.maxEntries {
		entries = entries[len(entries)-h.maxEntries:]
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(h.dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(h.Path(), data, 0o644)
}

// SaveAll replaces all entries and persists to disk.
func (h *History[T]) SaveAll(entries []T) error {
	if len(entries) > h.maxEntries {
		entries = entries[len(entries)-h.maxEntries:]
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(h.dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(h.Path(), data, 0o644)
}

// Clear removes all history by deleting the history file.
func (h *History[T]) Clear() error {
	err := os.Remove(h.Path())
	if os.IsNotExist(err) {
		return nil // Already cleared
	}

	return err
}

// Len returns the number of entries currently stored.
func (h *History[T]) Len() (int, error) {
	entries, err := h.Load()
	if err != nil {
		return 0, err
	}

	return len(entries), nil
}

// Last returns the most recent entry, if any.
// Returns the zero value and false if there are no entries.
func (h *History[T]) Last() (T, bool, error) {
	entries, err := h.Load()
	if err != nil {
		var zero T

		return zero, false, err
	}

	if len(entries) == 0 {
		var zero T

		return zero, false, nil
	}

	return entries[len(entries)-1], true, nil
}
