package gerrit

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var transport = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	DisableCompression:  true,
	MaxIdleConns:        100, // 调整为100以允许更多的连接复用
	IdleConnTimeout:     60 * time.Second,
	TLSHandshakeTimeout: 15 * time.Second,
	TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS10, InsecureSkipVerify: false},
}

// DefaultClient is the default http client for got requests.
var DefaultClient = &http.Client{
	Transport: transport,
	Timeout:   30 * time.Second,
}

// 定义认证类型的常量
const (
	AuthTypeCookie = "cookie"
	AuthTypeDigest = "digest"
	AuthTypeBasic  = "basic"
)

// AuthMethod 接口定义了各种认证方法
type AuthMethod interface {
	ApplyAuthentication(req *http.Request)
}

// CookieAuth 实现了基于Cookie的认证
type CookieAuth struct {
	Username string
	Password string
}

func (c *CookieAuth) ApplyAuthentication(req *http.Request) {
	// 注意：在生产环境中，应确保使用HTTPS和设置HttpOnly、Secure属性
	req.AddCookie(&http.Cookie{
		Name:  c.Username,
		Value: c.Password,
	})
}

// BasicAuth 实现了基本认证
type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) ApplyAuthentication(req *http.Request) {
	req.SetBasicAuth(b.Username, b.Password)
}

// DigestAuth 实现了摘要认证（示例代码，需要根据实际需求完成具体实现）
type DigestAuth struct {
	Username string
	Password string
}

func (d *DigestAuth) ApplyAuthentication(req *http.Request) {
	// TODO: 实现摘要认证逻辑
}

type Requester struct {
	// client is the HTTP client used to communicate with the API.
	client *http.Client

	// baseURL is the base URL of the Gerrit instance for API requests.
	baseURL *url.URL

	// Gerrit service for authentication.
	username, password, authType string
}

func (r *Requester) NewRequest(ctx context.Context, method, endpoint string, opt interface{}) (*http.Request, error) {
	hasAuth := false

	if len(r.authType) != 0 && len(r.username) != 0 && len(r.password) != 0 {
		hasAuth = true
	}

	// If there is a "/" at the start, remove it.
	urlStr := strings.TrimPrefix(endpoint, "/")

	baseURL := r.baseURL.String()

	// If we are authenticated, let's apply the "a/" prefix,
	if hasAuth {
		u, _ := url.Parse(baseURL)
		baseURL = fmt.Sprintf("%s://%s/a%s", u.Scheme, u.Host, u.Path)
	}

	urlStr = baseURL + urlStr

	if method == http.MethodGet {
		u, err := addOptions(urlStr, opt)
		if err != nil {
			return nil, err
		}
		urlStr = u
	}

	//log.Printf("Requesting %s %s", method, urlStr)

	req, err := http.NewRequestWithContext(ctx, method, urlStr, nil)

	if err != nil {
		return nil, err
	}

	if opt != nil && (method == http.MethodPost || method == http.MethodPut) {
		if reflect.TypeOf(opt).String() == "string" {
			req.Body = io.NopCloser(bytes.NewBuffer([]byte(opt.(string))))

			req.Header.Add("Content-Type", "plain/text;charset=UTF-8")
		} else {
			buf, err := json.Marshal(opt)
			//log.Printf("buf: %+v", buf)
			if err != nil {
				return nil, err
			}
			req.Body = io.NopCloser(bytes.NewBuffer(buf))

			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Apply Authentication
	if hasAuth {
		switch r.authType {
		case AuthTypeCookie:
			cookieAuth := &CookieAuth{
				Username: r.username,
				Password: r.password,
			}
			cookieAuth.ApplyAuthentication(req)

		case AuthTypeDigest:
			digestAuth := &DigestAuth{
				Username: r.username,
				Password: r.password,
			}
			digestAuth.ApplyAuthentication(req)

		default:
			basicAuth := &BasicAuth{
				Username: r.username,
				Password: r.password,
			}
			basicAuth.ApplyAuthentication(req)
		}
	}

	// Request compact JSON
	// See https://gerrit-review.googlesource.com/Documentation/rest-api.html#output
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (r *Requester) Do(req *http.Request, v interface{}) (*http.Response, error) {
	isText := false
	if _, ok := v.(*string); ok {
		req.Header.Set("Accept", "text/plain")
		isText = true
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return resp, err
	}

	err = CheckResponse(resp)

	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further.
		return resp, err
	}

	if v != nil {
		defer resp.Body.Close()

		if w, ok := v.(io.Writer); ok {
			if _, err := io.Copy(w, resp.Body); err != nil {
				return nil, err
			}
		} else {
			if isText {
				var body []byte
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					// even though there was an error, we still return the response
					// in case the caller wants to inspect it further
					return resp, err
				}
				body = RemoveMagicPrefixLine(body)

				w := v.(*string)
				*w = strings.Trim(string(body), "\"\n")

			} else {
				if _, err := io.CopyN(io.Discard, resp.Body, 5); err != nil {
					return resp, err
				}
				if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
					return resp, err
				}
			}
		}
	}

	return resp, err
}

func (r *Requester) Call(ctx context.Context, method, u string, opt interface{}, v interface{}) (*http.Response, error) {
	req, err := r.NewRequest(ctx, method, u, opt)
	if err != nil {
		return nil, err
	}

	resp, err := r.Do(req, v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetAuth 用于设置不同类型的认证方式。
// authType: 认证类型，可以是 "basic"、"digest" 或 "cookie"。
// username: 用户名。
// password: 密码。
func (r *Requester) SetAuth(authType, username, password string) {
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

	r.authType = authType
	r.username = username
	r.password = password
}
