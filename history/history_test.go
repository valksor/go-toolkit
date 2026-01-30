package history

import (
	"testing"
	"time"
)

// TestEntry is a simple test type for history entries.
type TestEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Value     int       `json:"value"`
}

func TestHistory_LoadEmpty(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	entries, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Load() returned %d entries, want 0", len(entries))
	}
}

func TestHistory_SaveAndLoad(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	entry := TestEntry{
		Timestamp: time.Now().Truncate(time.Second),
		Message:   "Test entry",
		Value:     42,
	}

	if err := h.Save(entry); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("Load() returned %d entries, want 1", len(loaded))
	}

	if loaded[0].Message != entry.Message {
		t.Errorf("Load() message = %q, want %q", loaded[0].Message, entry.Message)
	}

	if loaded[0].Value != entry.Value {
		t.Errorf("Load() value = %d, want %d", loaded[0].Value, entry.Value)
	}
}

func TestHistory_MaxEntries(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json", WithMaxEntries(5))

	// Save 10 entries
	for i := range 10 {
		entry := TestEntry{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			Message:   "Entry",
			Value:     i,
		}
		if err := h.Save(entry); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Should only have last 5
	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 5 {
		t.Errorf("Load() returned %d entries, want 5 (max)", len(loaded))
	}

	// Verify they are the last 5 (values 5-9)
	for i, entry := range loaded {
		expectedValue := i + 5
		if entry.Value != expectedValue {
			t.Errorf("Load()[%d].Value = %d, want %d", i, entry.Value, expectedValue)
		}
	}
}

func TestHistory_DefaultMaxEntries(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	// Save 15 entries
	for i := range 15 {
		entry := TestEntry{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			Value:     i,
		}
		if err := h.Save(entry); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Should only have last 10 (default)
	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != DefaultMaxEntries {
		t.Errorf("Load() returned %d entries, want %d (default max)", len(loaded), DefaultMaxEntries)
	}
}

func TestHistory_Clear(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	// Save an entry
	entry := TestEntry{
		Timestamp: time.Now(),
		Message:   "Test",
	}
	if err := h.Save(entry); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Clear
	if err := h.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// Verify empty
	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 0 {
		t.Errorf("Load() after Clear() returned %d entries, want 0", len(loaded))
	}
}

func TestHistory_ClearNonExistent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "nonexistent.json")

	// Clear should not error for non-existent file
	if err := h.Clear(); err != nil {
		t.Errorf("Clear() on non-existent file should not error, got: %v", err)
	}
}

func TestHistory_Path(t *testing.T) {
	h := New[TestEntry]("/some/dir", "myfile.json")

	path := h.Path()
	want := "/some/dir/myfile.json"

	if path != want {
		t.Errorf("Path() = %q, want %q", path, want)
	}
}

func TestHistory_Len(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	// Initially empty
	length, err := h.Len()
	if err != nil {
		t.Fatalf("Len() error = %v", err)
	}
	if length != 0 {
		t.Errorf("Len() = %d, want 0", length)
	}

	// Add entries
	for i := range 3 {
		if err := h.Save(TestEntry{Value: i}); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	length, err = h.Len()
	if err != nil {
		t.Fatalf("Len() error = %v", err)
	}
	if length != 3 {
		t.Errorf("Len() = %d, want 3", length)
	}
}

func TestHistory_Last(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json")

	// Empty - should return false
	_, ok, err := h.Last()
	if err != nil {
		t.Fatalf("Last() error = %v", err)
	}
	if ok {
		t.Error("Last() should return false for empty history")
	}

	// Add entries
	for i := range 3 {
		if err := h.Save(TestEntry{Value: i}); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	last, ok, err := h.Last()
	if err != nil {
		t.Fatalf("Last() error = %v", err)
	}
	if !ok {
		t.Error("Last() should return true for non-empty history")
	}
	if last.Value != 2 {
		t.Errorf("Last().Value = %d, want 2", last.Value)
	}
}

func TestHistory_SaveAll(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	h := New[TestEntry](tmpDir, "test_history.json", WithMaxEntries(5))

	entries := []TestEntry{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	if err := h.SaveAll(entries); err != nil {
		t.Fatalf("SaveAll() error = %v", err)
	}

	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 3 {
		t.Errorf("Load() returned %d entries, want 3", len(loaded))
	}

	// SaveAll with more than max
	bigEntries := make([]TestEntry, 10)
	for i := range bigEntries {
		bigEntries[i] = TestEntry{Value: i}
	}

	if err := h.SaveAll(bigEntries); err != nil {
		t.Fatalf("SaveAll() error = %v", err)
	}

	loaded, err = h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 5 {
		t.Errorf("Load() returned %d entries, want 5 (max)", len(loaded))
	}
}

func TestWithMaxEntries_Zero(t *testing.T) {
	tmpDir := t.TempDir()

	// Zero should use default
	h := New[TestEntry](tmpDir, "test.json", WithMaxEntries(0))

	// Save more than default
	for i := range 15 {
		if err := h.Save(TestEntry{Value: i}); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	loaded, err := h.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != DefaultMaxEntries {
		t.Errorf("Load() returned %d entries, want %d (default)", len(loaded), DefaultMaxEntries)
	}
}
