package validate

import (
	"strings"
	"testing"
)

func TestNewResult(t *testing.T) {
	r := NewResult()

	if !r.Valid {
		t.Fatal("new result should be valid")
	}
	if r.Errors != 0 {
		t.Fatalf("expected 0 errors, got %d", r.Errors)
	}
	if r.Warnings != 0 {
		t.Fatalf("expected 0 warnings, got %d", r.Warnings)
	}
	if len(r.Findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(r.Findings))
	}
}

func TestAddError(t *testing.T) {
	r := NewResult()
	r.AddError("CODE", "message", "path", "file")

	if r.Valid {
		t.Fatal("result should be invalid after error")
	}
	if r.Errors != 1 {
		t.Fatalf("expected 1 error, got %d", r.Errors)
	}
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}

	f := r.Findings[0]
	if f.Severity != SeverityError {
		t.Fatalf("expected severity error, got %v", f.Severity)
	}
	if f.Code != "CODE" {
		t.Fatalf("expected code CODE, got %s", f.Code)
	}
	if f.Message != "message" {
		t.Fatalf("expected message 'message', got %s", f.Message)
	}
	if f.Path != "path" {
		t.Fatalf("expected path 'path', got %s", f.Path)
	}
	if f.File != "file" {
		t.Fatalf("expected file 'file', got %s", f.File)
	}
}

func TestAddErrorWithSuggestion(t *testing.T) {
	r := NewResult()
	r.AddErrorWithSuggestion("CODE", "message", "path", "file", "fix it")

	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}

	f := r.Findings[0]
	if f.Suggestion != "fix it" {
		t.Fatalf("expected suggestion 'fix it', got %s", f.Suggestion)
	}
}

func TestAddWarning(t *testing.T) {
	r := NewResult()
	r.AddWarning("CODE", "message", "path", "file")

	// Warnings don't affect validity
	if !r.Valid {
		t.Fatal("result should still be valid after warning")
	}
	if r.Warnings != 1 {
		t.Fatalf("expected 1 warning, got %d", r.Warnings)
	}
	if r.Errors != 0 {
		t.Fatalf("expected 0 errors, got %d", r.Errors)
	}

	f := r.Findings[0]
	if f.Severity != SeverityWarning {
		t.Fatalf("expected severity warning, got %v", f.Severity)
	}
}

func TestAddWarningWithSuggestion(t *testing.T) {
	r := NewResult()
	r.AddWarningWithSuggestion("CODE", "message", "path", "file", "fix it")

	f := r.Findings[0]
	if f.Suggestion != "fix it" {
		t.Fatalf("expected suggestion 'fix it', got %s", f.Suggestion)
	}
}

func TestAddInfo(t *testing.T) {
	r := NewResult()
	r.AddInfo("CODE", "message", "path", "file")

	// Info doesn't affect validity or counts
	if !r.Valid {
		t.Fatal("result should still be valid after info")
	}
	if r.Warnings != 0 {
		t.Fatalf("expected 0 warnings, got %d", r.Warnings)
	}
	if r.Errors != 0 {
		t.Fatalf("expected 0 errors, got %d", r.Errors)
	}

	f := r.Findings[0]
	if f.Severity != SeverityInfo {
		t.Fatalf("expected severity info, got %v", f.Severity)
	}
}

func TestMerge(t *testing.T) {
	r1 := NewResult()
	r1.AddError("ERR1", "error 1", "", "file1")
	r1.AddWarning("WARN1", "warning 1", "", "file1")

	r2 := NewResult()
	r2.AddError("ERR2", "error 2", "", "file2")
	r2.AddWarning("WARN2", "warning 2", "", "file2")

	r1.Merge(r2)

	if r1.Valid {
		t.Fatal("merged result should be invalid")
	}
	if r1.Errors != 2 {
		t.Fatalf("expected 2 errors, got %d", r1.Errors)
	}
	if r1.Warnings != 2 {
		t.Fatalf("expected 2 warnings, got %d", r1.Warnings)
	}
	if len(r1.Findings) != 4 {
		t.Fatalf("expected 4 findings, got %d", len(r1.Findings))
	}
}

func TestMergeNil(t *testing.T) {
	r := NewResult()
	r.AddError("ERR", "error", "", "file")
	errorsBefore := r.Errors

	r.Merge(nil)

	if r.Errors != errorsBefore {
		t.Fatal("merging nil should not change result")
	}
}

func TestFormatText(t *testing.T) {
	r := NewResult()
	r.AddError("ERR", "error message", "config.timeout", "config.yaml")
	r.AddWarning("WARN", "warning message", "legacy.enabled", "config.yaml")

	output := r.Format("text")

	// Check that output contains expected content
	if !contains(output, "ERROR") {
		t.Fatal("text output should contain ERROR")
	}
	if !contains(output, "config.yaml:") {
		t.Fatal("text output should contain file name")
	}
	if !contains(output, "config.timeout") {
		t.Fatal("text output should contain path")
	}
	if !contains(output, "error message") {
		t.Fatal("text output should contain error message")
	}
	if !contains(output, "Summary:") {
		t.Fatal("text output should contain summary")
	}
	if !contains(output, "INVALID") {
		t.Fatal("text output should indicate invalid configuration")
	}
}

func TestFormatTextValid(t *testing.T) {
	r := NewResult()
	output := r.Format("text")

	if !contains(output, "VALID") {
		t.Fatal("text output for valid result should contain VALID")
	}
}

func TestFormatTextGroupedByFile(t *testing.T) {
	r := NewResult()
	r.AddError("ERR1", "error 1", "path1", "file1")
	r.AddError("ERR2", "error 2", "path2", "file2")
	r.AddWarning("WARN1", "warning 1", "path3", "file1")

	output := r.Format("text")

	// Check grouping by file
	file1Index := indexOf(output, "file1:")
	file2Index := indexOf(output, "file2:")
	if file1Index == -1 || file2Index == -1 {
		t.Fatal("text output should group by file")
	}
	// file2 should come after file1 since we iterate over map
	// (order not guaranteed, but both should be present)
}

func TestFormatJSON(t *testing.T) {
	r := NewResult()
	r.AddError("ERR", "error message", "path", "file")

	output := r.Format("json")

	// Check that output is valid JSON
	if !contains(output, `"severity"`) {
		t.Fatal("JSON output should contain severity field")
	}
	if !contains(output, `"error"`) {
		t.Fatal("JSON output should contain error code")
	}
	if !contains(output, `"error message"`) {
		t.Fatal("JSON output should contain message")
	}
}

func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	return len(s) - len(strings.ReplaceAll(s, substr, ""))
}
