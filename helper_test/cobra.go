// Package helper_test provides shared testing utilities for all valksor Go projects.
package helper_test

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestContext provides test context for Cobra command tests.
type TestContext struct {
	T         *testing.T
	TmpDir    string
	StdoutBuf *bytes.Buffer
	StderrBuf *bytes.Buffer
	RootCmd   *cobra.Command
}

// NewTestContext creates a test context for command testing.
// It sets up a temporary directory and captures output.
func NewTestContext(t *testing.T) *TestContext {
	t.Helper()

	tmpDir := t.TempDir()

	// Create stdout/stderr buffers
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	// Create a test root command
	rootCmd := createTestRootCommand(stdoutBuf, stderrBuf)

	// Set working directory
	t.Chdir(tmpDir)

	return &TestContext{
		T:         t,
		TmpDir:    tmpDir,
		StdoutBuf: stdoutBuf,
		StderrBuf: stderrBuf,
		RootCmd:   rootCmd,
	}
}

// createTestRootCommand creates a minimal root command for testing.
func createTestRootCommand(stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "Test command for testing",
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetContext(context.Background())

	return cmd
}

// ExecuteCommand executes a command with the given arguments.
func ExecuteCommand(cmd *cobra.Command, args ...string) error {
	ctx := context.Background()
	cmd.SetContext(ctx)
	cmd.SetArgs(args)

	return cmd.Execute()
}

// ExecuteCommandWithContext executes a command with a custom context.
func ExecuteCommandWithContext(ctx context.Context, cmd *cobra.Command, args ...string) error {
	cmd.SetContext(ctx)
	cmd.SetArgs(args)

	return cmd.Execute()
}

// AssertOutputContains fails the test if the output doesn't contain the substring.
func AssertOutputContains(t *testing.T, buf *bytes.Buffer, substr string) {
	t.Helper()

	output := buf.String()
	if !strings.Contains(output, substr) {
		t.Errorf("output does not contain %q\nGot:\n%s", substr, output)
	}
}

// AssertOutputNotContains fails the test if the output contains the substring.
func AssertOutputNotContains(t *testing.T, buf *bytes.Buffer, substr string) {
	t.Helper()

	output := buf.String()
	if strings.Contains(output, substr) {
		t.Errorf("output should not contain %q\nGot:\n%s", substr, output)
	}
}

// AssertOutputEquals fails the test if the output doesn't match exactly.
func AssertOutputEquals(t *testing.T, buf *bytes.Buffer, expected string) {
	t.Helper()

	output := buf.String()
	if output != expected {
		t.Errorf("output mismatch\nGot:\n%s\nWant:\n%s", output, expected)
	}
}

// AssertStdoutContains is a helper for TestContext.
func (tc *TestContext) AssertStdoutContains(substr string) {
	AssertOutputContains(tc.T, tc.StdoutBuf, substr)
}

// AssertStderrContains is a helper for TestContext.
func (tc *TestContext) AssertStderrContains(substr string) {
	AssertOutputContains(tc.T, tc.StderrBuf, substr)
}

// AssertStdoutNotContains is a helper for TestContext.
func (tc *TestContext) AssertStdoutNotContains(substr string) {
	AssertOutputNotContains(tc.T, tc.StdoutBuf, substr)
}

// Execute executes the root command with arguments.
func (tc *TestContext) Execute(args ...string) error {
	return ExecuteCommand(tc.RootCmd, args...)
}

// ExecuteWithContext executes the root command with a custom context.
func (tc *TestContext) ExecuteWithContext(ctx context.Context, args ...string) error {
	return ExecuteCommandWithContext(ctx, tc.RootCmd, args...)
}

// StdoutString returns the stdout buffer as a string.
func (tc *TestContext) StdoutString() string {
	return tc.StdoutBuf.String()
}

// StderrString returns the stderr buffer as a string.
func (tc *TestContext) StderrString() string {
	return tc.StderrBuf.String()
}

// ResetOutput resets the stdout and stderr buffers.
func (tc *TestContext) ResetOutput() {
	tc.StdoutBuf.Reset()
	tc.StderrBuf.Reset()
}

// AddSubCommand adds a subcommand to the root command.
func (tc *TestContext) AddSubCommand(cmd *cobra.Command) {
	tc.RootCmd.AddCommand(cmd)
}

// WithRegisteredSubCommands adds common subcommands to the root command.
func (tc *TestContext) WithRegisteredSubCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		tc.RootCmd.AddCommand(cmd)
	}
}

// CreateFile creates a file in the test directory.
func (tc *TestContext) CreateFile(relativePath, content string) {
	fullPath := filepath.Join(tc.TmpDir, relativePath)
	WriteFile(tc.T, fullPath, content)
}

// AssertFileExists fails if the file doesn't exist.
func (tc *TestContext) AssertFileExists(relativePath string) {
	AssertFileExists(tc.T, filepath.Join(tc.TmpDir, relativePath))
}

// AssertFileNotExists fails if the file exists.
func (tc *TestContext) AssertFileNotExists(relativePath string) {
	AssertFileNotExists(tc.T, filepath.Join(tc.TmpDir, relativePath))
}

// AssertFileContains fails if the file doesn't contain the content.
func (tc *TestContext) AssertFileContains(relativePath, content string) {
	AssertFileContains(tc.T, filepath.Join(tc.TmpDir, relativePath), content)
}
