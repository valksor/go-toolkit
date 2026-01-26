package licensing

import (
	"context"
	"testing"
)

func TestGetProjectLicense(t *testing.T) {
	result := GetProjectLicense()

	if result == "" {
		t.Error("GetProjectLicense() returned empty string")
	}

	// Check it contains expected license keywords
	expectedTerms := []string{"BSD 3-Clause", "Copyright", "SIA Valksor"}
	for _, term := range expectedTerms {
		if !contains(result, term) {
			t.Errorf("License text missing expected term: %s", term)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

func TestGetDependencyLicenses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ctx := context.Background()

	libs, err := GetDependencyLicenses(ctx, ".")
	if err != nil {
		t.Fatalf("GetDependencyLicenses() error = %v", err)
	}

	if len(libs) == 0 {
		t.Error("GetDependencyLicenses() returned no licenses")
	}

	// Check that at least some known dependencies are found
	knownDeps := []string{"github.com/spf13/cobra", "github.com/google/uuid"}
	for _, dep := range knownDeps {
		found := false
		for _, lib := range libs {
			if lib.Path == dep || contains(lib.Path, dep) {
				found = true

				break
			}
		}
		if !found {
			t.Logf("Warning: known dependency not found: %s", dep)
		}
	}

	// Check that all entries have non-empty paths
	for i, lib := range libs {
		if lib.Path == "" {
			t.Errorf("libs[%d].Path is empty", i)
		}
		if lib.Unknown && lib.License == "" {
			t.Errorf("libs[%d] has Unknown=true but empty License", i)
		}
	}
}
