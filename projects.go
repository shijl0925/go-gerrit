package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Project struct {
	Raw      *ProjectInfo
	gerrit   *Gerrit
	Base     string
	Branches *BranchService
	Tags     *TagService
	Commits  *CommitService
}

// ProjectInfo entity contains information about a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#project-info
type ProjectInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Parent      string            `json:"parent,omitempty"`
	Description string            `json:"description,omitempty"`
	State       string            `json:"state,omitempty"`
	Branches    map[string]string `json:"branches,omitempty"`
	WebLinks    []WebLinkInfo     `json:"web_links,omitempty"`
}

// ProjectInput entity contains information for the creation of a new project.
type ProjectInput struct {
	Name                             string                       `json:"name,omitempty"`
	Parent                           string                       `json:"parent,omitempty"`
	Description                      string                       `json:"description,omitempty"`
	PermissionsOnly                  bool                         `json:"permissions_only"`
	CreateEmptyCommit                bool                         `json:"create_empty_commit"`
	SubmitType                       string                       `json:"submit_type,omitempty"`
	Branches                         []string                     `json:"branches,omitempty"`
	Owners                           []string                     `json:"owners,omitempty"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements,omitempty"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by,omitempty"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target,omitempty"`
	UseContentMerge                  string                       `json:"use_content_merge,omitempty"`
	RequireChangeID                  string                       `json:"require_change_id,omitempty"`
	MaxObjectSizeLimit               string                       `json:"max_object_size_limit,omitempty"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values,omitempty"`
}

// DeleteOptionsInfo entity contains information for the deletion of a project.
type DeleteOptionsInfo struct {
	Force    bool `json:"force"`
	Preserve bool `json:"preserve"`
}

// HeadInput entity contains information for setting HEAD for a project.
type HeadInput struct {
	Ref string `json:"ref"`
}

// ThemeInfo entity describes a theme.
type ThemeInfo struct {
	CSS    string `type:"css,omitempty"`
	Header string `type:"header,omitempty"`
	Footer string `type:"footer,omitempty"`
}

// ReflogEntryInfo entity describes an entry in a reflog.
type ReflogEntryInfo struct {
	OldID   string        `json:"old_id"`
	NewID   string        `json:"new_id"`
	Who     GitPersonInfo `json:"who"`
	Comment string        `json:"comment"`
}

// ProjectParentInput entity contains information for setting a project parent.
type ProjectParentInput struct {
	Parent        string `json:"parent"`
	CommitMessage string `json:"commit_message,omitempty"`
}

// InheritedBooleanInfo entity represents a boolean value that can also be inherited.
type InheritedBooleanInfo struct {
	Value           bool   `json:"value"`
	ConfiguredValue string `json:"configured_value"`
	InheritedValue  bool   `json:"inherited_value,omitempty"`
}

// MaxObjectSizeLimitInfo entity contains information about the max object size limit of a project.
type MaxObjectSizeLimitInfo struct {
	Value           string `json:"value,omitempty"`
	ConfiguredValue string `json:"configured_value,omitempty"`
	InheritedValue  string `json:"inherited_value,omitempty"`
}

// ConfigParameterInfo entity describes a project configuration parameter.
type ConfigParameterInfo struct {
	DisplayName string   `json:"display_name,omitempty"`
	Description string   `json:"description,omitempty"`
	Warning     string   `json:"warning,omitempty"`
	Type        string   `json:"type"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	// TODO: 5 fields are missing here, because the documentation seems to be fucked up
	// See https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#config-parameter-info
}

// ProjectDescriptionInput entity contains information for setting a project description.
type ProjectDescriptionInput struct {
	Description   string `json:"description,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
}

// ConfigInfo entity contains information about the effective project configuration.
type ConfigInfo struct {
	Description                      string                         `json:"description,omitempty"`
	UseContributorAgreements         InheritedBooleanInfo           `json:"use_contributor_agreements,omitempty"`
	UseContentMerge                  InheritedBooleanInfo           `json:"use_content_merge,omitempty"`
	UseSignedOffBy                   InheritedBooleanInfo           `json:"use_signed_off_by,omitempty"`
	CreateNewChangeForAllNotInTarget InheritedBooleanInfo           `json:"create_new_change_for_all_not_in_target,omitempty"`
	RequireChangeID                  InheritedBooleanInfo           `json:"require_change_id,omitempty"`
	EnableSignedPush                 InheritedBooleanInfo           `json:"enable_signed_push,omitempty"`
	MaxObjectSizeLimit               MaxObjectSizeLimitInfo         `json:"max_object_size_limit"`
	SubmitType                       string                         `json:"submit_type"`
	State                            string                         `json:"state,omitempty"`
	Commentlinks                     map[string]string              `json:"commentlinks"`
	Theme                            ThemeInfo                      `json:"theme,omitempty"`
	PluginConfig                     map[string]ConfigParameterInfo `json:"plugin_config,omitempty"`
	Actions                          map[string]ActionInfo          `json:"actions,omitempty"`
}

// ConfigInput entity describes a new project configuration.
type ConfigInput struct {
	Description                      string                       `json:"description,omitempty"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements,omitempty"`
	UseContentMerge                  string                       `json:"use_content_merge,omitempty"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by,omitempty"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target,omitempty"`
	EnableSignedPush                 string                       `json:"enable_signed_push,omitempty"`
	RequireSignedPush                string                       `json:"require_signed_push,omitempty"`
	RejectImplicitMerges             string                       `json:"reject_implicit_merges,omitempty"`
	RequireChangeID                  string                       `json:"require_change_id,omitempty"`
	MaxObjectSizeLimit               interface{}                  `json:"max_object_size_limit,omitempty"`
	SubmitType                       string                       `json:"submit_type,omitempty"`
	State                            string                       `json:"state,omitempty"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values,omitempty"`
}

type ProjectBaseOptions struct {
	// Limit the number of projects to be included in the results.
	Limit int `url:"n,omitempty"`

	// Skip the given number of projects from the beginning of the list.
	Skip int `url:"S,omitempty"`
}

type ProjectOptions struct {
	ProjectBaseOptions

	// Limit the results to the projects having the specified branch and include the sha1 of the branch in the results.
	Branch string `url:"b,omitempty"`

	// Include project description in the results.
	Description bool `url:"d,omitempty"`

	// Limit the results to those projects that start with the specified prefix.
	Prefix string `url:"p,omitempty"`

	// Limit the results to those projects that match the specified regex.
	// Boundary matchers '^' and '$' are implicit.
	// For example: the regex 'test.*' will match any projects that start with 'test' and regex '.*test' will match any project that end with 'test'.
	Regex string `url:"r,omitempty"`

	// Skip the given number of projects from the beginning of the list.
	// Skip string `url:"S,omitempty"`

	// Limit the results to those projects that match the specified substring.
	Substring string `url:"m,omitempty"`

	// Get projects inheritance in a tree-like format.
	// This option does not work together with the branch option.
	Tree bool `url:"t,omitempty"`

	// Get projects with specified type: ALL, CODE, PERMISSIONS.
	Type string `url:"type,omitempty"`

	// Get all projects with the given state.
	State string `url:"state,omitempty"`
}

type ProjectService struct {
	gerrit *Gerrit
}

func NewProject(gerrit *Gerrit, projectName string) *Project {
	obj := &Project{
		Raw:    new(ProjectInfo),
		gerrit: gerrit,
		Base:   projectName,
	}

	obj.Branches = &BranchService{gerrit: gerrit, project: obj}
	obj.Tags = &TagService{gerrit: gerrit, project: obj}
	obj.Commits = &CommitService{gerrit: gerrit, project: obj}

	return obj
}

// List lists the projects accessible by the caller.
// This is the same as using the ls-projects command over SSH, and accepts the same options as query parameters.
// The entries in the map are sorted by project name.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-projects
func (s *ProjectService) List(ctx context.Context, opt *ProjectOptions) (map[string]ProjectInfo, *http.Response, error) {
	v := make(map[string]ProjectInfo)
	resp, err := s.gerrit.Requester.Call(ctx, "GET", "projects/", opt, &v)

	return v, resp, err
}

// Get retrieves a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-project
func (s *ProjectService) Get(ctx context.Context, projectName string) (*Project, *http.Response, error) {
	project := NewProject(s.gerrit, projectName)

	resp, err := project.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// Create creates a new project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-project
func (s *ProjectService) Create(ctx context.Context, projectName string, input *ProjectInput) (*Project, *http.Response, error) {
	obj := NewProject(s.gerrit, projectName)
	return obj.Create(ctx, input)
}

// Delete deletes a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-project
func (s *ProjectService) Delete(ctx context.Context, projectName string, input *DeleteOptionsInfo) (bool, *http.Response, error) {
	obj := NewProject(s.gerrit, projectName)
	return obj.Delete(ctx, input)
}

func (p *Project) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("projects/%s", url.QueryEscape(p.Base))
	return p.gerrit.Requester.Call(ctx, "GET", u, nil, p.Raw)
}

func (p *Project) Create(ctx context.Context, input *ProjectInput) (*Project, *http.Response, error) {
	u := fmt.Sprintf("projects/%s", url.QueryEscape(p.Base))
	resp, err := p.gerrit.Requester.Call(ctx, "PUT", u, input, nil)

	if err != nil {
		return nil, resp, err
	}

	resp, err = p.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

func (p *Project) Delete(ctx context.Context, input *DeleteOptionsInfo) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/delete-project~delete", url.QueryEscape(p.Base))
	resp, err := p.gerrit.Requester.Call(ctx, "POST", u, input, nil)

	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}

// GetDescription retrieves the description of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-project-description
func (p *Project) GetDescription(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/description", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// SetDescription sets the description of a project.
// The new project description must be provided in the request body inside a ProjectDescriptionInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-project-description
func (p *Project) SetDescription(ctx context.Context, input *ProjectDescriptionInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/description", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// DeleteDescription deletes the description of a project.
// The request body does not need to include a ProjectDescriptionInput entity if no commit message is specified.
//
// Please note that some proxies prohibit request bodies for DELETE requests.
// In this case, if you want to specify a commit message, use PUT to delete the description.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-project-description
func (p *Project) DeleteDescription(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/description", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// GetParent retrieves the name of a project’s parent project.
// For the All-Projects root project an empty string is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-project-parent
func (p *Project) GetParent(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/parent", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// SetParent sets the parent project for a project.
// The new name of the parent project must be provided in the request body inside a ProjectParentInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-project-parent
func (p *Project) SetParent(ctx context.Context, input *ProjectParentInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/parent", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// GetHEAD retrieves for a project the name of the branch to which HEAD points.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-head
func (p *Project) GetHEAD(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/HEAD", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// SetHEAD sets HEAD for a project.
// The new ref to which HEAD should point must be provided in the request body inside a HeadInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-head
func (p *Project) SetHEAD(ctx context.Context, input *HeadInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("projects/%s/HEAD", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// GetConfig gets some configuration information about a project.
// Note that this config info is not simply the contents of project.config;
// it generally contains fields that may have been inherited from parent projects.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-config
func (p *Project) GetConfig(ctx context.Context) (*ConfigInfo, *http.Response, error) {
	v := new(ConfigInfo)
	u := fmt.Sprintf("projects/%s/config", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// SetConfig sets the configuration of a project.
// The new configuration must be provided in the request body as a ConfigInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-config
func (p *Project) SetConfig(ctx context.Context, input *ConfigInput) (*ConfigInfo, *http.Response, error) {
	v := new(ConfigInfo)
	u := fmt.Sprintf("projects/%s/config", url.QueryEscape(p.Base))

	resp, err := p.gerrit.Requester.Call(ctx, "PUT", u, input, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}
