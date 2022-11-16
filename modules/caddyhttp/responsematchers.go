package caddyhttp

import "net/http"

type ResponseMatcher struct {
	StatusCode []int       `json:"status_code,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}
