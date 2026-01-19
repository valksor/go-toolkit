package helper_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// CreateTempGitRepo creates an initialized git repository in a temporary directory.
// It configures user.email and user.name, and creates an initial commit.
// Returns the path to the repository root.
func CreateTempGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Initialize git repo
	if err := runGit(t, dir, "init"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	// Configure git user (required for commits)
	mustRunGit(t, dir, "config", "user.email", "test@example.com")
	mustRunGit(t, dir, "config", "user.name", "Test User")

	// Create initial commit (many operations require at least one commit)
	WriteFile(t, filepath.Join(dir, "README.md"), "# Test Repository\n")
	mustRunGit(t, dir, "add", ".")
	mustRunGit(t, dir, "commit", "-m", "initial commit")

	return dir
}

// CreateTempGitRepoWithBranch creates a git repo and switches to the specified branch.
func CreateTempGitRepoWithBranch(t *testing.T, branch string) string {
	t.Helper()
	dir := CreateTempGitRepo(t)
	mustRunGit(t, dir, "checkout", "-b", branch)

	return dir
}

// CreateTempGitRepoInDir initializes a git repository in an existing directory.
func CreateTempGitRepoInDir(t *testing.T, dir string) {
	t.Helper()
	if err := runGit(t, dir, "init"); err != nil {
		t.Skipf("git not available: %v", err)
	}
	mustRunGit(t, dir, "config", "user.email", "test@example.com")
	mustRunGit(t, dir, "config", "user.name", "Test User")
}

// WriteFileAndCommit writes a file and commits it.
func WriteFileAndCommit(t *testing.T, dir, relativePath, content, message string) {
	t.Helper()
	WriteFile(t, filepath.Join(dir, relativePath), content)
	mustRunGit(t, dir, "add", relativePath)
	mustRunGit(t, dir, "commit", "-m", message)
}

// RunGit runs a git command and returns its output.
func RunGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\nOutput: %s", args, err, output)
	}

	return string(output)
}

// runGit runs a git command and returns any error.
func runGit(t *testing.T, dir string, args ...string) error {
	t.Helper()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE=2020-01-01T00:00:00Z",
		"GIT_COMMITTER_DATE=2020-01-01T00:00:00Z",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("git %v failed: %s", args, output)
	}

	return err
}

// mustRunGit runs a git command and fails the test if it errors.
func mustRunGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	if err := runGit(t, dir, args...); err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
}

// GitAdd stages a file for commit.
func GitAdd(t *testing.T, dir, path string) {
	t.Helper()
	mustRunGit(t, dir, "add", path)
}

// GitCommit creates a commit with the given message.
func GitCommit(t *testing.T, dir, message string) {
	t.Helper()
	mustRunGit(t, dir, "commit", "-m", message)
}

// GitCheckout switches to the specified branch.
func GitCheckout(t *testing.T, dir, branch string) {
	t.Helper()
	mustRunGit(t, dir, "checkout", branch)
}

// GitCreateBranch creates and switches to a new branch.
func GitCreateBranch(t *testing.T, dir, branch string) {
	t.Helper()
	mustRunGit(t, dir, "checkout", "-b", branch)
}

// GitCurrentBranch returns the current branch name.
func GitCurrentBranch(t *testing.T, dir string) string {
	t.Helper()
	return RunGit(t, dir, "rev-parse", "--abbrev-ref", "HEAD")
}

// GitStatus returns the git status output.
func GitStatus(t *testing.T, dir string) string {
	t.Helper()
	return RunGit(t, dir, "status", "--porcelain")
}

// CreateCheckpoint creates a git checkpoint for testing.
func CreateCheckpoint(t *testing.T, repoDir, taskID, message string) string {
	t.Helper()

	// Make a change first
	testFile := filepath.Join(repoDir, "checkpoint-test.txt")
	if err := os.WriteFile(testFile, []byte(message), 0o644); err != nil {
		t.Fatalf("Write checkpoint test file: %v", err)
	}

	// Stage and commit with task ID prefix
	mustRunGit(t, repoDir, "add", ".")
	commitMsg := "[" + taskID + "] " + message
	mustRunGit(t, repoDir, "commit", "-m", commitMsg)

	// Get commit hash
	output := RunGit(t, repoDir, "rev-parse", "HEAD")

	return strings.TrimSpace(output)
}

// GetCheckpointCount returns the number of checkpoints (commits) for a task.
func GetCheckpointCount(t *testing.T, repoDir, taskID string) int {
	t.Helper()

	output := RunGit(t, repoDir, "log", "--oneline", "--grep=^\\["+taskID+"\\]")
	lines := strings.Split(strings.TrimSpace(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}

	return count
}

// AssertBranchExists fails the test if the branch doesn't exist.
func AssertBranchExists(t *testing.T, repoDir, branch string) {
	t.Helper()

	branches := RunGit(t, repoDir, "branch", "--list", branch)
	if strings.TrimSpace(branches) == "" {
		t.Errorf("branch %q does not exist", branch)
	}
}

// AssertBranchNotExists fails the test if the branch exists.
func AssertBranchNotExists(t *testing.T, repoDir, branch string) {
	t.Helper()

	branches := RunGit(t, repoDir, "branch", "--list", branch)
	if strings.TrimSpace(branches) != "" {
		t.Errorf("branch %q exists but should not", branch)
	}
}

// AssertWorktreeExists fails the test if the worktree doesn't exist.
func AssertWorktreeExists(t *testing.T, repoDir, worktreePath string) {
	t.Helper()

	worktrees := RunGit(t, repoDir, "worktree", "list", "--porcelain")
	if !strings.Contains(worktrees, worktreePath) {
		t.Errorf("worktree %q does not exist", worktreePath)
	}
}

// AssertWorktreeNotExists fails the test if the worktree exists.
func AssertWorktreeNotExists(t *testing.T, repoDir, worktreePath string) {
	t.Helper()

	worktrees := RunGit(t, repoDir, "worktree", "list", "--porcelain")
	if strings.Contains(worktrees, worktreePath) {
		t.Errorf("worktree %q exists but should not", worktreePath)
	}
}

// AssertCurrentBranch fails the test if current branch doesn't match.
func AssertCurrentBranch(t *testing.T, repoDir, expected string) {
	t.Helper()

	current := GetCurrentBranch(t, repoDir)
	if current != expected {
		t.Errorf("current branch = %q, want %q", current, expected)
	}
}

// AssertFileInCommit fails the test if the file is not in the given commit.
func AssertFileInCommit(t *testing.T, repoDir, commit, filePath string) {
	t.Helper()

	// Check if file exists in commit
	output := RunGit(t, repoDir, "ls-tree", commit, filePath)
	if strings.TrimSpace(output) == "" {
		t.Errorf("file %q not found in commit %s", filePath, commit)
	}
}

// GetCommitCount returns the total number of commits in the repo.
func GetCommitCount(t *testing.T, repoDir string) int {
	t.Helper()

	output := RunGit(t, repoDir, "rev-list", "--count", "HEAD")
	count := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(output), "%d", &count); err == nil {
		return count
	}

	return 0
}

// GitCreateBranchOnly creates a new branch at the current HEAD without switching.
func GitCreateBranchOnly(t *testing.T, repoDir, branch string) {
	t.Helper()
	mustRunGit(t, repoDir, "branch", branch)
}

// GitCreateAndCheckoutBranch creates and switches to a new branch.
func GitCreateAndCheckoutBranch(t *testing.T, repoDir, branch string) {
	t.Helper()
	mustRunGit(t, repoDir, "checkout", "-b", branch)
}

// GetCurrentBranch returns the current branch name.
func GetCurrentBranch(t *testing.T, repoDir string) string {
	t.Helper()

	return strings.TrimSpace(RunGit(t, repoDir, "branch", "--show-current"))
}

// GitGetBaseBranch returns the base branch of the repo.
func GitGetBaseBranch(t *testing.T, repoDir string) string {
	t.Helper()

	// Try to get the default branch name
	branches := RunGit(t, repoDir, "branch", "--format=%(refname:short)")
	lines := strings.Split(strings.TrimSpace(branches), "\n")

	// Common default branch names
	defaults := []string{"main", "master", "develop"}
	for _, def := range defaults {
		for _, line := range lines {
			if strings.TrimSpace(line) == def {
				return def
			}
		}
	}

	// Fall back to first branch
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "main"
}
