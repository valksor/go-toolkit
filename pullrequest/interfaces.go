package pullrequest

import (
	"context"

	"github.com/valksor/go-toolkit/workunit"
)

// PRCreator creates pull requests (for GitHub-like providers).
type PRCreator interface {
	CreatePullRequest(ctx context.Context, opts PullRequestOptions) (*PullRequest, error)
}

// PRFetcher retrieves pull request details and diffs.
type PRFetcher interface {
	FetchPullRequest(ctx context.Context, number int) (*PullRequest, error)
	FetchPullRequestDiff(ctx context.Context, number int) (*PullRequestDiff, error)
}

// PRCommenter posts comments to pull requests.
type PRCommenter interface {
	AddPullRequestComment(ctx context.Context, number int, body string) (*workunit.Comment, error)
}

// PRCommentFetcher retrieves existing comments from a PR/MR.
type PRCommentFetcher interface {
	FetchPullRequestComments(ctx context.Context, number int) ([]workunit.Comment, error)
}

// PRCommentUpdater updates existing comments on a PR/MR.
type PRCommentUpdater interface {
	UpdatePullRequestComment(ctx context.Context, number int, commentID string, body string) (*workunit.Comment, error)
}

// PRReviewer submits formal reviews to pull requests.
// For GitHub: Creates a PR review with APPROVED/REQUEST_CHANGES.
// For GitLab: Creates MR notes (no formal review API).
type PRReviewer interface {
	SubmitReview(ctx context.Context, opts SubmitReviewOptions) (*ReviewSubmission, error)
}

// BranchLinker links work units to git branches.
type BranchLinker interface {
	LinkBranch(ctx context.Context, workUnitID, branch string) error
	UnlinkBranch(ctx context.Context, workUnitID, branch string) error
	GetLinkedBranch(ctx context.Context, workUnitID string) (string, error)
}
