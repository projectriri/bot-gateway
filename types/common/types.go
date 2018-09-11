package common

import (
	"net/http"
)

type HTTPRequest struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}
