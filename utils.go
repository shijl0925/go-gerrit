package gerrit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func SetBaseURL(urlStr string) (*url.URL, error) {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Make sure the given URL end with a slash
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	// Update the base URL of the client.
	return baseURL, nil
}

var magicPrefix = []byte(")]}'\n")

func RemoveMagicPrefixLine(body []byte) []byte {
	if bytes.HasPrefix(body, magicPrefix) {
		body = body[5:]
	}

	if bytes.Equal(body, []byte("ok")) || bytes.Equal(body, []byte("ok\n")) {
		return []byte("\"ok\"\n")
	}

	return body
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: %s", data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

func addOptions(s string, opt interface{}) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}