package display

import (
	"testing"
)

func TestInitColors(t *testing.T) {
	tests := []struct {
		name     string
		noColor  bool
		envSet   bool
		expected bool
	}{
		{
			name:     "colors enabled by default",
			noColor:  false,
			envSet:   false,
			expected: true,
		},
		{
			name:     "no-color flag disables",
			noColor:  true,
			envSet:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state
			SetColorsEnabled(true)

			if tt.envSet {
				t.Setenv("NO_COLOR", "1")
			}

			InitColors(tt.noColor)

			if ColorsEnabled() != tt.expected {
				t.Errorf("After InitColors(%v), ColorsEnabled() = %v, want %v",
					tt.noColor, ColorsEnabled(), tt.expected)
			}
		})
	}
}

func TestColorsEnabled(t *testing.T) {
	SetColorsEnabled(true)
	if !ColorsEnabled() {
		t.Error("ColorsEnabled() = false, want true after SetColorsEnabled(true)")
	}

	SetColorsEnabled(false)
	if ColorsEnabled() {
		t.Error("ColorsEnabled() = true, want false after SetColorsEnabled(false)")
	}
}

func TestSetColorsEnabled(t *testing.T) {
	SetColorsEnabled(true)
	if !ColorsEnabled() {
		t.Error("SetColorsEnabled(true) failed")
	}

	SetColorsEnabled(false)
	if ColorsEnabled() {
		t.Error("SetColorsEnabled(false) failed")
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		enabled  bool
		contains string
	}{
		{
			name:     "with colors enabled",
			text:     "test",
			enabled:  true,
			contains: "\033[",
		},
		{
			name:     "with colors disabled",
			text:     "test",
			enabled:  false,
			contains: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetColorsEnabled(tt.enabled)

			result := colorize(tt.text, red)

			if !contains(result, tt.contains) {
				t.Errorf("colorize() = %v, want to contain %v", result, tt.contains)
			}
		})
	}
}

func TestColorFunctions(t *testing.T) {
	SetColorsEnabled(true)

	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"Success", Success},
		{"Error", Error},
		{"Warning", Warning},
		{"Info", Info},
		{"Muted", Muted},
		{"Bold", Bold},
		{"Cyan", Cyan},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn("test")

			if result == "" {
				t.Errorf("%s() returned empty string", tt.name)
			}

			if !contains(result, "test") {
				t.Errorf("%s() = %v, should contain 'test'", tt.name, result)
			}
		})
	}
}

func TestPrefixFunctions(t *testing.T) {
	SetColorsEnabled(true)

	tests := []struct {
		name string
		fn   func() string
	}{
		{"SuccessPrefix", SuccessPrefix},
		{"ErrorPrefix", ErrorPrefix},
		{"WarningPrefix", WarningPrefix},
		{"InfoPrefix", InfoPrefix},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()

			if result == "" {
				t.Errorf("%s() returned empty string", tt.name)
			}
		})
	}
}

func TestMessageFunctions(t *testing.T) {
	SetColorsEnabled(true)

	tests := []struct {
		name string
		fn   func(string, ...any) string
	}{
		{"SuccessMsg", SuccessMsg},
		{"ErrorMsg", ErrorMsg},
		{"WarningMsg", WarningMsg},
		{"InfoMsg", InfoMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn("test message: %s", "value")

			if result == "" {
				t.Errorf("%s() returned empty string", tt.name)
			}

			if !contains(result, "test message:") {
				t.Errorf("%s() = %v, should contain formatted message", tt.name, result)
			}
		})
	}
}

func TestNoColorEnvVar(t *testing.T) {
	// Set NO_COLOR
	t.Setenv("NO_COLOR", "1")
	InitColors(false)

	if ColorsEnabled() {
		t.Error("NO_COLOR environment variable should disable colors")
	}

	// Clean up for other tests
	SetColorsEnabled(true)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
