// Package workunit provides types for representing tasks from various work management systems.
package workunit

import "time"

// WorkUnit represents a task from any provider.
type WorkUnit struct {
	ID          string
	ExternalID  string // Provider-specific ID
	Provider    string // Provider name
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Labels      []string
	Assignees   []Person
	Comments    []Comment
	Attachments []Attachment
	Subtasks    []string
	Metadata    map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Source      SourceInfo

	// Naming fields for branch/commit customization
	ExternalKey string // User-facing key (e.g., "FEATURE-123") for branches/commits
	TaskType    string // Task type (e.g., "feature", "fix", "task")
	Slug        string // URL-safe title slug for branch names

	// Agent configuration from task source
	AgentConfig *AgentConfig // Per-task agent configuration (optional)

	// Budget configuration from task source (optional)
	Budget *BudgetConfig
}

// SourceInfo tracks where the work unit came from.
type SourceInfo struct {
	Type      string    // Provider type
	Reference string    // Original reference
	SyncedAt  time.Time // Last sync time
}

// StepAgentConfig holds agent configuration for a specific workflow step.
type StepAgentConfig struct {
	Name string            // Agent name or alias
	Env  map[string]string // Step-specific env vars
	Args []string          // Step-specific CLI args
}

// AgentConfig holds per-task agent configuration from the task source.
type AgentConfig struct {
	Name  string                     // Agent name or alias (e.g., "glm", "claude")
	Env   map[string]string          // Inline environment variables
	Args  []string                   // CLI arguments
	Steps map[string]StepAgentConfig // Per-step agent overrides
}

// BudgetConfig defines cost/token budgets for a task.
type BudgetConfig struct {
	MaxTokens int     // Maximum total tokens for the task (0 = unlimited)
	MaxCost   float64 // Maximum total cost for the task (0 = unlimited)
	Currency  string  // Currency code (e.g., "USD")
	OnLimit   string  // warn | pause | stop
	WarningAt float64 // Warning threshold (0-1, e.g., 0.8)
}

// Status represents work unit status.
type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusReview     Status = "review"
	StatusDone       Status = "done"
	StatusClosed     Status = "closed"
)

// Priority represents work unit priority.
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// String returns priority as string.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "normal"
	}
}

// Person represents a user/assignee.
// ID is the provider-specific unique identifier (e.g., GitHub login, GitLab username).
// It is used for matching in PR comments and other provider interactions.
// Name is the display name.
// Email is the email address (optional, may not be available from all providers).
type Person struct {
	ID    string // Provider username/login (used for matching in PR comments)
	Name  string // Display name
	Email string
}

// PersonNames extracts names from a Person slice.
// If a person has a name, that is used; otherwise falls back to their ID.
// Duplicate persons (by ID) are deduplicated in the result.
func PersonNames(persons []Person) []string {
	if len(persons) == 0 {
		return []string{}
	}

	// Deduplicate by ID while preserving order
	seen := make(map[string]bool, len(persons))
	var names []string
	for _, p := range persons {
		if !seen[p.ID] {
			seen[p.ID] = true
			if p.Name != "" {
				names = append(names, p.Name)
			} else {
				names = append(names, p.ID)
			}
		}
	}

	return names
}

// Comment represents a comment on a work unit.
type Comment struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    Person
	ID        string
	Body      string
}

// Attachment represents a file attachment.
type Attachment struct {
	CreatedAt   time.Time
	ID          string
	Name        string
	URL         string
	ContentType string
	Size        int64
}

// ListOptions configures list operations.
type ListOptions struct {
	Status   Status
	Labels   []string
	Limit    int
	Offset   int
	OrderBy  string
	OrderDir string // asc, desc
}

// CreateWorkUnitOptions for creating a work unit.
type CreateWorkUnitOptions struct {
	CustomFields  map[string]any
	Title         string
	Description   string
	ParentID      string
	Labels        []string
	Assignees     []string
	Priority      Priority
	DependencyIDs []string // Work unit IDs this unit depends on (predecessors)
}

// ProjectStructure represents a complete project/epic with hierarchical tasks.
type ProjectStructure struct {
	ID          string         // Provider project ID
	Title       string         // Project/epic title
	Description string         // Optional description
	Source      string         // Provider name
	URL         string         // Provider URL
	Tasks       []*ProjectTask // All tasks in the project (flat list)
	Metadata    map[string]any // Provider-specific metadata
}

// ProjectTask represents a task within a project structure.
type ProjectTask struct {
	*WorkUnit

	ParentID string // Parent task ID (if nested)
	Depth    int    // Hierarchy depth (0 = top level)
	Position int    // Position within parent
}
