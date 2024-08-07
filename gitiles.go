package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Gitiles struct {
	Requester *Requester
}

func (gs *Gitiles) SetBasicAuth(username, password string) {
	gs.Requester.SetAuth("basic", username, password)
}

func (gs *Gitiles) SetDigestAuth(username, password string) {
	gs.Requester.SetAuth("digest", username, password)
}

func (gs *Gitiles) SetCookieAuth(username, password string) {
	gs.Requester.SetAuth("cookie", username, password)
}

func NewGitilesClient(gitilesURL string, httpClient *http.Client) (*Gitiles, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second, // 设置超时时间
		}
	}

	r := &Requester{client: httpClient}

	if baseURL, err := SetBaseURL(gitilesURL); err != nil {
		return nil, err
	} else {
		r.baseURL = baseURL
	}

	gitiles := &Gitiles{Requester: r}

	return gitiles, nil
}

type GitilesCommitsOptions struct {
	// The n parameter can be used to limit the returned results.
	// If the n query parameter is supplied and additional changes exist that match the query beyond the end, the last change object has a _more_changes: true JSON field set.
	Limit int `url:"n,omitempty"`

	// The S or start query parameter can be supplied to skip a number of changes from the list.
	Start string `url:"s,omitempty"`
}

type GitilesPersonInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Time  string `json:"time"`
}

type GitilesDiffInfo struct {
	Type    string `json:"type"`
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	OldMode int    `json:"old_mode"`
	NewMode int    `json:"new_mode"`
	OldID   string `json:"old_id"`
	NewID   string `json:"new_id"`
}

type GitilesCommitInfo struct {
	Commit    string            `json:"commit"`
	Tree      string            `json:"tree"`
	Parents   []string          `json:"parents"`
	Author    GitilesPersonInfo `json:"author"`
	Committer GitilesPersonInfo `json:"committer"`
	Message   string            `json:"message"`
	TreeDiff  []GitilesDiffInfo `json:"tree_diff,omitempty"`
}

type GitilesLogs struct {
	Log      []GitilesCommitInfo `json:"log"`
	Previous string              `json:"previous,omitempty"`
	Next     string              `json:"next,omitempty"`
}

type GitilesRef struct {
	Value  string `json:"value"`
	Target string `json:"target,omitempty"`
	Peeled string `json:"peeled,omitempty"`
}

func (gs *Gitiles) GetCommit(ctx context.Context, project, commitID string) (*GitilesCommitInfo, *http.Response, error) {
	v := new(GitilesCommitInfo)
	u := fmt.Sprintf("%s/+/%s?format=JSON", project, commitID)

	resp, err := gs.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

func (gs *Gitiles) GetRefLogs(ctx context.Context, project, Ref string, opt *GitilesCommitsOptions) (*GitilesLogs, *http.Response, error) {
	v := new(GitilesLogs)
	u := fmt.Sprintf("%s/+log/%s/?format=JSON", project, Ref)

	resp, err := gs.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

func (gs *Gitiles) GetRefs(ctx context.Context, project string) (map[string]GitilesRef, *http.Response, error) {
	v := make(map[string]GitilesRef)
	u := fmt.Sprintf("%s/+refs?format=JSON", project)

	resp, err := gs.Requester.Call(ctx, "GET", u, nil, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// func (gs *Gitiles) DownloadFile(ctx context.Context, project, Ref, path string) (string, *http.Response, error) {
// 	v := new(string)
// 	u := fmt.Sprintf("%s/+/%s/%s?format=TEXT", project, Ref, path)

// 	resp, err := gs.Requester.Call(ctx, "GET", u, nil, v)
// 	if err != nil {
// 		return "", resp, err
// 	}
// 	return *v, resp, nil
// }
