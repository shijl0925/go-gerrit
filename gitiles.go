package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

type GitilesService struct {
	gerrit *Gerrit
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

type GitilesCommits struct {
	Log      []GitilesCommitInfo `json:"log"`
	Previous string              `json:"previous,omitempty"`
	Next     string              `json:"next,omitempty"`
}

func (s *GitilesService) GetCommit(ctx context.Context, project, commitID string) (*GitilesCommitInfo, *http.Response, error) {
	v := new(GitilesCommitInfo)
	u := fmt.Sprintf("plugins/gitiles/%s/+/%s", project, commitID)

	resp, err := s.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

func (s *GitilesService) GetCommits(ctx context.Context, project, Ref string, opt *GitilesCommitsOptions) (*GitilesCommits, *http.Response, error) {
	v := new(GitilesCommits)
	u := fmt.Sprintf("plugins/gitiles/%s/+log/%s/", project, Ref)

	resp, err := s.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}