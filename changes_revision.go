package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DiffInfo entity contains information about the diff of a file in a revision.
type DiffInfo struct {
	MetaA           DiffFileMetaInfo  `json:"meta_a,omitempty"`
	MetaB           DiffFileMetaInfo  `json:"meta_b,omitempty"`
	ChangeType      string            `json:"change_type"`
	IntralineStatus string            `json:"intraline_status,omitempty"`
	DiffHeader      []string          `json:"diff_header"`
	Content         []DiffContent     `json:"content"`
	WebLinks        []DiffWebLinkInfo `json:"web_links,omitempty"`
	Binary          bool              `json:"binary,omitempty"`
}

// BlameInfo entity stores the commit metadata with the row coordinates where it applies.
type BlameInfo struct {
	Author    string      `json:"author"`
	ID        string      `json:"id"`
	Time      int         `json:"time"`
	CommitMsg string      `json:"commit_msg"`
	Ranges    []RangeInfo `json:"ranges"`
}

// RangeInfo entity stores the coordinates of a range.
type RangeInfo struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// RelatedChangesInfo entity contains information about related changes.
type RelatedChangesInfo struct {
	Changes []RelatedChangeAndCommitInfo `json:"changes"`
}

// FileInfo entity contains information about a file in a patch set.
type FileInfo struct {
	Status        string `json:"status,omitempty"`
	Binary        bool   `json:"binary,omitempty"`
	OldPath       string `json:"old_path,omitempty"`
	LinesInserted int    `json:"lines_inserted,omitempty"`
	LinesDeleted  int    `json:"lines_deleted,omitempty"`
	SizeDelta     int    `json:"size_delta"`
	Size          int    `json:"size"`
}

// ActionInfo entity describes a REST API call the client can make to manipulate a resource.
// These are frequently implemented by plugins and may be discovered at runtime.
type ActionInfo struct {
	Method  string `json:"method,omitempty"`
	Label   string `json:"label,omitempty"`
	Title   string `json:"title,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// CommitInfo entity contains information about a commit.
type CommitInfo struct {
	Commit    string        `json:"commit,omitempty"`
	Parents   []CommitInfo  `json:"parents"`
	Author    GitPersonInfo `json:"author"`
	Committer GitPersonInfo `json:"committer"`
	Subject   string        `json:"subject"`
	Message   string        `json:"message"`
	WebLinks  []WebLinkInfo `json:"web_links,omitempty"`
}

// MergeableInfo entity contains information about the mergeability of a change.
type MergeableInfo struct {
	SubmitType    string   `json:"submit_type"`
	Mergeable     bool     `json:"mergeable"`
	MergeableInto []string `json:"mergeable_into,omitempty"`
}

// DiffOptions specifies the parameters for GetDiff call.
type DiffOptions struct {
	// If the intraline parameter is specified, intraline differences are included in the diff.
	Intraline bool `url:"intraline,omitempty"`

	// The base parameter can be specified to control the base patch set from which the diff
	// should be generated.
	Base string `url:"base,omitempty"`

	// The integer-valued request parameter parent can be specified to control the parent commit number
	// against which the diff should be generated. This is useful for supporting review of merge commits.
	// The value is the 1-based index of the parent’s position in the commit object.
	Parent int `url:"parent,omitempty"`

	// If the weblinks-only parameter is specified, only the diff web links are returned.
	WeblinksOnly bool `url:"weblinks-only,omitempty"`

	// The ignore-whitespace parameter can be specified to control how whitespace differences are reported in the result. Valid values are NONE, TRAILING, CHANGED or ALL.
	IgnoreWhitespace string `url:"ignore-whitespace,omitempty"`

	// The context parameter can be specified to control the number of lines of surrounding context in the diff.
	// Valid values are ALL or number of lines.
	Context string `url:"context,omitempty"`
}

// CommitOptions specifies the parameters for GetCommit call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-commit
type CommitOptions struct {
	// Adding query parameter links (for example /changes/.../commit?links) returns a CommitInfo with the additional field web_links.
	Weblinks bool `url:"links,omitempty"`
}

// DescriptionInput entity contains information for setting a description.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#description-input
type DescriptionInput struct {
	Description string `json:"description"`
}

// MergableOptions specifies the parameters for GetMergable call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-mergeable
type MergableOptions struct {
	// If the other-branches parameter is specified, the mergeability will also be checked for all other branches.
	OtherBranches bool `url:"other-branches,omitempty"`
}

// FilesOptions specifies the parameters for ListFiles and ListFilesReviewed calls.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-files
type FilesOptions struct {
	// The request parameter q changes the response to return a list of all files (modified or unmodified)
	// that contain that substring in the path name. This is useful to implement suggestion services
	// finding a file by partial name.
	Q string `url:"q,omitempty"`

	// The base parameter can be specified to control the base patch set from which the list of files
	// should be generated.
	//
	// Note: This option is undocumented.
	Base string `url:"base,omitempty"`

	// The integer-valued request parameter parent changes the response to return a list of the files
	// which are different in this commit compared to the given parent commit. This is useful for
	// supporting review of merge commits. The value is the 1-based index of the parent’s position
	// in the commit object.
	Parent int `url:"parent,omitempty"`
}

// PatchOptions specifies the parameters for GetPatch call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-patch
type PatchOptions struct {
	// Adding query parameter zip (for example /changes/.../patch?zip) returns the patch as a single file inside of a ZIP archive.
	// Clients can expand the ZIP to obtain the plain text patch, avoiding the need for a base64 decoding step.
	// This option implies download.
	Zip bool `url:"zip,omitempty"`

	// Query parameter download (e.g. /changes/.../patch?download) will suggest the browser save the patch as commitsha1.diff.base64, for later processing by command line tools.
	Download bool `url:"download,omitempty"`

	// If the path parameter is set, the returned content is a diff of the single file that the path refers to.
	Path string `url:"path,omitempty"`
}

// GetRevisionCommit retrieves a parsed commit of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-commit
func (c *Change) GetRevisionCommit(ctx context.Context, revisionID string, opt *CommitOptions) (*CommitInfo, *http.Response, error) {
	//        if reflect.TypeOf(revisionID).String() == "int" {
	//                revisionID = strconv.Itoa(revisionID.(int))
	//        }
	//revisionID = revisionID.(string)

	v := new(CommitInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/commit", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionDescription retrieves the description of a patch set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-description
func (c *Change) GetRevisionDescription(ctx context.Context, revisionID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/revisions/%s/description", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetRevisionDescription sets the description of a patch set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-description
func (c *Change) SetRevisionDescription(ctx context.Context, revisionID string, input *DescriptionInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/revisions/%s/description", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "PUT", u, input, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// GetRevisionMergeDiff Returns the list of commits that are being integrated into a target branch by a merge commit.
//
//        By default the first parent is assumed to be uninteresting. By using the parent option another parent can be set
//        as uninteresting (parents are 1-based).
//        The list of commits is returned as a list of CommitInfo entities.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-merge-list
func (c *Change) GetRevisionMergeDiff(ctx context.Context, revisionID string) (*[]CommitInfo, *http.Response, error) {
	v := new([]CommitInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/mergelist", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionActions retrieves revision actions of the revision of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-revision-actions
func (c *Change) GetRevisionActions(ctx context.Context, revisionID string) (map[string]ActionInfo, *http.Response, error) {
	v := make(map[string]ActionInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/actions", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionReview retrieves a review of a revision.
//
// As response a ChangeInfo entity with detailed labels and detailed accounts is returned that describes the review of the revision.
// The revision for which the review is retrieved is contained in the revisions field.
// In addition the current_revision field is set if the revision for which the review is retrieved is the current revision of the change.
// Please note that the returned labels are always for the current patch set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-review
func (c *Change) GetRevisionReview(ctx context.Context, revisionID string) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/review", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetRevisionReview sets a review on a revision.
// The review must be provided in the request body as a ReviewInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-review
func (c *Change) SetRevisionReview(ctx context.Context, revisionID string, input *ReviewInput) (*ReviewResult, *http.Response, error) {
	v := new(ReviewResult)
	u := fmt.Sprintf("changes/%s/revisions/%s/review", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionRelatedChanges retrieves related changes of a revision.
// Related changes are changes that either depend on, or are dependencies of the revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-related-changes
func (c *Change) GetRevisionRelatedChanges(ctx context.Context, revisionID string) (*RelatedChangesInfo, *http.Response, error) {
	v := new(RelatedChangesInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/related", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// RebaseRevision rebase a revision.
// The review must be provided in the request body as a RebaseInput entity.
// If the revision cannot be rebased, e.g. due to conflicts, the response is “409 Conflict”
// and the error message is contained in the response body.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#rebase-revision
func (c *Change) RebaseRevision(ctx context.Context, revisionID string, input *RebaseInput) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/rebase", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SubmitRevision Submit a revision. As response a ChangeInfo entity is returned that describes the submitted/merged change.
// As response a ChangeInfo entity is returned that describes the submitted/merged change.
// If the revision cannot be submitted, e.g. because the submit rule doesn’t allow submitting the revision or
// the revision is not the current revision, the response is “409 Conflict” and the error message is contained in the response body.
// and the error message is contained in the response body.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#submit-revision
func (c *Change) SubmitRevision(ctx context.Context, revisionID string) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/submit", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionPatch Gets the formatted patch for one revision.
// The formatted patch is returned as text encoded inside base64.
// Adding query parameter zip (for example /changes/.../patch?zip) returns the patch as a single file inside of a ZIP archive.
// Clients can expand the ZIP to obtain the plain text patch, avoiding the need for a base64 decoding step.
// This option implies download.
//
// Query parameter download (e.g. /changes/.../patch?download) will suggest the browser save the patch as commitsha1.diff.base64, for later processing by command line tools.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-patch
func (c *Change) GetRevisionPatch(ctx context.Context, revisionID string, opt *PatchOptions) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/patch", c.Base, revisionID)
	return c.gerrit.Requester.Call(ctx, "GET", u, opt, nil)
}

// GetRevisionMergeable gets the method the server will use to submit (merge) the change and an indicator if the change is currently mergeable.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-mergeable
func (c *Change) GetRevisionMergeable(ctx context.Context, revisionID string, opt *MergableOptions) (*MergeableInfo, *http.Response, error) {
	v := new(MergeableInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/mergeable", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionSubmitType gets the method the server will use to submit (merge) the change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-submit-type
func (c *Change) GetRevisionSubmitType(ctx context.Context, revisionID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/revisions/%s/submit_type", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// TestRevisionSubmitType tests the submit_type Prolog rule in the project, or the one given.
//
// Request body may be either the Prolog code as text/plain or a RuleInput object.
// The query parameter filters may be set to SKIP to bypass parent project filters while testing a project-specific rule.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#test-submit-type
func (c *Change) TestRevisionSubmitType(ctx context.Context, revisionID string, input *RuleInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/revisions/%s/test.submit_type", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// TestRevisionSubmitRule tests the submit_rule Prolog rule in the project, or the one given.
//
// Request body may be either the Prolog code as text/plain or a RuleInput object.
// The query parameter filters may be set to SKIP to bypass parent project filters while testing a project-specific rule.
//
// The response is a list of SubmitRecord entries describing the permutations that satisfy the tested submit rule.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#test-submit-rule
func (c *Change) TestRevisionSubmitRule(ctx context.Context, revisionID string, input *RuleInput) (*[]SubmitRecord, *http.Response, error) {
	v := new([]SubmitRecord)
	u := fmt.Sprintf("changes/%s/revisions/%s/test.submit_rule", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListRevisionDrafts lists the draft comments of a revision that belong to the calling user.
// Returns a map of file paths to lists of CommentInfo entries.
// The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-drafts
func (c *Change) ListRevisionDrafts(ctx context.Context, revisionID string) (map[string][]CommentInfo, *http.Response, error) {
	v := make(map[string][]CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// CreateRevisionDraft creates a draft comment on a revision.
// The new draft comment must be provided in the request body inside a CommentInput entity.
//
// As response a CommentInfo entity is returned that describes the draft comment.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#create-draft
func (c *Change) CreateRevisionDraft(ctx context.Context, revisionID string, input *CommentInput) (*CommentInfo, *http.Response, error) {
	v := new(CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionDraft retrieves a draft comment of a revision that belongs to the calling user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-draft
func (c *Change) GetRevisionDraft(ctx context.Context, revisionID, draftID string) (*CommentInfo, *http.Response, error) {
	v := new(CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/%s", c.Base, revisionID, draftID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// UpdateRevisionDraft updates a draft comment on a revision.
// The new draft comment must be provided in the request body inside a CommentInput entity.
//
// As response a CommentInfo entity is returned that describes the draft comment.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#update-draft
func (c *Change) UpdateRevisionDraft(ctx context.Context, revisionID, draftID string, input *CommentInput) (*CommentInfo, *http.Response, error) {
	v := new(CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/%s", c.Base, revisionID, draftID)

	resp, err := c.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteRevisionDraft deletes a draft comment from a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-draft
func (c *Change) DeleteRevisionDraft(ctx context.Context, revisionID, draftID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/%s", c.Base, revisionID, draftID)
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// ListRevisionComments lists the published comments of a revision.
// As result a map is returned that maps the file path to a list of CommentInfo entries.
// The entries in the map are sorted by file path and only include file (or inline) comments.
// Use the Get Change Detail endpoint to retrieve the general change message (or comment).
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-comments
func (c *Change) ListRevisionComments(ctx context.Context, revisionID string) (map[string][]CommentInfo, *http.Response, error) {
	v := make(map[string][]CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/comments/", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionComment retrieves a published comment of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-comment
func (c *Change) GetRevisionComment(ctx context.Context, revisionID, commentID string) (*CommentInfo, *http.Response, error) {
	v := new(CommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/comments/%s", c.Base, revisionID, commentID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteRevisionComment deletes a published comment of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-comment
func (c *Change) DeleteRevisionComment(ctx context.Context, revisionID, commentID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/comments/%s", c.Base, revisionID, commentID)
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// ListRevisionRobotComments Lists the robot comments of a revision.
// Return a map that maps the file path to a list of RobotCommentInfo entries. The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-robot-comments
func (c *Change) ListRevisionRobotComments(ctx context.Context, revisionID string) (map[string][]RobotCommentInfo, *http.Response, error) {
	v := make(map[string][]RobotCommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/robotcomments/", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionRobotComments retrieves a robot comment of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-robot-comment
func (c *Change) GetRevisionRobotComments(ctx context.Context, revisionID, commentID string) (*RobotCommentInfo, *http.Response, error) {
	v := new(RobotCommentInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/robotcomments/%s", c.Base, revisionID, commentID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListRevisionFiles lists the files that were modified, added or deleted in a revision.
// As result a map is returned that maps the file path to a list of FileInfo entries.
// The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-files
func (c *Change) ListRevisionFiles(ctx context.Context, revisionID string, opt *FilesOptions) (map[string]FileInfo, *http.Response, error) {
	v := make(map[string]FileInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/files/", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetRevisionFileContent gets the content of a file from a certain revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-content
func (c *Change) GetRevisionFileContent(ctx context.Context, revisionID, fileID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/content", c.Base, revisionID, url.PathEscape(fileID))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// GetRevisionFileContentType gets the content type of a file from a certain revision.
// This is nearly the same as GetContent.
// But if only the content type is required, callers should use HEAD to avoid downloading the encoded file contents.
//
// For further documentation see GetContent.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-content
func (c *Change) GetRevisionFileContentType(ctx context.Context, revisionID, fileID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/content", c.Base, revisionID, url.PathEscape(fileID))
	return c.gerrit.Requester.Call(ctx, "HEAD", u, nil, nil)
}

func (c *Change) DownloadRevisionFileContent(ctx context.Context, revisionID, fileID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/download", c.Base, revisionID, url.PathEscape(fileID))
	return c.gerrit.Requester.Call(ctx, "GET", u, nil, nil)
}

// GetRevisionFileDiff gets the diff of a file from a certain revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-diff
func (c *Change) GetRevisionFileDiff(ctx context.Context, revisionID, fileID string, opt *DiffOptions) (*DiffInfo, *http.Response, error) {
	v := new(DiffInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/diff", c.Base, revisionID, url.PathEscape(fileID))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

func (c *Change) GetRevisionFileBlame(ctx context.Context, revisionID, fileID string) (*[]BlameInfo, *http.Response, error) {
	v := new([]BlameInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/blame", c.Base, revisionID, url.PathEscape(fileID))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListRevisionFilesReviewed lists the files that were modified, added or deleted in a revision.
// Unlike ListFiles, the response of ListFilesReviewed is a list of the paths the caller
// has marked as reviewed. Clients that also need the FileInfo should make two requests.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-files
func (c *Change) ListRevisionFilesReviewed(ctx context.Context, revisionID string, opt *FilesOptions) ([]string, *http.Response, error) {
	v := new([]string)
	u := fmt.Sprintf("changes/%s/revisions/%s/files/", c.Base, revisionID)

	o := struct {
		FilesOptions

		// The request parameter reviewed changes the response to return a list of the paths the caller has marked as reviewed.
		Reviewed bool `url:"reviewed,omitempty"`
	}{
		Reviewed: true,
	}
	if opt != nil {
		o.FilesOptions = *opt
	}
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, o, v)

	if err != nil {
		return nil, resp, err
	}
	return *v, resp, nil
}

// SetRevisionFileReviewed marks a file of a revision as reviewed by the calling user.
//
// If the file was already marked as reviewed by the calling user the response is “200 OK”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-reviewed
func (c *Change) SetRevisionFileReviewed(ctx context.Context, revisionID, fileID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/reviewed", c.Base, revisionID, url.PathEscape(fileID))
	return c.gerrit.Requester.Call(ctx, "PUT", u, nil, nil)
}

// DeleteRevisionFileReviewed deletes the reviewed flag of the calling user from a file of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-reviewed
func (c *Change) DeleteRevisionFileReviewed(ctx context.Context, revisionID, fileID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/reviewed", c.Base, revisionID, url.PathEscape(fileID))
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// CherryPickRevision publishes a draft revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#cherry-pick
func (c *Change) CherryPickRevision(ctx context.Context, revisionID string, input *CherryPickInput) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/revisions/%s/cherrypick", c.Base, revisionID)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}