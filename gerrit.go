package gerrit

import (
	"log"
	"net/http"
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

func NewClient(gerritURL string, httpClient *http.Client) (*Gerrit, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	r := &Requester{client: httpClient}

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

// SetAuth 用于设置不同类型的认证方式。
// authType: 认证类型，可以是 "basic"、"digest" 或 "cookie"。
// username: 用户名。
// password: 密码。
func (g *Gerrit) SetAuth(authType, username, password string) {
	// 参数验证
	if authType == "" || username == "" || password == "" {
		// 根据实际情况，这里可以记录日志、抛出异常或返回错误
		log.Fatal("authType, username, and password cannot be empty")
		return
	}

	// 对authType值进行校验，确保其为允许的值之一
	allowedAuthTypes := []string{"basic", "digest", "cookie"}

	found := false
	for _, allowedType := range allowedAuthTypes {
		if authType == allowedType {
			found = true
			break
		}
	}
	if !found {
		// 根据实际情况，这里可以记录日志、抛出异常或返回错误
		log.Fatal("Unsupported authType")
		return
	}

	g.Requester.authType = authType
	g.Requester.username = username
	g.Requester.password = password
}

func (g *Gerrit) SetBasicAuth(username, password string) {
	g.SetAuth("basic", username, password)
}

func (g *Gerrit) SetDigestAuth(username, password string) {
	g.SetAuth("digest", username, password)
}

func (g *Gerrit) SetCookieAuth(username, password string) {
	g.SetAuth("cookie", username, password)
}
