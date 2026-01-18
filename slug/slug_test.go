package slug

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name   string
		title  string
		want   string
		maxLen int
	}{
		{
			name:   "simple title",
			title:  "Add user authentication",
			maxLen: 50,
			want:   "add-user-authentication",
		},
		{
			name:   "with special characters",
			title:  "Fix bug #123: login fails!",
			maxLen: 50,
			want:   "fix-bug-123-login-fails",
		},
		{
			name:   "with diacritics",
			title:  "Résumé parsing für Änderungen",
			maxLen: 50,
			want:   "resume-parsing-fur-anderungen",
		},
		{
			name:   "truncate at word boundary",
			title:  "This is a very long title that should be truncated",
			maxLen: 20,
			want:   "this-is-a-very-long",
		},
		{
			name:   "truncate mid-word when necessary",
			title:  "Supercalifragilisticexpialidocious",
			maxLen: 15,
			want:   "supercalifragil",
		},
		{
			name:   "empty string",
			title:  "",
			maxLen: 50,
			want:   "",
		},
		{
			name:   "only special chars",
			title:  "!@#$%^&*()",
			maxLen: 50,
			want:   "",
		},
		{
			name:   "underscores to hyphens",
			title:  "user_authentication_module",
			maxLen: 50,
			want:   "user-authentication-module",
		},
		{
			name:   "multiple spaces",
			title:  "Add    multiple   spaces",
			maxLen: 50,
			want:   "add-multiple-spaces",
		},
		{
			name:   "no max length",
			title:  "No limit on length",
			maxLen: 0,
			want:   "no-limit-on-length",
		},
		{
			name:   "negative max length treated as no limit",
			title:  "No limit on length",
			maxLen: -1,
			want:   "no-limit-on-length",
		},
		{
			name:   "leading and trailing spaces",
			title:  "  trimmed  ",
			maxLen: 50,
			want:   "trimmed",
		},
		{
			name:   "numbers preserved",
			title:  "Version 2.0 Release",
			maxLen: 50,
			want:   "version-20-release",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.title, tt.maxLen)
			if got != tt.want {
				t.Errorf("Slugify(%q, %d) = %q, want %q", tt.title, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestTruncateAtWordBoundary(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "no truncation needed",
			input:  "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "truncate at hyphen",
			input:  "first-second-third",
			maxLen: 12,
			want:   "first-second",
		},
		{
			name:   "hyphen too early",
			input:  "a-verylongword",
			maxLen: 10,
			want:   "a-verylong",
		},
		{
			name:   "no hyphen at all",
			input:  "verylongword",
			maxLen: 8,
			want:   "verylong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateAtWordBoundary(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateAtWordBoundary(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func BenchmarkSlugify(b *testing.B) {
	title := "This is a fairly long title with special characters: #123 & diacritics: résumé"
	b.ResetTimer()
	for range b.N {
		_ = Slugify(title, 50)
	}
}
