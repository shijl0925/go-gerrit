package gerrit

import (
	"net/http"
	"time"
)

type Gerrit struct {
	Requester *Requester

	Access   *AccessService
	Gitiles  *GitilesService
	Projects *ProjectService
	Changes  *ChangeService
	Accounts *AccountsService
	Groups   *GroupsService
	Config   *ConfigService
}

func NewClient(gerritURL string) (*Gerrit, error) {
	client := &http.Client{
		Timeout: 15 * time.Second, // 设置超时时间
	}

	r := &Requester{client: client}

	if baseURL, err := SetBaseURL(gerritURL); err != nil {
		return nil, err
	} else {
		r.baseURL = baseURL
	}

	gerrit := &Gerrit{Requester: r}

	gerrit.Access = &AccessService{gerrit: gerrit}
	gerrit.Gitiles = &GitilesService{gerrit: gerrit}
	gerrit.Projects = &ProjectService{gerrit: gerrit}
	gerrit.Changes = &ChangeService{gerrit: gerrit}
	gerrit.Accounts = &AccountsService{gerrit: gerrit}
	gerrit.Groups = &GroupsService{gerrit: gerrit}
	gerrit.Config = &ConfigService{gerrit: gerrit}

	return gerrit, nil
}

func (g *Gerrit) SetBasicAuth(username, password string) {
	g.Requester.authType = "basic"
	g.Requester.username = username
	g.Requester.password = password
}

func (g *Gerrit) SetDigestAuth(username, password string) {
	g.Requester.authType = "digest"
	g.Requester.username = username
	g.Requester.password = password
}

func (g *Gerrit) SetCookieAuth(username, password string) {
	g.Requester.authType = "cookie"
	g.Requester.username = username
	g.Requester.password = password
}
