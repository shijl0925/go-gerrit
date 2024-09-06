package gerrit

import (
	"bytes"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
		errorResponse.Message = string(data)
	}

	return errorResponse
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()

	return origURL.String(), nil
}
