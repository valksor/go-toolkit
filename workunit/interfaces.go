package workunit

import (
	"context"
	"io"
)

// Reader fetches work units from a provider.
type Reader interface {
	Fetch(ctx context.Context, id string) (*WorkUnit, error)
}

// Identifier parses and validates references.
type Identifier interface {
	Parse(input string) (string, error)
	Match(input string) bool
}

// Lister enumerates work units.
type Lister interface {
	List(ctx context.Context, opts ListOptions) ([]*WorkUnit, error)
}

// AttachmentDownloader downloads attachments.
type AttachmentDownloader interface {
	DownloadAttachment(ctx context.Context, workUnitID, attachmentID string) (io.ReadCloser, error)
}

// CommentFetcher retrieves comments.
type CommentFetcher interface {
	FetchComments(ctx context.Context, workUnitID string) ([]Comment, error)
}

// Commenter adds comments to work units.
type Commenter interface {
	AddComment(ctx context.Context, workUnitID string, body string) (*Comment, error)
}

// StatusUpdater changes work unit status.
type StatusUpdater interface {
	UpdateStatus(ctx context.Context, workUnitID string, status Status) error
}

// LabelManager manages labels on work units.
type LabelManager interface {
	AddLabels(ctx context.Context, workUnitID string, labels []string) error
	RemoveLabels(ctx context.Context, workUnitID string, labels []string) error
}

// ReadOnlyProvider is the minimum interface for a provider.
type ReadOnlyProvider interface {
	Reader
	Identifier
}

// BidirectionalProvider supports read and write operations.
type BidirectionalProvider interface {
	Reader
	Identifier
	Commenter
	StatusUpdater
}

// WorkUnitCreator creates new work units (for Wrike, GitHub issues, etc.)
type WorkUnitCreator interface {
	CreateWorkUnit(ctx context.Context, opts CreateWorkUnitOptions) (*WorkUnit, error)
}

// SubtaskFetcher retrieves subtasks for a work unit.
type SubtaskFetcher interface {
	FetchSubtasks(ctx context.Context, workUnitID string) ([]*WorkUnit, error)
}

// ParentFetcher retrieves the parent task for a work unit.
// This is useful when working on a subtask and needing parent context.
type ParentFetcher interface {
	FetchParent(ctx context.Context, workUnitID string) (*WorkUnit, error)
}

// ProjectFetcher retrieves entire project/epic structures with all tasks.
// Different providers implement this differently:
// - Wrike: Fetch folder/project with all descendants
// - Jira: Fetch epic with all child stories
// - GitHub: Fetch project board with all cards (if API available)
// - Asana: Fetch portfolio/project with tasks.
type ProjectFetcher interface {
	FetchProject(ctx context.Context, reference string) (*ProjectStructure, error)
}

// DependencyCreator creates dependencies between work units.
type DependencyCreator interface {
	// CreateDependency creates a dependency where predecessorID must complete before successorID.
	CreateDependency(ctx context.Context, predecessorID, successorID string) error
}

// DependencyFetcher retrieves dependencies for a work unit.
type DependencyFetcher interface {
	// GetDependencies returns the IDs of work units that the given work unit depends on.
	GetDependencies(ctx context.Context, workUnitID string) ([]string, error)
}
