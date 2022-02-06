package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

// NewPostRequestWithBody ...
func NewPostRequestWithBody(endpoint string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	switch t := body.(type) {
	case []byte:
		req = httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(t))
	case string:
		req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(t))
	default:
		bodyBytes, err := json.Marshal(t)
		if err != nil {
			panic(err)
		}
		req = httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	rec := httptest.NewRecorder()
	return req, rec
}

// NewGetRequest ...
func NewGetRequest(endpoint string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	rec := httptest.NewRecorder()
	return req, rec
}

// NewDeleteRequest ...
func NewDeleteRequest(endpoint string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	rec := httptest.NewRecorder()
	return req, rec
}
