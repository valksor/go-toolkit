// Package capability provides provider capability discovery.
package capability

import (
	"github.com/valksor/go-toolkit/pullrequest"
	"github.com/valksor/go-toolkit/snapshot"
	"github.com/valksor/go-toolkit/workunit"
)

// Capability identifies provider capabilities.
type Capability string

const (
	CapRead               Capability = "read"
	CapList               Capability = "list"
	CapDownloadAttachment Capability = "download_attachment"
	CapFetchComments      Capability = "fetch_comments"
	CapComment            Capability = "comment"
	CapUpdateStatus       Capability = "update_status"
	CapManageLabels       Capability = "manage_labels"
	CapSnapshot           Capability = "snapshot"
	CapCreatePR           Capability = "create_pr"
	CapLinkBranch         Capability = "link_branch"
	CapCreateWorkUnit     Capability = "create_work_unit"
	CapFetchSubtasks      Capability = "fetch_subtasks"
	CapFetchParent        Capability = "fetch_parent"
	CapFetchPR            Capability = "fetch_pr"
	CapPRComment          Capability = "pr_comment"
	CapFetchPRComments    Capability = "fetch_pr_comments"
	CapUpdatePRComment    Capability = "update_pr_comment"
	CapCreateDependency   Capability = "create_dependency"
	CapFetchDependencies  Capability = "fetch_dependencies"
	CapFetchProject       Capability = "fetch_project"
)

// CapabilitySet is a set of capabilities.
type CapabilitySet map[Capability]bool

// Has checks if capability is present.
func (cs CapabilitySet) Has(c Capability) bool {
	return cs[c]
}

// Infer uses type assertions to determine capabilities of a provider.
func Infer(p any) CapabilitySet {
	caps := make(CapabilitySet)

	if _, ok := p.(workunit.Reader); ok {
		caps[CapRead] = true
	}
	if _, ok := p.(workunit.Lister); ok {
		caps[CapList] = true
	}
	if _, ok := p.(workunit.AttachmentDownloader); ok {
		caps[CapDownloadAttachment] = true
	}
	if _, ok := p.(workunit.CommentFetcher); ok {
		caps[CapFetchComments] = true
	}
	if _, ok := p.(workunit.Commenter); ok {
		caps[CapComment] = true
	}
	if _, ok := p.(workunit.StatusUpdater); ok {
		caps[CapUpdateStatus] = true
	}
	if _, ok := p.(workunit.LabelManager); ok {
		caps[CapManageLabels] = true
	}
	if _, ok := p.(snapshot.Snapshotter); ok {
		caps[CapSnapshot] = true
	}
	if _, ok := p.(pullrequest.PRCreator); ok {
		caps[CapCreatePR] = true
	}
	if _, ok := p.(pullrequest.BranchLinker); ok {
		caps[CapLinkBranch] = true
	}
	if _, ok := p.(workunit.WorkUnitCreator); ok {
		caps[CapCreateWorkUnit] = true
	}
	if _, ok := p.(workunit.SubtaskFetcher); ok {
		caps[CapFetchSubtasks] = true
	}
	if _, ok := p.(workunit.ParentFetcher); ok {
		caps[CapFetchParent] = true
	}
	if _, ok := p.(pullrequest.PRFetcher); ok {
		caps[CapFetchPR] = true
	}
	if _, ok := p.(pullrequest.PRCommenter); ok {
		caps[CapPRComment] = true
	}
	if _, ok := p.(pullrequest.PRCommentFetcher); ok {
		caps[CapFetchPRComments] = true
	}
	if _, ok := p.(pullrequest.PRCommentUpdater); ok {
		caps[CapUpdatePRComment] = true
	}
	if _, ok := p.(workunit.DependencyCreator); ok {
		caps[CapCreateDependency] = true
	}
	if _, ok := p.(workunit.DependencyFetcher); ok {
		caps[CapFetchDependencies] = true
	}
	if _, ok := p.(workunit.ProjectFetcher); ok {
		caps[CapFetchProject] = true
	}

	return caps
}
