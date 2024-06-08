package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// ListGroupMembersOptions specifies the different options for the ListGroupMembers call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#group-members
type ListGroupMembersOptions struct {
	// To resolve the included groups of a group recursively and to list all members the parameter recursive can be set.
	// Members from included external groups and from included groups which are not visible to the calling user are ignored.
	Recursive bool `url:"recursive,omitempty"`
}

// MembersInput entity contains information about accounts that should be added as members to a group or that should be deleted from the group
type MembersInput struct {
	OneMember string   `json:"_one_member,omitempty"`
	Members   []string `json:"members,omitempty"`
}

// ListMembers lists the direct members of a Gerrit internal group.
// The entries in the list are sorted by full name, preferred email and id.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#group-members
func (g *Group) ListMembers(ctx context.Context, opt *ListGroupMembersOptions) (*[]AccountInfo, *http.Response, error) {
	v := new([]AccountInfo)
	u := fmt.Sprintf("groups/%s/members/", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetMember retrieves a group member.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-member
func (g *Group) GetMember(ctx context.Context, accountID string) (*AccountInfo, *http.Response, error) {
	v := new(AccountInfo)
	u := fmt.Sprintf("groups/%s/members/%s", g.Base, accountID)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddMember adds a user as member to a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#add-group-member
func (g *Group) AddMember(ctx context.Context, accountID string) (*AccountInfo, *http.Response, error) {
	v := new(AccountInfo)
	u := fmt.Sprintf("groups/%s/members/%s", g.Base, accountID)

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddMembers adds one or several users to a Gerrit internal group.
// The users to be added to the group must be provided in the request body as a MembersInput entity.
//
// As response a list of detailed AccountInfo entities is returned that describes the group members that were specified in the MembersInput.
// An AccountInfo entity is returned for each user specified in the input, independently of whether the user was newly added to the group or whether the user was already a member of the group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#_add_group_members
func (g *Group) AddMembers(ctx context.Context, input *MembersInput) (*[]AccountInfo, *http.Response, error) {
	v := new([]AccountInfo)
	u := fmt.Sprintf("groups/%s/members", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteMember deletes a user from a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#delete-group-member
func (g *Group) DeleteMember(ctx context.Context, accountID string) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s/members/%s", g.Base, accountID)
	return g.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// DeleteMembers delete one or several users from a Gerrit internal group.
// The users to be deleted from the group must be provided in the request body as a MembersInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#delete-group-members
func (g *Group) DeleteMembers(ctx context.Context, input *MembersInput) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s/members.delete", g.Base)
	return g.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}