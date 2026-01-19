// Package helper_test provides shared testing utilities for all valksor Go projects.
package helper_test

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// WriteFile creates a file with the given content, creating parent directories as needed.
func WriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

// CreateFile creates a file with content in a directory.
func CreateFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	WriteFile(t, path, content)

	return path
}

// CreateDir creates a directory, including parent directories.
func CreateDir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", path, err)
	}
}

// TempDir creates a temporary directory and returns its path.
// Uses t.Cleanup() for automatic cleanup.
func TempDir(t *testing.T) string {
	t.Helper()

	return t.TempDir()
}

// ReadFile reads a file and returns its contents.
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}

	return string(data)
}

// AssertFileExists fails the test if the file does not exist.
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

// AssertFileNotExists fails the test if the file exists.
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected file to not exist: %s", path)
	}
}

// AssertFileContent fails the test if the file content doesn't match.
func AssertFileContent(t *testing.T, path, expected string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	if string(data) != expected {
		t.Errorf("file %s:\n  got:  %q\n  want: %q", path, string(data), expected)
	}
}

// AssertFileContains fails the test if the file doesn't contain the substring.
func AssertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("file %s does not contain %q", path, substr)
	}
}

// AssertEqual fails the test if got != want.
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertNotEqual fails the test if got == want.
func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v, should not equal %v", got, want)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("expected true: %s", msg)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("expected false: %s", msg)
	}
}

// AssertNil fails the test if value is not nil.
func AssertNil(t *testing.T, value any) {
	t.Helper()
	if value != nil && !reflect.ValueOf(value).IsNil() {
		t.Errorf("expected nil, got %v", value)
	}
}

// AssertNotNil fails the test if value is nil.
func AssertNotNil(t *testing.T, value any) {
	t.Helper()
	if value == nil || reflect.ValueOf(value).IsNil() {
		t.Errorf("expected non-nil value")
	}
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

// AssertErrorContains fails the test if err is nil or doesn't contain substr.
func AssertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", substr)

		return
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("error %q does not contain %q", err.Error(), substr)
	}
}

// AssertDeepEqual fails the test if got and want are not deeply equal.
func AssertDeepEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// AssertJSONEqual compares two values by marshaling them to JSON.
func AssertJSONEqual(t *testing.T, got, want any) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal got: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to marshal want: %v", err)
	}
	if string(gotJSON) != string(wantJSON) {
		t.Errorf("JSON mismatch:\n  got:  %s\n  want: %s", gotJSON, wantJSON)
	}
}

// AssertContains fails the test if slice doesn't contain item.
func AssertContains[T comparable](t *testing.T, slice []T, item T) {
	t.Helper()
	for _, s := range slice {
		if s == item {
			return
		}
	}
	t.Errorf("slice does not contain %v", item)
}

// AssertNotContains fails the test if slice contains item.
func AssertNotContains[T comparable](t *testing.T, slice []T, item T) {
	t.Helper()
	for _, s := range slice {
		if s == item {
			t.Errorf("slice should not contain %v", item)

			return
		}
	}
}

// AssertLen fails the test if slice length doesn't match expected.
func AssertLen[T any](t *testing.T, slice []T, expected int) {
	t.Helper()
	if len(slice) != expected {
		t.Errorf("len(slice) = %d, want %d", len(slice), expected)
	}
}

// Logger returns a logger for testing that writes to stderr.
func Logger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// DiscardLogger returns a logger that discards all output.
func DiscardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

// BoolPtr returns a pointer to a bool value.
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns a pointer to a string value.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int value.
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to a float64 value.
func Float64Ptr(f float64) *float64 {
	return &f
}

// TimePtr returns a pointer to a time.Time value.
func TimePtr(tm time.Time) *time.Time {
	return &tm
}
