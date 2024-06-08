package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Tag struct {
	Raw     *TagInfo
	project *Project
	gerrit  *Gerrit
	Base    string
}

type TagInput struct {
	//Ref      string `json:"ref"`
	Revision string `json:"revision,omitempty"`
	Message  string `json:"message,omitempty"`
}

type TagInfo struct {
	Ref       string        `json:"ref"`
	Revision  string        `json:"revision"`
	Object    string        `json:"object"`
	Message   string        `json:"message"`
	Tagger    GitPersonInfo `json:"tagger"`
	Created   *Timestamp    `json:"created,omitempty"`
	CanDelete bool          `json:"can_delete"`
	WebLinks  []WebLinkInfo `json:"web_links,omitempty"`
}

// DeleteTagsInput entity for delete tags.
type DeleteTagsInput struct {
	Tags []string `json:"tags"`
}

// TagOptions specifies the parameters to the tag API endpoints.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#tag-options
type TagOptions struct {
	// Limit the number of tags to be included in the results.
	Limit int `url:"n,omitempty"`

	// Skip the given number of tags from the beginning of the list.
	Skip int `url:"S,omitempty"`

	// Substring limits the results to those projects that match the specified substring.
	Substring string `url:"m,omitempty"`

	// Limit the results to those tags that match the specified regex.
	// Boundary matchers '^' and '$' are implicit.
	// For example: the regex 't*' will match any tags that start with 'test' and regex '*t' will match any tags that end with 'test'.
	Regex string `url:"r,omitempty"`
}

type TagService struct {
	gerrit  *Gerrit
	project *Project
}

type ITagService interface {
	List(ctx context.Context, opt *TagOptions) (*[]TagInfo, *http.Response, error)
	Get(ctx context.Context, tagID string) (*Tag, *http.Response, error)
	Create(ctx context.Context, tagID string, input *TagInput) (*Tag, *http.Response, error)
	Delete(ctx context.Context, tagID string) (bool, *http.Response, error)
	BulkDelete(ctx context.Context, input *DeleteTagsInput) (bool, *http.Response, error)
}

//func NewTagService() ITagService {
//        return &TagService{}
//}

//func NewTag(gerrit *Gerrit, projectName, tagID string) *Tag {
//        return &Tag{
//                Raw:     new(TagInfo),
//                gerrit:  gerrit,
//                project: &Project{Base: projectName, gerrit: gerrit},
//                Base:    tagID,
//        }
//}

//func NewTagService(gerrit *Gerrit, projectName string) *TagService {
//        return &TagService{
//                gerrit:  gerrit,
//                project: &Project{Base: projectName, gerrit: gerrit},
//        }
//}

// List lists the tags of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-tags
func (s *TagService) List(ctx context.Context, opt *TagOptions) (*[]TagInfo, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/tags/", url.QueryEscape(s.project.Base))

	v := &[]TagInfo{}
	resp, err := s.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	return v, resp, err
}

// Get retrieves a tag of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-tag
func (s *TagService) Get(ctx context.Context, tagID string) (*Tag, *http.Response, error) {
	tag := Tag{Raw: new(TagInfo), gerrit: s.gerrit, project: s.project, Base: tagID}

	resp, err := tag.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return &tag, resp, nil
}

// Create create a tag of a project
//
// Gerrit API docs:https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-tag
func (s *TagService) Create(ctx context.Context, tagID string, input *TagInput) (*Tag, *http.Response, error) {
	obj := Tag{Raw: new(TagInfo), gerrit: s.gerrit, project: s.project, Base: tagID}
	return obj.Create(ctx, input)
}

// Delete delete a tag of a project
//
// Gerrit API docs:https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-tag
func (s *TagService) Delete(ctx context.Context, tagID string) (bool, *http.Response, error) {
	obj := Tag{Raw: new(TagInfo), gerrit: s.gerrit, project: s.project, Base: tagID}
	return obj.Delete(ctx)
}

// BulkDelete delete one or more tags.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-tags
func (s *TagService) BulkDelete(ctx context.Context, input *DeleteTagsInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/tags:delete", url.QueryEscape(s.project.Base))
	resp, err := s.gerrit.Requester.Call(ctx, "POST", u, input, nil)

	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}

func (t *Tag) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(t.project.Base), url.QueryEscape(t.Base))
	return t.gerrit.Requester.Call(ctx, "GET", u, nil, t.Raw)
}

func (t *Tag) Create(ctx context.Context, input *TagInput) (*Tag, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(t.project.Base), url.QueryEscape(t.Base))
	resp, err := t.gerrit.Requester.Call(ctx, "PUT", u, input, nil)

	if err != nil {
		return nil, resp, err
	}

	resp, err = t.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}
	return t, resp, nil
}

func (t *Tag) Delete(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(t.project.Base), url.QueryEscape(t.Base))
	resp, err := t.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)

	if err != nil {
		return false, resp, err
	}

	return true, resp, nil
}