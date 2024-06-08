package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// GroupsService contains Group related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html
type GroupsService struct {
	gerrit *Gerrit
}

type Group struct {
	Raw    *GroupInfo
	gerrit *Gerrit
	Base   string
}

// GroupAuditEventInfo entity contains information about an audit event of a group.
type GroupAuditEventInfo struct {
	// TODO Member AccountInfo OR GroupInfo `json:"member"`
	Type string      `json:"type"`
	User AccountInfo `json:"user"`
	Date Timestamp   `json:"date"`
}

// GroupInfo entity contains information about a group.
// This can be a Gerrit internal group, or an external group that is known to Gerrit.
type GroupInfo struct {
	ID          string           `json:"id"`
	Name        string           `json:"name,omitempty"`
	URL         string           `json:"url,omitempty"`
	Options     GroupOptionsInfo `json:"options"`
	Description string           `json:"description,omitempty"`
	GroupID     int              `json:"group_id,omitempty"`
	Owner       string           `json:"owner,omitempty"`
	OwnerID     string           `json:"owner_id,omitempty"`
	CreatedOn   *Timestamp       `json:"created_on,omitempty"`
	MoreGroups  bool             `json:"_more_groups,omitempty"`
	Members     []AccountInfo    `json:"members,omitempty"`
	Includes    []GroupInfo      `json:"includes,omitempty"`
}

// GroupInput entity contains information for the creation of a new internal group.
type GroupInput struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	VisibleToAll bool   `json:"visible_to_all,omitempty"`
	OwnerID      string `json:"owner_id,omitempty"`
}

// GroupOptionsInfo entity contains options of the group.
type GroupOptionsInfo struct {
	VisibleToAll bool `json:"visible_to_all,omitempty"`
}

// GroupOptionsInput entity contains new options for a group.
type GroupOptionsInput struct {
	VisibleToAll bool `json:"visible_to_all,omitempty"`
}

// GroupsInput entity contains information about groups that should be included into a group or that should be deleted from a group.
type GroupsInput struct {
	OneGroup string   `json:"_one_group,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

// ListGroupsOptions specifies the different options for the ListGroups call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#list-groups
type ListGroupsOptions struct {
	// Group Options
	// Options fields can be obtained by adding o parameters, each option requires more lookups and slows down the query response time to the client so they are generally disabled by default.
	// Optional fields are:
	//        INCLUDES: include list of directly included groups.
	//        MEMBERS: include list of direct group members.
	Options []string `url:"o,omitempty"`

	// Check if a group is owned by the calling user
	// By setting the option owned and specifying a group to inspect with the option q, it is possible to find out, if this group is owned by the calling user.
	// If the group is owned by the calling user, the returned map contains this group. If the calling user doesn’t own this group an empty map is returned.
	Owned string `url:"owned,omitempty"`
	Group string `url:"q,omitempty"`

	// Group Limit
	// The /groups/ URL also accepts a limit integer in the n parameter. This limits the results to show n groups.
	Limit int `url:"n,omitempty"`

	// The /groups/ URL also accepts a start integer in the S parameter. The results will skip S groups from group list.
	Skip int `url:"S,omitempty"`

	// Limit the results to those groups that match the specified substring.
	Substring string `url:"m,omitempty"`

	// Limit the results to those groups that match the specified regex.
	Regex string `url:"r,omitempty"`
}

// List lists the groups accessible by the caller.
// This is the same as using the ls-groups command over SSH, and accepts the same options as query parameters.
// The entries in the map are sorted by group name.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#list-groups
func (s *GroupsService) List(ctx context.Context, opt *ListGroupsOptions) (map[string]GroupInfo, *http.Response, error) {
	v := make(map[string]GroupInfo)
	resp, err := s.gerrit.Requester.Call(ctx, "GET", "groups/", opt, &v)
	return v, resp, err
}

// Get retrieves a group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group
func (s *GroupsService) Get(ctx context.Context, groupID string) (*Group, *http.Response, error) {
	group := &Group{Raw: new(GroupInfo), gerrit: s.gerrit, Base: groupID}

	resp, err := group.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// Create creates a new Gerrit internal group.
// In the request body additional data for the group can be provided as GroupInput.
//
// As response the GroupInfo entity is returned that describes the created group.
// If the group creation fails because the name is already in use the response is “409 Conflict”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#create-group
func (s *GroupsService) Create(ctx context.Context, groupID string, input *GroupInput) (*Group, *http.Response, error) {
	obj := Group{Raw: new(GroupInfo), gerrit: s.gerrit, Base: groupID}
	return obj.Create(ctx, input)
}

func (g *Group) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s", g.Base)
	return g.gerrit.Requester.Call(ctx, "GET", u, nil, g.Raw)
}

func (g *Group) Create(ctx context.Context, input *GroupInput) (*Group, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}

	resp, err = g.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}

	g.Base = v.ID
	return g, resp, nil
}

// GetDetail retrieves a group with the direct members and the directly included groups.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-detail
func (g *Group) GetDetail(ctx context.Context) (*GroupInfo, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s/detail", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetName retrieves the name of a group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-name
func (g *Group) GetName(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("groups/%s/name", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// Rename renames a Gerrit internal group.
// The new group name must be provided in the request body.
//
// As response the new group name is returned.
// If renaming the group fails because the new name is already in use the response is “409 Conflict”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#rename-group
func (g *Group) Rename(ctx context.Context, name string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("groups/%s/name", g.Base)

	input := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// GetDescription retrieves the description of a group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-description
func (g *Group) GetDescription(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("groups/%s/description", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetDescription sets the description of a Gerrit internal group.
// The new group description must be provided in the request body.
//
// As response the new group description is returned.
// If the description was deleted the response is “204 No Content”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#set-group-description
func (g *Group) SetDescription(ctx context.Context, description string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("groups/%s/description", g.Base)

	input := struct {
		Description string `json:"description"`
	}{
		Description: description,
	}

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// DeleteDescription deletes the description of a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#delete-group-description
func (g *Group) DeleteDescription(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("groups/%s/description", g.Base)
	return g.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// GetOptions retrieves the options of a group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-options
func (g *Group) GetOptions(ctx context.Context) (*GroupOptionsInfo, *http.Response, error) {
	v := new(GroupOptionsInfo)
	u := fmt.Sprintf("groups/%s/options", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetOptions sets the options of a Gerrit internal group.
// The new group options must be provided in the request body as a GroupOptionsInput entity.
//
// As response the new group options are returned as a GroupOptionsInfo entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#set-group-options
func (g *Group) SetOptions(ctx context.Context, input *GroupOptionsInput) (*GroupOptionsInfo, *http.Response, error) {
	v := new(GroupOptionsInfo)
	u := fmt.Sprintf("groups/%s/options", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetOwner retrieves the owner group of a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-owner
func (g *Group) GetOwner(ctx context.Context) (*GroupInfo, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s/owner", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetOwner sets the owner group of a Gerrit internal group.
// The new owner group must be provided in the request body.
// The new owner can be specified by name, by group UUID or by the legacy numeric group ID.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#set-group-owner
func (g *Group) SetOwner(ctx context.Context, owner string) (*GroupInfo, *http.Response, error) {
	v := new(GroupInfo)
	u := fmt.Sprintf("groups/%s/owner", g.Base)

	input := struct {
		Owner string `json:"owner"`
	}{
		Owner: owner,
	}

	resp, err := g.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetAuditLog gets the audit log of a Gerrit internal group.
// The returned audit events are sorted by date in reverse order so that the newest audit event comes first.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-audit-log
func (g *Group) GetAuditLog(ctx context.Context) (*[]GroupAuditEventInfo, *http.Response, error) {
	v := new([]GroupAuditEventInfo)
	u := fmt.Sprintf("groups/%s/log.audit", g.Base)

	resp, err := g.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}