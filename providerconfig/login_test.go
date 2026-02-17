package providerconfig

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "*******"},
		{"short", "*******"},
		{"12345678", "*******"},
		{"123456789", "1234...6789"},
		{"abcdefghijklmnop", "abcd...mnop"},
	}

	for _, tt := range tests {
		got := MaskToken(tt.input)
		if got != tt.want {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPrintTokenHelp(t *testing.T) {
	var buf bytes.Buffer
	lc := LoginConfig{
		Name:        "TestProvider",
		HelpURL:     "https://example.com/tokens",
		HelpSteps:   "Settings → API → Create Token",
		Scopes:      "read, write",
		TokenPrefix: "tp_",
	}

	PrintTokenHelp(&buf, lc)
	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("TestProvider Token Setup")) {
		t.Error("Expected provider name in output")
	}
	if !bytes.Contains([]byte(output), []byte("https://example.com/tokens")) {
		t.Error("Expected help URL in output")
	}
	if !bytes.Contains([]byte(output), []byte("Settings → API → Create Token")) {
		t.Error("Expected help steps in output")
	}
	if !bytes.Contains([]byte(output), []byte("read, write")) {
		t.Error("Expected scopes in output")
	}
	if !bytes.Contains([]byte(output), []byte("tp_")) {
		t.Error("Expected token prefix in output")
	}
}

// mockConfigManager implements ConfigManager for testing.
type mockConfigManager struct {
	config     Config
	exists     bool
	path       string
	readErr    error
	writeErr   error
	writeCalls []Config
}

func (m *mockConfigManager) Read(_ context.Context) (Config, error) {
	if m.readErr != nil {
		return Config{}, m.readErr
	}
	return m.config, nil
}

func (m *mockConfigManager) Write(_ context.Context, cfg Config) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	m.writeCalls = append(m.writeCalls, cfg)
	m.config = cfg
	return nil
}

func (m *mockConfigManager) Path() string {
	return m.path
}

func (m *mockConfigManager) Exists() bool {
	return m.exists
}

func TestRunLogin_NewConfig(t *testing.T) {
	mgr := &mockConfigManager{
		exists: false,
		path:   ".crealfy/test.yaml",
	}

	lc := LoginConfig{
		Name:       "Test",
		TokenField: "token",
	}

	// Simulate user entering a token
	input := bytes.NewBufferString("my-secret-token\n")
	output := &bytes.Buffer{}

	result, err := RunLogin(context.Background(), mgr, lc, LoginOptions{
		Stdin:  input,
		Stdout: output,
	})

	if err != nil {
		t.Fatalf("RunLogin error: %v", err)
	}

	if result.Cancelled {
		t.Error("Expected login to succeed, got cancelled")
	}

	if len(mgr.writeCalls) != 1 {
		t.Fatalf("Expected 1 write call, got %d", len(mgr.writeCalls))
	}

	savedToken := mgr.writeCalls[0].GetString("token")
	if savedToken != "my-secret-token" {
		t.Errorf("Expected token 'my-secret-token', got %q", savedToken)
	}
}

func TestRunLogin_Cancelled(t *testing.T) {
	mgr := &mockConfigManager{
		exists: false,
		path:   ".crealfy/test.yaml",
	}

	lc := LoginConfig{
		Name:       "Test",
		TokenField: "token",
	}

	// Simulate user pressing enter without input (cancel)
	input := bytes.NewBufferString("\n")
	output := &bytes.Buffer{}

	result, err := RunLogin(context.Background(), mgr, lc, LoginOptions{
		Stdin:  input,
		Stdout: output,
	})

	if err != nil {
		t.Fatalf("RunLogin error: %v", err)
	}

	if !result.Cancelled {
		t.Error("Expected login to be cancelled")
	}

	if len(mgr.writeCalls) != 0 {
		t.Error("Expected no write calls when cancelled")
	}
}

func TestRunLogin_ExistingToken_NoOverride(t *testing.T) {
	existingConfig := NewConfig().Set("token", "existing-token")
	mgr := &mockConfigManager{
		exists: true,
		config: existingConfig,
		path:   ".crealfy/test.yaml",
	}

	lc := LoginConfig{
		Name:       "Test",
		TokenField: "token",
	}

	// Simulate user declining to override
	input := bytes.NewBufferString("n\n")
	output := &bytes.Buffer{}

	result, err := RunLogin(context.Background(), mgr, lc, LoginOptions{
		Stdin:  input,
		Stdout: output,
	})

	if err != nil {
		t.Fatalf("RunLogin error: %v", err)
	}

	if !result.Cancelled {
		t.Error("Expected login to be cancelled when user declines override")
	}
}

func TestRunLogin_Force(t *testing.T) {
	existingConfig := NewConfig().Set("token", "old-token")
	mgr := &mockConfigManager{
		exists: true,
		config: existingConfig,
		path:   ".crealfy/test.yaml",
	}

	lc := LoginConfig{
		Name:       "Test",
		TokenField: "token",
	}

	// Simulate entering new token with force option
	input := bytes.NewBufferString("new-token\n")
	output := &bytes.Buffer{}

	result, err := RunLogin(context.Background(), mgr, lc, LoginOptions{
		Force:  true,
		Stdin:  input,
		Stdout: output,
	})

	if err != nil {
		t.Fatalf("RunLogin error: %v", err)
	}

	if result.Cancelled {
		t.Error("Expected login to succeed with force option")
	}

	if len(mgr.writeCalls) != 1 {
		t.Fatalf("Expected 1 write call, got %d", len(mgr.writeCalls))
	}

	savedToken := mgr.writeCalls[0].GetString("token")
	if savedToken != "new-token" {
		t.Errorf("Expected token 'new-token', got %q", savedToken)
	}
}

func TestDetectExistingToken_FromEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_TOKEN", "env-token-value")
	defer os.Unsetenv("TEST_TOKEN")

	mgr := &mockConfigManager{exists: false}
	lc := LoginConfig{
		EnvVar:     "TEST_TOKEN",
		TokenField: "token",
	}

	source, masked := DetectExistingToken(context.Background(), mgr, lc)

	if source != "TEST_TOKEN environment variable" {
		t.Errorf("Expected source 'TEST_TOKEN environment variable', got %q", source)
	}

	if masked != "env-...alue" {
		t.Errorf("Expected masked 'env-...alue', got %q", masked)
	}
}

func TestDetectExistingToken_FromConfig(t *testing.T) {
	existingConfig := NewConfig().Set("token", "config-token-value")
	mgr := &mockConfigManager{
		exists: true,
		config: existingConfig,
	}
	lc := LoginConfig{
		EnvVar:     "NONEXISTENT_VAR",
		TokenField: "token",
	}

	source, masked := DetectExistingToken(context.Background(), mgr, lc)

	if source != "config file" {
		t.Errorf("Expected source 'config file', got %q", source)
	}

	if masked != "conf...alue" {
		t.Errorf("Expected masked 'conf...alue', got %q", masked)
	}
}

func TestDetectExistingToken_NotFound(t *testing.T) {
	mgr := &mockConfigManager{exists: false}
	lc := LoginConfig{
		EnvVar:     "NONEXISTENT_VAR",
		TokenField: "token",
	}

	source, masked := DetectExistingToken(context.Background(), mgr, lc)

	if source != "" || masked != "" {
		t.Errorf("Expected empty result, got source=%q masked=%q", source, masked)
	}
}

func TestDetectExistingToken_EnvVarReference(t *testing.T) {
	// Token is a reference like ${VAR}, should be treated as not found
	existingConfig := NewConfig().Set("token", "${WRIKE_TOKEN}")
	mgr := &mockConfigManager{
		exists: true,
		config: existingConfig,
	}
	lc := LoginConfig{
		EnvVar:     "NONEXISTENT_VAR",
		TokenField: "token",
	}

	source, masked := DetectExistingToken(context.Background(), mgr, lc)

	if source != "" || masked != "" {
		t.Errorf("Expected empty result for env var reference, got source=%q masked=%q", source, masked)
	}
}

func TestLoginConfig_DefaultFields(t *testing.T) {
	// Verify LoginConfig can be created with all fields
	lc := LoginConfig{
		Name:        "Provider",
		EnvVar:      "PROVIDER_TOKEN",
		TokenField:  "api_key",
		HelpURL:     "https://example.com",
		HelpSteps:   "Step 1, Step 2",
		Scopes:      "read, write",
		TokenPrefix: "pk_",
	}

	if lc.Name != "Provider" {
		t.Error("Name not set correctly")
	}
	if lc.TokenPrefix != "pk_" {
		t.Error("TokenPrefix not set correctly")
	}
}

// Integration test with real file system
func TestRunLogin_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".crealfy", "test.yaml")

	// Create a real file-based config manager for integration test
	mgr := &fileTestManager{
		path:   configPath,
		exists: false,
	}

	lc := LoginConfig{
		Name:       "IntegrationTest",
		TokenField: "token",
	}

	input := bytes.NewBufferString("integration-token\n")
	output := &bytes.Buffer{}

	result, err := RunLogin(context.Background(), mgr, lc, LoginOptions{
		Stdin:  input,
		Stdout: output,
	})

	if err != nil {
		t.Fatalf("RunLogin error: %v", err)
	}

	if result.Cancelled {
		t.Error("Expected login to succeed")
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}
}

// fileTestManager is a simple file-based config manager for integration testing.
type fileTestManager struct {
	path   string
	exists bool
	config Config
}

func (m *fileTestManager) Read(_ context.Context) (Config, error) {
	return m.config, nil
}

func (m *fileTestManager) Write(_ context.Context, cfg Config) error {
	// Create directory
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Write a placeholder file
	return os.WriteFile(m.path, []byte("token: test"), 0o600)
}

func (m *fileTestManager) Path() string {
	return m.path
}

func (m *fileTestManager) Exists() bool {
	return m.exists
}
