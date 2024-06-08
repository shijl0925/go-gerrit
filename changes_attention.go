package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// AttentionSetInfo entity contains details of users that are in the attention set.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#attention-set-info
type AttentionSetInfo struct {
	// AccountInfo entity.
	Account AccountInfo `json:"account"`
	// The timestamp of the last update.
	LastUpdate Timestamp `json:"last_update"`
	// The reason of for adding or removing the user.
	Reason string `json:"reason"`
}

// Doc: https://gerrit-review.googlesource.com/Documentation/user-notify.html#recipient-types
type RecipientType string

// AttentionSetInput entity contains details for adding users to the attention
// set and removing them from it.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#attention-set-input
type AttentionSetInput struct {
	User          string                       `json:"user,omitempty"`
	Reason        string                       `json:"reason"`
	Notify        string                       `json:"notify,omitempty"`
	NotifyDetails map[RecipientType]NotifyInfo `json:"notify_details,omitempty"`
}

// GetAttentionSet returns all users that are currently in the attention set. As response a list of AttentionSetInfo entity is returned.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-attention-set
func (c *Change) GetAttentionSet(ctx context.Context) (*[]AttentionSetInfo, *http.Response, error) {
	v := new([]AttentionSetInfo)
	u := fmt.Sprintf("changes/%s/attention", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddAttention adds a single user to the attention set of a change.
// AttentionSetInput.Input must be provided
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#add-to-attention-set
func (c *Change) AddAttention(ctx context.Context, input *AttentionSetInput) (*AccountInfo, *http.Response, error) {
	v := new(AccountInfo)
	u := fmt.Sprintf("changes/%s/attention", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// RemoveAttention deletes a single user from the attention set of a change.
// AttentionSetInput.Input must be provided
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#remove-from-attention-set
func (c *Change) RemoveAttention(ctx context.Context, accountID string, input *AttentionSetInput) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/attention/%s/delete", c.Base, accountID)
	return c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}