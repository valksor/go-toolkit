// Package snapshot provides types for capturing source content.
package snapshot

import "context"

// Snapshot contains captured source content (read-only copy).
type Snapshot struct {
	Type    string         // directory, file, or provider name
	Ref     string         // original reference
	Files   []SnapshotFile // for directories or multi-file snapshots
	Content string         // for single files
}

// SnapshotFile represents a single file in a snapshot.
type SnapshotFile struct {
	Path    string
	Content string
}

// Snapshotter captures source content for storage.
type Snapshotter interface {
	Snapshot(ctx context.Context, id string) (*Snapshot, error)
}
