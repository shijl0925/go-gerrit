package gerrit

import (
	"bytes"
	"context"
	"crypto/tls"
	"time"

	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var transport = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	DisableCompression:  true,
	MaxIdleConns:        10,
	IdleConnTimeout:     30 * time.Second,
	TLSHandshakeTimeout: 10 * time.Second,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
}

// defaultClient is the default http client for got requests.
var defaultClient = &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second,
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

	// If we are authenticated, let's apply the "a/" prefix,
	if hasAuth {
		urlStr = "a/" + urlStr
	}

	urlStr = r.baseURL.String() + urlStr

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
		case "cookie":
			req.AddCookie(&http.Cookie{Name: r.username, Value: r.password})
		case "digest":
			// todo
		default:
			req.SetBasicAuth(r.username, r.password)
		}
	}

	// Request compact JSON
	// See https://gerrit-review.googlesource.com/Documentation/rest-api.html#output
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (r *Requester) Do(req *http.Request, v interface{}) (*http.Response, error) {
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
			var body []byte
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				// even though there was an error, we still return the response
				// in case the caller wants to inspect it further
				return resp, err
			}
			body = RemoveMagicPrefixLine(body)
			//log.Println(string(body))

			err = json.Unmarshal(body, v)
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
