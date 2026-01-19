// Package helper_test provides shared testing utilities for all valksor Go projects.
package helper_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// T panics if err is non-nil, returning v. Useful for test setup.
// Usage: f := must.T(os.ReadFile("file.txt")).
func T[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

// Eq panics if got != want. Useful for assertions in test setup.
func Eq[T comparable](got, want T) {
	if got != want {
		panic(fmt.Sprintf("got %v, want %v", got, want))
	}
}

// EqFatal calls t.Fatal if got != want.
func EqFatal[T comparable](t *testing.T, got, want T, msgAndArgs ...interface{}) {
	t.Helper()
	if got != want {
		t.Fatal(append([]interface{}{fmt.Sprintf("got %v, want %v", got, want)}, msgAndArgs...)...)
	}
}

// NoError calls t.Fatal if err is non-nil.
func NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatal(append([]interface{}{err}, msgAndArgs...)...)
	}
}

// PanicHandler catches panics and reports them as test failures.
// Usage: defer PanicHandler(t).
func PanicHandler(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Fatalf("panic recovered: %v", r)
	}
}

// TempFile creates a temporary file with the given content.
// Returns the file path and a cleanup function.
func TempFile(t *testing.T, content string) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	return path, func() {
		// Cleanup is handled by t.TempDir()
	}
}
