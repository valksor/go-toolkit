// Package pullrequest provides types for VCS pull/merge requests.
package pullrequest

import "time"

// PullRequestOptions for creating a PR.
type PullRequestOptions struct {
	Title        string
	Body         string
	SourceBranch string
	TargetBranch string
	Labels       []string
	Reviewers    []string
	Draft        bool
}

// PullRequest represents a pull/merge request.
type PullRequest struct {
	ID         string
	URL        string
	Title      string
	State      string
	Number     int
	Body       string
	HeadSHA    string    // Commit SHA of the head branch
	HeadBranch string    // Name of the head branch
	BaseBranch string    // Name of the base branch
	Author     string    // Author username
	CreatedAt  time.Time // Creation time
	UpdatedAt  time.Time // Last update time
	Labels     []string  // PR labels
	Assignees  []string  // Assignee usernames
}

// PullRequestDiff contains PR diff information.
type PullRequestDiff struct {
	URL        string     // URL to view the diff
	BaseBranch string     // Base branch name
	HeadBranch string     // Head branch name
	Files      []FileDiff // Files changed
	Patch      string     // Full diff in unified format
	Additions  int        // Total lines added
	Deletions  int        // Total lines deleted
	Commits    int        // Number of commits
}

// FileDiff represents a single file's changes.
type FileDiff struct {
	Path      string // File path
	Mode      string // "added", "modified", "deleted", "renamed"
	Patch     string // Unified diff for this file
	Additions int    // Lines added
	Deletions int    // Lines deleted
}

// ReviewEvent represents the type of review action.
type ReviewEvent string

const (
	// ReviewEventApprove approves the PR.
	ReviewEventApprove ReviewEvent = "APPROVE"
	// ReviewEventRequestChanges requests changes to the PR.
	ReviewEventRequestChanges ReviewEvent = "REQUEST_CHANGES"
	// ReviewEventComment posts a comment without approval/rejection.
	ReviewEventComment ReviewEvent = "COMMENT"
)

// SubmitReviewOptions contains options for submitting a PR review.
type SubmitReviewOptions struct {
	PRNumber int             // PR/MR number
	Event    ReviewEvent     // APPROVE, REQUEST_CHANGES, COMMENT
	Summary  string          // Overall review summary
	Comments []ReviewComment // Per-line comments (optional)
}

// ReviewComment represents a per-line comment in a review.
type ReviewComment struct {
	Path      string // File path
	Line      int    // Line number in the diff
	Body      string // Comment body
	Side      string // "LEFT" (old) or "RIGHT" (new) for GitHub; ignored for GitLab
	StartLine int    // For multi-line comments (0 if single line)
}

// ReviewSubmission contains the result of submitting a review.
type ReviewSubmission struct {
	ID             string // Review/note ID
	URL            string // URL to the review
	CommentsPosted int    // Number of comments posted
}

// Note: PR comments use workunit.Comment - see interfaces.go
