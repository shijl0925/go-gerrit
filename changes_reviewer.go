package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// ReviewerInfo entity contains information about a reviewer and its votes on a change.
type ReviewerInfo struct {
	AccountInfo
	Approvals map[string]string `json:"approvals"`
}

//type SuggestedReviewerQueryOptions struct {
//        QueryOptions
//}

// SuggestedReviewerInfo entity contains information about a reviewer that can be added to a change (an account or a group).
type SuggestedReviewerInfo struct {
	Account AccountInfo   `json:"account,omitempty"`
	Group   GroupBaseInfo `json:"group,omitempty"`
}

// ReviewerResult entity describes the result of adding a reviewer to a change.
type ReviewerResult struct {
	Input     string         `json:"input,omitempty"`
	Reviewers []ReviewerInfo `json:"reviewers,omitempty"`
	CCS       []AccountInfo  `json:"ccs,omitempty"`
	Removed   []AccountInfo  `json:"removed,omitempty"`
	Error     string         `json:"error,omitempty"`
	Confirm   bool           `json:"confirm,omitempty"`
}

// DeleteVoteInput entity contains options for the deletion of a vote.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-vote-input
type DeleteVoteInput struct {
	Label         string                `json:"label,omitempty"`
	Notify        string                `json:"notify,omitempty"`
	NotifyDetails map[string]NotifyInfo `json:"notify_details"`
}

// ListReviewers lists the reviewers of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-reviewers
func (c *Change) ListReviewers(ctx context.Context) (*[]ReviewerInfo, *http.Response, error) {
	v := new([]ReviewerInfo)
	u := fmt.Sprintf("changes/%s/reviewers/", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// SuggestReviewers suggest the reviewers for a given query q and result limit n.
// If result limit is not passed, then the default 10 is used.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#suggest-reviewers
func (c *Change) SuggestReviewers(ctx context.Context, opt *QueryOptions) (*[]SuggestedReviewerInfo, *http.Response, error) {
	v := new([]SuggestedReviewerInfo)
	u := fmt.Sprintf("changes/%s/suggest_reviewers", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetReviewer retrieves a reviewer of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-reviewer
func (c *Change) GetReviewer(ctx context.Context, accountID string) (*[]ReviewerInfo, *http.Response, error) {
	v := new([]ReviewerInfo)
	u := fmt.Sprintf("changes/%s/reviewers/%s", c.Base, accountID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddReviewer adds one user or all members of one group as reviewer to the change.
// The reviewer to be added to the change must be provided in the request body as a ReviewerInput entity.
//
// As response an ReviewerResult entity is returned that describes the newly added reviewers.
// If a group is specified, adding the group members as reviewers is an atomic operation.
// This means if an error is returned, none of the members are added as reviewer.
// If a group with many members is added as reviewer a confirmation may be required.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#add-reviewer
func (c *Change) AddReviewer(ctx context.Context, input *ReviewerInput) (*ReviewerResult, *http.Response, error) {
	v := new(ReviewerResult)
	u := fmt.Sprintf("changes/%s/reviewers", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteReviewer deletes a reviewer from a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-reviewer
func (c *Change) DeleteReviewer(ctx context.Context, accountID string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/reviewers/%s", c.Base, accountID)
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// ListVotes lists the votes for a specific reviewer of the change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-votes
func (c *Change) ListVotes(ctx context.Context, accountID string) (map[string]int, *http.Response, error) {
	v := make(map[string]int)
	u := fmt.Sprintf("changes/%s/reviewers/%s/votes/", c.Base, accountID)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteVote deletes a single vote from a change. Note, that even when the
// last vote of a reviewer is removed the reviewer itself is still listed on
// the change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-vote
func (c *Change) DeleteVote(ctx context.Context, accountID string, label string) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/reviewers/%s/votes/%s'", c.Base, accountID, label)
	return c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}