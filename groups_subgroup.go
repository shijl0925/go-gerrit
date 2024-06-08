package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// ListSubgroups lists the directly subgroups of a group.
// The entries in the list are sorted by group name and UUID.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#list-subgroups
func (g *Group) ListSubgroups(ctx context.Context) (*[]GroupInfo, *http.Response, error) {
	v := new([]GroupInfo)
	u := fmt.Sprintf("groups/%s/groups/", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetSubGroup retrieves a subgroup.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-subgroup
func (g *Group) GetSubGroup(ctx context.Context, groupID string) (*GroupInfo, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s/groups/%s", g.Base, groupID)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddSubgroup adds an internal or external group as subgroup to a Gerrit internal group
// External groups must be specified using the UUID.
//
// As response a GroupInfo entity is returned that describes the subgroup.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#add-subgroup
func (g *Group) AddSubgroup(ctx context.Context, groupID string) (*GroupInfo, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s/groups/%s", g.Base, groupID)

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddSubgroups adds one or several groups as subgroups to a Gerrit internal group.
// The subgroups to be added must be provided in the request body as a GroupsInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#add-subgroups
func (g *Group) AddSubgroups(ctx context.Context, input *GroupsInput) (*[]GroupInfo, *http.Response, error) {
	v := new([]GroupInfo)
	u := fmt.Sprintf("groups/%s/groups", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "POST", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// RemoveSubgroup removes a subgroup from a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#remove-subgroup
func (g *Group) RemoveSubgroup(ctx context.Context, groupID string) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s/groups/%s", g.Base, groupID)
	return g.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// RemoveSubgroups removes one or several subgroups from a Gerrit internal group.
// The groups to be deleted from the group must be provided in the request body as a GroupsInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#remove-subgroup
func (g *Group) RemoveSubgroups(ctx context.Context, groupID string, input *GroupsInput) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s/groups.delete", groupID)
	return g.gerrit.Requester.Call(ctx, "POST", u, input, nil)
}