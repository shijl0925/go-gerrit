package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// EditInfo entity contains information about a change edit.
type EditInfo struct {
	Commit             CommitInfo           `json:"commit"`
	BasePatchSetNumber int                  `json:"base_patch_set_number"`
	Ref                string               `json:"ref"`
	BaseRevision       string               `json:"base_revision"`
	Fetch              map[string]FetchInfo `json:"fetch"`
	Files              map[string]FileInfo  `json:"files,omitempty"`
}

// EditFileInfo entity contains additional information of a file within a change edit.
type EditFileInfo struct {
	WebLinks []WebLinkInfo `json:"web_links,omitempty"`
}

// ChangeEditDetailOptions specifies the parameters to the ChangesService.GetChangeEditDetails.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-detail
type ChangeEditDetailOptions struct {
	// When request parameter list is provided the response also includes the file list.
	List bool `url:"list,omitempty"`
	// When base request parameter is provided the file list is computed against this base revision.
	Base bool `url:"base,omitempty"`
	// When request parameter download-commands is provided fetch info map is also included.
	DownloadCommands bool `url:"download-commands,omitempty"`
}

// GetEditDetails Retrieves the details of the change edit done by the caller to the given change.
// As response an EditInfo entity is returned that describes the change edit, or “204 No Content” when change edit doesn’t exist for this change.
// Change edits are stored on special branches and there can be max one edit per user per change.
// Edits aren’t tracked in the database.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-detail
func (c *Change) GetEditDetails(ctx context.Context, opt *ChangeEditDetailOptions) (*EditInfo, *http.Response, error) {
	v := new(EditInfo)
	u := fmt.Sprintf("changes/%s/edit", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ChangeFileContentInChangeEdit put content of a file to a change edit.
//
// When change edit doesn’t exist for this change yet it is created.
// When file content isn’t provided, it is wiped out for that file.
// As response “204 No Content” is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#put-edit-file
func (c *Change) ChangeFileContentInChangeEdit(ctx context.Context, filePath, content string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit/%s", c.Base, url.QueryEscape(filePath))
	return c.gerrit.Requester.Call(ctx, "PUT", u, content, nil)
}

// RestoreChangeEdit restores file content or renames files in change edit.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#post-edit
func (c *Change) RestoreChangeEdit(ctx context.Context, input *RestoreChangeEditInput) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit", c.Base)
	return c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}

// RenameChangeEdit renames files in change edit.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#post-edit
func (c *Change) RenameChangeEdit(ctx context.Context, input *RenameChangeEditInput) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit", c.Base)
	return c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}

// RetrieveCommitMessageFromChangeEdit retrieves commit message from change edit.
// The commit message is returned as base64 encoded string.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-message
func (c *Change) RetrieveCommitMessageFromChangeEdit(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/edit:message", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// ChangeCommitMessageInChangeEdit modify commit message.
// The request body needs to include a ChangeEditMessageInput entity.
// If a change edit doesn’t exist for this change yet, it is created.
// As response “204 No Content” is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#put-change-edit-message
func (c *Change) ChangeCommitMessageInChangeEdit(ctx context.Context, input *ChangeEditMessageInput) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit:message", c.Base)
	return c.gerrit.Requester.Call(ctx, "PUT", u, input, nil)
}

// DeleteFileInChangeEdit deletes a file from a change edit.
// This deletes the file from the repository completely.
// This is not the same as reverting or restoring a file to its previous contents.
// When change edit doesn’t exist for this change yet it is created.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-edit-file
func (c *Change) DeleteFileInChangeEdit(ctx context.Context, filePath string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit/%s", c.Base, url.QueryEscape(filePath))
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// RetrieveFileContentFromChangeEdit retrieves content of a file from a change edit.
//
// The content of the file is returned as text encoded inside base64.
// The Content-Type header will always be text/plain reflecting the outer base64 encoding.
// A Gerrit-specific X-FYI-Content-Type header can be examined to find the server detected content type of the file.
//
// When the specified file was deleted in the change edit “204 No Content” is returned.
// If only the content type is required, callers should use HEAD to avoid downloading the encoded file contents.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-file
func (c *Change) RetrieveFileContentFromChangeEdit(ctx context.Context, filePath string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/edit/%s", c.Base, url.QueryEscape(filePath))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// RetrieveFileMetaFromChangeEdit retrieves meta data of a file from a change edit.
// Currently only web links are returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-meta-data
func (c *Change) RetrieveFileMetaFromChangeEdit(ctx context.Context, filePath string) (*EditFileInfo, *http.Response, error) {
	v := new(EditFileInfo)
	u := fmt.Sprintf("changes/%s/edit/%s/meta", c.Base, url.QueryEscape(filePath))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// PublishChangeEdit promotes change edit to a regular patch set.
// As response “204 No Content” is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#publish-edit
func (c *Change) PublishChangeEdit(ctx context.Context, input *PublishChangeEditInput) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit:publish", c.Base)
	return c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}

// RebaseChangeEdit rebase change edit on top of latest patch set.
// When change was rebased on top of latest patch set, response “204 No Content” is returned.
// When change edit is already based on top of the latest patch set, the response “409 Conflict” is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#rebase-edit
func (c *Change) RebaseChangeEdit(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit:rebase", c.Base)
	return c.gerrit.Requester.Call(ctx, "POST", u, nil, nil)
}

// DeleteChangeEdit deletes change edit.
// As response “204 No Content” is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-edit
func (c *Change) DeleteChangeEdit(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/edit", c.Base)
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}