// Package validate provides a generic validation result system.
//
// The validation system captures findings with severity levels (error, warning, info),
// supports grouping by source file, and outputs in multiple formats (text, JSON).
//
// Usage:
//
//	result := validate.NewResult()
//	result.AddError("INVALID_VALUE", "Value must be positive", "config.timeout", "config.yaml")
//	result.AddWarning("DEPRECATED", "This field is deprecated", "legacy.enabled", "config.yaml")
//
//	if !result.Valid {
//	    fmt.Println(result.Format("text"))
//	}
package validate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Severity indicates the importance of a validation finding.
type Severity string

const (
	// SeverityError indicates a critical issue that must be fixed.
	SeverityError Severity = "error"
	// SeverityWarning indicates a non-critical issue that should be addressed.
	SeverityWarning Severity = "warning"
	// SeverityInfo indicates informational feedback.
	SeverityInfo Severity = "info"
)

// Finding represents a single validation issue.
type Finding struct {
	Severity   Severity `json:"severity"`
	Code       string   `json:"code"`                 // e.g., "INVALID_VALUE"
	Message    string   `json:"message"`              // Human-readable message
	Path       string   `json:"path,omitempty"`       // Config path, e.g., "server.port"
	File       string   `json:"file,omitempty"`       // Source file
	Suggestion string   `json:"suggestion,omitempty"` // How to fix
}

// Result holds all validation findings.
type Result struct {
	Findings []Finding `json:"findings"`
	Errors   int       `json:"errors"`
	Warnings int       `json:"warnings"`
	Valid    bool      `json:"valid"`
}

// NewResult creates an empty validation result.
func NewResult() *Result {
	return &Result{
		Valid:    true,
		Findings: make([]Finding, 0),
	}
}

// AddError adds an error finding.
func (r *Result) AddError(code, message, path, file string) {
	r.addFinding(SeverityError, code, message, path, file, "")
}

// AddErrorWithSuggestion adds an error finding with a fix suggestion.
func (r *Result) AddErrorWithSuggestion(code, message, path, file, suggestion string) {
	r.addFinding(SeverityError, code, message, path, file, suggestion)
}

// AddWarning adds a warning finding.
func (r *Result) AddWarning(code, message, path, file string) {
	r.addFinding(SeverityWarning, code, message, path, file, "")
}

// AddWarningWithSuggestion adds a warning finding with a fix suggestion.
func (r *Result) AddWarningWithSuggestion(code, message, path, file, suggestion string) {
	r.addFinding(SeverityWarning, code, message, path, file, suggestion)
}

// AddInfo adds an informational finding.
func (r *Result) AddInfo(code, message, path, file string) {
	r.addFinding(SeverityInfo, code, message, path, file, "")
}

func (r *Result) addFinding(severity Severity, code, message, path, file, suggestion string) {
	finding := Finding{
		Severity:   severity,
		Code:       code,
		Message:    message,
		Path:       path,
		File:       file,
		Suggestion: suggestion,
	}
	r.Findings = append(r.Findings, finding)

	switch severity {
	case SeverityError:
		r.Errors++
		r.Valid = false
	case SeverityWarning:
		r.Warnings++
	case SeverityInfo:
		// Info level findings don't affect validation result
	}
}

// Merge combines another result into this one.
func (r *Result) Merge(other *Result) {
	if other == nil {
		return
	}
	r.Findings = append(r.Findings, other.Findings...)
	r.Errors += other.Errors
	r.Warnings += other.Warnings
	if other.Errors > 0 {
		r.Valid = false
	}
}

// Format returns the result in the specified format.
// Supported formats: "json", "text" (default).
func (r *Result) Format(format string) string {
	switch format {
	case "json":
		return r.formatJSON()
	default:
		return r.formatText()
	}
}

func (r *Result) formatJSON() string {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal result: %s"}`, err)
	}

	return string(data)
}

func (r *Result) formatText() string {
	var sb strings.Builder

	// Group findings by file
	byFile := make(map[string][]Finding)
	for _, f := range r.Findings {
		file := f.File
		if file == "" {
			file = "(general)"
		}
		byFile[file] = append(byFile[file], f)
	}

	// Print findings grouped by file
	for file, findings := range byFile {
		sb.WriteString(file + ":\n")
		for _, f := range findings {
			severityStr := strings.ToUpper(string(f.Severity))
			sb.WriteString(fmt.Sprintf("  %s [%s] %s: %s\n", severityStr, f.Code, f.Path, f.Message))
			if f.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("    Suggestion: %s\n", f.Suggestion))
			}
		}
		sb.WriteString("\n")
	}

	// Print summary
	if r.Errors == 0 && r.Warnings == 0 {
		sb.WriteString("Configuration is VALID\n")
	} else {
		sb.WriteString(fmt.Sprintf("Summary: %d error(s), %d warning(s)\n", r.Errors, r.Warnings))
		if r.Valid {
			sb.WriteString("Configuration is VALID (with warnings)\n")
		} else {
			sb.WriteString("Configuration is INVALID\n")
		}
	}

	return sb.String()
}
