package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Branch struct {
	Raw     *BranchInfo
	project *Project
	gerrit  *Gerrit
	Base    string
}

// BranchInfo entity contains information about a branch.
type BranchInfo struct {
	Ref       string        `json:"ref"`
	Revision  string        `json:"revision"`
	CanDelete bool          `json:"can_delete"`
	WebLinks  []WebLinkInfo `json:"web_links,omitempty"`
}

// BranchInput entity contains information for the creation of a new branch.
type BranchInput struct {
	Ref      string `json:"ref,omitempty"`
	Revision string `json:"revision,omitempty"`
}

// DeleteBranchesInput entity contains information about branches that should be deleted.
type DeleteBranchesInput struct {
	Branches []string `json:"branches"`
}

// BranchOptions specifies the parameters to the branch API endpoints.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#branch-options
type BranchOptions struct {
	// Limit the number of branches to be included in the results.
	Limit int `url:"n,omitempty"`

	// Skip the given number of branches from the beginning of the list.
	Skip int `url:"S,omitempty"`

	// Substring limits the results to those projects that match the specified substring.
	Substring string `url:"m,omitempty"`

	// Limit the results to those branches that match the specified regex.
	// Boundary matchers '^' and '$' are implicit.
	// For example: the regex 't*' will match any branches that start with 'test' and regex '*t' will match any branches that end with 'test'.
	Regex string `url:"r,omitempty"`
}

type MergeOptions struct {
	Source         string `url:"source"`
	SourceBranch   string `url:"source_branch,omitempty"`
	Strategy       string `url:"strategy,omitempty"`
	AllowConflicts bool   `url:"allow_conflicts,omitempty"`
}

type BranchService struct {
	gerrit  *Gerrit
	project *Project
}

type IBranchService interface {
	List(ctx context.Context, opt *BranchOptions) (*[]BranchInfo, *http.Response, error)
	Get(ctx context.Context, branchID string) (*Branch, *http.Response, error)
	Create(ctx context.Context, branchID string, input *BranchInput) (*Branch, *http.Response, error)
	Delete(ctx context.Context, branchID string) (bool, *http.Response, error)
	BulkDelete(ctx context.Context, input *DeleteBranchesInput) (bool, *http.Response, error)
}

//func NewBranchService() IBranchService {
//        return &BranchService{}
//}

//func NewBranch(gerrit *Gerrit, projectName, branchID string) *Branch {
//        return &Branch{
//                Raw:     new(BranchInfo),
//                gerrit:  gerrit,
//                project: &Project{Base: projectName, gerrit: gerrit},
//                Base:    branchID,
//        }
//}

// List lists the branches of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-branches
func (s *BranchService) List(ctx context.Context, opt *BranchOptions) (*[]BranchInfo, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/branches/", url.QueryEscape(s.project.Base))

	v := &[]BranchInfo{}
	resp, err := s.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	return v, resp, err
}

// Get retrieves a branch of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-branch
func (s *BranchService) Get(ctx context.Context, branchID string) (*Branch, *http.Response, error) {
	branch := Branch{Raw: new(BranchInfo), gerrit: s.gerrit, project: s.project, Base: branchID}

	resp, err := branch.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return &branch, resp, nil
}

// Create creates a new branch.
// In the request body additional data for the branch can be provided as BranchInput.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-branch
func (s *BranchService) Create(ctx context.Context, branchID string, input *BranchInput) (*Branch, *http.Response, error) {
	obj := Branch{Raw: new(BranchInfo), gerrit: s.gerrit, project: s.project, Base: branchID}
	return obj.Create(ctx, input)
}

// Delete deletes a branch.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-branch
func (s *BranchService) Delete(ctx context.Context, branchID string) (bool, *http.Response, error) {
	obj := Branch{Raw: new(BranchInfo), gerrit: s.gerrit, project: s.project, Base: branchID}
	return obj.Delete(ctx)
}

// BulkDelete delete one or more branches.
// If some branches could not be deleted, the response is “409 Conflict” and the error message is contained in the response body.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-branches
func (s *BranchService) BulkDelete(ctx context.Context, input *DeleteBranchesInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/branches:delete", url.QueryEscape(s.project.Base))
	resp, err := s.gerrit.Requester.Call(ctx, "POST", u, input, nil)

	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}

func (b *Branch) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("projects/%s/branches/%s", url.QueryEscape(b.project.Base), url.QueryEscape(b.Base))
	return b.gerrit.Requester.Call(ctx, "GET", u, nil, b.Raw)
}

func (b *Branch) Create(ctx context.Context, input *BranchInput) (*Branch, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/branches/%s", url.QueryEscape(b.project.Base), url.QueryEscape(b.Base))
	resp, err := b.gerrit.Requester.Call(ctx, "PUT", u, input, nil)

	if err != nil {
		return nil, resp, err
	}

	resp, err = b.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return b, resp, nil
}

func (b *Branch) Delete(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/branches/%s", url.QueryEscape(b.project.Base), url.QueryEscape(b.Base))
	resp, err := b.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)

	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}

// GetContent gets the content of a file from the HEAD revision of a certain branch.
// The content is returned as base64 encoded string.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-content
func (b *Branch) GetContent(ctx context.Context, fileID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/branches/%s/files/%s/content",
		url.QueryEscape(b.project.Base),
		url.QueryEscape(b.Base),
		url.QueryEscape(fileID))

	resp, err := b.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// GetMergeableInformation Gets whether the source is mergeable with the target branch.
// The source query parameter is required, which can be anything that could be resolved to a commit,
// and is visible to the caller. See examples of the source attribute in MergeInput.
// Also takes an optional parameter strategy, which can be recursive, resolve, simple-two-way-in-core, ours or theirs,
// default will use project settings.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-mergeable-info
func (b *Branch) GetMergeableInformation(ctx context.Context, opt *MergeOptions) (*MergeableInfo, *http.Response, error) {
	v := new(MergeableInfo)
	u := fmt.Sprintf("projects/%s/branches/%s/mergeable",
		url.QueryEscape(b.project.Base),
		url.QueryEscape(b.Base))

	resp, err := b.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetReflog gets the reflog of a certain branch.
// The caller must be project owner.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-reflog
func (b *Branch) GetReflog(ctx context.Context) (*[]ReflogEntryInfo, *http.Response, error) {
	v := new([]ReflogEntryInfo)
	u := fmt.Sprintf("projects/%s/branches/%s/reflog",
		url.QueryEscape(b.project.Base),
		url.QueryEscape(b.Base))

	resp, err := b.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}