package gerrit

import (
	"context"
	"net/http"
)

// AccessService contains Access Right related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html
type AccessService struct {
	gerrit *Gerrit
}

// AccessSectionInfo describes the access rights that are assigned on a ref.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html#access-section-info
type AccessSectionInfo struct {
	Permissions map[string]PermissionInfo `json:"permissions"`
}

// PermissionInfo entity contains information about an assigned permission.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html#permission-info
type PermissionInfo struct {
	Label     string                        `json:"label,omitempty"`
	Exclusive bool                          `json:"exclusive"`
	Rules     map[string]PermissionRuleInfo `json:"rules"`
}

// PermissionRuleInfo entity contains information about a permission rule that is assigned to group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html#permission-rule-info
type PermissionRuleInfo struct {
	// TODO Possible values for action: ALLOW, DENY or BLOCK, INTERACTIVE and BATCH
	Action string `json:"action"`
	Force  bool   `json:"force"`
	Min    int    `json:"min"`
	Max    int    `json:"max"`
}

// ProjectAccessInfo entity contains information about the access rights for a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html#project-access-info
type ProjectAccessInfo struct {
	Revision       string                       `json:"revision"`
	InheritsFrom   ProjectInfo                  `json:"inherits_from"`
	Local          map[string]AccessSectionInfo `json:"local"`
	IsOwner        bool                         `json:"is_owner"`
	OwnerOf        []string                     `json:"owner_of"`
	CanUpload      bool                         `json:"can_upload"`
	CanAdd         bool                         `json:"can_add"`
	CanAddTags     bool                         `json:"can_add_tags"`
	ConfigVisible  bool                         `json:"config_visible"`
	Groups         map[string]GroupInfo         `json:"groups"`
	ConfigWebLinks []string                     `json:"configWebLinks"`
}

// ListAccessRightsOptions specifies the parameters to the AccessService.ListAccessRights.
type ListAccessRightsOptions struct {
	// The projects for which the access rights should be returned must be specified as project options.
	// The project can be specified multiple times.
	Project []string `url:"project,omitempty"`
}

// ListAccessRights lists the access rights for projects.
// As result a map is returned that maps the project name to ProjectAccessInfo entities.
// The entries in the map are sorted by project name.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-access.html#list-access
func (s *AccessService) ListAccessRights(ctx context.Context, opt *ListAccessRightsOptions) (map[string]ProjectAccessInfo, *http.Response, error) {
	v := make(map[string]ProjectAccessInfo)
	u := "access/"

	resp, err := s.gerrit.Requester.Call(ctx, "GET", u, opt, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}