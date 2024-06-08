package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Commit struct {
	Raw     *CommitInfo
	project *Project
	gerrit  *Gerrit
	Base    string
}

type CommitService struct {
	gerrit  *Gerrit
	project *Project
}

// Get retrieves a commit of a project.
// The commit must be visible to the caller.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-commit
func (s *CommitService) Get(ctx context.Context, commitID string) (*Commit, *http.Response, error) {
	commit := Commit{Raw: new(CommitInfo), gerrit: s.gerrit, project: s.project, Base: commitID}
	resp, err := commit.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return &commit, resp, nil
}

func (c *Commit) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("projects/%s/commits/%s", url.QueryEscape(c.project.Base), c.Base)
	return c.gerrit.Requester.Call(ctx, "GET", u, nil, c.Raw)
}

// GetIncludeIn Retrieves the branches and tags in which a change is included.
// Branches that are not visible to the calling user according to the project’s read permissions are filtered out from the result.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-included-in
func (c *Commit) GetIncludeIn(ctx context.Context) (*IncludedInInfo, *http.Response, error) {
	v := new(IncludedInInfo)
	u := fmt.Sprintf("projects/%s/commits/%s/in", url.QueryEscape(c.project.Base), c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetContent gets the content of a file from a certain commit.
// The content is returned as base64 encoded string.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-content-from-commit
func (c *Commit) GetContent(ctx context.Context, fileID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/commits/%s/files/%s/content",
		url.QueryEscape(c.project.Base),
		c.Base,
		url.QueryEscape(fileID))

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// ListFiles gets the files that were modified, added or deleted in a commit.
// As result a map is returned that maps the file path to a FileInfo entry. The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-files
func (c *Commit) ListFiles(ctx context.Context) (map[string]FileInfo, *http.Response, error) {
	v := make(map[string]FileInfo)
	u := fmt.Sprintf("projects/%s/commits/%s/files/", url.QueryEscape(c.project.Base), c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}